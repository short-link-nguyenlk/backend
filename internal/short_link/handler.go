package short_link

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(ShortLink{
		OriginalURL: req.OriginalURL,
		ShortCode:   "abc123", // This should be generated
	})

	if err != nil {
		http.Error(w, "Create failed", http.StatusInternalServerError)

		return
	}

	// TODO: Implement service layer call
	response := ShortLinkResponse{
		ID:          id,
		OriginalURL: req.OriginalURL,
		ShortURL:    "http://localhost:8080/abc123",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	// Extract code from the path, removing the leading slash
	code := r.URL.Path[len("/api/v1/shortlinks/"):]
	if code == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)

		return
	}

	shortLink, err := h.repo.FindByCode(code)
	if err != nil {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	// TODO: Implement service layer call
	response := ShortLinkResponse{
		ID:          shortLink.ID,
		OriginalURL: shortLink.OriginalURL,
		ShortURL:    "http://localhost:8080/" + shortLink.ShortCode,
		CreatedAt:   shortLink.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
