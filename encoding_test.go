package uhttp

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dunv/ulog"
	"github.com/stretchr/testify/require"
)

func setupEncodingTest(t *testing.T, enableBrotli, enableGzip, enableDeflate bool) *UHTTP {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	opts := []UhttpOption{}
	if !enableBrotli {
		opts = append(opts, WithBrotliCompression(false, 5))
	}
	if !enableGzip {
		opts = append(opts, WithGzipCompression(false, 5))
	}
	if !enableDeflate {
		opts = append(opts, WithDeflateCompression(false, 5))
	}

	u := NewUHTTP(opts...)
	handler := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"hello": "world"}
		}),
	)
	u.Handle("/test", handler)
	return u
}

func TestEncodingBrotli(t *testing.T) {
	u := setupEncodingTest(t, true, true, true)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "br, gzip, deflate")
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "br", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"hello": "world"}`, string(body))

}

func TestEncodingNoBrotli(t *testing.T) {
	u := setupEncodingTest(t, false, true, true)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "br, gzip, deflate")
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"hello": "world"}`, string(body))

}

func TestEncodingGzip(t *testing.T) {
	u := setupEncodingTest(t, true, true, true)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"hello": "world"}`, string(body))
}

func TestEncodingDeflate(t *testing.T) {
	u := setupEncodingTest(t, true, true, true)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "deflate")
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "deflate", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"hello": "world"}`, string(body))
}

func TestEncodingNoEncoding(t *testing.T) {
	u := setupEncodingTest(t, false, false, false)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "deflate")
	w := httptest.NewRecorder()
	u.ServeMux().ServeHTTP(w, req)
	res := w.Result()
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"hello": "world"}`, string(body))
}
