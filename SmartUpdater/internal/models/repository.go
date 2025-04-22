package models

import (
	"time"

	"gorm.io/gorm"
)

// Repository represents a GitHub repository being monitored
type Repository struct {
	gorm.Model
	Owner           string    `gorm:"not null"`           // GitHub repository owner
	Name            string    `gorm:"not null"`           // GitHub repository name
	LastUpdated     time.Time                             // Last dependency update time
	UpdateFrequency string    `gorm:"default:'@daily'"`   // Cron expression for update frequency
	LastCommit      time.Time                             // Last commit time
	CISuccessRate   float64                               // CI pipeline success rate
	Active          bool      `gorm:"default:true"`       // Whether the repository is actively monitored
}

// CommitHistory tracks commit information for analysis
type CommitHistory struct {
	gorm.Model
	RepositoryID uint      `gorm:"not null"`              // Reference to Repository
	SHA          string    `gorm:"not null"`              // Commit SHA
	Timestamp    time.Time `gorm:"not null"`              // Commit timestamp
	CIStatus     string                                   // CI pipeline status
}

// UpdateHistory tracks dependency update attempts
type UpdateHistory struct {
	gorm.Model
	RepositoryID uint      `gorm:"not null"`              // Reference to Repository
	StartTime    time.Time `gorm:"not null"`              // Update start time
	EndTime     time.Time                                // Update completion time
	Status      string    `gorm:"not null"`              // Update status (success/failed)
	Error       string                                   // Error message if failed
} 