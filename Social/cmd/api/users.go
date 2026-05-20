package main

import (
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request){
	idParam := chi.URLParam(r,"userID")
	userId, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.badRequestError(w,r,err)
		return
	}

	ctx := r.Context()

	user, err := app.store.Users.GetById(ctx, userId)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w,r,err)
			return
		default: 
			app.internalServerError(w,r,err)
			return
		}
	}

	data := map[string]any {
		"status":"true",
		"message":user,
		"env":app.config.env,
	}

	if err := writeJson(w,http.StatusOK, data); err != nil {
		app.internalServerError(w,r,err)
	}

}