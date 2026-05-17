package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request){

	//Most basic way of sending json response without using any library

	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{
	// "status":"true",
	// "message":"Server is working properly"
	// }`))
	data:= map[string]string{
		"status":"true",
		"message":"Server is working properly", 
		"env": app.config.env,
	}

	err:=  writeJson(w, http.StatusOK, data);
	if err != nil {
		app.internalServerError(w,r,err)
	}
}