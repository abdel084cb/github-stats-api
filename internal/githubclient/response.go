package githubclient

type GithubResponse struct {
	Stars    int64   `json:"stargazers_count"`
	Language *string `json:"language"`
}
