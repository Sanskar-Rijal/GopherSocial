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

type userkey string

const UserCtx userkey = "user"

// middleware to fetch post and put it into the context of the request so that we can use it in the handlers
func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postId, err := strconv.ParseInt(idParam, 10, 64)

		//if error comes user is sending "abc" instead of 123 numbers etc
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context() // we will use this context to create new context and insert the post

		post, err := app.store.Posts.GetById(ctx, postId)

		if err != nil {
			switch {
			//if not found in database we send 404 error
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			//else return internal server error
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		//we never mutate context, we always create new context from scratch
		ctx = context.WithValue(ctx, PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// middleware to fetch userId from url
func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idParam := chi.URLParam(r, "userID")
		userId, err := strconv.ParseInt(idParam, 10, 64)

		//if error comes then user is sending incorrect format
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		//getting the actual context
		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		//creating new context and adding the user property in it
		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
