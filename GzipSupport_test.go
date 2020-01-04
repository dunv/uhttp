package uhttp

import (
	"net/http"
	"testing"
)

func TestGzipResponse(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessage(w, r, "testResponse")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"testResponse"}`)

	ExecuteHandlerWithGzipResponse(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, t)

}

func TestGzipRequestAndResponse(t *testing.T) {
	handler := Handler{
		PostHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, ParsedModel(r))
		}),
		PostModel: map[string]string{},
	}
	requestBody := []byte(`{"request":"gzipped"}`)
	expectedResponseBody := []byte(`{"request":"gzipped"}`)

	ExecuteHandlerWithGzipRequestAndResponse(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, t)

}
