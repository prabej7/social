package main

import (
	"net/http"

	"github.com/prabej7/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=100"`
	PostId  int64  `json:"post_id" validate:"required"`
}

// createCommentHandler godoc
//
//	@Summary		Create a comment
//	@Description	Creates a comment for a post.
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			comment		body		CreateCommentPayload	true	"Comment payload"
//	@Success		200			{object}	store.Comment
//	@Failure		500			{object}	map[string]string
//	@Router			/comments	[post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var commentPayload = ctx.Value(BodyKey).(CreateCommentPayload)

	post, err := app.store.Posts.GetById(ctx, commentPayload.PostId)

	// we could also check for user existence

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := app.getUserFromContext(r, userContextKey)

	var comment store.Comment = store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
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
