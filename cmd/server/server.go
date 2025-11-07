package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/axellelanca/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RunServerCmd represents the 'run-server' Cobra command.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Start the URL shortener API server and background processes.",
	Long:  `This command initializes the database, configures the APIs, starts background click workers and the URL monitor, then starts the HTTP server.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatal("Error: global configuration is not loaded (cmd2.Cfg is nil)")
		}

		// Database init + migrations
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to open SQLite database: %v", err)
		}
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("AutoMigrate error: %v", err)
		}

		// Repositories
		linkRepo := repository.NewLinkRepository(db)
		clickRepo := repository.NewClickRepository(db)
		log.Println("Repositories initialized.")

		// Services
		linkService := services.NewLinkService(linkRepo)
		clickService := services.NewClickService(clickRepo)
		log.Println("Domain services initialized.")

		// Click events channel + workers (use models.ClickEvent)
		clickEvents := make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		api.ClickEventsChannel = clickEvents
		workers.StartClickWorkers(cfg.Analytics.WorkerCount, clickEvents, clickRepo)
		log.Printf("Click event channel initialized with buffer %d. Started %d click worker(s).",
			cfg.Analytics.BufferSize, cfg.Analytics.WorkerCount)

		// URL monitor
		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()
		log.Printf("URL monitor started with interval %v.", monitorInterval)

		// Router and routes
		router := gin.Default()
		api.RegisterRoutes(router, linkService, clickService, clickEvents)
		log.Println("API routes configured.")

		// HTTP server
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// Start server
		go func() {
			log.Printf("HTTP server listening on %s", serverAddr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Server error: %v", err)
			}
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown signal received. Stopping server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		log.Println("Shutting down... giving workers time to finish.")
		time.Sleep(5 * time.Second)
		log.Println("Server stopped cleanly.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
