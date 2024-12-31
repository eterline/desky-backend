package sys

import (
	"net/http"
	"strconv"
	"time"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/services/system"
	"github.com/eterline/desky-backend/internal/utils"
	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const (
	WS_Message_delay = 5
)

type SysHandlerGroup struct {
	Sys *system.HostInfoService
	WS  websocket.Upgrader
}

func Init() *SysHandlerGroup {
	log = logger.ReturnEntry().Logger

	return &SysHandlerGroup{
		Sys: system.NewHostInfoService(),

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

func (s *SysHandlerGroup) SystemdUnits(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.sys.host-info"

	r.ParseForm()

	pageNumber, _ := strconv.Atoi(r.FormValue("page"))
	perPage, _ := strconv.Atoi(r.FormValue("count"))

	list, err := system.UnitsList()
	if err != nil {
		return op, err
	}

	w.Header().Add("All-Count", strconv.Itoa(len(list)))

	if pageNumber != 0 && perPage != 0 {

		cuted := utils.CutList(list, pageNumber*perPage, pageNumber*perPage+perPage)

		list = []system.SystemdUnit{}

		for _, d := range cuted {
			if d.UnitFile != "" {
				list = append(list, d)
			}
		}
	}

	return op, handlers.WriteJSON(w, http.StatusOK, list)
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
