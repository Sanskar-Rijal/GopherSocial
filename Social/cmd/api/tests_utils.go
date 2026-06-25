package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/store"
	"social/internal/store/cache"
	"testing"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T)*application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	return &application{
		logger: logger,
		store: mockStore,
		cacheStorage: mockCacheStore,
	}
}

func executeRequests(req *http.Request, mux http.Handler) *httptest.ResponseRecorder{
		//creating  fake https response writer 
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr 
}