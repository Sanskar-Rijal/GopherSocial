package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	//500 error
	app.logger.Errorw("internal server error","method",r.Method,"path",r.URL.Path,"error", err.Error())

	writeJson(w, http.StatusInternalServerError, map[string]string{
		"status":  "false",
		"message": "internal server errors",
		"env":     app.config.env,
	})
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	//400 error
	app.logger.Warnf("bad Requset error","method",r.Method,"path",r.URL.Path,"error", err.Error())

	writeJson(w, http.StatusBadRequest, map[string]string{
		"status":  "false",
		"message": err.Error(),
		"env":     app.config.env,
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	// 404 error

	app.logger.Warnf("not Found error","method",r.Method,"path",r.URL.Path,"error", err.Error())


	writeJson(w, http.StatusNotFound, map[string]string{
		"status":  "false",
		"message": err.Error(),
		"env":     app.config.env,
	})
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("unAuthorizedError error","method",r.Method,"path",r.URL.Path,"error", err)


	writeJson(w, http.StatusUnauthorized, map[string]string{
		"status":  "false",
		"message": err.Error(),
		"env":     app.config.env,
	})
}

func (app *application) methodNotAllowedError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("method not allowed: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	app.logger.Warnf("methodNotAllowed Error","method",r.Method,"path",r.URL.Path,"error", err.Error())


	writeJson(w, http.StatusMethodNotAllowed, map[string]string{
		"status":  "false",
		"message": err.Error(),
		"env":     app.config.env,
	})
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("Conflict error","method",r.Method,"path",r.URL.Path,"error", err.Error())


	writeJson(w, http.StatusConflict, map[string]string{
		"status":  "false",
		"message": err.Error(),
		"env":     app.config.env,
	})
}
