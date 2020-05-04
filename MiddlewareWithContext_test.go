package uhttp

import (
	"net/http"
	"testing"
)

func TestWithContextMiddleware(t *testing.T) {
	u := NewUHTTP()
	AddContext("testKey", map[string]string{"addedContext": "testAddedContext"})

	handler := Handler{
		GetHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				u.Render(w, r, r.Context().Value("testKey"))
			})
		},
	}
	expectedResponseBody := []byte(`{"addedContext":"testAddedContext"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
