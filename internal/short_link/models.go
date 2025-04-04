package short_link

import (
	"time"

	"gorm.io/gorm"
)

type ShortLink struct {
	gorm.Model
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

type ShortLinkResponse struct {
	ID          uint      `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}
