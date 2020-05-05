package uhttp

import (
	"net/http"
	"testing"
)

func TestGzipResponse(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"msg": "testResponse"}
	}))
	expectedResponseBody := []byte(`{"msg":"testResponse"}`)

	ExecuteHandlerWithGzipResponse(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)

}

func TestGzipRequestAndResponse(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithPostModel(map[string]string{}, func(r *http.Request, model interface{}, ret *int) interface{} {
		return model
	}))

	requestBody := []byte(`{"request":"gzipped"}`)
	expectedResponseBody := []byte(`{"request":"gzipped"}`)

	ExecuteHandlerWithGzipRequestAndResponse(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)

}
