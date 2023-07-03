package uhttp_test

import (
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
)

func TestParsePostModel(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithPostModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	executeHandler(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseGetModel(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithGetModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	executeHandler(handler, http.MethodGet, http.StatusOK, requestBody, expectedResponseBody, u, t)
}

func TestParseDeleteModel(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithDeleteModel(
		map[string]string{},
		func(r *http.Request, model interface{}, ret *int) interface{} {
			return model
		},
	))

	requestBody := []byte(`{"test":"test"}`)
	expectedResponseBody := []byte(`{"test":"test"}`)

	executeHandler(handler, http.MethodDelete, http.StatusOK, requestBody, expectedResponseBody, u, t)
}
