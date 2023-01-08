package uhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dunv/ulog"
	"github.com/stretchr/testify/require"
)

func TestCORS(t *testing.T) {
	ulog.SetWriter(io.Discard, nil)
	u := NewUHTTP()
	handler1 := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": "world"}
		}),
	)
	u.Handle("/test", handler1)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "my-header")
	require.NoError(t, err)
	u.ServeMux().ServeHTTP(w, req)

	require.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "my-header")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "true")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Max-Age"), "86400")

}

func TestNoCORS(t *testing.T) {
	ulog.SetWriter(io.Discard, nil)
	u := NewUHTTP(
		WithCORS(""),
	)
	handler1 := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": "world"}
		}),
	)
	u.Handle("/test", handler1)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Access-Control-Request-Headers", "my-header")
	require.NoError(t, err)
	u.ServeMux().ServeHTTP(w, req)

	require.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "")
	require.Equal(t, w.Header().Get("Access-Control-Allow-Max-Age"), "")

}
