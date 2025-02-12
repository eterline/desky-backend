package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	sshlander "github.com/eterline/desky-backend/internal/services/ssh-lander"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-ping/ping"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type SSHRepository interface {
	AddHost(username string, host string, port uint16, osType string, privateKeyUsage bool, password string, key string) error
	Delete(id int) error
	QueryAll() ([]models.SSHCredentialsT, error)
	QueryById(id int) (*models.SSHCredentialsT, error)
}

type SSHLanderControllers struct {
	ctx context.Context

	term    sshlander.TerminalType
	repoSSH SSHRepository
	websock *websocket.Upgrader

	logging *logrus.Logger

	sshMu sync.Mutex
}

func Init(ctx context.Context, repo SSHRepository) *SSHLanderControllers {
	return &SSHLanderControllers{
		ctx:     ctx,
		repoSSH: repo,
		websock: &websocket.Upgrader{
			HandshakeTimeout:  10 * time.Second,
			EnableCompression: true,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
		},

		term: sshlander.XtermColored,

		logging: logger.ReturnEntry().Logger,
	}
}

// ============================================================================================

func (mc *SSHLanderControllers) ListHosts(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.list-hosts"

	sshList, err := mc.repoSSH.QueryAll()
	if err != nil {
		return op, err
	}

	if handler.ListIsEmpty(w, sshList) {
		return op, nil
	}

	resultList := make([]models.SSHInstanceObject, len(sshList))

	for idx, sshInst := range sshList {

		hostString := fmt.Sprintf(
			"%s@%s:%v",
			sshInst.Username,
			sshInst.Host,
			sshInst.Port,
		)

		resultList[idx] = models.SSHInstanceObject{
			ID:            int(sshInst.ID),
			HostString:    hostString,
			PrivateKeyUse: sshInst.Security.PrivateKeyUse,
		}
	}

	return op, handler.WriteJSON(w, http.StatusOK, resultList)
}

func (mc *SSHLanderControllers) AppendHost(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.append-host"

	form := new(models.RequestFormSSH)
	if err := handler.DecodeRequest(r, form); err != nil {
		return op, err
	}

	if err := handler.Validate(form); err != nil {
		return op, err
	}

	if err := mc.repoSSH.AddHost(
		form.User, form.Host, form.Port,
		form.System,
		form.PrivateKeyUse,
		form.Password, form.PrivateKey,
	); err != nil {
		return op, err
	}

	response := models.ResponseCreateSSH{
		PrivateKeyUse: form.PrivateKeyUse,
		Target:        fmt.Sprintf("%s@%s:%v", form.User, form.Host, form.Port),
	}

	return op, handler.WriteJSON(w, http.StatusCreated, response)
}

func (mc *SSHLanderControllers) DeleteHost(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.delete-host"

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	if err := mc.repoSSH.Delete(q.GetInt("id")); err != nil {
		return op, err
	}

	return op, handler.StatusOK(w, "host deleted")
}

// ============================================================================================

func (mc *SSHLanderControllers) TestHosts(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.test-hosts"

	sshList, err := mc.repoSSH.QueryAll()
	if err != nil {
		return op, err
	}

	if handler.ListIsEmpty(w, sshList) {
		return op, nil
	}

	result := getPingedList(sshList)

	return op, handler.WriteJSON(w, http.StatusOK, result)
}

func getPingedList(hostsData []models.SSHCredentialsT) []models.SSHTestObject {
	var wg sync.WaitGroup
	testedList := make([]models.SSHTestObject, len(hostsData))

	for idx, credential := range hostsData {

		wg.Add(1)

		go func() {
			var attempt bool

			defer func() {
				testedList[idx] = models.SSHTestObject{
					ID:        int(credential.ID),
					Available: attempt,
				}
				wg.Done()
			}()

			pinger, err := ping.NewPinger(credential.Host)
			if err != nil {
				attempt = false
				return
			}

			pinger.Count = 1
			attempt = pinger.Run() == nil
			return
		}()
	}

	wg.Wait()
	return testedList
}

// ============================================================================================

func (mc *SSHLanderControllers) ConnectionWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sshlander.connection[WS]"

	mc.sshMu.Lock()
	defer mc.sshMu.Unlock()

	if !websocket.IsWebSocketUpgrade(r) {
		return op, handler.StatusOK(w, "websocket connection only")
	}

	// ================================================================

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	baseCredentials, err := mc.repoSSH.QueryById(q.GetInt("id"))
	if err != nil {
		return op, err
	}

	// ================================================================

	wsConn, err := mc.websock.Upgrade(w, r, nil)
	if err != nil {
		return op, err
	}

	socket := handler.NewSocketWithContext(mc.ctx, wsConn, mc.logging)

	xTerm, err := mc.initTerm(socket.UUID(), baseCredentials, ssh.InsecureIgnoreHostKey())
	if err != nil {
		socket.Exit()
		return op, err
	}

	return op, mc.wsSSH(socket, xTerm)
}

func (mc *SSHLanderControllers) initTerm(
	uuid uuid.UUID,
	creds sshlander.SessionCredentials,
	callback ssh.HostKeyCallback,
) (*sshlander.TerminalSession, error) {

	sshSession, err := sshlander.NewClientSession(creds, callback, uuid)
	if err != nil {
		return nil, err
	}

	sshTerm, err := sshlander.ConnectTerminal(mc.ctx, sshSession, mc.term)
	if err != nil {
		sshSession.CloseDial()
		return nil, err
	}

	mc.logging.Infof(
		"ssh terminal connected: %s uuid %s:",
		sshSession.Instance(), uuid.String(),
	)

	return sshTerm, nil
}

func (mc *SSHLanderControllers) wsSSH(socket *handler.WsHandlerWrap, term *sshlander.TerminalSession) error {

	go func() {

		defer func() {
			mc.logging.Infof("stdout closed: %s", term.UUID())
			mc.logging.Infof("ssh terminal closed: %s", term.UUID())
			socket.Exit()
			mc.logging.Infof("ws closed: %s", term.UUID())
		}()

		// res := new(models.SSHResponseWS)

		for batch := range term.TerminalRead() {
			// res.Line = string(batch)
			socket.WriteBase64(batch)
		}

	}()

	data := new(models.SSHRequestWS)

	for msg := range socket.AwaitMessage(
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
	) {

		data = new(models.SSHRequestWS)
		if err := json.Unmarshal(msg.Body, data); err != nil {
			mc.logging.Error(err)
			continue
		}

		if err := term.SendTerminal(data.Command); err != nil {
			socket.WriteJSON(models.SSHResponseWS{
				Line: fmt.Sprintf("stdin error: %v", err),
			})
			mc.logging.Error(err)
			return err
		}

		mc.logging.Infof(
			"ssh sent command: '%s' uuid: %s ",
			data.Command,
			term.SSHSession.UUID(),
		)
	}

	if err := term.SendTerminal(sshlander.ExitCommand); err != nil {
		mc.logging.Errorf("send 'exit' error: %v uuid: %s", err, term.UUID())
	}

	return nil
}
