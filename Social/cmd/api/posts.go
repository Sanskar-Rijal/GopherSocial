package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)


//getting post from context 
func getPostFromContext(r *http.Request) *store.Post{
	post, _ := r.Context().Value(PostCtx).(*store.Post)
	return post
}


type CreatePostPayload struct {
	Content string `json:"content" validate:"required,max=500"`
	Title string   `json:"title" validate:"required,max=1000"`
	Tags []string `json:"tags"`
}

//Create Post
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request){

	var userId int64 = 1 

	var payload CreatePostPayload

	err := readJson(w,r, &payload)

	if err != nil {
		app.badRequestError(w,r,err)
		return 
	}

	//validating the payload sent by the user
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w,r,err)
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
		app.internalServerError(w,r,err)
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

	// idParam := chi.URLParam(r, "postID")
	// postId, err := strconv.ParseInt(idParam, 10, 64)

	// //if error comes user is sending "abc" instead of 123 numbers etc
	// if err != nil {
	// 	app.badRequestError(w,r,err)
	// 	return
	// }

	// ctx := r.Context()

	// post, err := app.store.Posts.GetById(ctx, postId)

	// if err != nil {
	// 	switch{
	// 		//if not found in database we send 404 error
	// 	case errors.Is(err, store.ErrNotFound):
	// 		app.notFoundError(w,r, err)
	// 	//else return internal server error 
	// 	default: 
	// 		writeJson(w, http.StatusInternalServerError, err.Error())
	// 	}
	// return 
	// }
	
	post := getPostFromContext(r)
	

	//fetch comments for the post 
	comments, err := app.store.Comments.GetCommentsFromPost(r.Context(), post.ID)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	post.Comments = comments

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


//Delete post 
func (app *application) deletePostHandler(w http.ResponseWriter, r * http.Request){
	idParam := chi.URLParam(r,"postID")
	postID, err := strconv.ParseInt(idParam,10,64)

	if err != nil {
		app.badRequestError(w, r, err)
	}

	ctx := r.Context()
	if err := app.store.Posts.Delete(ctx,postID); err != nil {
		switch{
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w,r,err)

		default:
			app.internalServerError(w,r,err)
		}
		return 
	}
	//post deleted successfully 
	writeJson(w,http.StatusNoContent, nil);
}



type UpdatePostPayload struct {
	Title *string `json:"title" validate:"omitempty,max=200"` //using *string, so that if user doesn't send data it becomes nil
	Content *string `json:"content" validate:"omitempty,max=500"`
}

//update post 
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request){
	//getting existing post first 
	existingPost := getPostFromContext(r)

	var payload UpdatePostPayload

	err :=  readJson(w,r, &payload)
	if err != nil {
		app.badRequestError(w,r,err)
		return 
	}
	//validating the payload sent by the user
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w,r,err)
		return
	}

	//checking whether user sent title and content or not and sending it to the databse
	if payload.Title != nil {
		existingPost.Title = *payload.Title
	}
	if payload.Content != nil {
		existingPost.Content = *payload.Content
	}
	ctx := r.Context()

	//now passing payload to the store 
	if err := app.store.Posts.Update(ctx,existingPost); err != nil {
		switch{
		case errors.Is(err, store.ErrConflict):
			app.conflictError(w,r,err) //409 someone else updated first
		default:
			app.internalServerError(w,r,err)
		}
		return 
	}

	data := map[string]any{
		"status":"true",
		"message":existingPost,
		"env":app.config.env,
	}
	if err := writeJson(w,http.StatusCreated, data); err != nil {
		app.internalServerError(w,r,err)
	}
}



