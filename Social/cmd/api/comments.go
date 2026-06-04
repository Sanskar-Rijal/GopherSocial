package main

import (
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
	PostID  int64  `json:"post_id" validate:"required"`
}

type addCommentHandlerResponse = SuccessResponse[store.Comment]

// CreatePost godoc
//
//	@Summary		Creates a Comment
//	@Description	Create  a comment on users post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	createPostHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		401		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/comments [post]
//
// Add Comments to the post
func (app *application) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	//to do via jwt later on
	var userId int64 = 1

	var payload CreateCommentPayload

	err := readJson(w, r, &payload)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	//validating the payload sent by the user
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	comment := &store.Comment{
		PostID:  payload.PostID,
		UserID:  userId,
		Content: payload.Content,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	data := map[string]any{
		"status":  "true",
		"message": comment,
		"env":     app.config.env,
	}

	if err := writeJson(w, http.StatusCreated, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeleteComment godoc
//
//	@Summary		Deletes comment from post
//	@Description	Delete a comment by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			commentID	path	int	true	"comment ID"
//	@Success		204			"No Content"
//	@Failure		404			{object}	ErrorResponseWrapper
//	@Failure		500			{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/comments/{commentID} [delete]
//
// Delete post
func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "commentID")
	commentID, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Comments.Delete(ctx, commentID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	//comment deleted successfully
	writeJson(w, http.StatusNoContent, nil)
}
