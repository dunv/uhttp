package uhttp

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dunv/ulog"
)

func ExecuteHandler(
	handler Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	t *testing.T,
) {
	// Suppress log-output
	ulog.SetWriter(bufio.NewWriter(nil), nil)

	ts := httptest.NewServer(handler.HandlerFunc())
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
		Body:   ioutil.NopCloser(bytes.NewReader(requestBody)),
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

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	res.Body.Close()

	expectedWithNewLine := append(expectedResponseBody, []byte("\n")...)

	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}
}

func ExecuteHandlerWithGzipResponse(
	handler Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	t *testing.T,
) {
	// Suppress log-output
	ulog.SetWriter(bufio.NewWriter(nil), nil)

	ts := httptest.NewServer(handler.HandlerFunc())
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
		Body:   ioutil.NopCloser(bytes.NewReader(requestBody)),
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

	response, err := GzipDecodeResponseBody(res)
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

func ExecuteHandlerWithGzipRequestAndResponse(
	handler Handler,
	method string,
	expectedStatus int,
	requestBody []byte,
	expectedResponseBody []byte,
	t *testing.T,
) {
	// Suppress log-output
	ulog.SetWriter(bufio.NewWriter(nil), nil)

	ts := httptest.NewServer(handler.HandlerFunc())
	defer ts.Close()

	url, err := url.Parse(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	encoded, err := GzipEncodeRequestBody(requestBody)
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

	response, err := GzipDecodeResponseBody(res)
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
