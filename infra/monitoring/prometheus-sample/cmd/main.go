package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ktpr1223214/til/infra/monitoring/prometheus-sample/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Prometheus interface {
	ObserveHTTPRequest(code string, method string, path string, start time.Time)
}

type Server struct {
	router *mux.Router

	prometheus Prometheus
}

func New(p Prometheus) *Server {
	srv := &Server{
		prometheus: p,
		router:     mux.NewRouter(),
	}
	srv.routes()

	return srv
}

func (s *Server) Run(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.router); err != nil {
		return err
	}
	return nil
}

// ServeHTTP これを定義することで、Server 自体が http.Handler interface を満たすことができる
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.Handle("/hello", s.observe(s.hello())).Methods(http.MethodGet)
	s.router.Handle("/world", s.observe(s.world())).Methods(http.MethodGet)
	s.router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
}

// handler
func (s *Server) hello() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
}

func (s *Server) world() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code
// cf. https://github.com/gorilla/handlers/blob/7e0847f9db758cdebd26c149d0ae9d5d0b9c98ce/handlers.go#L46
type responseLogger struct {
	w      http.ResponseWriter
	status int
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	return size, err
}

// middleware
func (s *Server) observe(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rl := &responseLogger{w: w, status: http.StatusOK}

		h.ServeHTTP(rl, r)

		status := strconv.Itoa(rl.status)
		s.prometheus.ObserveHTTPRequest(status, r.Method, r.URL.Path, start)
	})
}

func main() {
	p := metrics.New()
	s := New(p)

	log.Fatal(s.Run(8000))
}
