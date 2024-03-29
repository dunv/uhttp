package uhttp_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dunv/uhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeHandler(
	handler uhttp.Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	u *uhttp.UHTTP,
	t *testing.T,
) {

	ts := httptest.NewServer(handler.HandlerFunc(u))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	c := http.Client{}
	request := &http.Request{
		Method: method,
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(requestBody)),
	}

	res, err := c.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	if res.StatusCode != expectedStatus {
		t.Errorf("did not return http %d (actual: %d)", expectedStatus, res.StatusCode)
		return
	}

	expectedWithNewLine := append(expectedResponseBody, []byte("\n")...)

	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}
}

func executeHandlerWithGzipResponse(
	handler uhttp.Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	u *uhttp.UHTTP,
	t *testing.T,
) {

	ts := httptest.NewServer(handler.HandlerFunc(u))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	c := http.Client{}
	request := &http.Request{
		Method: method,
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(requestBody)),
		Header: http.Header{"Accept-Encoding": []string{"gzip"}},
	}

	res, err := c.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	if res.StatusCode != expectedStatus {
		t.Errorf("did not return http %d (actual: %d)", expectedStatus, res.StatusCode)
		return
	}

	response, err := uhttp.DecodeResponseBody(res)
	if err != nil {
		t.Error(err)
		return
	}

	expectedWithNewLine := append(expectedResponseBody, []byte("\n")...)
	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}
}

func executeHandlerWithGzipRequestAndResponse(
	handler uhttp.Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	u *uhttp.UHTTP,
	t *testing.T,
) {

	ts := httptest.NewServer(handler.HandlerFunc(u))
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	encoded, err := uhttp.GzipEncodeRequestBody(requestBody)
	if err != nil {
		t.Error(err)
		return
	}

	c := http.Client{}
	request := &http.Request{
		Method: method,
		URL:    url,
		Body:   encoded,
		Header: http.Header{
			"Accept-Encoding":  []string{"gzip"},
			"Content-Encoding": []string{"gzip"},
		},
	}

	res, err := c.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	if res.StatusCode != expectedStatus {
		t.Errorf("did not return http %d (actual: %d)", expectedStatus, res.StatusCode)
		return
	}

	response, err := uhttp.DecodeResponseBody(res)
	if err != nil {
		t.Error(err)
		return
	}

	expectedWithNewLine := append(expectedResponseBody, []byte("\n")...)
	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}
}

func RequireHTTPBodyJSONEq(
	t *testing.T,
	handlerFunc http.HandlerFunc,
	method string,
	url string,
	values url.Values,
	expected string,
) {
	actual := assert.HTTPBody(handlerFunc, method, url, values)
	require.JSONEq(t, expected, actual)
}

// Right out of testify but with headers
func RequireHTTPBodyAndHeader(
	t *testing.T,
	handler http.HandlerFunc,
	method string,
	url string,
	values url.Values,
	expectedBody string,
	expectedHeader http.Header,
) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	require.NoError(t, err)
	handler(w, req)

	require.JSONEq(t, expectedBody, w.Body.String())

	actualHeaders := w.Header()

	for key, expected := range expectedHeader {
		actual := actualHeaders.Values(key)
		require.ElementsMatch(t, expected, actual)
	}
}

func RequireHTTPHeader(
	t *testing.T,
	handler http.HandlerFunc,
	method string,
	url string,
	values url.Values,
	expectedHeader http.Header,
) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	require.NoError(t, err)
	handler(w, req)

	actualHeaders := w.Header()
	for k, header := range actualHeaders {
		fmt.Println("header", k, header)
	}
	for key, expected := range expectedHeader {
		actual := actualHeaders.Values(key)
		require.ElementsMatch(t, expected, actual)
	}
}

func RequireHTTPBodyAndNotHeader(
	t *testing.T,
	handler http.HandlerFunc,
	method string,
	url string,
	values url.Values,
	expectedBody string,
	bannedHeaders []string,
) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url+"?"+values.Encode(), nil)
	require.NoError(t, err)
	handler(w, req)

	require.JSONEq(t, expectedBody, w.Body.String())

	actualHeaders := w.Header()

	for _, bannedHeader := range bannedHeaders {
		if val := actualHeaders.Get(bannedHeader); val != "" {
			t.Errorf("discovered banned header: %s", bannedHeader)
			t.FailNow()
		}
	}
}

func Run(
	t *testing.T,
	u *uhttp.UHTTP,
	method string,
	url string,
	header map[string]string,
) (int, string, http.Header, *http.Response) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)
	for key, val := range header {
		req.Header.Set(key, val)
	}
	u.ServeMux().ServeHTTP(w, req)
	return w.Result().StatusCode, w.Body.String(), w.Header(), w.Result()
}
