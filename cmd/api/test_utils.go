package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/na1tto/go-social/internal/auth"
	repository "github.com/na1tto/go-social/internal/store"
	"github.com/na1tto/go-social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := repository.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := auth.NewMockAuthenticator()

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
