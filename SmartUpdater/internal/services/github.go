package services

import (
	"context"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	client *github.Client
	ctx    context.Context
}

func NewGitHubService(token string) *GitHubService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubService{
		client: client,
		ctx:    ctx,
	}
}

func (s *GitHubService) GetCommitHistory(owner, repo string, since time.Time) ([]*github.RepositoryCommit, error) {
	opts := &github.CommitsListOptions{
		Since: since,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := s.client.Repositories.ListCommits(s.ctx, owner, repo, opts)
		if err != nil {
			logrus.Errorf("Error fetching commits for %s/%s: %v", owner, repo, err)
			return nil, err
		}
		allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allCommits, nil
}

func (s *GitHubService) GetCIStatus(owner, repo, ref string) (string, error) {
	statuses, _, err := s.client.Repositories.GetCombinedStatus(s.ctx, owner, repo, ref, nil)
	if err != nil {
		logrus.Errorf("Error fetching CI status for %s/%s@%s: %v", owner, repo, ref, err)
		return "", err
	}

	return string(statuses.GetState()), nil
}

func (s *GitHubService) CreatePullRequest(owner, repo, title, body, head, base string) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Head:  github.String(head),
		Base:  github.String(base),
	}

	pr, _, err := s.client.PullRequests.Create(s.ctx, owner, repo, newPR)
	if err != nil {
		logrus.Errorf("Error creating pull request for %s/%s: %v", owner, repo, err)
		return nil, err
	}

	return pr, nil
} 