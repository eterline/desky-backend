package handlers

import (
	"net/http"
	"path/filepath"
)

func (h *ApiHandler) fsProtection(dir string) http.Handler {
	fs := http.FileServer(http.Dir(filepath.Join(h.WebDir, dir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path[0] == '.' {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.StripPrefix("/"+dir+"/", fs).ServeHTTP(w, r)
	})
}

func (h *ApiHandler) serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	http.ServeFile(w, r, filepath.Join(h.WebDir, filePath))
}
