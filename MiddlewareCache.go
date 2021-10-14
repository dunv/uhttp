package uhttp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dunv/uhttp/cache"
	"github.com/dunv/ulog"
)

const CACHE_HEADER = "X-UHTTP-CACHE"
const CACHE_HEADER_AGE_HUMAN_READABLE = "X-UHTTP-CACHE-AGE-HUMAN-READABLE"
const CACHE_HEADER_AGE_MS = "X-UHTTP-CACHE-AGE-MS"

// This middleware provides a per-handler cache
// It will cache the original response to the client based on
// - "relevant" headers
// - queryParams
// - requestBody
func cacheMiddleware(u *UHTTP, handler Handler) func(next http.HandlerFunc) http.HandlerFunc {
	var c *cache.Cache

	// only register cache once (this make the "HandlerFunc" callable more than once)
	u.cacheLock.Lock()
	if registeredCache, ok := u.cache[handler.opts.HandlerPattern]; ok {
		c = registeredCache
	} else {
		c = cache.NewCache(handler.opts.CacheMaxAge, u.opts.cachePersistDifferentEncodings)
		ulog.PanicIfError(u.registerCache(handler.opts.HandlerPattern, c))

		if u.opts.cacheExposeHandlers {
			u.Handle(fmt.Sprintf("/uhttp/cache/clear%s", handler.opts.HandlerPattern), specificCacheClearHandler(u, c))
		}

		if handler.opts.CacheAutomaticUpdatesInterval > 0 {
			// Run automatic refresher
			go func() {
				f := handler.handlerFuncExcludeMiddlewareByName(u, handler.opts.CacheAutomaticUpdatesSkipMiddleware)
				for {
					r, err := http.NewRequest(http.MethodGet, NO_LOG_MAGIC_URL_FORCE_CACHE, nil)
					if err != nil {
						ulog.Errorf("this error should never happen (%s)", err)
						time.Sleep(handler.opts.CacheAutomaticUpdatesInterval)
						continue
					}
					r.Header.Set(handler.opts.CacheBypassHeader, "true")
					noopWriter := &noopResponseWriter{}
					f(noopWriter, r)

					if noopWriter.statusCode != http.StatusOK {
						u.opts.log.Errorf("could not populate cache of %s. statusCode:%d body:%s", handler.opts.HandlerPattern, noopWriter.statusCode, strings.TrimSpace(noopWriter.body))
					}
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
					r:     r,
					w:     w,
					cache: c,
				}, r)
				return
			}

			if entry, ok, key := c.Get(ExtractAndRestoreRequestBody(r), r.URL.RawQuery); ok {
				if time.Since(entry.UpdatedOn()) < handler.opts.CacheMaxAge {
					u.renderCacheEntry(w, r, entry)
					return
				}
				c.Delete(key)
			}

			next.ServeHTTP(&cachingResponseWriter{
				r:     r,
				w:     w,
				cache: c,
			}, r)
		}
	}
}

// a response writer whch updates the cache as soon as a response is sent to the client
type cachingResponseWriter struct {
	r            *http.Request
	w            http.ResponseWriter
	cache        *cache.Cache
	wroteHeader  bool
	responseBody []byte
}

// a response writer whch updates the cache as soon as a response is sent to the client
func (w *cachingResponseWriter) Header() http.Header {

	return w.w.Header()
}

// a response writer whch updates the cache as soon as a response is sent to the client
// Write always needs to be called AFTER WriteHeader, otherwise we get an error
func (w *cachingResponseWriter) Write(data []byte) (int, error) {
	w.responseBody = append(w.responseBody, data...)
	return w.w.Write(data)
}

// a response writer whch updates the cache as soon as a response is sent to the client
func (w *cachingResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.w.WriteHeader(code)
	}
}

func (w *cachingResponseWriter) Close(model interface{}, statusCode int) error {
	w.cache.Set(
		ExtractAndRestoreRequestBody(w.r), w.r.URL.RawQuery, w.r.Header.Clone(),
		model, w.responseBody, w.w.Header().Clone(), statusCode,
	)
	return nil
}

// a response writer whch updates the cache as soon as a response is sent to the client
func (w *cachingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// a response writer which does nothing (used for automatically updating the cache in the background)
// it can simulate an actual call which discards the anwer to the client
type noopResponseWriter struct {
	body       string
	statusCode int
}

// a response writer which does nothing (used for automatically updating the cache in the background)
func (w *noopResponseWriter) Header() http.Header { return http.Header{} }

// a response writer which does nothing (used for automatically updating the cache in the background)
func (w *noopResponseWriter) Write(data []byte) (int, error) {
	w.body = string(data)
	return 0, nil
}

// a response writer which does nothing (used for automatically updating the cache in the background)
func (w *noopResponseWriter) WriteHeader(statusCode int) { w.statusCode = statusCode }
