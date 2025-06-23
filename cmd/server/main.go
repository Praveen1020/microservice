package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metric to count total HTTP requests, labeled by path
var httpRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "golang_microservice_http_requests_total",
		Help: "Total number of HTTP requests handled by the service.",
	},
	[]string{"path"},
)

// Middleware to collect Prometheus metrics for each request
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}

// Handlers
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("live"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Application handler wrapped with metrics collection
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from GKE Golang Microservice! (v2 with Probes & Metrics)")
	})
	http.Handle("/", metricsMiddleware(mainHandler))

	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Health and probe endpoints
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/healthz/live", livenessHandler)
	http.HandleFunc("/healthz/ready", readinessHandler)

	// Read port from environment variable (fallback to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}