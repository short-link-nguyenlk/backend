package short_link

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"short-link/internal/config"
	"time"

	"github.com/spaolacci/murmur3"
	"gorm.io/gorm"
)

var base62Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Charset[0])
	}

	result := ""
	for num > 0 {
		remainder := num % 62
		num /= 62
		result = string(base62Charset[remainder]) + result
	}

	return result
}

func (h *Handler) GenerateHash(URL string, seed int64) string {
	hash, _ := murmur3.Sum128WithSeed([]byte(URL), uint32(seed))
	return base62Encode(hash)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	maxRetries := 3

	var shortCode string
	for i := range maxRetries {
		if i == 0 {
			shortCode = h.GenerateHash(req.OriginalURL, 0)
		} else {
			shortCode = h.GenerateHash(req.OriginalURL, time.Now().UnixNano())
		}
		existingShortCode, err := h.repo.FindByCode(shortCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}

		if existingShortCode != nil && existingShortCode.OriginalURL == req.OriginalURL {
			response := ShortLinkResponse{
				ID:          existingShortCode.ID,
				OriginalURL: existingShortCode.OriginalURL,
				ShortURL:    fmt.Sprintf("%s:%s/%s", config.GetConfig().BaseURL, config.GetConfig().ServerPort, existingShortCode.ShortCode),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)

			return
		}

		if i == maxRetries-1 {
			http.Error(w, "Short code collision", http.StatusInternalServerError)

			return
		}
	}

	id, err := h.repo.Create(ShortLink{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode, // This should be generated
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// TODO: Implement service layer call
	response := ShortLinkResponse{
		ID:          id,
		OriginalURL: req.OriginalURL,
		ShortURL:    fmt.Sprintf("%s:%s/%s", config.GetConfig().BaseURL, config.GetConfig().ServerPort, shortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetByCode(w http.ResponseWriter, r *http.Request) {
	// Extract code from the path, removing the leading slash
	code := r.URL.Path[len("/"):]
	if code == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)

		return
	}

	shortLink, err := h.repo.FindByCode(code)
	if err != nil {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortLink.OriginalURL, http.StatusFound)
}
