package uhttp

import (
	"net/http"
	"testing"
)

func TestWithContextMiddleware(t *testing.T) {
	u := NewUHTTP()

	ctxKey := ContextKey("testKey")

	err := u.AddContext(ctxKey, map[string]string{"addedContext": "testAddedContext"})
	if err != nil {
		t.Error(err)
		return
	}

	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return r.Context().Value(ctxKey)
	}))

	expectedResponseBody := []byte(`{"addedContext":"testAddedContext"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
