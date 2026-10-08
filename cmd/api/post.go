package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prabej7/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJson(w, r, &payload); err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := 1

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		writeJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, http.StatusCreated, post)
}

func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := app.store.Posts.GetById(r.Context(), parsedId)

	if err != nil {
		switch {
		case errors.Is(err, store.NotFoundError):
			writeJsonError(w, http.StatusNotFound, err.Error())
		default:
			writeJsonError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJson(w, http.StatusOK, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

}
