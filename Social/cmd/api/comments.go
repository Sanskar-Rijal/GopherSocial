package main

import (
	"net/http"
	"social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
	PostID int64 `json:"post_id" validate:"required"`
}

//Add Comments to the post 
func(app *application) addCommentHandler(w http.ResponseWriter, r *http.Request){
	//to do via jwt later on 
	var userId int64 = 1
	 
	var payload CreateCommentPayload

	err := readJson(w,r, &payload)

	if err !=nil {
		app.badRequestError(w,r,err)
		return
	}

	//validating the payload sent by the user
	if err := Validate.Struct(payload); err != nil{
		app.badRequestError(w,r,err)
		return 
	}
	comment := &store.Comment{
		PostID: payload.PostID,
		UserID: userId,
		Content: payload.Content,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx,comment); err != nil {
		app.internalServerError(w,r,err)
		return
	}
	data := map[string]any{
		"status":"true",
		"message":comment,
		"env":app.config.env,
	}

	if err := writeJson(w,http.StatusCreated,data); err != nil {
		app.internalServerError(w,r,err)
		return
	}
}

