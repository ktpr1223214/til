package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type WebhookPayload struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []Alert `json:"alerts"`
	// grouping に使った Labels
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
}

// Alert alertmanager webhook メッセージの Alert の詳細を記載した構造体
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	// https://github.com/prometheus/alertmanager/issues/1903
	Fingerprint string `json:"fingerprint"`
}

func webhook(w http.ResponseWriter, r *http.Request) {
	var aml WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&aml); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("incident:")
	log.Println(len(aml.Alerts))
	log.Println(aml)
}

func ticket(w http.ResponseWriter, r *http.Request) {
	var aml WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&aml); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("ticket:")
	log.Println(len(aml.Alerts))
	log.Println(aml)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", webhook)
	mux.HandleFunc("/ticket", ticket)

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}
