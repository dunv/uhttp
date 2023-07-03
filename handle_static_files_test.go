package uhttp_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dunv/uhttp"
)

func setupSinglePage(t *testing.T) *uhttp.UHTTP {
	u := uhttp.NewUHTTP()

	tempDir, err := os.MkdirTemp(os.TempDir(), "uhttpTest")
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	defer os.RemoveAll(tempDir)

	f, err := os.OpenFile(filepath.Join(tempDir, "index.html"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
	}
	_, err = f.Write([]byte("<html></html>"))
	if err != nil {
		f.Close()
		t.Error(err)
		t.FailNow()
		return nil
	}
	f.Close()

	f2, err := os.OpenFile(filepath.Join(tempDir, "main.css"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	_, err = f2.Write([]byte(".test{ font-weight:bold;}"))
	if err != nil {
		f.Close()
		t.Error(err)
		t.FailNow()
		return nil
	}
	f2.Close()

	// Test for serving index.html when requesting
	err = u.RegisterStaticFilesHandler(tempDir)
	if err != nil {
		t.Errorf("could not register staticFilesHandler (%s)", err)
		t.FailNow()
		return nil
	}

	return u
}

func TestSinglePageAppHandlerReturnIndex(t *testing.T) {
	u := setupSinglePage(t)

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	uhttp.StaticFilesHandler(u)(w, req)
	res := w.Result()
	response, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("did not return http %d (actual: %d)", http.StatusOK, res.StatusCode)
		return
	}
	expectedWithNewLine := []byte(`<html></html>`)
	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}

}

func TestSinglePageAppHandlerReturnActualFile(t *testing.T) {
	u := setupSinglePage(t)

	req := httptest.NewRequest("GET", "http://example.com/main.css", nil)
	w := httptest.NewRecorder()
	uhttp.StaticFilesHandler(u)(w, req)
	res := w.Result()
	response, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Errorf("did not return http %d (actual: %d)", http.StatusOK, res.StatusCode)
		return
	}
	expectedWithNewLine := []byte(`.test{ font-weight:bold;}`)
	if !bytes.Equal(expectedWithNewLine, response) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, response)
		return
	}

}

func TestSinglePageAppHandlerReturnActualFileGzip(t *testing.T) {
	u := setupSinglePage(t)

	req := httptest.NewRequest("GET", "http://example.com/main.css", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	uhttp.StaticFilesHandler(u)(w, req)
	res := w.Result()
	responseDecoded, err := uhttp.DecodeResponseBody(res)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("did not return http %d (actual: %d)", http.StatusOK, res.StatusCode)
		return
	}
	expectedWithNewLine := []byte(`.test{ font-weight:bold;}`)

	if !bytes.Equal(expectedWithNewLine, responseDecoded) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, responseDecoded)
		return
	}

}

func TestSinglePageAppHandlerReturnActualFileBrotli(t *testing.T) {
	u := setupSinglePage(t)

	req := httptest.NewRequest("GET", "http://example.com/main.css", nil)
	req.Header.Add("Accept-Encoding", "br")
	w := httptest.NewRecorder()
	uhttp.StaticFilesHandler(u)(w, req)
	res := w.Result()
	responseDecoded, err := uhttp.DecodeResponseBody(res)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("did not return http %d (actual: %d)", http.StatusOK, res.StatusCode)
		return
	}
	expectedWithNewLine := []byte(`.test{ font-weight:bold;}`)

	if !bytes.Equal(expectedWithNewLine, responseDecoded) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, responseDecoded)
		return
	}

}

func TestSinglePageAppHandlerReturnActualFileDeflate(t *testing.T) {
	u := setupSinglePage(t)

	req := httptest.NewRequest("GET", "http://example.com/main.css", nil)
	req.Header.Add("Accept-Encoding", "deflate")
	w := httptest.NewRecorder()
	uhttp.StaticFilesHandler(u)(w, req)
	res := w.Result()
	responseDecoded, err := uhttp.DecodeResponseBody(res)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("did not return http %d (actual: %d)", http.StatusOK, res.StatusCode)
		return
	}
	expectedWithNewLine := []byte(`.test{ font-weight:bold;}`)

	if !bytes.Equal(expectedWithNewLine, responseDecoded) {
		t.Errorf("expected does not match actual (expected: '%s', actual: '%s')", expectedWithNewLine, responseDecoded)
		return
	}

}
