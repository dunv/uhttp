package uhttp

import (
	"net/http"
	"testing"
)

func TestSelectMethodNotFound(t *testing.T) {
	u := NewUHTTP()
	// All success cases are already tested by other tests
	handler := Handler{
		GetHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				u.RenderMessage(w, r, "test")
			})
		},
	}
	expectedResponseBody := []byte(`{"error":"method not allowed"}`)
	ExecuteHandler(handler, http.MethodPost, http.StatusMethodNotAllowed, nil, expectedResponseBody, u, t)
}
