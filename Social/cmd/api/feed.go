package main

import (
	"net/http"
	"social/internal/store"
)

type getUserFeedHandlerResponse = SuccessResponse[[]store.PostWithMetaData]

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed using algorithm
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	getUserFeedHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	//pagination, filter, sorting
	feedQuery := &store.PaginatedQuery{
		Limit:  10, //Default values
		Offset: 0,
		Sort:   "desc",
	}

	if err := feedQuery.Parse(r); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	//Validate
	if err := Validate.Struct(feedQuery); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	//var userID int64 = 1 //get from JWT later
	userID := getselfFromContext(r)

	posts, err := app.store.Posts.GetUserFeed(ctx, userID.ID, feedQuery)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := map[string]any{
		"status":  "true",
		"message": posts,
		"env":     app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
