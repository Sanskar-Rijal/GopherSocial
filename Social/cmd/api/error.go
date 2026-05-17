package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error){
	//500 error
	log.Printf("internal server error: %s path: %s error: %s",r.Method, r.URL.Path, err.Error());
	writeJson(w, http.StatusInternalServerError, map[string]string{
		"status":"false",
		"message":"internal server errors",
		"env":app.config.env,
	})
}

func (app * application) badRequestError(w http.ResponseWriter, r *http.Request, err error){
	//400 error 
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL, err.Error())

	writeJson(w, http.StatusBadRequest, map[string]string{
		"status":"false",
		"message":err.Error(),
		"env":app.config.env,
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r* http.Request, err error){
	// 404 error 
		log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL, err.Error())

		writeJson(w, http.StatusNotFound, map[string]string{
			"status":"false",
			"message":err.Error(),
			"env":app.config.env,
		})
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r* http.Request, err error){
		log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL, err.Error())

		writeJson(w,http.StatusUnauthorized,map[string]string{
			"status":"false",
			"message":err.Error(),
			"env":app.config.env,
		} )
}

func (app *application) methodNotAllowedError(w http.ResponseWriter, r *http.Request, err error){
	    log.Printf("method not allowed: %s path: %s error: %s",r.Method, r.URL.Path, err.Error())

		writeJson(w, http.StatusMethodNotAllowed, map[string]string {
			"status":"false",
			"message":err.Error(),
			"env":app.config.env,
		})
}