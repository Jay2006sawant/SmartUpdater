package services

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/models"
	"gorm.io/gorm"
)

type SchedulerService struct {
	db           *gorm.DB
	githubSvc    *GitHubService
	cron         *cron.Cron
	activeJobs   map[uint]cron.EntryID
}

func NewSchedulerService(db *gorm.DB, githubSvc *GitHubService) *SchedulerService {
	return &SchedulerService{
		db:         db,
		githubSvc:  githubSvc,
		cron:      cron.New(cron.WithSeconds()),
		activeJobs: make(map[uint]cron.EntryID),
	}
}

func (s *SchedulerService) Start() {
	s.cron.Start()
	s.scheduleAllRepositories()
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

func (s *SchedulerService) scheduleAllRepositories() {
	var repos []models.Repository
	if err := s.db.Where("active = ?", true).Find(&repos).Error; err != nil {
		logrus.Errorf("Error fetching active repositories: %v", err)
		return
	}

	for _, repo := range repos {
		s.scheduleRepository(&repo)
	}
}

func (s *SchedulerService) scheduleRepository(repo *models.Repository) {
	// Remove existing job if any
	if existingID, ok := s.activeJobs[repo.ID]; ok {
		s.cron.Remove(existingID)
		delete(s.activeJobs, repo.ID)
	}

	// Schedule new job
	entryID, err := s.cron.AddFunc(repo.UpdateFrequency, func() {
		s.processUpdate(repo)
	})

	if err != nil {
		logrus.Errorf("Error scheduling repository %s/%s: %v", repo.Owner, repo.Name, err)
		return
	}

	s.activeJobs[repo.ID] = entryID
}

func (s *SchedulerService) processUpdate(repo *models.Repository) {
	updateHistory := &models.UpdateHistory{
		RepositoryID: repo.ID,
		StartTime:    time.Now(),
		Status:      "in_progress",
	}
	s.db.Create(updateHistory)

	// Analyze commit patterns
	commits, err := s.githubSvc.GetCommitHistory(repo.Owner, repo.Name, time.Now().AddDate(0, -1, 0))
	if err != nil {
		s.finishUpdate(updateHistory, "failed", err.Error())
		return
	}

	// Calculate optimal time based on commit patterns
	if len(commits) > 0 {
		lastCommit := commits[0]
		repo.LastCommit = lastCommit.GetCommit().GetAuthor().GetDate()
		
		// Get CI status
		status, err := s.githubSvc.GetCIStatus(repo.Owner, repo.Name, lastCommit.GetSHA())
		if err != nil {
			logrus.Warnf("Error fetching CI status: %v", err)
		} else {
			// Update CI success rate
			isSuccess := status == "success"
			repo.CISuccessRate = (repo.CISuccessRate*0.7 + float64(btoi(isSuccess))*0.3)
		}
	}

	repo.LastUpdated = time.Now()
	s.db.Save(repo)

	s.finishUpdate(updateHistory, "completed", "")
}

func (s *SchedulerService) finishUpdate(history *models.UpdateHistory, status, errorMsg string) {
	history.EndTime = time.Now()
	history.Status = status
	history.Error = errorMsg
	s.db.Save(history)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
} 