package uhttp

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dunv/ulog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`, map[string][]string{})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`, map[string][]string{})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`, map[string][]string{CACHE_HEADER: {"true"}})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`, map[string][]string{CACHE_HEADER: {"true"}})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 1}`, map[string][]string{CACHE_HEADER: {"true"}})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`, map[string][]string{CACHE_HEADER: {"true"}})
}

func TestCacheNoCacheWhenNotOK(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter := 0
	handler := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter++
			return fmt.Errorf("counter:%d", counter)
		}),
	)
	u.Handle("/cachedHandler", handler)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:1"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:2"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:3"}`)
}

func TestCacheForceCacheWhenNotOK(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter := 0
	handler := NewHandler(
		WithCache(10*time.Second),
		WithCacheFailedRequests(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter++
			return fmt.Errorf("counter:%d", counter)
		}),
	)
	u.Handle("/cachedHandler", handler)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:1"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:1"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil, `{"error": "counter:1"}`)
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
	u := NewUHTTP()
	u.ExposeCacheHandlers()

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
	u := NewUHTTP()
	u.ExposeCacheHandlers()

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

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", map[string][]string{"path": {"/cachedHandler1"}}, `{"deletedEntries": 1}`)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler1", nil, `{"counter1": 2}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler2", nil, `{"counter2": 1}`)
}

func TestExposeCacheManagementNotAvailable(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithAutomaticCacheUpdates(200*time.Millisecond, nil, nil),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all1": "ok"}
		}),
	)
	u.Handle("/cache1", handler1)
	require.HTTPError(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.HTTPError(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil)
	require.HTTPError(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", nil)
}

func TestExposeCacheManagement(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	u.ExposeCacheHandlers()

	longResponseLength := 100000
	longResponse := make([]byte, longResponseLength)
	for i := 0; i < longResponseLength; i++ {
		longResponse[i] = 'a'
	}

	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all1": "ok"}
		}),
	)
	handler2 := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"longerResponse": string(longResponse)}
		}),
	)

	u.Handle("/cache1", handler1)
	u.Handle("/cache2", handler2)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 0, "sizeInBytes": 0}, "/cache2": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache1", nil, `{"all1": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache1":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:false bodyBr:false bodyGzip:false bodyDeflate:false }"},"/cache2":{}}`)

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 1, "sizeInBytes": 16}, "/cache2": {"entries": 0, "sizeInBytes": 0},  "total": {"entries": 1, "sizeInBytes": 16}}`)

	// populate second
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache2", nil, fmt.Sprintf(`{"longerResponse": "%s"}`, longResponse))

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 1, "sizeInBytes": 16}, "/cache2": {"entries": 1, "sizeInBytes": 16},  "total": {"entries": 2, "sizeInBytes": 32}}`)

	// clear first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", map[string][]string{"path": {"/cache1"}}, `{"deletedEntries": 1}`)

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 0, "sizeInBytes": 0}, "/cache2": {"entries": 1, "sizeInBytes": 16},  "total": {"entries": 1, "sizeInBytes": 16}}`)

	// clear second
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", map[string][]string{"path": {"/cache2"}}, `{"deletedEntries": 1}`)

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 0, "sizeInBytes": 0}, "/cache2": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache1", nil, `{"all1": "ok"}`)

	// populate second
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache2", nil, fmt.Sprintf(`{"longerResponse": "%s"}`, longResponse))

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 1, "sizeInBytes": 16}, "/cache2": {"entries": 1, "sizeInBytes": 16},  "total": {"entries": 2, "sizeInBytes": 32}}`)

	// clear all
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", nil, `{"deletedEntries": 2}`)

	// check result
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache1": {"entries": 0, "sizeInBytes": 0}, "/cache2": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)
}

func TestCacheEncodings(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	u.ExposeCacheHandlers()

	handler := NewHandler(
		WithCache(10*time.Second),
		WithCachePersistEncodings(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)

	u.Handle("/cache", handler)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache", nil, `{"all": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:true bodyBr:true bodyGzip:true bodyDeflate:true }"}}`)
}

func TestCacheEncodingsNoBrotli(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP(WithBrotliCompression(false, 5))
	u.ExposeCacheHandlers()

	handler := NewHandler(
		WithCache(10*time.Second),
		WithCachePersistEncodings(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)

	u.Handle("/cache", handler)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache", nil, `{"all": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:true bodyBr:false bodyGzip:true bodyDeflate:true }"}}`)
}

func TestCacheEncodingsNoGzip(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP(WithGzipCompression(false, 5))
	u.ExposeCacheHandlers()

	handler := NewHandler(
		WithCache(10*time.Second),
		WithCachePersistEncodings(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)

	u.Handle("/cache", handler)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache", nil, `{"all": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:true bodyBr:true bodyGzip:false bodyDeflate:true }"}}`)
}

func TestCacheEncodingsNoDeflate(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP(WithDeflateCompression(false, 5))
	u.ExposeCacheHandlers()

	handler := NewHandler(
		WithCache(10*time.Second),
		WithCachePersistEncodings(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)

	u.Handle("/cache", handler)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache", nil, `{"all": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:true bodyBr:true bodyGzip:true bodyDeflate:false }"}}`)
}

func TestCacheNoEncodings(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	u.ExposeCacheHandlers()

	handler := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)

	u.Handle("/cache", handler)

	// check initial size
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0}, "total": {"entries": 0, "sizeInBytes": 0}}`)

	// populate first
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cache", nil, `{"all": "ok"}`)

	// check result
	body := assert.HTTPBody(u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil)
	require.Contains(t, body, `{"/cache":{"3333366435656263353433363533346536316431366536336464666361333237":"{ updated:`)
	require.Contains(t, body, `statusCode:200 model:true bodyPlain:false bodyBr:false bodyGzip:false bodyDeflate:false }"}}`)
}

