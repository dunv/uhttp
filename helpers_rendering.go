package uhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dunv/uhttp/cache"
)

// HelperMethod for rendering a JSON model
func (u *UHTTP) Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	u.rawRenderWithStatusCode(w, r, http.StatusOK, model)
}

// HelperMethod for rendering a JSON model with statusCode in the response
func (u *UHTTP) RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	u.rawRenderWithStatusCode(w, r, statusCode, model)
}

// HelperMethod for rendering an error as JSON while automatically setting a 400 statusCode
func (u *UHTTP) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	u.RenderErrorWithStatusCode(w, r, http.StatusBadRequest, err, true)
}

// HelperMethod for rendering an error as JSON with defining a custom statusCode
func (u *UHTTP) RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error, logOut bool) {
	if err != nil {
		u.rawRenderWithStatusCode(w, r, statusCode, NewHttpErrorResponse(err))
		if logOut {
			u.opts.log.LogWithLevelf(u.opts.handlerErrorLogLevel, "[uri: %s] err: %s", r.RequestURI, err.Error())
		}
	} else {
		u.opts.log.Panic("Error to be rendered is nil")
	}
}

// Internal helperMethod with is used for ALL rendering throughout uhttp
// Takes care of encoding responses
func (u *UHTTP) rawRenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	encoding := u.determineEncoding(r, statusCode)

	// Write header
	w.Header().Set(HEADER_CONTENT_ENCODING, encoding)
	w.WriteHeader(statusCode)

	// Prepare body-writer
	ew := u.encodingWriter(w, encoding)
	defer ew.Close()

	// Write body
	err := json.NewEncoder(ew).Encode(model)
	if err != nil {
		u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err encoding http response (%s)", err)
		return
	}

	// If we are in a cache-run: give the cache all info
	if crw, ok := w.(*cachingResponseWriter); ok {
		crw.Close(model, statusCode)
		return
	}
}

func (u *UHTTP) renderCacheEntry(handler Handler, w http.ResponseWriter, r *http.Request, entry cache.CacheEntry) {
	_ = AddLogOutput(w, "cached", "true")
	w.Header().Add(CACHE_HEADER, "true")
	w.Header().Add(CACHE_HEADER_AGE_HUMAN_READABLE, time.Since(entry.UpdatedOn()).String())
	w.Header().Add(CACHE_HEADER_AGE_MS, strconv.FormatInt(time.Since(entry.UpdatedOn()).Milliseconds(), 10))

	if handler.opts.cachePersistEncodings {
		encoding := u.determineEncoding(r, entry.ResponseStatusCode())

		// Write
		w.Header().Set(HEADER_CONTENT_ENCODING, encoding)
		w.WriteHeader(entry.ResponseStatusCode())

		var body []byte

		switch encoding {
		case ENCODING_PLAIN:
			body = entry.ResponseBodyPlain()
		case ENCODING_BROTLI:
			body = entry.ResponseBodyBrotli()
		case ENCODING_GZIP:
			body = entry.ResponseBodyGzip()
		case ENCODING_DEFLATE:
			body = entry.ResponseBodyDeflate()
		}
		_, _ = w.Write(body)
		return
	}

	u.rawRenderWithStatusCode(w, r, entry.ResponseStatusCode(), entry.ResponseModel())
}
