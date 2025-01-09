package sys

import (
	"net/http"
	"strconv"
	"time"

	"github.com/eterline/desky-backend/internal/api/handlers"
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
	Sys   HostService
	WS    websocket.Upgrader
	Cache Cacher
}

func Init(hs HostService, ch Cacher) *SysHandlerGroup {
	log = logger.ReturnEntry().Logger
	// conf := configuration.GetConfig()

	return &SysHandlerGroup{
		Sys:   hs,
		Cache: ch,

		WS: websocket.Upgrader{
			HandshakeTimeout:  10 * time.Second,
			EnableCompression: true,
		},
	}
}

func (s *SysHandlerGroup) HostInfo(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.sys.host-info"

	return op, handlers.WriteJSON(w, http.StatusOK, s.Sys.HostInfo())
}

func (s *SysHandlerGroup) HostStatsWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.sys.host-stats-WS"

	connection, err := s.WS.Upgrade(w, r, nil)
	defer connection.Close()
	if err != nil {
		return op, err
	}

	connection.SetCloseHandler(func(code int, text string) error {
		log.Infof("websocket closed: %s - code: %d, reason: %s", connection.RemoteAddr(), code, text)
		return nil
	})

	ticker := time.NewTicker(time.Second * WS_Message_delay)
	defer ticker.Stop()

	infoGet := func() StatsResponse {
		return StatsResponse{
			RAM:  s.Sys.RAMInfo(),
			CPU:  s.Sys.CPUInfo(),
			Temp: s.Sys.Temperatures(),
			Load: s.Sys.Load(),
		}
	}

	connection.WriteJSON(infoGet())

	for {
		select {

		case <-ticker.C:

			if err := connection.WriteJSON(infoGet()); err != nil {
				switch e := err.(type) {
				case *websocket.CloseError:
					log.Infof("websocket connection: %s - closed", connection.RemoteAddr())
					return op, nil
				default:
					log.Errorf("websocket error: %s", e.Error())
					return op, e
				}
			}
		}

	}
}

func (s *SysHandlerGroup) TtyWS(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.sys.tty-WS"

	connection, err := s.WS.Upgrade(w, r, nil)
	defer connection.Close()
	if err != nil {
		return op, err
	}

	log.Infof("websocket connection: %s - opened", connection.RemoteAddr())

	for {
		_, msgContent, err := connection.ReadMessage()
		if err != nil {

			switch e := err.(type) {

			case *websocket.CloseError:
				log.Infof("websocket connection: %s - closed", connection.RemoteAddr())
				return op, nil

			default:
				log.Errorf("websocket tty error: %s", err.Error())
				return op, e
			}
		}

		resp, err := system.HandleCommand(msgContent)
		if err != nil {
			log.Errorf("websocket tty error: %s", err.Error())
		} else {
			log.Infof("command: '%s' - executed by request: %s", resp.Command, connection.RemoteAddr())
		}

		connection.WriteJSON(resp)
	}
}

func (s *SysHandlerGroup) SystemdUnits(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.sys.host-info"

	r.ParseForm()
	pageNumber, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("count"))

	if data := s.Cache.GetValue(op); data != nil && !s.Cache.OlderThanAndExists(op, time.Second*30) {

		list := data.([]system.SystemdUnit)

		paginatedList := paginateSystemdUnits(list, pageNumber, perPage)
		w.Header().Add("All-Count", strconv.Itoa(len(list)))

		return op, handlers.WriteJSON(w, http.StatusOK, paginatedList)
	}

	list, err := system.UnitsList()
	if err != nil {
		return op, err
	}

	s.Cache.PushValue(op, list)

	paginatedList := paginateSystemdUnits(list, pageNumber, perPage)
	w.Header().Add("All-Count", strconv.Itoa(len(list)))

	return op, handlers.WriteJSON(w, http.StatusOK, paginatedList)
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
	op = "handlers.sys.unit-command"

	qStr, err := handlers.QueryURLParameters(r, "unit", "command")
	if err != nil {
		return op, err
	}

	unit, err := system.UnitInstance(qStr["unit"])
	if err != nil {

		if err == system.ErrUnitNotFound {
			return op, handlers.NewErrorResponse(
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
		return op, handlers.NewErrorResponse(
			http.StatusBadRequest,
			ErrUnknownUnitCommand,
		)
	}

	if err == nil {
		err = handlers.StatusOK(w, "command successfully completed")
	}

	return op, err
}
