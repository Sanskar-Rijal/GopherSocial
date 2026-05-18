package main

import (
	"encoding/json"
	"net/http"
	 "github.com/go-playground/validator/v10"
)


var Validate *validator.Validate

//it runs whenever the program initializes 
func init(){
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

//Write json data
func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

//Reading json data
//we want to read from the request and write into the response 
func readJson(w http.ResponseWriter, r *http.Request, data any) error {

	// Go allows underscores in numbers for readability
	// same number — just easier to read

	//total size of body allowed 
	maxBytes := 1_048_576 //1mb = 1024*1024 = 1,048,576 bytes
	//for 5 mb -> 5 * 1048_576 = 5,242,880 bytes
	r.Body = http.MaxBytesReader(w,r.Body,int64(maxBytes))
	decoder := json.NewDecoder(r.Body)

	//suppose i want only post_id but frontend sends me user_id comment_id and all things
	//let's remove them 
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
