package main

import (
	"net/http"
	"social/internal/store"
)

type CreatePostPayload struct {
	Content string `json:"content"`
	Title string   `json:"title"`
	Tags []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request){

	var userId int64 = 1 

	var payload CreatePostPayload

	err := readJson(w,r, &payload)

	if err != nil {
		writeJson(w,http.StatusBadRequest, err.Error())
		return 
	}

	post := &store.Post{
		Content: payload.Content,
		Title: payload.Title,
		Tags : payload.Tags,
		//TODO : we will get user id using JWT later on
		UserId : userId,
	}
	
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJson(w,http.StatusInternalServerError, err.Error())
		return 
	}

	data := map[string]any{
		"status":"true",
		"message":post,
		"env":app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		errorJson(w, http.StatusInternalServerError, err.Error())
		return
	}

}