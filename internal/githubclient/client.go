package githubclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abdelbassat/github-stats-api/internal/apperrors"
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

func (c *Client) GetGitHubRepositories(ctx context.Context, username string) ([]GitHubResponse, error) {
	const perPage = 100
	var allRepositories []GitHubResponse

	for page := 1; ; page++ {
		endpoint := fmt.Sprintf(
			"%s/users/%s/repos?per_page=%d&page=%d",
			c.baseURL,
			username,
			perPage,
			page,
		)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, apperrors.ErrUpstreamUnavailable
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				return nil, apperrors.ErrUserNotFound
			}

			return nil, apperrors.ErrUpstreamUnavailable
		}

		var repositories []GitHubResponse
		err = json.NewDecoder(resp.Body).Decode(&repositories)
		resp.Body.Close()

		if err != nil {
			return nil, apperrors.ErrUpstreamUnavailable
		}

		allRepositories = append(allRepositories, repositories...)

		if len(repositories) < perPage {
			return allRepositories, nil
		}
	}
}
