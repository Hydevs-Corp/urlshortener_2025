package repository

import (
	"log"

	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	// Configure GORM logger to ignore ErrRecordNotFound logs so the console
	// isn't cluttered with expected "record not found" messages.
	newLogger := logger.New(
		log.New(log.Writer(), "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("failed to connect database (%s): %v", path, err)
	}

	db = database
	log.Printf("✅ SQLite connected at %s", path)
	return db
}
