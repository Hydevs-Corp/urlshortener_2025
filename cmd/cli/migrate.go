package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// dbPathFlag permettra d'overrider le chemin vers la DB via le flag --db
var dbPathFlag string

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Charger la configuration : priorité au flag --db, sinon config, sinon par défaut
		dbPath := dbPathFlag
		if dbPath == "" {
			// Tente de charger la config depuis le fichier
			if cfg, err := config.LoadConfig(); err == nil && cfg != nil && cfg.Database.Name != "" {
				dbPath = cfg.Database.Name
			} else {
				// Valeur par défaut si aucune config
				dbPath = "url_shortener.db"
			}
		}

		// Initialiser la connexion à la BDD
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: impossible d'ouvrir la base de données '%s' : %v", dbPath, err)
		}

		// Récupérer la connexion SQL pour la fermer proprement
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// Assure la fermeture propre de la connexion après la migration
		defer func() {
			if cerr := sqlDB.Close(); cerr != nil {
				log.Printf("⚠️  WARN: erreur lors de la fermeture de la DB : %v", cerr)
			}
		}()

		// Exécuter les migrations automatiques de GORM pour tous les modèles
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("✗ FATAL: échec des migrations : %v", err)
		}

		// Message final de succès
		fmt.Println("✓ Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	// Ajouter la commande migrate au RootCmd
	cmd2.RootCmd.AddCommand(MigrateCmd)

	// Déclare un flag optionnel --db pour overrider le chemin de la base de données si nécessaire
	MigrateCmd.Flags().StringVar(&dbPathFlag, "db", "", "Chemin vers le fichier SQLite (overrides config)")
}
