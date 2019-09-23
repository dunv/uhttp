package uhttp

import (
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Header)

	actual := fmt.Sprintf("%s", greeting)
	expected := `{"hallo":"welt"}` + "\n"

	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Could not produce valid JSON response. Expected: %s, Actual: %s", expected, actual)
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
		log.Fatal(err)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Did not set Content-Type as expeceted")
	}
}

// TODO: write more tests
