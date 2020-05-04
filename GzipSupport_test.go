package uhttp

import (
	"net/http"
	"testing"
)

func TestGzipResponse(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				u.RenderMessage(w, r, "testResponse")
			})
		},
	}
	expectedResponseBody := []byte(`{"msg":"testResponse"}`)

	ExecuteHandlerWithGzipResponse(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)

}

func TestGzipRequestAndResponse(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		PostHandler: func(u *UHTTP) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				u.Render(w, r, ParsedModel(r))
			})
		},
		PostModel: map[string]string{},
	}
	requestBody := []byte(`{"request":"gzipped"}`)
	expectedResponseBody := []byte(`{"request":"gzipped"}`)

	ExecuteHandlerWithGzipRequestAndResponse(handler, http.MethodPost, http.StatusOK, requestBody, expectedResponseBody, u, t)

}
