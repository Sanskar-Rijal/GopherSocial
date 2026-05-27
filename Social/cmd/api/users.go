package main

import (
	"errors"
	"fmt"
	"net/http"
	"social/internal/store"

	"github.com/go-chi/chi/v5"
)

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(UserCtx).(*store.User)
	return user
}

type getUserByIdHandlerResponse = SuccessResponse[store.User]

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	getUserByIdHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	// idParam := chi.URLParam(r,"userID")
	// userId, err := strconv.ParseInt(idParam, 10, 64)

	// if err != nil {
	// 	app.badRequestError(w,r,err)
	// 	return
	// }

	// ctx := r.Context()

	// user, err := app.store.Users.GetById(ctx, userId)

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, store.ErrNotFound):
	// 		app.notFoundError(w,r,err)
	// 		return
	// 	default:
	// 		app.internalServerError(w,r,err)
	// 		return
	// 	}
	// }

	user := getUserFromContext(r)

	data := map[string]any{
		"status":  "true",
		"message": user,
		"env":     app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}

}

// FollowUser godoc
//
//	@Summary		Follow a user profile
//	@Description	You’re free to follow any user.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	//this is the user you want to follow
	user := getUserFromContext(r)

	//this is you, you are followiing someone else so you are the follower and the
	//user in the url is the one you are following
	var followerId int64 = 1 //To get later from jwt

	//you cannot follow yourself
	if user.ID == followerId {
		app.badRequestError(w, r, fmt.Errorf("You cannot follow yourself"))
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.FollowUser(ctx, followerId, user.ID); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.badRequestError(w, r, fmt.Errorf("You have already followed this user"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJson(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnFollowUser godoc
//
//	@Summary		UnFollow a user profile
//	@Description	You’re free to unfollow any user.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [post]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	//this is the user we want to unfollow
	user := getUserFromContext(r)

	//this is you, you are unfollowing someone else so you were the follower
	var followerId int64 = 1 //To get later from jwt

	//you cannot unfolow yourself
	if user.ID == followerId {
		app.badRequestError(w, r, fmt.Errorf("You cannot unfollow yourself"))
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.UnFollowUser(ctx, followerId, user.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type getFollowersHandlerResponse = SuccessResponse[[]store.Cuser]

// getFollowers godoc
//
//	@Summary		GetFollowers
//	@Description	You can see all the followers of the user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		204		{object}	getFollowersHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/getfollowers [get]
func (app *application) getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	ctx := r.Context()

	followers, err := app.store.Followers.GetFollowers(ctx, user.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := map[string]any{
		"status":  "true",
		"message": followers,
		"env":     app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type getFollowingHandlerResponse = SuccessResponse[[]store.Cuser]

// getFollowing godoc
//
//	@Summary		GetFollowing
//	@Description	You can see all followers of any user, including yourself.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		204		{object}	getFollowersHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/getfollowing [get]
func (app *application) getFollowingHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	ctx := r.Context()

	following, err := app.store.Followers.GetFollowing(ctx, user.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	data := map[string]any{
		"status":  "true",
		"message": following,
		"env":     app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [post]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	//getting token from the url
	token := chi.URLParam(r, "token")

	ctx := r.Context()

	if err := app.store.Users.ActivateUser(ctx, token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	data := map[string]any{
		"env":     app.config.env,
		"message": "User Activated Successfully",
		"status":  "true",
	}

	if err := writeJson(w, http.StatusNoContent, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
