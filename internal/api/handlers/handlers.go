package handlers

import (
	"net/http"

	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

type ApiHandler struct {
	WebDir, FrontFile string
	log               *logrus.Logger

	_ struct{}
}

func Init() *ApiHandler {
	return &ApiHandler{
		WebDir:    "./web",
		FrontFile: "index.html",
		log:       logger.ReturnEntry().Logger,
	}
}

// Main frontend handlers ===============================

func (h *ApiHandler) Front(w http.ResponseWriter, r *http.Request) {
	h.serveFile(w, r, h.FrontFile)
}

func (h *ApiHandler) Static(w http.ResponseWriter, r *http.Request) {
	h.fsProtection("static").ServeHTTP(w, r)
}

func (h *ApiHandler) Assets(w http.ResponseWriter, r *http.Request) {
	h.fsProtection("assets").ServeHTTP(w, r)
}

// =============================== API handlers ===============================

// AppsMenu handlers

func (h *ApiHandler) AppsList(w http.ResponseWriter, r *http.Request) {

}

func (h *ApiHandler) DeleteFromAppsList(w http.ResponseWriter, r *http.Request) {

}

func (h *ApiHandler) CreateInAppsList(w http.ResponseWriter, r *http.Request) {

}
