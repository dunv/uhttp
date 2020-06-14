package uhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: add filters for logging (i.e. do not log everything, or only user etc)
// TODO: make statistics trackable -> prometheus?

func init() {
	// Make expected output (which is only for info, not for debugging) more readable
	ulog.AddSkipFunctions(
		"github.com/dunv/uhttp.RenderError",
		"github.com/dunv/uhttp/helpers.RenderError",
		"github.com/dunv/uhttp.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp.RenderMessage",
		"github.com/dunv/uhttp/helpers.RenderMessage",
		"github.com/dunv/uhttp.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp/helpers.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderMessageWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.rawRenderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.rawRenderErrorWithStatusCode",
	)
	ulog.AddReplaceFunction("github.com/dunv/uhttp.addLoggingMiddleware.func1.1", "uhttp.Log")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.(*UHTTP).ListenAndServe", "uhttp.ListenAndServe")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.(*UHTTP).ListenAndServe.func1", "uhttp.ListenAndServe")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.(*UHTTP).Handle", "uhttp.Handle")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.(*UHTTP).RegisterStaticFilesHandler", "uhttp.HandleStatic")
}

type UHTTP struct {
	opts           *uhttpOptions
	requestContext map[ContextKey]interface{}
	metrics        map[string]interface{}
}

func NewUHTTP(opts ...UhttpOption) *UHTTP {
	mergedOpts := &uhttpOptions{
		cors:                    "*",
		log:                     ulog.NewUlog(),
		gzipCompressionLevel:    4,
		encodingErrorLogLevel:   ulog.LEVEL_ERROR,
		parseModelErrorLogLevel: ulog.LEVEL_ERROR,
		serveMux:                http.NewServeMux(),
		address:                 "0.0.0.0:8080",
		readTimeout:             30 * time.Second,
		readHeaderTimeout:       30 * time.Second,
		writeTimeout:            30 * time.Second,
		idleTimeout:             30 * time.Second,
		enableMetrics:           false,
		metricsPath:             "/metrics",
	}
	for _, opt := range opts {
		opt.apply(mergedOpts)
	}

	metrics := map[string]interface{}{}
	if mergedOpts.enableMetrics {
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
	}

	u := &UHTTP{
		opts:           mergedOpts,
		requestContext: map[ContextKey]interface{}{},
		metrics:        metrics,
	}

	return u
}

func (u *UHTTP) ServeMux() *http.ServeMux {
	return u.opts.serveMux
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
	handlerFunc := handler.HandlerFunc(u)

	if handler.opts.Get != nil {
		u.opts.log.Infof("Registered http GET %s", pattern)
	} else if handler.opts.Post != nil {
		u.opts.log.Infof("Registered http POST %s", pattern)
	} else if handler.opts.Delete != nil {
		u.opts.log.Infof("Registered http DELETE %s", pattern)
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
		mux := http.NewServeMux()
		mux.Handle(u.opts.metricsPath, promhttp.Handler())
		metricsServer = &http.Server{
			Handler:           mux,
			Addr:              u.opts.metricsSocket,
			ReadTimeout:       u.opts.readTimeout,
			ReadHeaderTimeout: u.opts.readHeaderTimeout,
			WriteTimeout:      u.opts.writeTimeout,
			IdleTimeout:       u.opts.idleTimeout,
		}
	}

	if !u.opts.enableTLS {
		if u.opts.enableMetrics {
			go func() {
				ulog.Infof("Serving metrics at %s", u.opts.metricsSocket)
				ulog.Fatal(metricsServer.ListenAndServe())
			}()
		}

		ulog.Infof("Serving at %s", u.opts.address)
		return srv.ListenAndServe()
	}

	if u.opts.enableMetrics {
		go func() {
			ulog.Infof("ServingTLS metrics at %s", u.opts.metricsSocket)
			ulog.Fatal(metricsServer.ListenAndServeTLS(*u.opts.tlsCertPath, *u.opts.tlsKeyPath))
		}()
	}
	ulog.Infof("ServingTLS at %s", u.opts.address)
	return srv.ListenAndServeTLS(*u.opts.tlsCertPath, *u.opts.tlsKeyPath)
}
