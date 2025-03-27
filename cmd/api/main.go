package main

import (
	"log"
	"net/http"
	"short-link/internal/handlers"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Initialize handlers
	shortLinkHandler := handlers.NewShortLinkHandler()

	// Define routes
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("GET /health", handleHealth)

	// Short link routes
	mux.HandleFunc("/api/v1/shortlinks", shortLinkHandler.Create)
	mux.HandleFunc("/api/v1/shortlinks/", shortLinkHandler.GetByCode)

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to Short Link API"}`))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
