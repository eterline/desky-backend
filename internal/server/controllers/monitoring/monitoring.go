package monitoring

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

type MonitorProvider interface {
	List() []models.SessionCredentials
	Pool() (chan any, context.CancelFunc)
}

type MonitoringControllers struct {
	monitor   MonitorProvider
	wsHandler *handler.WebSocketHandler
	ctx       context.Context
}

func Init(ctx context.Context, m MonitorProvider, compress bool) *MonitoringControllers {

	log = logger.ReturnEntry().Logger

	return &MonitoringControllers{
		ctx:     ctx,
		monitor: m,

		wsHandler: handler.NewWebSocketHandler(ctx, &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: compress,
		}),
	}
}

func (mc *MonitoringControllers) Monitor(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "agents.monitor"

	if websocket.IsWebSocketUpgrade(r) {
		return mc.MonitorWS(w, r)
	}

	monitorList := mc.monitor.List()

	if handler.ListIsEmpty(w, monitorList) {
		return op, err
	}

	return op, handler.WriteJSON(w, http.StatusOK, monitorList)
}

func (mc *MonitoringControllers) MonitorWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "agents.monitor[WS]"

	sock, err := mc.wsHandler.HandleConnect(w, r)
	if err != nil {
		log.Error(err)
		return op, err
	}
	defer sock.Exit()

	sock.AwaitClose(websocket.CloseNormalClosure, websocket.CloseGoingAway)

	wr := sock.InitWebSocketWriting(false)
	defer wr.CloseWriting()

	monitor, stop := mc.monitor.Pool()
	defer stop()

	for {
		select {

		case <-sock.SessionDone():
			return op, nil

		case data := <-monitor:
			json.NewEncoder(wr).Encode(data)
		}
	}
}
