package githubclient

type GitHubResponse struct {
	Stars    int64   `json:"stargazers_count"`
	Language *string `json:"language"`
}
