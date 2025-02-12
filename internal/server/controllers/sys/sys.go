package sys

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/eterline/desky-backend/internal/models"
	"github.com/eterline/desky-backend/internal/services/router/handler"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const (
	WS_Message_delay = 5
)

type HostService interface {
	HostInfo() *system.HostInfo
	RAMInfo() *system.RAMInfo
	CPUInfo() *system.CPUInfo
	Temperatures() []system.SensorInfo
	Load() *system.AverageLoad
}

type Cacher interface {
	GetValue(key any) any
	PushValue(key, value any)
	OlderThanAndExists(key any, duration time.Duration) bool
}

type SysHandlerGroup struct {
	HostService
	websock *websocket.Upgrader
	ctx     context.Context
}

func Init(ctx context.Context, hs HostService) *SysHandlerGroup {
	log = logger.ReturnEntry().Logger

	return &SysHandlerGroup{
		HostService: hs,
		ctx:         ctx,
		websock: &websocket.Upgrader{
			HandshakeTimeout:  10 * time.Second,
			EnableCompression: true,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
		},
	}
}

// HostInfo godoc
//
//	@Tags			system
//	@Summary		HostInfo
//	@Description	host information
//
//	@Produce		json
//	@Success		200	{object}	system.HostInfo
//	@Failure		500	{object}	handler.APIErrorResponse	"Internal server error"
//	@Router			/system/info [get]
func (s *SysHandlerGroup) Stats(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sys.stats"
	if websocket.IsWebSocketUpgrade(r) {
		return s.StatsWS(w, r)
	}
	return op, handler.WriteJSON(w, http.StatusOK, s.HostInfo())
}

func (s *SysHandlerGroup) StatsWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sys.stats[WS]"

	conn, err := s.websock.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		handler.InternalErrorResponse()
		return op, err
	}

	infoGet := func() models.StatsResponse {
		return models.StatsResponse{
			RAM:  s.RAMInfo(),
			CPU:  s.CPUInfo(),
			Temp: s.Temperatures(),
			Load: s.Load(),
		}
	}

	so := handler.NewSocketWithContext(s.ctx, conn, log)
	so.AwaitClose(websocket.CloseNormalClosure, websocket.CloseGoingAway)

	go func(so *handler.WsHandlerWrap) {

		ticker := time.NewTicker(time.Second * WS_Message_delay)
		defer ticker.Stop()

		for {
			select {

			case <-so.Done():
				so.Exit()
				return

			case <-ticker.C:
				data := infoGet()
				conn.WriteJSON(data)
			}
		}

	}(so)

	return
}

// HostStatsWS godoc
//
//	@Tags			system
//	@Summary		HostStatsWS
//	@Description	host information ws interval update = 5s
//
//	@Produce		json
//	@Success		200	{object}	StatsResponse
//	@Failure		500	{object}	handler.APIErrorResponse	"Internal server error"
//	@Router			/system/stats [get]
// func (s *SysHandlerGroup) HostStatsWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
// 	op = "handler.sys.host-stats-WS"

// 	connection, err := s.websock.Upgrade(w, r, nil)
// 	defer connection.Close()
// 	if err != nil {
// 		return op, err
// 	}

// 	connection.SetCloseHandler(func(code int, text string) error {
// 		log.Infof("websocket closed: %s - code: %d, reason: %s", connection.RemoteAddr(), code, text)
// 		return nil
// 	})

// 	ticker := time.NewTicker(time.Second * WS_Message_delay)
// 	defer ticker.Stop()

// 	infoGet := func() models.StatsResponse {
// 		return models.StatsResponse{
// 			RAM:  s.RAMInfo(),
// 			CPU:  s.CPUInfo(),
// 			Temp: s.Temperatures(),
// 			Load: s.Load(),
// 		}
// 	}

// 	connection.WriteJSON(infoGet)

// 	for {
// 		select {

// 		case <-ticker.C:

// 			if err := connection.WriteJSON(infoGet()); err != nil {
// 				switch e := err.(type) {
// 				case *websocket.CloseError:
// 					log.Infof("websocket connection: %s - closed", connection.RemoteAddr())
// 					return op, nil
// 				default:
// 					log.Errorf("websocket error: %s", e.Error())
// 					return op, e
// 				}
// 			}
// 		}

