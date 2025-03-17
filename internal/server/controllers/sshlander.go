package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/handler"
	sshlander "github.com/eterline/desky-backend/internal/services/ssh-lander"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/eterline/desky-backend/pkg/net-wait-go-forked/wait"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	MicroByteChunk    = 128   // 128 B
	SmallByteChunk    = 512   // 512 B
	UsualByteChunk    = 1024  // 1 KB
	MiddleByteChunk   = 2048  // 2 KB
	BigByteChunk      = 4092  // 4 KB
	ExtremeByteChunk  = 8192  // 8 KB
	OverflowByteChunk = 16384 // 16 KB
)

type SSHRepository interface {
	AddHost(username string, host string, port uint16, osType string, privateKeyUsage bool, password string, key string) error
	Delete(id int) error
	QueryAll() ([]models.SSHCredentialsT, error)
	QueryById(id int) (*models.SSHCredentialsT, error)
}

type SSHLanderControllers struct {
	ctx context.Context

	term      sshlander.TerminalType
	repoSSH   SSHRepository
	wsHandler *handler.WebSocketHandler

	logging *logrus.Logger

	sshMu sync.Mutex
}

func InitSSHlander(ctx context.Context, repo SSHRepository) *SSHLanderControllers {
	return &SSHLanderControllers{
		ctx:     ctx,
		repoSSH: repo,
		logging: logger.ReturnEntry().Logger,
		term:    sshlander.XtermColored,

		wsHandler: handler.NewWebSocketHandler(ctx, &websocket.Upgrader{
			ReadBufferSize:    MiddleByteChunk,
			WriteBufferSize:   MiddleByteChunk,
			EnableCompression: true,
		}),
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
	var mu sync.Mutex
	testedList := make([]models.SSHTestObject, len(hostsData))

	for idx, credential := range hostsData {

		wg.Add(1)

		go func() {
			var attempt bool = false

			defer func() {
				mu.Lock()
				testedList[idx] = models.SSHTestObject{
					ID:        int(credential.ID),
					Available: attempt,
				}
				mu.Unlock()

				wg.Done()
			}()

			attempt = wait.New(
				wait.WithDeadline(5*time.Second),
				wait.WithProto("tcp"),
			).Do([]string{credential.Socket()})

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

	socket, err := mc.wsHandler.HandleConnect(w, r)
	if err != nil {
		return op, err
	}
	defer socket.Exit()

	wsBase64Writer := socket.InitWebSocketBase64Writing()
	defer wsBase64Writer.CloseWriting()

	term, err := mc.initTerm(socket.ID, baseCredentials, ssh.InsecureIgnoreHostKey())
	if err != nil {
		fmt.Fprintf(wsBase64Writer, "ssh connection error: %v", err)
		socket.Exit()
		return op, err
	}
	defer term.CloseDial()

	return op, mc.wsSSH(wsBase64Writer, socket, term)
}

func (mc *SSHLanderControllers) initTerm(
	uuid uuid.UUID,
	creds sshlander.SessionCredentials,
	callback ssh.HostKeyCallback,
) (*sshlander.TerminalSession, error) {

	session, err := sshlander.NewClientSession(creds, callback, uuid)
	if err != nil {
		return nil, err
	}

	term, err := sshlander.ConnectTerminal(mc.ctx, session, mc.term)
	if err != nil {
		session.CloseDial()
		return nil, err
	}

	mc.logging.Infof("ssh terminal connected: %s uuid: %s", term.Instance(), uuid.String())

	return term, nil
}

func (mc *SSHLanderControllers) wsSSH(writer *handler.WebSocketBase64Writing, socket *handler.WebSocketSession, term *sshlander.TerminalSession) error {

	defer mc.logging.Infof("ws uuid: %s ssh closed", term.UUID())

	go func() {

		defer func() {
			defer term.WriteExit()
			mc.logging.Infof("ws uuid: %s terminal session end", term.UUID())
		}()

		if err := term.FromTerminalBytes(writer, UsualByteChunk); err != nil {
			mc.logging.Infof(
				"ws uuid: %s terminal write error: %v",
				term.UUID(), err,
			)
		}
	}()

	data := new(models.SSHRequestWS)

	for msg := range socket.AwaitMessage(
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
	) {

		// prepare message
		data = new(models.SSHRequestWS)
		if err := json.Unmarshal(msg.Body, data); err != nil {
			mc.logging.Error(err)
			continue
		}

		//
		if err := term.Send([]byte(data.Command)); err != nil {
			mc.logging.Error(err)
			return err
		}

		mc.logging.Infof(
			"ssh sent command: '%s' uuid: %s ",
			data.Command,
			term.SSHSession.UUID(),
		)
	}

	return nil
}
