package uhttp_test

import (
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
)

func TestGzipResponse(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"msg": "testResponse"}
	}))
	expectedResponseBody := []byte(`{"msg":"testResponse"}`)

	executeHandlerWithGzipResponse(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestGzipRequestAndResponse(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithPostModel(map[string]string{}, func(r *http.Request, model interface{}, ret *int) interface{} {
		return model
	}))

	requestBody := []byte(`{"request":"gzipped"}`)
	expectedResponseBody := []byte(`{"request":"gzipped"}`)

	executeHandlerWithGzipRequestAndResponse(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)
}
