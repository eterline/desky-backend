package handler

import (
	"fmt"
	"net/http"
)

// Basic type response constructors

// Create new response
func NewResponse(code int, message string) *APIResponse {
	return &APIResponse{
		StatusCode: code,
		Message:    message,
	}
}

// Create OK (200) response
func StatusOK(w http.ResponseWriter, message string) error {
	resp := NewResponse(http.StatusOK, message)
	return WriteJSON(w, resp.StatusCode, resp)
}

// Create StatusCreated (201) response
func StatusCreated(w http.ResponseWriter, message string) error {
	resp := NewResponse(http.StatusCreated, message)
	return WriteJSON(w, resp.StatusCode, resp)
}

// Error type response constructors

// Create new error type response
func NewErrorResponse(code int, err error) *APIErrorResponse {
	return &APIErrorResponse{
		APIResponse: APIResponse{
			StatusCode: code,
			Message:    err.Error(),
		},
	}
}

// Create error StatusUnprocessableEntity (422) response:
// DataErrors <= map[string]error
func InvalidRequestDataResponse(errors DataErrors) *APIErrorResponse {
	return &APIErrorResponse{
		APIResponse: APIResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errors,
		},
	}
}

// Create error StatusBadRequest (406) response
func InvalidContentTypeResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusNotAcceptable,
		fmt.Errorf("uncorrect 'Content-Type' header, must be: JSON"),
	)
}

// Create error StatusBadRequest (400) response
func ErrorBadRequest() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusBadRequest,
		fmt.Errorf("uncorrect request data"),
	)
}

// Create error StatusBadRequest (500) response
func InternalErrorResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusInternalServerError,
		fmt.Errorf("internal server error"),
	)
}

// Create error StatusUnauthorized (401) response.
func UnauthorizedErrorResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusUnauthorized,
		fmt.Errorf("unauthorized request"),
	)
}

// Create error StatusUnauthorized (403) response.
func ForbiddenRequestResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusForbidden,
		fmt.Errorf("forbidden request"),
	)
}

func NotFoundPageResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusNotFound,
		fmt.Errorf("route controller not found"),
	)
}

func NoContentResponse() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusNoContent,
		fmt.Errorf("content not exists"),
	)
}

func StatusUnauthorized() *APIErrorResponse {
	return NewErrorResponse(
		http.StatusUnauthorized,
		fmt.Errorf("incorrect credentials"),
	)
}

func BadRequestParam(param string) *APIErrorResponse {
	return NewErrorResponse(
		http.StatusBadRequest,
		fmt.Errorf("bad parameter: %v", param),
	)
}
