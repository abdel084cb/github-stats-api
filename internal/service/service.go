package service

import (
	"context"
	"regexp"

	"github.com/abdelbassat/github-stats-api/internal/apperrors"
	"github.com/abdelbassat/github-stats-api/internal/githubclient"
	"github.com/abdelbassat/github-stats-api/internal/model"
	"github.com/abdelbassat/github-stats-api/internal/repository"
)

type Service struct {
	repository   *repository.Repository
	githubClient *githubclient.Client
}

func NewService(repository *repository.Repository, githubClient *githubclient.Client) *Service {
	return &Service{
		repository:   repository,
		githubClient: githubClient,
	}
}

const maxUsernameLength = 39

var usernamePattern *regexp.Regexp = regexp.MustCompile(
	`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`,
)

func isValidUsername(username string) bool {
	if username == "" {
		return false
	}

	if len(username) > maxUsernameLength {
		return false
	}

	return usernamePattern.MatchString(username)
}

func (s *Service) GetStatsByUsername(ctx context.Context, username string) (stats model.Stats, err error) {

	if !isValidUsername(username) {
		return model.Stats{}, apperrors.ErrInvalidUsername
	}

	return model.Stats{}, nil
}