// 	}
// }

// func (s *SysHandlerGroup) TtyWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
// 	op = "handler.sys.tty-WS"

// 	connection, err := s.WS.Upgrade(w, r, nil)
// 	defer connection.Close()
// 	if err != nil {
// 		return op, err
// 	}

// 	log.Infof("websocket connection: %s - opened", connection.RemoteAddr())

// 	for {
// 		_, msgContent, err := connection.ReadMessage()
// 		if err != nil {

// 			switch e := err.(type) {

// 			case *websocket.CloseError:
// 				log.Infof("websocket connection: %s - closed", connection.RemoteAddr())
// 				return op, nil

// 			default:
// 				log.Errorf("websocket tty error: %v", err)
// 				return op, e
// 			}
// 		}

// 		resp, err := system.HandleCommand(msgContent)
// 		if err != nil {
// 			log.Errorf("websocket tty error: %v", err)
// 		} else {
// 			log.Infof("command: '%s' - executed by request: %s", resp.Command, connection.RemoteAddr())
// 		}

// 		connection.WriteJSON(resp)
// 	}
// }

// SystemdUnits godoc
//
//	@Tags			system
//	@Summary		SystemdUnits
//	@Description	units systemd list
//	@Param			page	query	string	false	"Page number for pagination (optional)"
//	@Param			count	query	string	false	"Number of items per page (optional)"
//	@Produce		json
//	@Success		200	{object}	[]system.SystemdUnit
//	@Failure		500	{object}	handler.APIErrorResponse	"Internal server error"
//	@Router			/system/systemd/status [get]
func (s *SysHandlerGroup) SystemdUnits(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.sys.host-info"

	r.ParseForm()
	pageNumber, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("count"))

	list, err := system.UnitsList()
	if err != nil {
		return op, err
	}

	paginatedList := paginateSystemdUnits(list, pageNumber, perPage)
	w.Header().Add("All-Count", strconv.Itoa(len(list)))

	return op, handler.WriteJSON(w, http.StatusOK, paginatedList)
}

func paginateSystemdUnits(list []system.SystemdUnit, pageNumber, perPage int) []system.SystemdUnit {
	if pageNumber == 0 || perPage == 0 {
		return list
	}

	start := pageNumber * perPage
	end := start + perPage
	if start >= len(list) {
		return []system.SystemdUnit{}
	}
	if end > len(list) {
		end = len(list)
	}

	cuted := list[start:end]
	filtered := []system.SystemdUnit{}
	for _, d := range cuted {
		if d.UnitFile != "" {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// UnitCommand godoc
//
//	@Tags			system
//	@Summary		UnitCommand
//	@Description	execute device command
//
//	@Produce		json
//	@Param			command	path		string	true	"systemd command"
//	@Param			service	path		string	true	"systemd service"
//	@Success		200		{object}	handler.APIResponse
//	@Failure		400		{object}	handler.APIErrorResponse	"Invalid parameters"
//	@Failure		500		{object}	handler.APIErrorResponse	"Internal server error"
//	@Failure		501		{object}	handler.APIErrorResponse	"Uninplemented"
//	@Router			/system/systemd/{service}/{command} [post]
func (s *SysHandlerGroup) UnitCommand(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.sys.unit-command"

	qStr, err := handler.QueryURLParameters(r, "unit", "command")
	if err != nil {
		return op, err
	}

	unit, err := system.UnitInstance(qStr["unit"])
	if err != nil {

		if err == system.ErrUnitNotFound {
			return op, handler.NewErrorResponse(
				http.StatusBadRequest,
				err,
			)
		}

		return op, err
	}

	switch qStr["command"] {

	case "stop":
		err = unit.Stop()
		break

	case "start":
		err = unit.Start()
		break

	case "restart":
		err = unit.Restart()
		break

	default:
		return op, handler.NewErrorResponse(
			http.StatusBadRequest,
			ErrUnknownUnitCommand,
		)
	}

	if err == nil {
		err = handler.StatusOK(w, "command successfully completed")
	}

	return op, err
}
