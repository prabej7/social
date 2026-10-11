package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type contextKey string

const (
	BodyKey contextKey = "Body"
)

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
			writeJson(w, http.StatusBadRequest, ValidationError{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), BodyKey, body)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
