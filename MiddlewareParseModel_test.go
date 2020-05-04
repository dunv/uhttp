package uhttp

import (
	"net/http"
	"testing"
)

func TestParsePostModel(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		PostHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { u.Render(w, r, ParsedModel(r)) })
		},
		PostModel: map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseGetModel(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { u.Render(w, r, ParsedModel(r)) })
		},
		GetModel: map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseDeleteModel(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		DeleteHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { u.Render(w, r, ParsedModel(r)) })
		},
		DeleteModel: map[string]string{},
	}
	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodDelete, http.StatusOK, requestBody, expectedResponseBody, u, t)
}
