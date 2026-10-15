package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prabej7/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload = r.Context().Value(BodyKey).(CreatePostPayload)

	userId := 1

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJson(w, http.StatusCreated, post)
}

func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.store.Posts.GetById(r.Context(), parsedId)

	if err != nil {
		switch {
		case errors.Is(err, store.NotFoundError):
			app.notFoundError(w, r, fmt.Errorf("post not found"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostId(r.Context(), parsedId)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
