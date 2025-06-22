package main

import (
	"fmt"
	"log"
	"net/http"
)

// healthzHandler is a simple handler that responds to health checks.
// Kubernetes will call this endpoint to check if the application is alive and ready.
// If it returns a 200-299 status code, the probe is considered successful.
// If it returns any other status code, the probe is considered a failure.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 200 OK.
	w.WriteHeader(http.StatusOK)
	// Optionally, write a body. "OK" is a common practice.
	w.Write([]byte("OK"))
}

func main() {
	// This is your existing handler for the main application logic.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from GKE Golang Microservice!")
	})

	// THIS IS THE NEW LINE: Register the health check handler for the /healthz path.
	http.HandleFunc("/healthz", healthzHandler)

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}