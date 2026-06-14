package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"social/internal/store"
	"strconv"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type postkey string

const PostCtx postkey = "post"

type userkey string

const UserCtx userkey = "user"

type selfKey string 
const SelfCtx selfKey = "self"


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

// Basic authentication middleware
func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//step-1 read authorization header
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				//no header
				app.unAuthorizedBasicError(w, r, fmt.Errorf("Authorization header is missing"))
				return
			}
			//step-2 Split header into parts

			// header looks like: "Basic YWRtaW46cGFzc3dvcmQxMjM="
			// after split:       ["Basic", "YWRtaW46cGFzc3dvcmQxMjM="]
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				//wrong format or auth type
				app.unAuthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}
			//step-3  decode base64

			// "YWRtaW46cGFzc3dvcmQxMjM=" → "admin:password123"
			decoded, err := base64.StdEncoding.DecodeString(parts[1])

			if err != nil {
				app.unAuthorizedBasicError(w, r, err)
				return
			}

			//step-4 check credentials
			//Getting expected credentials from config
			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			//splitting the decoded string  by :
			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				//wrong username or password
				app.unAuthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			//step-5 everything ok, allow the request to go through
			next.ServeHTTP(w, r)

		})
	}
}


//protect middleware 
func (app *application) protect(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){

		//step -1 get token from auth header 
		authHeader := r.Header.Get("Authorization")

		if authHeader == ""{
			app.unAuthorizedError(w,r,fmt.Errorf("Authorization header is missing"))
			return 
		}

		//step-2 split the header into parts 
		parts := strings.Split(authHeader," ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unAuthorizedError(w,r,fmt.Errorf("Invalid authrorization Header"))
			return 
		}

		token := parts[1]

		//step-3 verify token 
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unAuthorizedError(w,r,fmt.Errorf("Invalid Token"))
			return 
		}

		//get user id from jwtToken 
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID := int64(claims["sub"].(float64))

		 ctx := r.Context()
		 //Get user from redis cache 
		// user,err := app.store.Users.GetById(ctx,userID)
		user , err := app.getUserFromCache(ctx, userID)
		
		if err != nil {
			switch  {
			case errors.Is(err, store.ErrConflict):
				app.unAuthorizedError(w,r,fmt.Errorf("User not found"))
				return
			default:
				app.internalServerError(w,r,err)
				return 
			}
		}

		//Store the. user in context 
		ctx = context.WithValue(ctx,SelfCtx, user)
		next.ServeHTTP(w,r.WithContext(ctx))
	} )
}

//function to get user and store it in cache 
func (app *application) getUserFromCache(ctx context.Context, userID int64) (*store.User, error){
	//1) if redis is disabled give data from database
	if !app.config.redisCfg.enabled{
		return app.store.Users.GetById(ctx,userID)
	}
	//2) If redis is enabled search for the data and return it
	user, err := app.cacheStorage.Users.GetUser(ctx,userID)
	if err != nil {
		return nil, err
	}

	//3) If user == nil then we must add it in db and fetch from db 
	if user == nil {
		user, err = app.store.Users.GetById(ctx, userID)
		if err != nil {
			return nil, err
		}
		//4) Store the data in cache 
		if err := app.cacheStorage.Users.SetUser(ctx, user); err != nil {
			return nil, err
		}
	}
	return user,nil
}

//middleware for authorization, Check post ownership 
func (app *application) checkPostOwnerShip(requiredRole string ,next http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get user from jwt 
		user  := getselfFromContext(r)
		post := getPostFromContext(r)

		//if the post belongs to user allow it
		if (user.ID == post.UserId){
			next.ServeHTTP(w,r)
			return 
		}
		ctx := r.Context()
		//the post doesn't belong to user so we check roles
		allowed, err  := app.checkRolePrecedence(ctx, user, requiredRole)

		if err != nil {
			app.internalServerError(w,r, err)
			return
		}

		if !allowed {
			app.forbiddenError(w,r, fmt.Errorf("You are not allowed to Perform this action"))
			return 
		}

		//if everything's right go to the next middleware 
		next.ServeHTTP(w,r)
	})
}

//middleware for checking roles 
func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string)(bool, error){
	role, err := app.store.Roles.GetByName(ctx, roleName)

	if err != nil {
		return false, err
	}
	//role.Level gives us what permission we need to delete the role
	//check if user has higher level role than specified role
	return user.Role.Level >= role.Level, nil 
}