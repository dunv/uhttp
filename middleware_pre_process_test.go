package uhttp_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
	"github.com/stretchr/testify/require"
)

func TestPreProcess(t *testing.T) {
	u := uhttp.NewUHTTP()

	originalMessage := "??"
	handler := uhttp.NewHandler(
		uhttp.WithPreProcess(func(ctx context.Context) error {
			originalMessage = "world"
			return nil
		}),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": originalMessage}
		}),
	)
	u.Handle("/test", handler)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/test", nil, `{"hello": "world"}`)
}

func TestPreError(t *testing.T) {
	u := uhttp.NewUHTTP()
	originalMessage := "??"
	handler := uhttp.NewHandler(
		uhttp.WithPreProcess(func(ctx context.Context) error {
			return errors.New("did not work")
		}),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": originalMessage}
		}),
	)
	u.Handle("/test", handler)

	require.HTTPError(t, u.ServeMux().ServeHTTP, http.MethodGet, "/test", nil)
}
