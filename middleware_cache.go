package uhttp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"path"
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
	if registeredCache, ok := u.cache[handler.opts.handlerPattern]; ok {
		c = registeredCache
	} else {
		// enable wanted encodings in cache
		c = cache.NewCache(handler.opts.cacheMaxAge, handler.opts.handlerPattern)

		ulog.PanicIfError(u.registerCache(handler.opts.handlerPattern, c))

		if handler.opts.cacheAutomaticUpdatesInterval > 0 {
			// Run automatic refresher
			go func() {
				f := handler.handlerFuncExcludeMiddlewareByName(u, handler.opts.cacheAutomaticUpdatesSkipMiddleware)
				for {
					// parameters are populated with a single empty set by default
					for _, paramSet := range handler.opts.cacheAutomaticUpdatesParameters {
						r, err := http.NewRequest(http.MethodGet, NO_LOG_MAGIC_URL_FORCE_CACHE, nil)
						if err != nil {
							ulog.Errorf("this error should never happen (%s)", err)
							continue
						}
						q := r.URL.Query()
						for paramKey, paramValue := range paramSet {
							q.Add(paramKey, paramValue)
						}
						r.URL.RawQuery = q.Encode()
						r.Header.Set(handler.opts.cacheBypassHeader, "true")

						noopWriter := &noopResponseWriter{}
						f(noopWriter, r.WithContext(context.WithValue(r.Context(), CtxKeyIsAutomaticCacheExecution, true)))
						if noopWriter.statusCode != http.StatusOK {
							u.opts.log.Errorf("could not populate cache of %s. statusCode:%d body:%s", handler.opts.handlerPattern, noopWriter.statusCode, strings.TrimSpace(noopWriter.body))
						}
					}

					time.Sleep(handler.opts.cacheAutomaticUpdatesInterval)
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

			bypassCache := r.Header.Get(handler.opts.cacheBypassHeader)
			if bypassCache == "true" {
				next.ServeHTTP(newCachingResponseWriter(u, handler, w, r, c), r)
				return
			}

			if entry, ok, key := c.Get(ExtractAndRestoreRequestBody(r), r.URL.RawQuery); ok {
				if time.Since(entry.UpdatedOn()) < handler.opts.cacheMaxAge {
					u.renderCacheEntry(handler, w, r, entry)
					return
				}
				c.Delete(key)
			}

			next.ServeHTTP(newCachingResponseWriter(u, handler, w, r, c), r)
		}
	}
}

// a response writer whch updates the cache as soon as a response is sent to the client
type cachingResponseWriter struct {
	u            *UHTTP
	h            Handler
	r            *http.Request
	w            http.ResponseWriter
	cache        *cache.Cache
	wroteHeader  bool
	responseBody []byte
	startTime    time.Time
}

func newCachingResponseWriter(u *UHTTP, h Handler, w http.ResponseWriter, r *http.Request, cache *cache.Cache) *cachingResponseWriter {
	if u.opts.logCacheRuns {
		if strings.Contains(r.URL.String(), NO_LOG_MAGIC_URL_FORCE_CACHE) {
			u.Log().Infof("Started automatic caching of %s", h.opts.handlerPattern)
		} else {
			u.Log().Infof("Started caching of %s by userRequest", h.opts.handlerPattern)
		}
	}

	return &cachingResponseWriter{
		u:         u,
		h:         h,
		w:         w,
		r:         r,
		cache:     cache,
		startTime: time.Now(),
	}
}

// a response writer whch updates the cache as soon as a response is sent to the client
func (w *cachingResponseWriter) Header() http.Header {
	return w.w.Header()
}

// a response writer whch updates the cache as soon as a response is sent to the client

func (w *cachingResponseWriter) Write(data []byte) (int, error) {
	w.responseBody = append(w.responseBody, data...)

	// the default implementation in net/http/server.go (line 1577 in go 1.17.2) writes the response-header as
	// soon as write is called, if there are no headers written yet
	if !w.wroteHeader {
		w.w.WriteHeader(http.StatusOK)
	}
	return w.w.Write(data)
}

// a response writer whch updates the cache as soon as a response is sent to the client
func (w *cachingResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		// copied straight out of the standard-library net/http/server.go
		caller := relevantCaller()
		w.u.opts.log.Warnf("superfluous response.WriteHeader call from %s (%s:%d). could happen if the responseWriter is used in a uhttp.Handler AND the function returns something non-nil", caller.Function, path.Base(caller.File), caller.Line)
		return
	}

	w.wroteHeader = true
	w.w.WriteHeader(code)
}

func (w *cachingResponseWriter) Close(model interface{}, statusCode int) {
	var err error
	var bodyPlain []byte
	var bodyBrotli []byte
	var bodyGzip []byte
	var bodyDeflate []byte

	if !w.h.opts.cacheFailedRequests && statusCode != http.StatusOK {
		return
	}

	if w.h.opts.cachePersistEncodings {
		bodyPlain, err = json.Marshal(model)
		if err != nil {
			w.u.Log().Errorf("could not encode model for caching (%s)", err)
		} else {
			if w.u.opts.enableBrotli {
				bodyBrotli, err = w.u.compressJSON(ENCODING_BROTLI, bodyPlain)
				if err != nil {
					w.u.Log().Errorf("could not compress JSON for caching (%s)", err)
				}
			}
			if w.u.opts.enableGzip {
				bodyGzip, err = w.u.compressJSON(ENCODING_GZIP, bodyPlain)
				if err != nil {
					w.u.Log().Errorf("could not compress JSON for caching (%s)", err)
				}
			}
			if w.u.opts.enableDeflate {
				bodyDeflate, err = w.u.compressJSON(ENCODING_DEFLATE, bodyPlain)
				if err != nil {
					w.u.Log().Errorf("could not compress JSON for caching (%s)", err)
				}
			}
		}
	}

	w.cache.Set(
		ExtractAndRestoreRequestBody(w.r), w.r.URL.RawQuery, w.r.Header.Clone(),
		model, w.w.Header().Clone(), statusCode,
		bodyPlain, bodyBrotli, bodyGzip, bodyDeflate,
	)

	if w.u.opts.logCacheRuns {
		if w.r.URL.String() == NO_LOG_MAGIC_URL_FORCE_CACHE {
			w.u.Log().Infof("Finished automatic caching of %s in %s", w.h.opts.handlerPattern, time.Since(w.startTime).String())
		} else {
			w.u.Log().Infof("Finished caching by userRequest of %s in %s", w.h.opts.handlerPattern, time.Since(w.startTime).String())
		}
	}
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
