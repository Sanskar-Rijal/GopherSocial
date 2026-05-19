package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)


type postkey string
const PostCtx postkey = "post"

//middleware to fetch post and put it into the context of the request so that we can use it in the handlers
func (app *application) postContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r * http.Request){
	idParam := chi.URLParam(r, "postID")
	postId, err := strconv.ParseInt(idParam, 10, 64)

	//if error comes user is sending "abc" instead of 123 numbers etc
	if err != nil {
		app.badRequestError(w,r,err)
		return
	}

	ctx := r.Context() // we will use this context to create new context and insert the post

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
	//we never mutate context, we always create new context from scratch
	ctx = context.WithValue(ctx,PostCtx,post)
	next.ServeHTTP(w,r.WithContext(ctx))
	})
}