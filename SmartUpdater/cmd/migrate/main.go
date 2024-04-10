package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/config"
	"github.com/smartupdater/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&models.Repository{},
		&models.CommitHistory{},
		&models.UpdateHistory{},
	)
	if err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	logrus.Info("Migrations completed successfully")
} 