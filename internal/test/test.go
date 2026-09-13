package test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/abdelbassat/github-stats-api/internal/model"
)

func newTestStats(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()
	return nil, nil
}

func insertStatsCacheable(t *testing.T, stats model.Stats) model.Stats {
	t.Helper()
	return model.Stats{}
}

func insertStatsNotCacheable(t *testing.T, stats model.Stats) model.Stats {
	t.Helper()
	return model.Stats{}
}

func readStatsByUsername(t *testing.T, username string) model.Stats {
	t.Helper()
	return model.Stats{}
}

func TestGetStatsByUsername(t *testing.T) {
	t.Run("cache_miss_and_caches_stats", func(t *testing.T) {

	})

	t.Run("cache_hit", func(t *testing.T) {

	})

	t.Run("cache_miss_expired", func(t *testing.T) {

	})

	t.Run("null_top_language", func(t *testing.T) {

	})

	t.Run("github_user_not_found_404", func(t *testing.T) {

	})

	t.Run("github_unavailable_503", func(t *testing.T) {

	})

	t.Run("github_rate_limit_503", func(t *testing.T) {

	})

	t.Run("invalid_username_400", func(t *testing.T) {
		t.Run("empty_username", func(t *testing.T) {

		})

		t.Run("invalid_characters", func(t *testing.T) {

		})

		t.Run("username_too_long", func(t *testing.T) {

		})

		t.Run("username_starting_with_hyphen", func(t *testing.T) {

		})

		t.Run("username_ending_with_hyphen", func(t *testing.T) {

		})
	})
}
