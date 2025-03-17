package controllers

import (
	"net/http"

	"github.com/eterline/desky-backend/internal/services/handler"
)

type FrontendHandlerGroup struct {
	FS, Storage *FilesHandlerGroup
	HTMLfile    string
}

func InitFronEnd() *FrontendHandlerGroup {
	return &FrontendHandlerGroup{
		HTMLfile: "index.html",

		FS:      InitFiles("./web"),
		Storage: InitFiles("./storage"),
	}
}

func (fh *FrontendHandlerGroup) HTML(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "frontend.html"

	fh.FS.ServeFile(w, r, fh.HTMLfile)
	return op, err
}

func (fh *FrontendHandlerGroup) Static(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "frontend.static"

	fh.FS.ServeDir("static").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) Assets(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "frontend.assets"

	fh.FS.ServeDir("assets").ServeHTTP(w, r)
	return op, err
}

func (fh *FrontendHandlerGroup) WallpaperHandle(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "frontend.storage"

	fh.Storage.ServeDir("wallpaper").ServeHTTP(w, r)
	return op, err
}

// func (fh *FrontendHandlerGroup) Login(w http.ResponseWriter, r *http.Request) (op string, err error) {
// 	op = "handler.frontend.login"

// 	form, err := authorization.DecodeCredentials(r)
// 	if err != nil {
// 		return op, err
// 	}

// 	token, err := fh.Auth.Token(authorization.NewPayload(form.GetUsername()))
// 	if err != nil {
// 		return op, err
// 	}

// 	if fh.Auth.IsValid(form) {
// 		return op, handler.WriteJSON(w, http.StatusAccepted, NewTokenResponse(token))
// 	}

// 	return op, handler.NewErrorResponse(
// 		http.StatusNotAcceptable,
// 		ErrUncorrectCredentials,
// 	)
// }

func AccessCheck(w http.ResponseWriter, r *http.Request) (op string, err error) {
	op = "handler.frontend.access-check"

	return op, handler.StatusOK(w, "accepted")
}
