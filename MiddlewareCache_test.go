package uhttp

import (
	"bufio"
	"net/http"
	"testing"
	"time"

	"github.com/dunv/ulog"
)

func TestCacheHit(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	)
	counter2 := 0
	handler2 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	)
	u.Handle("/cachedHandler1", handler1)
	u.Handle("/cachedHandler2", handler2)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
}

func TestCacheExpiry(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	)
	counter2 := 0
	handler2 := NewHandler(
		WithCache(1*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	)
	u.Handle("/cachedHandler1", handler1)
	u.Handle("/cachedHandler2", handler2)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
	time.Sleep(2 * time.Second)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 2}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 2}`)
}

func TestCacheClear(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP(WithExposeCacheHandlers())
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	)
	counter2 := 0
	handler2 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	)
	u.Handle("/cachedHandler1", handler1)
	u.Handle("/cachedHandler2", handler2)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1":1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2":1}`)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", nil, `{"deletedEntries":2}`)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1":2}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2":2}`)
}

func TestCacheClearSpecific(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP(WithExposeCacheHandlers())
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	)
	counter2 := 0
	handler2 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	)
	u.Handle("/cachedHandler1", handler1)
	u.Handle("/cachedHandler2", handler2)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear/cachedHandler1", nil, `{"deletedEntries": 1}`)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 2}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
}

func TestCacheAutomatic(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithAutomaticCacheUpdates(200*time.Millisecond, nil),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	)
	counter2 := 0
	handler2 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	)
	u.Handle("/cachedHandler1", handler1)
	u.Handle("/cachedHandler2", handler2)

	// wait for initial update to have run through
	time.Sleep(50 * time.Millisecond)

	// cache should be initialized automatically for handler1
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`, map[string][]string{
		CACHE_HEADER: {"true"},
	})

	// no cache should be present for handler2
	RequireHTTPBodyAndNotHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`, []string{CACHE_HEADER})

	// wait for automatic update
	time.Sleep(200 * time.Millisecond)

	// automatic update should have happened in the background
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 2}`, map[string][]string{
		CACHE_HEADER: {"true"},
	})

	// cache for handler2 is still the old one
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`, map[string][]string{
		CACHE_HEADER: {"true"},
	})
}
