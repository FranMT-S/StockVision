package services_test

import (
	"net/http"
	"net/http/httptest"
)

func initMockServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	return ts
}
