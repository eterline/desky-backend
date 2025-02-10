package sshlander

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	sshlander "github.com/eterline/desky-backend/internal/services/ssh-lander"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-ping/ping"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type SSHRepository interface {
	AddHost(username string, host string, port uint16, osType string, privateKeyUsage bool, password string, key string) error
	Delete(id int) error
	QueryAll() ([]models.SSHCredentialsT, error)
	QueryById(id int) (*models.SSHCredentialsT, error)
}

type SSHLanderControllers struct {
	ctx     context.Context
	repoSSH SSHRepository
	websock *websocket.Upgrader

	logging *logrus.Logger
}

func Init(
	ctx context.Context,
	repo SSHRepository,
) *SSHLanderControllers {
	return &SSHLanderControllers{
		ctx:     ctx,
		repoSSH: repo,
		websock: &websocket.Upgrader{
			HandshakeTimeout:  10 * time.Second,
			EnableCompression: true,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
		},

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

	q, err := handler.ParseURLParameters(r, handler.NumOpts("id"))
	if err != nil {
		return op, err
	}

	if !websocket.IsWebSocketUpgrade(r) {
		return op, handler.StatusOK(w, "websocket connection only")
	}

	credentials, err := mc.repoSSH.QueryById(q.GetInt("id"))
	if err != nil {
		return op, err
	}

	conn, err := mc.websock.Upgrade(w, r, nil)
	if err != nil {
		return op, err
	}

	socket := handler.NewSocketWithContext(mc.ctx, conn)
	mc.logging.Warnf("ssh external connection: %s. user: %s", credentials.Host, credentials.Username)
	return op, mc.ProcessSSH(socket, credentials)
}

func (mc *SSHLanderControllers) ProcessSSH(sock *handler.WebSocketHandle, data *models.SSHCredentialsT) error {

	ssh := sshlander.New(data.Username)

	if data.Security.PrivateKeyUse {
		ssh.SetupAuth(sshlander.PrivateKeyMethod, data.Security.PrivateKey)
	} else {
		ssh.SetupAuth(sshlander.PasswordMethod, data.Security.Password)
	}

	if err := ssh.Connect(mc.ctx, data.Socket()); err != nil {
		mc.logging.Errorf("ssh connection error: %s", err.Error())
		sock.WriteCloseError(err)
		return err
	}

	sock.AwaitClose(
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
	)

	go func(
		socket *handler.WebSocketHandle,
		ssh *sshlander.SSHLanderService,
		data *models.SSHCredentialsT,
	) {

		defer func() {
			socket.Exit()
			ssh.Exit()
			mc.logging.Infof("ssh connection closed: %s", data.Host)
		}()

		msg := socket.AwaitMessage(new(models.SSHSessionRequest))

		for {
			select {

			case <-socket.Done():
				return

			case m := <-msg:
				val, ok := m.(*models.SSHSessionRequest)
				if !ok {
					continue
				}

				mc.logging.Infof("ssh command: %s. host: %s", val.Command, data.Host)

				response := ssh.SendCommand(val.Command)
				for output := range response {
					socket.WriteJSON(models.SSHSessionResponse{
						Response: output,
					})
				}
			}
		}

	}(sock, ssh, data)

	return nil
}
