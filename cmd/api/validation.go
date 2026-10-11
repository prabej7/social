package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type contextKey string

const (
	BodyKey contextKey = "Body"
)

func getValidationErrorMessage(err validator.FieldError) string {
	var msg string

	switch err.Tag() {
	case "required":
		msg = fmt.Sprintf("%s is required", strings.ToLower(err.Field()))

	case "email":
		msg = fmt.Sprintf("%s must be a valid email", strings.ToLower(err.Field()))

	case "min":
		msg = fmt.Sprintf("%s must be at least %s characters", strings.ToLower(err.Field()), err.Param())
	default:
		msg = fmt.Sprintf("%s is invalid", strings.ToLower(err.Field()))
	}

	return msg
}

func buildValidationErrorMessage(validationErrors validator.ValidationErrors) map[string][]string {
	validationErrorsMap := make(map[string][]string)

	for _, err := range validationErrors {
		validationErrorsMap[strings.ToLower(err.Field())] = append(validationErrorsMap[strings.ToLower(err.Field())], getValidationErrorMessage(err))
	}

	return validationErrorsMap
}

func validatePayload[T any](next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJson(w, http.StatusBadRequest, ValidationError{
				Success: false,
				Message: "invalid payload",
			})
			return
		}

		if err := validator.New().Struct(body); err != nil {

			validationErrors, ok := err.(validator.ValidationErrors)

			if !ok {
				writeJson(w, http.StatusBadRequest, ValidationError{
					Success: false,
					Message: err.Error(),
				})
				return
			}

			validationErrorsMsg := buildValidationErrorMessage(validationErrors)
			writeJson(w, http.StatusBadRequest, ValidationError{
				Success: false,
				Message: "ValidationError",
				Errors:  validationErrorsMsg,
			})
			return

		}

		ctx := context.WithValue(r.Context(), BodyKey, body)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
