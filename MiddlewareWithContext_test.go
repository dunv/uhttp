package uhttp

import (
	"net/http"
	"testing"
)

func TestWithContextMiddleware(t *testing.T) {
	u := NewUHTTP()

	err := u.AddContext("testKey", map[string]string{"addedContext": "testAddedContext"})
	if err != nil {
		t.Error(err)
		return
	}

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
