package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request){

	ctx := r.Context()

	var userID int64 = 1 //get from JWT later

	posts, err := app.store.Posts.GetUserFeed(ctx, userID)

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