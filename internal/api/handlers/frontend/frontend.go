package frontend

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/api/handlers"
	"github.com/eterline/desky-backend/internal/api/handlers/files"
	"github.com/eterline/desky-backend/internal/services/authorization"
)

type Authorization interface {
	Token(authorization.Payload) (string, error)
	IsValid(authorization.AuthForm) bool
}

type FrontendHandlerGroup struct {
	FS, Storage *files.FilesHandlerGroup
	HTMLfile    string
	Auth        Authorization
}

func Init(service Authorization) *FrontendHandlerGroup {
	return &FrontendHandlerGroup{
		Auth:     service,
		HTMLfile: "index.html",

		FS:      files.Init("./web"),
		Storage: files.Init("./storage"),
	}
}

func (fh *FrontendHandlerGroup) HTML(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.html"

	fh.FS.ServeFile(w, r, fh.HTMLfile)
	return op, err
}

func (fh *FrontendHandlerGroup) Static(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.static"

	fh.FS.ServeDir("static").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) Assets(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.assets"

	fh.FS.ServeDir("assets").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) WallpaperHandle(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.storage"

	fh.Storage.ServeDir("wallpaper").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) Login(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.login"

	form, err := authorization.DecodeCredentials(r)
	if err != nil {
		return op, err
	}

	token, err := fh.Auth.Token(authorization.NewPayload(form.GetUsername()))
	if err != nil {
		return op, err
	}

	if fh.Auth.IsValid(form) {
		return op, handlers.WriteJSON(w, http.StatusAccepted, NewTokenResponse(token))
	}

	return op, handlers.NewErrorResponse(
		http.StatusNotAcceptable,
		ErrUncorrectCredentials,
	)
}

func AccessCheck(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handlers.frontend.access-check"

	return op, handlers.StatusOK(w, "accepted")
}
