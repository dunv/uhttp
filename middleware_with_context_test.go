package uhttp_test

import (
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
)

func TestWithContextMiddleware(t *testing.T) {
	u := uhttp.NewUHTTP()

	ctxKey := uhttp.ContextKey("testKey")

	err := u.AddContext(ctxKey, map[string]string{"addedContext": "testAddedContext"})
	if err != nil {
		t.Error(err)
		return
	}

	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return r.Context().Value(ctxKey)
	}))

	expectedResponseBody := []byte(`{"addedContext":"testAddedContext"}`)

	executeHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
