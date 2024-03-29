package uhttp_test

import (
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
	"github.com/stretchr/testify/require"
)

func TestSelectMethodNotAllowed(t *testing.T) {
	u := uhttp.NewUHTTP()
	// All success cases are already tested by other tests
	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"test": "test"}
	}))
	u.Handle("/test", handler)
	statusCode, body, _, _ := Run(t, u, http.MethodPost, "/test", nil)
	require.Equal(t, http.StatusMethodNotAllowed, statusCode)
	require.Contains(t, body, `{"error":"method not allowed"}`)
}

func TestRecover(t *testing.T) {
	u := uhttp.NewUHTTP()
	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		panic("handlerPanic")
	}))
	u.Handle("/panic", handler)
	statusCode, body, _, _ := Run(t, u, http.MethodGet, "/panic", nil)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Contains(t, body, `{"error":"internal server error"}`)
}

func TestRecoverWithStackTrace(t *testing.T) {
	u := uhttp.NewUHTTP(uhttp.WithSendPanicInfoToClient(true))
	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		panic("handlerPanic")
	}))
	u.Handle("/panic", handler)
	statusCode, body, _, _ := Run(t, u, http.MethodGet, "/panic", nil)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Contains(t, body, `{"error":"panic: handlerExecution (handlerPanic) stackTrace: goroutine`)
}
