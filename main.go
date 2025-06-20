package main

import (
	"fmt"
	"log"
	"net/http"
	"short-link/internal/config"
	"short-link/internal/short_link"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	config := config.GetConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database failed to start: %v", err)

		return
	}

	// Auto-migrate the ShortLink table
	log.Println("Running database migrations...")
	err = db.AutoMigrate(&short_link.ShortLink{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully!")

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shortLinkHandler.GetByCode(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	log.Println("Server starting on :8080")

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ServerPort),
		Handler: mux,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
