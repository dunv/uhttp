package uhttp

import (
	"net/http"
	"testing"
)

func TestParsePostModel(t *testing.T) {
	handler := Handler{
		PostHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Render(w, r, ParsedModel(r)) }),
		PostModel:   map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, NewUHTTP(), t)
}

func TestParseGetModel(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Render(w, r, ParsedModel(r)) }),
		GetModel:   map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, requestBody, expectedResponseBody, NewUHTTP(), t)
}

func TestParseDeleteModel(t *testing.T) {
	handler := Handler{
		DeleteHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Render(w, r, ParsedModel(r)) }),
		DeleteModel:   map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodDelete, http.StatusOK, requestBody, expectedResponseBody, NewUHTTP(), t)
}
