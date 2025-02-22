package application

import (
	"net/http"
	"sync"

	"github.com/eterline/desky-backend/internal/services/cache"
	"github.com/eterline/desky-backend/internal/services/router/handler"
)

type AppLanguage string

const (
	LangEN AppLanguage = "EN"
	LangRU AppLanguage = "RU"
)

type ApplicationSettings struct {
	Language   string `json:"language"`
	Background string `json:"background"`
	Auth       bool   `json:"auth"`

	sync.Mutex `json:"-"`
}

func (s *ApplicationSettings) SetAuth(value bool) {

	s.Lock()
	defer s.Unlock()

	s.Auth = value
}

func (s *ApplicationSettings) SetLanguage(l AppLanguage) {

	s.Lock()
	defer s.Unlock()

	s.Language = string(l)
}

func (s *ApplicationSettings) SetBG(value string) {

	s.Lock()
	defer s.Unlock()

	s.Background = value
}

func (s *ApplicationSettings) SettingHandler(w http.ResponseWriter, r *http.Request) {
	handler.WriteJSON(w, http.StatusOK, s)
}

func (s *ApplicationSettings) ThemeHandler(w http.ResponseWriter, r *http.Request) {
	handler.WriteJSON(w, http.StatusOK, ThemeStorage)
}

func (s *ApplicationSettings) WriteBG(w http.ResponseWriter, r *http.Request) {

	image, ok := cache.GetEntry().GetValue("bg").([]byte)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(200)
	w.Write(image)
}
