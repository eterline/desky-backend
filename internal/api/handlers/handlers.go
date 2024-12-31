package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eterline/desky-backend/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = nil

// Initialize custom type handler with error processing
func InitController(handle APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if log == nil {
			log = logger.ReturnEntry().Logger
		}

		op, err := handle(w, r)

		log.Debugf("requested controller: %s", op)

		if err != nil {

			switch e := err.(type) {

			case *APIErrorResponse:
				WriteJSON(w, e.StatusCode, e)
				break

			default:
				errDefault := InternalErrorResponse()
				WriteJSON(w, errDefault.StatusCode, errDefault)
				break
			}

			log.Errorf(
				"API Error - path: %s | controller: %s | error: %s",
				r.URL.Path, op, err.Error(),
			)
		}
	}
}

func QueryURLParameters(r *http.Request, params ...string) (map[string]string, error) {

	data := make(map[string]string, len(params))

	for _, param := range params {

		str := chi.URLParam(r, param)

		if str == "" || str == " " {
			return nil, NewErrorResponse(
				http.StatusNotAcceptable,
				ErrEmptyParameter(param),
			)
		}

		data[param] = str
	}

	return data, nil
}

func QueryURLNumeredParameters(r *http.Request, params ...string) (map[string]int, error) {

	stringPs, err := QueryURLParameters(r, params...)
	if err != nil {
		return nil, err
	}

	numPs := make(map[string]int, len(params))

	for _, param := range params {

		num, err := strconv.Atoi(stringPs[param])

		if err != nil {
			return nil, NewErrorResponse(
				http.StatusNotAcceptable,
				ErrInterpretationToNumber(param),
			)
		}

		numPs[param] = num
	}

	return numPs, nil
}

// Decode JSON request body to data structure
func DecodeRequest(r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return InvalidContentTypeResponse()
	}
	return json.NewDecoder(r.Body).Decode(&v)
}

// Send JSON object response
func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}
