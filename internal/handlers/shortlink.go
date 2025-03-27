package handlers

import (
	"encoding/json"
	"net/http"
	"short-link/internal/models"
)

type ShortLinkHandler struct {
	// TODO: Add service layer dependency
}

func NewShortLinkHandler() *ShortLinkHandler {
	return &ShortLinkHandler{}
}

func (h *ShortLinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement service layer call
	response := models.ShortLinkResponse{
		ID:          "1",
		OriginalURL: req.OriginalURL,
		ShortURL:    "http://localhost:8080/abc123",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ShortLinkHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	// Extract code from the path, removing the leading slash
	code := r.URL.Path[len("/api/v1/shortlinks/"):]
	if code == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// TODO: Implement service layer call
	response := models.ShortLinkResponse{
		ID:          "1",
		OriginalURL: "https://example.com",
		ShortURL:    "http://localhost:8080/" + code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
