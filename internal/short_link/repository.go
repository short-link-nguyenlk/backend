package short_link

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(shortLink ShortLink) (uint, error) {
	if err := r.db.Create(&shortLink).Error; err != nil {
		return shortLink.ID, err
	}
	return shortLink.ID, nil
}

func (r *Repository) FindByCode(code string) (*ShortLink, error) {
	var shortLink ShortLink
	if err := r.db.Where("short_code = ?", code).First(&shortLink).Error; err != nil {
		return nil, err
	}
	return &shortLink, nil
}
