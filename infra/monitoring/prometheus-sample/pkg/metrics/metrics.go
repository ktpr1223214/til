package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	reqCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)

	reqErrCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "How many HTTP requests failed, partitioned by status code and HTTP method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)

	reqDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds",
			Help: "How fast requests processed, partitioned by status code and HTTP method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)

	// TODO: batch 系のバリエーションもサンプルとして考えてみる
)

type Prometheus struct{}

func New() *Prometheus {
	return &Prometheus{}
}

func (p *Prometheus) ObserveHTTPRequest(code string, method string, path string, start time.Time) {
	reqCnt.WithLabelValues(code, method, path).Inc()
	reqDuration.WithLabelValues(code, method, path).Observe(time.Since(start).Seconds())
}

func (p *Prometheus) ObserveHTTPError(code string, method string, path string) {
	reqErrCnt.WithLabelValues(code, method, path).Inc()
}
