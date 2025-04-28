package services

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/models"
	"gorm.io/gorm"
)

// SchedulerService manages the scheduling of dependency updates
type SchedulerService struct {
	db            *gorm.DB
	githubService *GitHubService
	cron          *cron.Cron
}

// NewSchedulerService creates a new scheduler service instance
func NewSchedulerService(db *gorm.DB, githubService *GitHubService) *SchedulerService {
	return &SchedulerService{
		db:            db,
		githubService: githubService,
		cron:          cron.New(cron.WithSeconds()),
	}
}

// Start initializes and starts the scheduler
func (s *SchedulerService) Start() error {
	// Schedule repository updates
	if err := s.scheduleRepositoryUpdates(); err != nil {
		return fmt.Errorf("failed to schedule repository updates: %w", err)
	}

	s.cron.Start()
	return nil
}

// Stop gracefully stops the scheduler
func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

// scheduleRepositoryUpdates schedules updates for all active repositories
func (s *SchedulerService) scheduleRepositoryUpdates() error {
	var repositories []models.Repository
	if err := s.db.Where("active = ?", true).Find(&repositories).Error; err != nil {
		return err
	}

	for _, repo := range repositories {
		if err := s.scheduleRepository(repo); err != nil {
			logrus.WithError(err).WithField("repository", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)).
				Error("Failed to schedule repository")
		}
	}

	return nil
}

// scheduleRepository schedules updates for a single repository
func (s *SchedulerService) scheduleRepository(repo models.Repository) error {
	_, err := s.cron.AddFunc(repo.UpdateFrequency, func() {
		ctx := context.Background()
		if err := s.processRepository(ctx, repo); err != nil {
			logrus.WithError(err).WithField("repository", fmt.Sprintf("%s/%s", repo.Owner, repo.Name)).
				Error("Failed to process repository")
		}
	})

	return err
}

// processRepository processes a single repository for updates
func (s *SchedulerService) processRepository(ctx context.Context, repo models.Repository) error {
	// Check for dependency updates
	updates, err := s.githubService.GetDependencyUpdates(ctx, repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		return nil
	}

	// Create update history record
	updateHistory := &models.UpdateHistory{
		RepositoryID: repo.ID,
		StartTime:    time.Now(),
		Status:       "pending",
	}

	if err := s.db.Create(updateHistory).Error; err != nil {
		return err
	}

	// Create pull request for updates
	branch := fmt.Sprintf("deps/update-%d", time.Now().Unix())
	title := fmt.Sprintf("Update dependencies: %s", time.Now().Format("2006-01-02"))
	body := fmt.Sprintf("Automated dependency updates:\n\n%s", formatUpdates(updates))

	_, err = s.githubService.CreatePullRequest(ctx, repo.Owner, repo.Name, title, body, branch)
	if err != nil {
		updateHistory.Status = "failed"
		updateHistory.Error = err.Error()
		s.db.Save(updateHistory)
		return err
	}

	updateHistory.Status = "success"
	updateHistory.EndTime = time.Now()
	return s.db.Save(updateHistory).Error
}

// formatUpdates formats the list of updates for the PR description
func formatUpdates(updates []string) string {
	var result string
	for _, update := range updates {
		result += fmt.Sprintf("- %s\n", update)
	}
	return result
} 