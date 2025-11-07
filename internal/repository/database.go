package repository

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// DB retourne l'instance GORM initialisée (panic si non init).
func DB() *gorm.DB {
	if db == nil {
		log.Fatal("database not initialized: call repository.ConnectDatabase() first")
	}
	return db
}

// ConnectDatabase ouvre la base SQLite grâce au chemin défini dans la config Viper.
// Clé attendue : database.name (fallback: url_shortener.db).
func ConnectDatabase() *gorm.DB {
	path := viper.GetString("database.name")
	if path == "" {
		path = "url_shortener.db"
	}

	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database (%s): %v", path, err)
	}

	db = database
	log.Printf("✅ SQLite connected at %s", path)
	return db
}
