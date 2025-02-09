package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

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

func Validate(v any) error {
	var validatorError validator.ValidationErrors

	err := validator.New().Struct(v)

	if err == nil {
		return nil
	}

	if errors.As(err, &validatorError) {
		invalidData := make(DataErrors)
		for _, value := range validatorError {
			invalidData[value.Field()] = valFieldErr(value)
		}
		return NewErrorResponse(
			http.StatusBadRequest,
			invalidData,
		)
	}

	return err
}

func valFieldErr(field validator.FieldError) error {

	switch field.Tag() {

	case "required":
		return errors.New("field must be filled")

	case "email":
		return errors.New("field must be email type")

	case "ip":
		return errors.New("field must be ip type")

	case "url":
		return errors.New("field must be url type")

	case "min":
		return errors.New("field value is too small")

	default:
		return field
	}
}

func ListIsEmpty[Type any](w http.ResponseWriter, list []Type) bool {
	if list == nil || len(list) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}
