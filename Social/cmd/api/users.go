package main

import (
	"fmt"
	"net/http"
	"social/internal/store"
)

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(UserCtx).(*store.User)
	return user
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request){
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

	data := map[string]any {
		"status":"true",
		"message":user,
		"env":app.config.env,
	}

	if err := writeJson(w,http.StatusOK, data); err != nil {
		app.internalServerError(w,r,err)
	}

}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request){
	//this is the user you want to follow 
	user := getUserFromContext(r)

	//this is you, you are followiing someone else so you are the follower and the 
	//user in the url is the one you are following
	var followerId int64 = 2; //To get later from jwt

	//you cannot follow yourself 
	if user.ID == followerId {
		app.badRequestError(w,r,fmt.Errorf("You cannot follow yourself"))
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.FollowUser(ctx,followerId, user.ID); err != nil {
		app.internalServerError(w,r,err)
		return 
	}

	if err :=writeJson(w, http.StatusNoContent,nil); err != nil {
		app.internalServerError(w,r,err)
		return
	}
}



func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request){
	//this is the user we want to unfollow 
	user := getUserFromContext(r)

	//this is you, you are unfollowing someone else so you were the follower 
	var followerId int64 =1; //To get later from jwt


	//you cannot unfolow yourself 
	if user.ID == followerId {
		app.badRequestError(w,r, fmt.Errorf("You cannot unfollow yourself"))
		return 
	}

	ctx := r.Context()

	if err := app.store.Followers.UnFollowUser(ctx,followerId, user.ID); err != nil {
		app.internalServerError(w,r,err)
		return 
	}

	if err := writeJson(w,http.StatusNoContent,nil); err != nil {
		app.internalServerError(w,r,err)
		return 
	}

}

func (app *application) getFollowersHandler(w http.ResponseWriter, r *http.Request){
	user := getUserFromContext(r)

	ctx := r.Context()

	followers, err := app.store.Followers.GetFollowers(ctx,user.ID)

	if err != nil {
		app.internalServerError(w,r,err)
		return 
	}

	data := map[string]any {
		"status":"true",
		"message":followers,
		"env":app.config.env,
	}

	if err := writeJson(w,http.StatusOK, data); err != nil {
		app.internalServerError(w,r,err)
		return
	}
}

func (app *application) getFollowingHandler(w http.ResponseWriter, r *http.Request){
	user := getUserFromContext(r)

	ctx := r.Context()

	following, err := app.store.Followers.GetFollowing(ctx,user.ID)

	if err != nil {
		app.internalServerError(w,r,err)
		return
	}

	data := map[string]any {
		"status":"true",
		"message":following, 
		"env":app.config.env,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w,r,err)
		return 
	}
}