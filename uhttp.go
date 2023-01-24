package uhttp

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/cache"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type UHTTP struct {
	opts           *uhttpOptions
	requestContext map[ContextKey]interface{}

	metricsServeMux *http.ServeMux
	metrics         map[string]interface{}

	// hold handle to all caches for calculating total and management
	cache     map[string]*cache.Cache
	cacheLock *sync.RWMutex
}

func NewUHTTP(opts ...UhttpOption) *UHTTP {
	mergedOpts := &uhttpOptions{
		cors:                    "*",
		log:                     zap.L(),
		encodingErrorLogLevel:   zapcore.ErrorLevel,
		parseModelErrorLogLevel: zapcore.ErrorLevel,
		handlerErrorLogLevel:    zapcore.ErrorLevel,
		sendPanicInfoToClient:   false,
		serveMux:                http.NewServeMux(),
		address:                 "0.0.0.0:8080",
		readTimeout:             30 * time.Second,
		readHeaderTimeout:       30 * time.Second,
		writeTimeout:            30 * time.Second,
		idleTimeout:             30 * time.Second,
		enableMetrics:           false,
		metricsPath:             "/metrics",
		enableGzip:              true,
		gzipCompressionLevel:    gzip.BestCompression,
		enableBrotli:            true,
		brotliCompressionLevel:  11,
		enableDeflate:           true,
		deflateCompressionLevel: flate.BestCompression,

		silentStaticFileRegistration: false,
		logHandlerCalls:              true,
		logHandlerErrors:             true,
		logHandlerRegistrations:      true,
		logStaticFileAccess:          true,

		cacheTTLEnforcerInterval: 30 * time.Second,
	}
	for _, opt := range opts {
		opt.apply(mergedOpts)
	}

	u := &UHTTP{
		opts:           mergedOpts,
		requestContext: map[ContextKey]interface{}{},
		cache:          map[string]*cache.Cache{},
		cacheLock:      &sync.RWMutex{},
	}

	if mergedOpts.enableMetrics {
		metrics := map[string]interface{}{}
		metrics[Metric_Requests_Total] = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "uhttp",
			Subsystem: "requests",
			Name:      "total",
			Help:      "request counters",
		}, []string{"method", "code", "handler"})

		metrics[Metric_Requests_Duration] = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "uhttp",
			Subsystem: "requests",
			Name:      "duration",
			Help:      "request durations",
			Buckets:   []float64{1, 100, 500, 1000, 5000, 10000, 60000},
		}, []string{"method", "code", "handler"})

		u.metrics = metrics
		u.metricsServeMux = http.NewServeMux()
	}

	return u
}

func (u *UHTTP) registerCache(pattern string, cache *cache.Cache) error {
	if _, ok := u.cache[pattern]; ok {
		return fmt.Errorf("cache for handler %s already exists", pattern)
	}
	u.cache[pattern] = cache
	return nil
}

func (u *UHTTP) ExposeCacheHandlers(middlewares ...Middleware) {
	// this handler is mainly exposed for testing
	u.Handle("/uhttp/cache/size", cacheSizeHandler(u, middlewares...))
	// these handlers are for maintenance
	u.Handle("/uhttp/cache/details", cacheDetailsHandler(u, middlewares...))
	u.Handle("/uhttp/cache/clear", cacheClearHandler(u, middlewares...))
}

func (u *UHTTP) Log() *zap.Logger {
	return u.opts.log
}

func (u *UHTTP) CORS() string {
	return u.opts.cors
}

func (u *UHTTP) ServeMux() *http.ServeMux {
	return u.opts.serveMux
}

func (u *UHTTP) MetricsServeMux() (*http.ServeMux, error) {
	if u.metricsServeMux != nil {
		return u.metricsServeMux, nil
	}
	return nil, errors.New("cannot get metricsServeMux if not initialized")
}

func (u *UHTTP) AddContext(key ContextKey, value interface{}) error {
	keys := ContextKeysFromMap(u.requestContext)
	if !uhelpers.SliceContainsItem(keys, key) {
		u.requestContext[key] = value
		return nil
	}
	return fmt.Errorf("contextKey %s already exists", key)
}

// Handle configuration
func (u *UHTTP) Handle(pattern string, handler Handler) {
	handler.opts.handlerPattern = pattern
	handlerFunc := handler.HandlerFunc(u)

	if u.opts.logHandlerRegistrations {
		if handler.opts.get != nil || handler.opts.getWithModel != nil {
			u.opts.log.Sugar().Infof("Registered http GET %s", pattern)
		} else if handler.opts.post != nil || handler.opts.postWithModel != nil {
			u.opts.log.Sugar().Infof("Registered http POST %s", pattern)
		} else if handler.opts.delete != nil || handler.opts.deleteWithModel != nil {
			u.opts.log.Sugar().Infof("Registered http DELETE %s", pattern)
		}
	}

	u.opts.serveMux.Handle(pattern, handlerFunc)
}

func (u *UHTTP) ListenAndServe() error {
	srv := &http.Server{
		Handler:           u.opts.serveMux,
		Addr:              u.opts.address,
		ReadTimeout:       u.opts.readTimeout,
		ReadHeaderTimeout: u.opts.readHeaderTimeout,
		WriteTimeout:      u.opts.writeTimeout,
		IdleTimeout:       u.opts.idleTimeout,
		ErrorLog:          u.opts.tlsErrorLogger,
	}

	var metricsServer *http.Server
	if u.opts.enableMetrics {
		u.metricsServeMux.Handle(u.opts.metricsPath, promhttp.Handler())
		metricsServer = &http.Server{
			Handler:           u.metricsServeMux,
			Addr:              u.opts.metricsSocket,
			ReadTimeout:       u.opts.readTimeout,
			ReadHeaderTimeout: u.opts.readHeaderTimeout,
			WriteTimeout:      u.opts.writeTimeout,
			IdleTimeout:       u.opts.idleTimeout,
		}
	}

	// Execute TTL for cache (a handler will never serve a cache which is too old, this routine only
	// makes sure that the cache size does not grow too much)
	go func() {
		for {
			u.cacheLock.RLock()
			for _, patternCache := range u.cache {
				keys := patternCache.Keys()
				for _, key := range keys {
					if entry, ok := patternCache.GetByKey(key); ok {
						if time.Since(entry.UpdatedOn()) > patternCache.MaxAge() {
							patternCache.Delete(key)
						}
					}
				}
			}
			u.cacheLock.RUnlock()
			time.Sleep(u.opts.cacheTTLEnforcerInterval)
		}
	}()

	if !u.opts.enableTLS {
		if u.opts.enableMetrics {
			go func() {
				u.Log().Sugar().Infof("Serving metrics at %s", u.opts.metricsSocket)
				u.Log().Sugar().Fatal(metricsServer.ListenAndServe())
			}()
		}

		u.Log().Sugar().Infof("Serving at %s", u.opts.address)
		return srv.ListenAndServe()
	}

	if u.opts.enableMetrics {
		go func() {
			u.Log().Sugar().Infof("ServingTLS metrics at %s", u.opts.metricsSocket)
			u.Log().Sugar().Fatal(metricsServer.ListenAndServeTLS(*u.opts.tlsCertPath, *u.opts.tlsKeyPath))
		}()
	}
	u.Log().Sugar().Infof("ServingTLS at %s", u.opts.address)
	return srv.ListenAndServeTLS(*u.opts.tlsCertPath, *u.opts.tlsKeyPath)
}
