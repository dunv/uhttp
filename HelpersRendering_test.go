package uhttp

import (
	"errors"
	"net/http"
	"testing"
)

func TestRender(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, map[string]string{"test": "test"})
		}),
	}
	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, t)
}

func TestRenderWithStatusCode(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderWithStatusCode(w, r, http.StatusCreated, map[string]string{"test": "test"})
		}),
	}
	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusCreated, nil, expectedResponseBody, t)
}

func TestRenderError(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderError(w, r, errors.New("testError"))
		}),
	}
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadRequest, nil, expectedResponseBody, t)
}

func TestRenderErrorWithStatusCode(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderErrorWithStatusCode(w, r, http.StatusBadGateway, errors.New("testError"))
		}),
	}
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadGateway, nil, expectedResponseBody, t)
}

func TestRenderMessage(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessage(w, r, "test")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, t)
}

func TestRenderMessageWithStatusCode(t *testing.T) {
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessageWithStatusCode(w, r, http.StatusConflict, "test")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusConflict, nil, expectedResponseBody, t)
}
