package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Content string `json:"content"`
	Title string   `json:"title"`
	Tags []string `json:"tags"`
}

//Create Post
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request){

	var userId int64 = 1 

	var payload CreatePostPayload

	err := readJson(w,r, &payload)

	if err != nil {
		app.internalServerError(w,r,err)
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
		app.internalServerError(w,r,err)
		return
	}

}


//GetPost by Id 
func (app *application) getPostByIdHandler(w http.ResponseWriter, r *http.Request){

	idParam := chi.URLParam(r, "postID")
	postId, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, err.Error())
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetById(ctx, postId)

	if err != nil {
		switch{
			//if not found in database we send 404 error
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w,r, err)
		//else return internal server error 
		default: 
			writeJson(w, http.StatusInternalServerError, err.Error())
		}
	return 
	}

	//send the data to the user 
	if err := writeJson(w, http.StatusOK, map[string]any{
		"status":"true",
		"message":post,
		"env":app.config.env,
	}); err != nil {
		writeJson(w, http.StatusInternalServerError, err.Error())
		return 
	}

}