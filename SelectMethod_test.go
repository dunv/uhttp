package uhttp

import (
	"net/http"
	"testing"
)

func TestSelectMethodNotFound(t *testing.T) {
	// All success cases are already tested by other tests
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessage(w, r, "test")
		}),
	}
	expectedResponseBody := []byte(`{"error":"method not allowed"}`)
	ExecuteHandler(handler, http.MethodPost, http.StatusMethodNotAllowed, nil, expectedResponseBody, NewUHTTP(), t)
}
