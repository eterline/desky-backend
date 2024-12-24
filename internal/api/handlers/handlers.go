package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eterline/desky-backend/pkg/logger"
)

type ApiHandleFunc func(http.ResponseWriter, *http.Request) (op string, err error)

func InitController(handle ApiHandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.ReturnEntry()

		op, err := handle(w, r)

		log.Debugf("requested controller: %s", op)

		if err != nil {

			var errResponse APIErrorResponse

			if apiError, ok := err.(*APIErrorResponse); ok {
				errResponse = *apiError
			} else {
				errResponse = InternalError()
			}

			WriteJSON(w, errResponse.StatusCode, errResponse)

			log.Errorf(
				"API Error - path: %s | controller: %s | error: %s",
				r.URL.Path, op, err.Error(),
			)
		}
	}
}

func InvalidRequestData(errors DataErrors) APIErrorResponse {
	return APIErrorResponse{
		APIResponse: APIResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errors,
		},
	}
}

func InvalidJSON() APIErrorResponse {
	return NewErrorResponse(
		http.StatusBadRequest,
		fmt.Errorf("uncorrect JSON data"),
	)
}

func InternalError() APIErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		fmt.Errorf("internal server error"),
	)
}

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}
