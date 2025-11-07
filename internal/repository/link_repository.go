package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// *DONE LinkRepository est une interface qui définit les méthodes d'accès aux données
// pour les opérations CRUD sur les liens.
// L'implémenter avec les méthodes nécessaires

type GormLinkRepository struct { // * Done
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *GormLinkRepository { // * Done
	return &GormLinkRepository{db: db}
}

func (r *GormLinkRepository) CreateLink(link *models.Link) error { // * Done
	return r.db.Create(link).Error
}

func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) { // * Done
	var link models.Link
	if err := r.db.Where("shortcode = ?", shortCode).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) { // * Done
	var links []models.Link
	if err := r.db.Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) { // * Done
	var count int64
	if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

type LinkRepository interface {
    CreateLink(link *models.Link) error
    GetLinkByShortCode(shortCode string) (*models.Link, error)
    GetLinkByID(id uint) (*models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}