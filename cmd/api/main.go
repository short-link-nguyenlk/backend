package main

import (
	"log"
	"net/http"
	"short-link/internal/short_link"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Connecting to database

	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database failed to start: %v", err)

		return
	}
	log.Printf("Database connected successfully %v", db.Name())

	shortLinkRepo := short_link.NewRepository(db)
	// Initialize handlers
	shortLinkHandler := short_link.NewHandler(*shortLinkRepo)

	// Định nghĩa handler cho POST /api/v1/shortlinks
	mux.HandleFunc("/api/v1/shortlinks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			shortLinkHandler.Create(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Định nghĩa handler cho GET /api/v1/shortlinks/{code}
	mux.HandleFunc("/api/v1/shortlinks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shortLinkHandler.GetByCode(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	log.Println("Server starting on :8080")

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
