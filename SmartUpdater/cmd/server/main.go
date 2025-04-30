package main

import (
	"fmt"
	"log"
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

	// Initialize logger
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
	githubService := services.NewGitHubService(cfg.GithubToken)
	schedulerService := services.NewSchedulerService(db, githubService)

	// Start scheduler
	if err := schedulerService.Start(); err != nil {
		logrus.Fatalf("Failed to start scheduler: %v", err)
	}
	defer schedulerService.Stop()

	// Initialize router
	router := mux.NewRouter()
	handler := api.NewHandler(db, githubService)
	handler.RegisterRoutes(router)

	// Add middleware
	router.Use(loggingMiddleware)

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logrus.Info("Shutting down server...")
		if err := server.Close(); err != nil {
			logrus.Errorf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	logrus.Infof("Server starting on port %d", cfg.ServerPort)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatalf("Server failed: %v", err)
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		}).Info("HTTP request")

		next.ServeHTTP(w, r)
	})
} 