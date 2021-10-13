package uhttp

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dunv/ulog"
)

func cacheMiddleware(u *UHTTP, handler Handler) func(next http.HandlerFunc) http.HandlerFunc {
	var c cache

	// only register cache once (this make the "HandlerFunc" callable more than once)
	u.cacheLock.Lock()
	if registeredCache, ok := u.cache[handler.opts.HandlerPattern]; ok {
		c = *registeredCache
	} else {
		c = cache{
			&sync.RWMutex{},
			handler.opts.CacheMaxAge,
			map[string]cacheEntry{},
		}
		ulog.PanicIfError(u.registerCache(handler.opts.HandlerPattern, &c))
		u.Handle(fmt.Sprintf("/uhttp/cache/clear%s", handler.opts.HandlerPattern), specificCacheClearHandler(u, &c))
		if handler.opts.CacheAutomaticUpdatesInterval > 0 {
			// Run automatic refresher
			go func() {
				f := handler.HandlerFunc(u)
				for {
					r, err := http.NewRequest(http.MethodGet, "forceCache", nil)
					if err != nil {
						ulog.Errorf("this error should never happen (%s)", err)
						continue
					}
					r.Header.Set(handler.opts.CacheBypassHeader, "true")
					f(&noopResponseWriter{}, r)
					time.Sleep(handler.opts.CacheAutomaticUpdatesInterval)
				}
			}()
		}
	}
	u.cacheLock.Unlock()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// only cache GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			bypassCache := r.Header.Get(handler.opts.CacheBypassHeader)
			if bypassCache == "true" {
				next.ServeHTTP(&cachingResponseWriter{
					r:              r,
					w:              w,
					handlerOptions: handler.opts,
					cache:          c,
				}, r)
				return
			}

			h, err := cacheHash(handler.opts, r)
			if err == nil {
				c.RLock()
				if entry, ok := c.data[h]; ok {
					if time.Since(entry.updatedOn) < handler.opts.CacheMaxAge {
						w.Header().Add("X-UHTTP-CACHE", "true")
						w.Header().Add("X-UHTTP-CACHE-AGE-HUMAN-READABLE", time.Since(entry.updatedOn).String())
						w.Header().Add("X-UHTTP-CACHE-AGE-MS", strconv.FormatInt(time.Since(entry.updatedOn).Milliseconds(), 10))
						u.RenderWithStatusCode(w, r, entry.statusCode, entry.data)
						c.RUnlock()
						return
					}
					c.RUnlock()

					c.Lock()
					delete(c.data, h)
					c.Unlock()
				} else {
					c.RUnlock()
				}
			}

			next.ServeHTTP(&cachingResponseWriter{
				r:              r,
				w:              w,
				handlerOptions: handler.opts,
				cache:          c,
			}, r)

		}
	}
}

type cache struct {
	*sync.RWMutex
	maxAge time.Duration
	data   map[string]cacheEntry
}

func (c cache) size() uint64 {
	total := uint64(0)
	c.RLock()
	for _, entry := range c.data {
		total += uint64(len(entry.data))
	}
	c.RUnlock()
	return total
}

type cacheEntry struct {
	updatedOn  time.Time
	data       json.RawMessage
	statusCode int
}

func cacheHash(opts handlerOptions, r *http.Request) (string, error) {
	// Request-Body
	body := ""
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		defer r.Body.Close()
		if r.Body != nil {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		body = string(bodyBytes)
	}

	// Request-Params
	params := r.URL.RawQuery

	// Relevant headers
	headers := []string{}
	for _, h := range opts.CacheRelevantHeaders {
		if val := r.Header.Get(h); val != "" {
			headers = append(headers, fmt.Sprintf("%s: %s", h, val))
		}
	}
	uniqueString := fmt.Sprintf("%s-%s-%s", body, params, strings.Join(headers, "-"))
	s := md5.Sum([]byte(uniqueString))
	return string(s[:]), nil
}

type cachingResponseWriter struct {
	r              *http.Request
	w              http.ResponseWriter
	handlerOptions handlerOptions
	headerWritten  bool
	cache          cache
}

func (w *cachingResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *cachingResponseWriter) Write(data []byte) (int, error) {
	h, err := cacheHash(w.handlerOptions, w.r)
	if err != nil {
		return 0, err
	}

	w.cache.Lock()
	w.cache.data[h] = cacheEntry{
		updatedOn:  time.Now(),
		data:       data,
		statusCode: http.StatusOK,
	}
	w.cache.Unlock()
	return w.w.Write(data)
}

func (w *cachingResponseWriter) WriteHeader(code int) {
	if !w.headerWritten {
		w.headerWritten = true
		w.w.WriteHeader(code)
	}
}

func (w *cachingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

var cacheSizeHandler = func(u *UHTTP) Handler {
	return NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			res := map[string]map[string]uint64{}
			totalSize := uint64(0)
			totalEntries := uint64(0)
			u.cacheLock.RLock()
			for pattern, cache := range u.cache {
				cache.RLock()
				size := cache.size()
				totalSize += size
				entries := uint64(len(cache.data))
				totalEntries += entries
				res[pattern] = map[string]uint64{
					"sizeInBytes": size,
					"entries":     entries,
				}
				cache.RUnlock()
			}
			u.cacheLock.RUnlock()
			res["total"] = map[string]uint64{
				"sizeInBytes":  totalSize,
				"totalEntries": totalEntries,
			}
			return res
		}),
	)
}

var cacheClearHandler = func(u *UHTTP) Handler {
	return NewHandler(
		WithPost(func(r *http.Request, ret *int) interface{} {
			deletedEntries := 0
			u.cacheLock.RLock()
			for _, c := range u.cache {
				c.Lock()
				for key := range c.data {
					delete(c.data, key)
					deletedEntries++
				}
				c.Unlock()
			}
			u.cacheLock.RUnlock()
			return map[string]int{
				"deletedEntries": deletedEntries,
			}
		}),
	)
}

var specificCacheClearHandler = func(u *UHTTP, c *cache) Handler {
	return NewHandler(
		WithPost(func(r *http.Request, ret *int) interface{} {
			deletedEntries := 0
			c.Lock()
			for key := range c.data {
				delete(c.data, key)
				deletedEntries++
			}
			c.Unlock()
			return map[string]int{
				"deletedEntries": deletedEntries,
			}
		}),
	)
}

type noopResponseWriter struct{}

func (w *noopResponseWriter) Header() http.Header { return http.Header{} }

func (w *noopResponseWriter) Write(data []byte) (int, error) { return 0, nil }

func (w *noopResponseWriter) WriteHeader(statusCode int) {}
