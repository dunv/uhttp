package uhttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dunv/uhttp"
	"github.com/stretchr/testify/require"
)

func TestErrorModel(t *testing.T) {
	u := uhttp.NewUHTTP()

	handler := uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			return errors.New("err from handler")
		}),
	)
	u.Handle("/test", handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	parsedErr, err := uhttp.ErrorFromHttpResponseBody(res.Body)
	require.NoError(t, err)
	require.Equal(t, "err from handler", parsedErr.Error())
}
