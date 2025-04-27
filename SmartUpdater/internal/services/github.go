package services

import (
	"context"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"
	"github.com/smartupdater/internal/models"
)

// GitHubService handles interactions with GitHub API
type GitHubService struct {
	client *github.Client
}

// NewGitHubService creates a new GitHub service instance
func NewGitHubService(token string) *GitHubService {
	client := github.NewClient(nil).WithAuthToken(token)
	return &GitHubService{client: client}
}

// GetRepositoryInfo fetches repository information from GitHub
func (s *GitHubService) GetRepositoryInfo(ctx context.Context, owner, repo string) (*models.Repository, error) {
	repository, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return &models.Repository{
		Owner:       owner,
		Name:        repo,
		LastUpdated: repository.GetUpdatedAt().Time,
		LastCommit:  repository.GetPushedAt().Time,
		Active:      true,
	}, nil
}

// GetCommitHistory fetches recent commit history for a repository
func (s *GitHubService) GetCommitHistory(ctx context.Context, owner, repo string, since time.Time) ([]*models.CommitHistory, error) {
	commits, _, err := s.client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		Since: since,
	})
	if err != nil {
		return nil, err
	}

	var commitHistory []*models.CommitHistory
	for _, commit := range commits {
		commitHistory = append(commitHistory, &models.CommitHistory{
			SHA:       commit.GetSHA(),
			Timestamp: commit.GetCommit().GetAuthor().GetDate().Time,
			CIStatus:  "unknown", // Will be updated by CI service
		})
	}

	return commitHistory, nil
}

// GetDependencyUpdates checks for available dependency updates
func (s *GitHubService) GetDependencyUpdates(ctx context.Context, owner, repo string) ([]string, error) {
	alerts, _, err := s.client.Dependabot.ListRepoAlerts(ctx, owner, repo, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch dependency alerts")
		return nil, err
	}

	var updates []string
	for _, alert := range alerts {
		if alert.GetState() == "open" {
			updates = append(updates, alert.GetDependency().GetPackage().GetName())
		}
	}

	return updates, nil
}

// CreatePullRequest creates a new pull request for dependency updates
func (s *GitHubService) CreatePullRequest(ctx context.Context, owner, repo, title, body, branch string) (*github.PullRequest, error) {
	pr, _, err := s.client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &branch,
		Base:  github.String("main"),
	})
	if err != nil {
		return nil, err
	}

	return pr, nil
} 