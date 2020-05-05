package uhttp

import (
	"net/http"
	"testing"
)

func TestSelectMethodNotFound(t *testing.T) {
	u := NewUHTTP()
	// All success cases are already tested by other tests
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"test": "test"}
	}))

	expectedResponseBody := []byte(`{"error":"method not allowed"}`)
	ExecuteHandler(handler, http.MethodPost, http.StatusMethodNotAllowed, nil, expectedResponseBody, u, t)
}
