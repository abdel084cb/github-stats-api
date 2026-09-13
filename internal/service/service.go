package service

import (
	"github.com/abdelbassat/github-stats-api/internal/githubclient"
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
