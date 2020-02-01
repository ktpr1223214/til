package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	helloCnt = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hello_total", Help: "Hello requested.",
		})

	worldCnt = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "world_total", Help: "World requested.",
		})
)

func hello(w http.ResponseWriter, r *http.Request) {
	helloCnt.Inc()
	w.Write([]byte("Hello"))
}

func world(w http.ResponseWriter, r *http.Request) {
	worldCnt.Inc()
	w.Write([]byte("World"))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/world", world)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8000", nil))
}
