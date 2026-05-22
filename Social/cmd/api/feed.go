package main

import (
	"net/http"
	"social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request){

	//pagination, filter, sorting 
	feedQuery := &store.PaginatedQuery{
		Limit: 10, //Default values
		Offset: 0,
		Sort: "desc",
	}

	if err := feedQuery.Parse(r); err != nil {
		app.badRequestError(w,r,err)
		return
	}

	//Validate 
	if err := Validate.Struct(feedQuery); err != nil {
		app.badRequestError(w,r,err)
		return
	}

	ctx := r.Context()
	var userID int64 = 1 //get from JWT later

	posts, err := app.store.Posts.GetUserFeed(ctx, userID, feedQuery)

	if err != nil {
		app.internalServerError(w,r,err)
		return 
	}
	
	data := map[string]any{
		"status":"true",
		"message":posts,
		"env":app.config.env,
	}

	if err := writeJson(w,http.StatusOK, data); err != nil {
		app.internalServerError(w,r,err)
		return
	}
}