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

	ExecuteHandler(handler, http.MethodPost, requestBody, expectedResponseBody, t)
}

func TestParseGetModel(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Render(w, r, ParsedModel(r)) }),
		GetModel:   map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodGet, requestBody, expectedResponseBody, t)
}

func TestParseDeleteModel(t *testing.T) {
	handler := Handler{
		DeleteHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Render(w, r, ParsedModel(r)) }),
		DeleteModel:   map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodDelete, requestBody, expectedResponseBody, t)
}
