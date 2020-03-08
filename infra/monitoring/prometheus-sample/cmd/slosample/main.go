package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	reqCnt = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The number of HTTP requests processed, labeled with status code, HTTP method and URL path.",
		},
		[]string{"code", "method", "path"},
	)

	reqDur = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "The latency of HTTP requests processed, labeled with status code, HTTP method and URL path.",
		},
		[]string{"code", "method", "path"},
	)
)

func respondPrettyJSON(w http.ResponseWriter, body interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// format
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		return errors.Wrap(err, "failed to encode response")
	}
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func err(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "err", http.StatusInternalServerError)
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
}

// WriteHeader 通常の WriteHeader メソッドを呼び出し、status code を保存
func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// Header 通常の Header メソッドを呼び出すだけ
func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

// Write 通常の Write メソッドを呼び出すだけ
func (l *responseLogger) Write(b []byte) (int, error) {
	return l.w.Write(b)
}

// instrumentRequest Prometheus 開示用に HTTP Request のメトリクスを計測する middleware
func instrumentRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// metrics パスのメトリクスは、Prometheus のクライアントが別途開示しているのでスキップ
		if r.URL.Path == "/metrics" {
			h.ServeHTTP(w, r)
			return
		}
		s := time.Now()

		rl := &responseLogger{w: w, status: http.StatusOK}

		h.ServeHTTP(rl, r)

		status := strconv.Itoa(rl.status)
		reqCnt.WithLabelValues(status, r.Method, r.URL.Path).Inc()
		// nanosecond を 1e9 で割って second
		reqDur.WithLabelValues(status, r.Method, r.URL.Path).Observe(float64(time.Since(s)) / float64(time.Second))
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	// cf. https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html
	type HealthCheck struct {
		Status *string `json:"status"`
		Output *string `json:"output,omitempty"`
	}
	hc := HealthCheck{
		Status: stringPtr("pass"),
	}
	if err := respondPrettyJSON(w, hc, http.StatusOK); err != nil {
		http.Error(w, "レスポンス書き込みに失敗しました", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/health", instrumentRequest(http.HandlerFunc(health)))
	mux.Handle("/error", instrumentRequest(http.HandlerFunc(err)))
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
