package workers

import (
	"log"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Nécessaire pour interagir avec le ClickRepository
)

// StartClickWorkers lance un pool de goroutines "workers" pour traiter les événements de clic.
// Chaque worker lira depuis le même 'clickEventsChan' et utilisera le 'clickRepo' pour la persistance.
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("▶ Démarrage de %d worker(s) pour le traitement des clicks...", workerCount)
	for i := 0; i < workerCount; i++ {
		// Lance chaque worker dans sa propre goroutine.
		// Le channel est passé en lecture seule (<-chan) pour renforcer l'immutabilité du channel à l'intérieur du worker.
		workerID := i + 1 //  id simple de compréhension
		log.Printf("▶ Worker %d lancé", workerID)
		go clickWorker(workerID, clickEventsChan, clickRepo)
	}
}

// clickWorker est la fonction exécutée par chaque goroutine worker.
// Elle tourne indéfiniment, lisant les événements de clic dès qu'ils sont disponibles dans le channel.
// Implémentation : conversion ClickEvent -> models.Click, validation minimale,
// persistance via clickRepo.CreateClick avec retry/backoff limité.
func clickWorker(workerID int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	const (
		maxRetries      = 3
		initialBackoffMs = 100
	)

	for event := range clickEventsChan { // Boucle qui lit les événements du channel
		// Conversion ClickEvent -> models.Click
		click := models.Click{
			LinkID:    event.LinkID,
			Timestamp: event.Timestamp,
			UserAgent: event.UserAgent,
			IP: event.IP,
		}

		// Validation minimale
		if click.LinkID == 0 {
			log.Printf("⚠️  Worker %d — événement de click invalide : LinkID=%d — UserAgent=%q — IP=%q", workerID, click.LinkID, click.UserAgent, click.IP)
			continue
		}

		// Persister le click avec retry/backoff simple (meilleure gestion de la surcharge)
		var err error
		backoff := time.Duration(initialBackoffMs) * time.Millisecond
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = clickRepo.CreateClick(&click)
			if err == nil {
				// Succès
				log.Printf("✓ Worker %d — click enregistré (LinkID=%d, ts=%s)", workerID, click.LinkID, click.Timestamp.Format(time.RFC3339))
				break
			}

			// Échec — loguer et préparer éventuellement un retry
			log.Printf("✗ Worker %d — échec enregistrement (LinkID=%d) — tentative %d/%d — erreur: %v", workerID, click.LinkID, attempt, maxRetries, err)

			if attempt < maxRetries {
				time.Sleep(backoff)
				backoff = backoff * 2
			}
		}

		if err != nil {
			// Après toutes les tentatives, on abandonne et on loggue l'erreur finale.
			log.Printf("‼️  Worker %d — abandon après %d tentatives pour LinkID=%d — erreur finale: %v", workerID, maxRetries, click.LinkID, err)
		}
	}
}
