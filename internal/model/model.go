package model

import "time"

type Stats struct {
	Username    string    `json:"username"`
	TotalRepos  int64     `json:"total_repos"`
	TotalStars  int64     `json:"total_stars"`
	TopLanguage *string   `json:"top_language"`
	CachedAt    time.Time `json:"cached_at"`
}
