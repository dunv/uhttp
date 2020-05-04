package uhttp

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dunv/ulog"
)

func TestRender(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, map[string]string{"test": "test"})
		}),
	}
	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestRenderWithStatusCode(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderWithStatusCode(w, r, http.StatusCreated, map[string]string{"test": "test"})
		}),
	}
	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusCreated, nil, expectedResponseBody, u, t)
}

func TestRenderError(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderError(w, r, errors.New("testError"))
		}),
	}
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadRequest, nil, expectedResponseBody, u, t)
}

func TestRenderErrorWithStatusCode(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderErrorWithStatusCode(w, r, http.StatusBadGateway, errors.New("testError"))
		}),
	}
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadGateway, nil, expectedResponseBody, u, t)
}

func TestRenderMessage(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessage(w, r, "test")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestRenderMessageWithStatusCode(t *testing.T) {
	u := NewUHTTP()
	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessageWithStatusCode(w, r, http.StatusConflict, "test")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusConflict, nil, expectedResponseBody, u, t)
}

func TestRenderMessageWithStatusCodeAndLogLevelOverride(t *testing.T) {
	u := NewUHTTP(WithEncodingErrorLogLevel(ulog.LEVEL_INFO))

	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RenderMessageWithStatusCode(w, r, http.StatusConflict, "test")
		}),
	}
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusConflict, nil, expectedResponseBody, u, t)
}
