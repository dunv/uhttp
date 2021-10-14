package uhttp

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dunv/uhttp/cache"
	"github.com/itchio/go-brotli/enc"
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
			u.opts.log.Errorf("renderError [path: %s] %s", r.RequestURI, err.Error())
		}
	} else {
		u.opts.log.Panic("Error to be rendered is nil")
	}
}

// Internal helperMethod with is used for ALL rendering throughout uhttp
// Takes care of encoding responses
func (u *UHTTP) rawRenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	// TODO: find a way of doing the logging per handler!
	var writer io.Writer
	var err error

	// The go-http-client implementation decodes gzip out-of-the-box, but only if it gets 200 OK
	// For now: use the same behavior here
	if statusCode == http.StatusOK {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if u.opts.enableBrotli && strings.Contains(acceptEncoding, "br") {
			w.Header().Add("Content-Encoding", "br")
			writer = enc.NewBrotliWriter(w, &enc.BrotliWriterOptions{Quality: u.opts.brotliCompressionLevel})
		} else if u.opts.enableGzip && strings.Contains(acceptEncoding, "gzip") {
			w.Header().Add("Content-Encoding", "gzip")
			writer, err = gzip.NewWriterLevel(w, u.opts.gzipCompressionLevel)
			if err != nil {
				u.opts.log.Panic(fmt.Errorf("could not initialize gzip writer (%s)", err))
			}
		} else if u.opts.enableDeflate && strings.Contains(acceptEncoding, "deflate") {
			w.Header().Add("Content-Encoding", "deflate")
			writer, err = flate.NewWriter(w, u.opts.deflateCompressionLevel)
			if err != nil {
				u.opts.log.Panic(fmt.Errorf("could not initialize deflate writer (%s)", err))
			}
		} else {
			writer = w
		}
	} else {
		writer = w
	}

	w.WriteHeader(statusCode)

	err = json.NewEncoder(writer).Encode(model)
	if err != nil {
		u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err encoding http response (%s)", err)
		return
	}

	switch typedWriter := writer.(type) {
	case *enc.BrotliWriter:
		err = typedWriter.Close()
		if err != nil {
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing brotli writer (%s)", err)
		}
	case *gzip.Writer:
		err = typedWriter.Close()
		if err != nil {
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing gzip writer (%s)", err)
		}
	case *flate.Writer:
		err = typedWriter.Close()
		if err != nil {
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing deflate writer (%s)", err)
		}
	}

	switch responseWriter := w.(type) {
	case *cachingResponseWriter:
		err = responseWriter.Close(model, statusCode)
		if err != nil {
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing cachingResponseWriter (%s)", err)
		}
	}
}

func (u *UHTTP) renderCacheEntry(w http.ResponseWriter, r *http.Request, entry cache.CacheEntry) {
	_ = AddLogOutput(w, "cached", "true")
	w.Header().Add(CACHE_HEADER, "true")
	w.Header().Add(CACHE_HEADER_AGE_HUMAN_READABLE, time.Since(entry.UpdatedOn()).String())
	w.Header().Add(CACHE_HEADER_AGE_MS, strconv.FormatInt(time.Since(entry.UpdatedOn()).Milliseconds(), 10))

	// if u.opts.cachePersistDifferentEncodings {
	// for k, v := range entry.ResponseHeader() {
	// w.Header().Set(k, strings.Join(v, ", "))
	// }
	// w.WriteHeader(entry.ResponseStatusCode())
	//
	// var err error
	// acceptEncoding := r.Header.Get("Accept-Encoding")
	// if u.opts.enableBrotli && strings.Contains(acceptEncoding, "br") {
	// _, err = w.Write(e.responseBodyBrotli)
	// // w.Header().Set("Content-Encoding", "br")
	// } else if u.opts.enableGzip && strings.Contains(acceptEncoding, "gzip") {
	// _, err = w.Write(e.responseBodyGzip)
	// // w.Header().Set("Content-Encoding", "gzip")
	// } else if u.opts.enableDeflate && strings.Contains(acceptEncoding, "deflate") {
	// _, err = w.Write(e.responseBodyDeflate)
	// // w.Header().Set("Content-Encoding", "deflate")
	// } else {
	// // writer = w
	// _, err = w.Write(e.responseBody)
	// }
	//
	// if err != nil {
	// u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err rendering cacheEntry (%s)", err)
	// }
	// return
	// }

	u.rawRenderWithStatusCode(w, r, entry.ResponseStatusCode(), entry.ResponseModel())
}
