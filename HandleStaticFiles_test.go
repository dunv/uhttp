package uhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSinglePageAppHandlerReturnIndex(t *testing.T) {
	u := NewUHTTP(WithServeMux(http.NewServeMux()))

	tempDir, err := ioutil.TempDir(os.TempDir(), "uhttpTest")
	if err != nil {
		t.Error(err)
		return
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
		return
	}
	f.Close()

	f2, err := os.OpenFile(filepath.Join(tempDir, "main.css"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
	}
	_, err = f2.Write([]byte(".test{ font-weight:bold;"))
	if err != nil {
		f.Close()
		t.Error(err)
		return
	}
	f2.Close()

	// Test for serving index.html when requesting
	err = u.RegisterStaticFilesHandler(tempDir)
	if err != nil {
		t.Errorf("could not register staticFilesHandler (%s)", err)
		return
	}

	req := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	staticFilesHandler(u)(w, req)
	res := w.Result()
	response, _ := ioutil.ReadAll(res.Body)

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
	u := NewUHTTP(WithServeMux(http.NewServeMux()))

	tempDir, err := ioutil.TempDir(os.TempDir(), "uhttpTest")
	if err != nil {
		t.Error(err)
		return
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
		return
	}
	f.Close()

	f2, err := os.OpenFile(filepath.Join(tempDir, "main.css"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Error(err)
	}
	_, err = f2.Write([]byte(".test{ font-weight:bold;}"))
	if err != nil {
		f.Close()
		t.Error(err)
		return
	}
	f2.Close()

	// Test for serving index.html when requesting
	err = u.RegisterStaticFilesHandler(tempDir)
	if err != nil {
		t.Errorf("could not register staticFilesHandler (%s)", err)
		return
	}
	req := httptest.NewRequest("GET", "http://example.com/main.css", nil)
	w := httptest.NewRecorder()
	staticFilesHandler(u)(w, req)
	res := w.Result()
	response, _ := ioutil.ReadAll(res.Body)

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
