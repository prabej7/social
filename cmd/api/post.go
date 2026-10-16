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

type UpdatePostPayload struct {
	Title   string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post := r.Context().Value(BodyKey).(UpdatePostPayload)
	fmt.Println(post.Tags)
	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updatedPost, err := app.store.Posts.Update(r.Context(), parsedId, store.Post{
		Title:   post.Title,
		Content: post.Content,
		Tags:    post.Tags,
	})

	if err != nil {
		if errors.Is(store.NotFoundError, err) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, updatedPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Posts.DeleteById(r.Context(), parsedId); err != nil {
		if errors.Is(store.NotFoundError, err) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, map[string]string{
		"message": "post deleted successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
