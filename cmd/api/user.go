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

var (
	userContextKey      contextKey = "user"
	paramUserContextKey contextKey = "paramUser"
)

// GetUser godoc
//
//	@Summary		Get user
//	@Description	Get a user by ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Faiure			500 {object} error
//	@Router			/user/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromContext(r, userContextKey)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		parsedId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, parsedId)

		if err != nil {
			switch {
			case errors.Is(store.NotFoundError, err):
				app.notFoundError(w, r, fmt.Errorf("user not found"))
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(r.Context(), paramUserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUserFromContext(r *http.Request, key contextKey) *store.User {
	user := r.Context().Value(key).(*store.User)

	return user
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by id
//	@Tags			users
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/user/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r, paramUserContextKey)

	user := app.getUserFromContext(r, userContextKey)

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, user.ID); err != nil {
		switch err {
		case store.ConflictError:
			app.conflictError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnFollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	FollowUser	true	"User to unfollow"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/user/unfollow [post]
func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := app.getUserFromContext(r, paramUserContextKey)
	user := app.getUserFromContext(r, userContextKey)

	if err := app.store.Followers.UnFollow(r.Context(), unFollowedUser.ID, user.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Get the authenticated user's feed
//	@Tags			users
//	@Produce		json
//	@Param			limit	query		int		false	"Number of posts to return"
//	@Param			offset	query		int		false	"Number of posts to skip"
//	@Param			sort	query		string	false	"Sort order"
//	@Success		200		{array}		store.Post
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/user/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.getUserFromContext(r, userContextKey)
	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(user.ID), fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Activate User godoc
//
//	@Summary		Activate a user
//	@Description	Activate a user account using an activation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path	string	true	"Activation token"
//	@Success		204
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/user/activate/{token} [get]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.NotFoundError:
			app.notFoundError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusAccepted)
}
