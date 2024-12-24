package handlers

import "fmt"

type APIResponse struct {
	StatusCode int `json:"code"`
	Message    any `json:"message"`
}

func NewResponse(code int, message any) APIResponse {
	return APIResponse{
		StatusCode: code,
		Message:    message,
	}
}

func (ar *APIResponse) String() string {
	return fmt.Sprintf("api response code: %d. message: %v", ar.StatusCode, ar.Message)
}

type APIErrorResponse struct {
	APIResponse
}

func NewErrorResponse(code int, err error) APIErrorResponse {
	return APIErrorResponse{
		APIResponse: APIResponse{
			StatusCode: code,
			Message:    err.Error(),
		},
	}
}

func (ar *APIErrorResponse) Error() string {
	return fmt.Sprintf("api error code: %d. action: %s. error: %v", ar.StatusCode, ar.Message)
}

type DataErrors map[string]error
