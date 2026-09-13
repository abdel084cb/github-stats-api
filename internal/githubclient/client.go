package githubclient

import (
	"context"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    "https://api.github.com",
	}
}

func (c *Client) GetGitHubRepositories(
	ctx context.Context,
	username string,
) ([]GitHubResponse, error) {

	return []GitHubResponse{}, nil
}
