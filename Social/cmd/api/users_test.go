package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T){

	app := newTestApplication(t)
	mux := app.mount()
	t.Run("Should not allow unaunthenticated Requst", func (t *testing.T){
		//check for 401 code 
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1",nil)

		if err != nil {
			t.Fatal(err)
		}
		// //creating  fake https response writer 
		// rr := httptest.NewRecorder()
		// mux.ServeHTTP(rr, req)
		rr := executeRequests(req, mux)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected the response code  to be %d  and we got %d ",http.StatusUnauthorized, rr.Code)
		}
	})
}