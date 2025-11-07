package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// LinkRepository est une interface qui définit les méthodes d'accès aux données
// pour les opérations CRUD sur les liens.
// L'implémenter avec les méthodes nécessaires

type GormLinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	return r.db.Create(link).Error
}

func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	err := r.db.Where("short_code = ?", shortCode).First(&link).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return the GORM error so callers can detect "not found" reliably
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &link, nil
}

func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	if err := r.db.Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GormLinkRepository) GetLinkByID(id uint) (*models.Link, error) {
	var link models.Link
	if err := r.db.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}


type LinkRepository interface {
	GetAllLinks() ([]models.Link, error)
    CreateLink(link *models.Link) error
    GetLinkByShortCode(shortCode string) (*models.Link, error)
    GetLinkByID(id uint) (*models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}