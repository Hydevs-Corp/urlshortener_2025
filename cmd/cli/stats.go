package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

var shortCodeFlag string


// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortCodeFlag == "" {
			fmt.Println("Erreur: le flag --code est requis.")
			os.Exit(1)
		}

		cfg, err := cmd2.LoadConfig()
		if err != nil {
			log.Fatalf("FATAL: Échec du chargement de la configuration: %v", err)
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}


		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Printf("WARNING: Échec de la fermeture de la base de données: %v", err)
			}
		}()

		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("Erreur: Aucun lien trouvé pour le code court '%s'.\n", shortCodeFlag)
			} else {
				log.Fatalf("FATAL: Échec de la récupération des statistiques du lien: %v", err)
			}
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

func init() {
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Le code court de l'URL pour laquelle récupérer les statistiques.")

	StatsCmd.MarkFlagRequired("code")

	cmd2.RootCmd.AddCommand(StatsCmd)
}
