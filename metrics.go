package uhttp

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	Metric_Requests_Total    string = "uhttp_requests_total"
	Metric_Requests_Duration string = "uhttp_requests_durations"
)

func HandleMetrics(metrics map[string]interface{}, method string, status int, uri string, duration time.Duration) error {
	counter := metrics[Metric_Requests_Total].(*prometheus.CounterVec)
	counterMetric, err := counter.GetMetricWith(prometheus.Labels{
		"method":  method,
		"code":    strconv.Itoa(status),
		"handler": uri,
	})
	if err != nil {
		return err
	}
	counterMetric.Inc()

	dur := metrics[Metric_Requests_Duration].(*prometheus.HistogramVec)
	durationMetric, err := dur.GetMetricWith(prometheus.Labels{
		"method":  method,
		"code":    strconv.Itoa(status),
		"handler": uri,
	})
	if err != nil {
		return err
	}
	durationMetric.Observe(float64(duration / time.Millisecond))
	return nil
}
