package models

import "time"

// Click représente un événement de clic sur un lien raccourci.
// GORM utilisera ces tags pour créer la table 'clicks'.
type Click struct {
	ID        uint      `gorm:"primaryKey"`        // Clé primaire
	LinkID    uint      `gorm:"index"`             // Clé étrangère vers la table 'links', indexée pour des requêtes efficaces
	Link      Link      `gorm:"foreignKey:LinkID"` // Relation GORM: indique que LinkID est une FK vers le champ ID de Link
	Timestamp time.Time // Horodatage précis du clic
	UserAgent string    `gorm:"size:255"` // User-Agent de l'utilisateur qui a cliqué
	IP string    `gorm:"size:50"`  // Adresse IP de l'utilisateur
}

type ClickEvent struct {
	LinkID    uint      // ID du lien ajouté par l'utilisateur
	Timestamp time.Time // Heure de l'event
	UserAgent string    // Referrer du navigateur
	IP string    // Adresse IP de l'utilisateur
}
