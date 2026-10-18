package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prabej7/social/internal/store"
)

var userContextKey contextKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(*store.User)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) GetUserContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		parsedId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetById(r.Context(), parsedId)

		if err != nil {
			switch {
			case errors.Is(store.NotFoundError, err):
				app.notFoundError(w, r, fmt.Errorf("user not found"))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
