package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// longURLFlag stocke la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.
		longURL, err := cmd.Flags().GetString("url")
		if err != nil || longURL == "" {
			log.Fatalf("FATAL: Le flag --url est requis et doit être une chaîne non vide.")
		}

		// TODO Validation basique du format de l'URL avec le package url et la fonction ParseRequestURI
		// si erreur, os.Exit(1)
		_, err = url.ParseRequestURI(longURL)
		if err != nil {
			log.Fatalf("FATAL: L'URL fournie n'est pas valide: %v", err)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg, err := cmd2.Cfg, nil
		if err != nil {
			log.Fatalf("FATAL: Échec du chargement de la configuration: %v", err)
		}

		// TODO : Initialiser la connexion à la base de données SQLite.

		sqlDB, err := repository.ConnectDatabase().DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		defer func() {
			sqlDB, err := repository.DB().DB()
			if err != nil {
				log.Printf("WARN: Échec de l'obtention de la base de données SQL pour la fermeture: %v", err)
				return
			}
			err = sqlDB.Close()
			if err != nil {
				log.Printf("WARN: Échec de la fermeture de la connexion à la base de données: %v", err)
			} else {
				log.Printf("✅ Connexion à la base de données fermée.")
			}
		}()
		
		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		linkRepo := repository.NewLinkRepository(repository.DB())
		linkService := services.NewLinkService(linkRepo)

		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		// os.Exit(1) si erreur

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	// Définir le flag --url pour la commande create.
	CreateCmd.Flags().StringP("url", "u", "", "L'URL longue à raccourcir")

	// Marquer le flag comme requis
	CreateCmd.MarkFlagRequired("url")

	// Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(CreateCmd)

}
