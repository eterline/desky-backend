package middlewares

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) Write(w []byte) (int, error) {
	return rw.ResponseWriter.Write(w)
}