func TestExposeCacheManagementMiddleware(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	requireQueryStringMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.RawQuery, "testMiddleware") {
				w.WriteHeader(http.StatusForbidden)
				_, err := w.Write([]byte(`{"err": "forbidden"}`))
				require.NoError(t, err)
				return
			}
			next.ServeHTTP(w, r)
		}
	}

	u := NewUHTTP()
	u.ExposeCacheHandlers(requireQueryStringMiddleware)

	handler := NewHandler(
		WithCache(10*time.Second),
		WithGet(func(r *http.Request, ret *int) interface{} {
			return map[string]string{"all": "ok"}
		}),
	)
	u.Handle("/cache", handler)

	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size", nil, `{"err": "forbidden"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details", nil, `{"err": "forbidden"}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear", nil, `{"err": "forbidden"}`)

	// working
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/size?testMiddleware", nil, `{"/cache": {"entries": 0, "sizeInBytes": 0},  "total": {"entries": 0, "sizeInBytes": 0}}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodGet, "/uhttp/cache/details?testMiddleware", nil, `{"/cache": {}}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear?testMiddleware", nil, `{"deletedEntries": 0}`)
	RequireHTTPBodyJSONEq(t, u.ServeMux().ServeHTTP, http.MethodPost, "/uhttp/cache/clear?path=cache&testMiddleware", nil, `{"deletedEntries": 0}`)
}

func TestCacheAutomatic(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter1 := 0
	handler1 := NewHandler(
		WithCache(10*time.Second),
		WithAutomaticCacheUpdates(200*time.Millisecond, nil, nil),
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

func TestCacheAutomaticWithParameters(t *testing.T) {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counterParam0 := 0
	counterParam1 := 0
	counterParam2 := 0
	handler := NewHandler(
		WithCache(10*time.Second),
		WithAutomaticCacheUpdates(200*time.Millisecond, nil, []map[string]string{
			{"param1": "param1"},
			{"param2": "param2"},
		}),
		WithGet(func(r *http.Request, ret *int) interface{} {
			if r.URL.Query().Get("param1") == "" && r.URL.Query().Get("param2") == "" {
				counterParam0++
			}
			if r.URL.Query().Get("param1") == "param1" {
				counterParam1++
			}
			if r.URL.Query().Get("param2") == "param2" {
				counterParam2++
			}

			return map[string]interface{}{
				"counterParam0": counterParam0,
				"counterParam1": counterParam1,
				"counterParam2": counterParam2,
			}
		}),
	)
	u.Handle("/cachedHandler", handler)

	// wait for initial update to have run through
	time.Sleep(50 * time.Millisecond)

	// cache should be initialized automatically. Since during the first cacheRun the one with param2 was not run yet -> counter == 0
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", url.Values{"param1": []string{"param1"}},
		`{"counterParam0": 0, "counterParam1": 1, "counterParam2": 0}`,
		map[string][]string{CACHE_HEADER: {"true"}})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", url.Values{"param2": []string{"param2"}},
		`{"counterParam0": 0, "counterParam1": 1, "counterParam2": 1}`,
		map[string][]string{CACHE_HEADER: {"true"}})

	// this request should not have been cached yet
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil,
		`{"counterParam0": 1, "counterParam1": 1, "counterParam2": 1}`,
		map[string][]string{})

	// wait for automatic update
	time.Sleep(200 * time.Millisecond)

	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", url.Values{"param1": []string{"param1"}},
		`{"counterParam0": 1, "counterParam1": 2, "counterParam2": 1}`,
		map[string][]string{CACHE_HEADER: {"true"}})
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", url.Values{"param2": []string{"param2"}},
		`{"counterParam0": 1, "counterParam1": 2, "counterParam2": 2}`,
		map[string][]string{CACHE_HEADER: {"true"}})

	// should still deliver the old response, as the regular cache-time is 10s
	RequireHTTPBodyAndHeader(t, u.ServeMux().ServeHTTP, http.MethodGet, "/cachedHandler", nil,
		`{"counterParam0": 1, "counterParam1": 1, "counterParam2": 1}`,
		map[string][]string{CACHE_HEADER: {"true"}})

}

func setupCacheEncodingTest(t *testing.T) *UHTTP {
	ulog.SetWriter(bufio.NewWriter(nil), nil)
	u := NewUHTTP()
	counter := 0
	handler := NewHandler(
		WithCache(10*time.Second),
		WithCachePersistEncodings(),
		WithGet(func(r *http.Request, ret *int) interface{} {
			counter++
			return map[string]int{"counter": counter}
		}),
	)
	u.Handle("/test", handler)
	return u
}

func TestCachePersistedEncodingPlain(t *testing.T) {
	u := setupCacheEncodingTest(t)
	_, _, _, res := Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": ""})
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
	_, _, _, res = Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": ""})
	body, err = decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
}

func TestCachePersistedEncodingGzip(t *testing.T) {
	u := setupCacheEncodingTest(t)
	_, _, _, res := Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": "gzip"})
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
	_, _, _, res = Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": "gzip"})
	body, err = decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
}

func TestCachePersistedEncodingDeflate(t *testing.T) {
	u := setupCacheEncodingTest(t)
	_, _, _, res := Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": "deflate"})
	body, err := decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "deflate", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
	_, _, _, res = Run(t, u, http.MethodGet, "/test", map[string]string{"Accept-Encoding": "deflate"})
	body, err = decodeResponseBody(res)
	require.NoError(t, err)
	require.Equal(t, "deflate", res.Header.Get("Content-Encoding"))
	require.JSONEq(t, `{"counter": 1}`, string(body))
}
