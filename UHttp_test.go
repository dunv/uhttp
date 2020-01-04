package uhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRendering(t *testing.T) {
	tmp := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, map[string]string{
				"hallo": "welt",
			})
		}),
	}

	ts := httptest.NewServer(tmp.HandlerFunc())
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

	actual := fmt.Sprintf("%s", greeting)
	expected := `{"hallo":"welt"}` + "\n"

	if strings.Compare(actual, expected) != 0 {
		t.Errorf("could not produce valid JSON response. Expected: %s, Actual: %s", expected, actual)
	}
}

func TestJSONResponse(t *testing.T) {
	tmp := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	ts := httptest.NewServer(tmp.HandlerFunc())
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
