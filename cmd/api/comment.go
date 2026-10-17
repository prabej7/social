package main

import (
	"net/http"

	"github.com/prabej7/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
	UserId  int64  `json:"user_id" validate:"required"`
	PostId  int64  `json:"post_id" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var commentPayload = ctx.Value(BodyKey).(CreateCommentPayload)

	post, err := app.store.Posts.GetById(ctx, commentPayload.PostId)

	// we could also check for user existence

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var comment store.Comment = store.Comment{
		PostID:  post.ID,
		UserID:  commentPayload.UserId,
		Content: commentPayload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), &comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
