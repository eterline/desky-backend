package sys

import (
	"context"
	"encoding/json"
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
	wsHandler *handler.WebSocketHandler
	ctx       context.Context
}

func Init(ctx context.Context, hs HostService) *SysHandlerGroup {
	log = logger.ReturnEntry().Logger

	return &SysHandlerGroup{
		HostService: hs,
		ctx:         ctx,
		wsHandler: handler.NewWebSocketHandler(ctx, &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}),
	}
}

func (s *SysHandlerGroup) Stats(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sys.stats"
	if websocket.IsWebSocketUpgrade(r) {
		return s.StatsWS(w, r)
	}
	return op, handler.WriteJSON(w, http.StatusOK, s.HostInfo())
}

func (s *SysHandlerGroup) StatsWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "sys.stats[WS]"

	infoGet := func() models.StatsResponse {
		return models.StatsResponse{
			RAM:  s.RAMInfo(),
			CPU:  s.CPUInfo(),
			Temp: s.Temperatures(),
			Load: s.Load(),
		}
	}

	socket, err := s.wsHandler.HandleConnect(w, r)
	defer socket.Exit()

	socket.AwaitClose(websocket.CloseNormalClosure, websocket.CloseGoingAway)
	wr := socket.InitWebSocketWriting(false)
	defer wr.CloseWriting()
	sockEnc := json.NewEncoder(wr)

	ticker := time.NewTicker(time.Second * WS_Message_delay)
	defer ticker.Stop()

	go sockEnc.Encode(infoGet())

	for {
		select {

		case <-socket.SessionDone():
			return

		case <-ticker.C:
			sockEnc.Encode(infoGet())
		}
	}
}

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
