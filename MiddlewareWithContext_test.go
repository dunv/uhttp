package uhttp

import (
	"net/http"
	"testing"
)

func TestWithContextMiddleware(t *testing.T) {
	AddContext("testKey", map[string]string{"addedContext": "testAddedContext"})

	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, r.Context().Value("testKey"))
		}),
	}
	expectedResponseBody := []byte(`{"addedContext":"testAddedContext"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, t)
}
