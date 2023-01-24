package uhttp

import (
	"errors"
	"net/http"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestRender(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"test": "test"}
	}))

	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestRenderWithStatusCode(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		*ret = http.StatusCreated
		return map[string]string{"test": "test"}
	}))

	expectedResponseBody := []byte(`{"test":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusCreated, nil, expectedResponseBody, u, t)
}

func TestRenderError(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return errors.New("testError")
	}))
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadRequest, nil, expectedResponseBody, u, t)
}

func TestRenderErrorWithStatusCode(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		*ret = http.StatusBadGateway
		return errors.New("testError")
	}))
	expectedResponseBody := []byte(`{"error":"testError"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusBadGateway, nil, expectedResponseBody, u, t)
}

func TestRenderMessageWithStatusCodeAndLogLevelOverride(t *testing.T) {
	u := NewUHTTP(WithEncodingErrorLogLevel(zapcore.InfoLevel))

	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		*ret = http.StatusConflict
		return map[string]string{"msg": "test"}
	}))
	expectedResponseBody := []byte(`{"msg":"test"}`)
	ExecuteHandler(handler, http.MethodGet, http.StatusConflict, nil, expectedResponseBody, u, t)
}
