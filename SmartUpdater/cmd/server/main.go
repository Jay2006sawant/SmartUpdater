package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/api"
	"github.com/smartupdater/internal/config"
	"github.com/smartupdater/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	githubSvc := services.NewGitHubService(cfg.GithubToken)
	schedulerSvc := services.NewSchedulerService(db, githubSvc)

	// Start the scheduler
	schedulerSvc.Start()
	defer schedulerSvc.Stop()

	// Setup HTTP router
	r := mux.NewRouter()
	handler := api.NewHandler(schedulerSvc, db)
	handler.SetupRoutes(r)

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: r,
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logrus.Infof("Starting server on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stop
	logrus.Info("Shutting down server...")
} 