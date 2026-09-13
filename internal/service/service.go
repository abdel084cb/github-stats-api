package service

import (
	"context"
	"errors"
	"regexp"
	"time"

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

func calculateStats(username string, repositories []githubclient.GitHubResponse) model.Stats {
	stats := model.Stats{
		Username:   username,
		TotalRepos: int64(len(repositories)),
	}

	languageCounts := make(map[string]int)
	maxLanguageCount := 0

	for _, repo := range repositories {
		stats.TotalStars += repo.Stars

		if repo.Language == nil {
			continue
		}

		language := *repo.Language
		languageCounts[language]++

		if languageCounts[language] > maxLanguageCount {
			maxLanguageCount = languageCounts[language]
			stats.TopLanguage = &language
		}
	}

	return stats
}

func (s *Service) GetStatsByUsername(ctx context.Context, username string) (model.Stats, error) {
	if !isValidUsername(username) {
		return model.Stats{}, apperrors.ErrInvalidUsername
	}

	cachedStats, err := s.repository.GetStatsByUsername(ctx, username)

	if err == nil {
		cacheAge := time.Since(cachedStats.CachedAt)

		if cacheAge >= 0 && cacheAge < time.Hour {
			return cachedStats, nil
		}
	} else if !errors.Is(err, apperrors.ErrCacheMiss) {
		return model.Stats{}, err
	}

	githubRepositories, err := s.githubClient.GetGitHubRepositories(ctx, username)
	if err != nil {
		return model.Stats{}, err
	}

	stats := calculateStats(username, githubRepositories)
	stats.CachedAt = time.Now().UTC().Truncate(time.Second)

	if err := s.repository.CacheStatsByUsername(ctx, stats); err != nil {
		return model.Stats{}, err
	}

	return stats, nil
}
