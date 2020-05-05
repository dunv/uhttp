package uhttp

import (
	"net/http"
	"testing"
)

func TestParsePostModel(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithPostModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseGetModel(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGetModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseDeleteModel(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithDeleteModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	ExecuteHandler(handler, http.MethodDelete, http.StatusOK, requestBody, expectedResponseBody, u, t)
}
