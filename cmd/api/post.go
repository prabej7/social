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

const postContextKey contextKey = "post"

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

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostContext(r)

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type UpdatePostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags" validate:"omitempty"`
}

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	post := getPostContext(r)

	postPayload := r.Context().Value(BodyKey).(UpdatePostPayload)

	if postPayload.Content == nil {
		postPayload.Content = &post.Content
	}

	if postPayload.Title == nil {
		postPayload.Title = &post.Title
	}

	updatedPost, err := app.store.Posts.Update(r.Context(), post.ID, store.Post{
		Title:   *postPayload.Title,
		Content: *postPayload.Content,
		Tags:    *postPayload.Tags,
		Version: post.Version,
	})

	if err != nil {
		if errors.Is(store.NotFoundError, err) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, updatedPost); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	postContext := getPostContext(r)

	if err := app.store.Posts.DeleteById(r.Context(), postContext.ID); err != nil {
		if errors.Is(store.NotFoundError, err) {
			app.notFoundError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "post deleted successfully",
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), postContextKey, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostContext(r *http.Request) *store.Post {
	post := r.Context().Value(postContextKey).(*store.Post)
	return post
}
