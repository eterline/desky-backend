package monitoring

import (
	"context"
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
	Pool() (chan models.FetchedResponse, context.CancelFunc)
}

type MonitoringControllers struct {
	monitor MonitorProvider
	websock *websocket.Upgrader
	ctx     context.Context
}

func Init(ctx context.Context, m MonitorProvider, compress bool) *MonitoringControllers {

	log = logger.ReturnEntry().Logger

	return &MonitoringControllers{
		monitor: m,
		websock: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: compress,
		},
		ctx: ctx,
	}
}

func (mc *MonitoringControllers) Monitor(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "agents.monitor"

	if websocket.IsWebSocketUpgrade(r) {
		return mc.MonitorWS(w, r)
	}

	return op, handler.WriteJSON(w, http.StatusOK, mc.monitor.List())
}

func (mc *MonitoringControllers) MonitorWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "agents.monitor[WS]"

	conn, err := mc.websock.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		handler.InternalErrorResponse()
		return op, err
	}

	so := handler.NewSocketWithContext(mc.ctx, conn)
	so.AwaitClose(websocket.CloseNormalClosure, websocket.CloseGoingAway)

	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, _, err := conn.ReadMessage()
	// 		if err != nil {
	// 			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

	go func(so *handler.WebSocketHandle, ctx context.Context) {

		mon, stop := mc.monitor.Pool()
		defer stop()

		for {
			select {

			case <-so.Done():
				so.Exit()
				return

			case data := <-mon:
				conn.WriteJSON(data)
			}
		}

	}(so, mc.ctx)

	return
}
