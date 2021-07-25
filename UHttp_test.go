package uhttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRendering(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hallo": "welt"}
	}))

	ts := httptest.NewServer(handler.HandlerFunc(u))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	actual := string(greeting)
	expected := `{"hallo":"welt"}` + "\n"

	if strings.Compare(actual, expected) != 0 {
		t.Errorf("could not produce valid JSON response. Expected: %s, Actual: %s", expected, actual)
	}
}

func TestJSONResponse(t *testing.T) {
	u := NewUHTTP()
	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hallo": "welt"}
	}))

	ts := httptest.NewServer(handler.HandlerFunc(u))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("did not set Content-Type as expeceted")
	}
}

// TODO: write more tests
