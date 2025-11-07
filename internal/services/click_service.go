package services

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// ClickService est une structure qui fournit des méthodes pour la logique métier des clics.
// Elle est composée d'un repository permettant la persistance des clics.
type ClickService struct {
	clickRepo repository.ClickRepository
}

// NewClickService crée et retourne une nouvelle instance de ClickService.
// C'est la fonction recommandée pour obtenir un service, assurant que toutes ses dépendances sont injectées.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

// RecordClick enregistre un nouvel événement de click dans la base de données.
// Cette méthode est appelée par le worker asynchrone.
func (s *ClickService) RecordClick(click *models.Click) error {
	// Vérifie que le service et son repository sont correctement initialisés
	if s == nil || s.clickRepo == nil {
		return fmt.Errorf("⚠️  service ClickService non initialisé")
	}

	// Vérifie que l'objet click fourni n'est pas nil
	if click == nil {
		return fmt.Errorf("⚠️  click fourni invalide : nil")
	}

	// Vérifie que le LinkID est valide (non nul)
	if click.LinkID == 0 {
		return fmt.Errorf("⚠️  LinkID invalide : %d", click.LinkID)
	}

	// Tente de persister le click via le repository et remonte l'erreur si échec
	if err := s.clickRepo.CreateClick(click); err != nil {
		return fmt.Errorf("✗ impossible d'enregistrer le click (LinkID=%d) : %w", click.LinkID, err)
	}
	return nil
}

// GetClicksCountByLinkID récupère le nombre total de clics pour un LinkID donné.
// Cette méthode pourrait être utilisée par le LinkService pour les statistiques, ou directement par l'API stats.
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
	// Vérifie que le service est initialisé
	if s == nil || s.clickRepo == nil {
		return 0, fmt.Errorf("⚠️  service ClickService non initialisé")
	}

	// Vérifie que le LinkID fourni est valide
	if linkID == 0 {
		return 0, fmt.Errorf("⚠️  LinkID invalide : %d", linkID)
	}

	// Appelle le repository pour obtenir le compte et propage l'erreur si nécessaire
	cnt, err := s.clickRepo.CountClicksByLinkID(linkID)
	if err != nil {
		return 0, fmt.Errorf("✗ impossible de récupérer le nombre de clics pour LinkID=%d : %w", linkID, err)
	}
	return cnt, nil
}
