package uhttp

import (
	"bufio"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dunv/ulog"
	"github.com/stretchr/testify/require"
)

func TestPreProcess(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()

	originalMessage := "??"
	handler := NewHandler(
		WithPreProcess(func(ctx context.Context) error {
			originalMessage = "world"
			return nil
		}),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": originalMessage}
		}),
	)
	u.Handle("/test", handler)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/test", nil, `{"hello": "world"}`)
}

func TestPreError(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	originalMessage := "??"
	handler := NewHandler(
		WithPreProcess(func(ctx context.Context) error {
			return errors.New("did not work")
		}),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": originalMessage}
		}),
	)
	u.Handle("/test", handler)

	require.HTTPError(t, u.ServeMux().ServeHTTP, http.MethodGet, "/test", nil)
}
