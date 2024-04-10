package models

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	Owner           string    `gorm:"not null"`
	Name            string    `gorm:"not null"`
	LastUpdated     time.Time
	UpdateFrequency string    `gorm:"default:'@daily'"`
	LastCommit      time.Time
	CISuccessRate   float64
	Active          bool      `gorm:"default:true"`
}

type CommitHistory struct {
	gorm.Model
	RepositoryID uint      `gorm:"not null"`
	SHA          string    `gorm:"not null"`
	Timestamp    time.Time `gorm:"not null"`
	CIStatus     string
}

type UpdateHistory struct {
	gorm.Model
	RepositoryID uint      `gorm:"not null"`
	StartTime    time.Time `gorm:"not null"`
	EndTime     time.Time
	Status      string    `gorm:"not null"`
	Error       string
} 