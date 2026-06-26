package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/auth"
	"social/internal/store"
	"social/internal/store/cache"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T)*application {
	t.Helper()

	// logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger: logger,
		store: mockStore,
		cacheStorage: mockCacheStore,
		authenticator: testAuth,
	}
}

func executeRequests(req *http.Request, mux http.Handler) *httptest.ResponseRecorder{
		//creating  fake https response writer 
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr 
}

func checkResponseCode(t *testing.T,expected, actual int){
	if expected != actual {
		t.Errorf("Expected the response code to be %d but the actual is %d",expected,actual)
	}
}