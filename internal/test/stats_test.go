package test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/abdelbassat/github-stats-api/internal/githubclient"
	"github.com/abdelbassat/github-stats-api/internal/handler"
	"github.com/abdelbassat/github-stats-api/internal/model"
	"github.com/abdelbassat/github-stats-api/internal/repository"
	"github.com/abdelbassat/github-stats-api/internal/router"
	"github.com/abdelbassat/github-stats-api/internal/service"

	_ "modernc.org/sqlite"
)

const (
	missingUsername     = "stats-test-missing-7b9e4a20c5d6f830"
	invalidUsernameJSON = `{"error":"invalid username"}`
	userNotFoundJSON    = `{"error":"user not found"}`
	upstreamUnavailable = `{"error":"upstream unavailable"}`
)

var (
	goLanguage   = "Go"
	rustLanguage = "Rust"

	// Actualizar si cambian los datos reales de GitHub.
	expectedGitHubStats = model.Stats{
		Username:    "abdel084cb",
		TotalRepos:  5,
		TotalStars:  0,
		TopLanguage: &goLanguage,
	}

	cachedStatsFixture = model.Stats{
		Username:    "cached-user",
		TotalRepos:  8,
		TotalStars:  24,
		TopLanguage: &rustLanguage,
	}

	expiredStatsFixture = model.Stats{
		Username:    "abdel084cb",
		TotalRepos:  1,
		TotalStars:  999,
		TopLanguage: &rustLanguage,
	}

	nullLanguageStatsFixture = model.Stats{
		Username:   "cached-null",
		TotalRepos: 2,
		TotalStars: 7,
	}
)

func newTestStats(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "stats.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})

	statsRepository := repository.NewRepository(db)
	if err := statsRepository.CreateDatabase(db); err != nil {
		t.Fatal(err)
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	gitHubClient := githubclient.NewClient(httpClient)
	statsService := service.NewService(statsRepository, gitHubClient)
	statsHandler := handler.NewHandler(statsService)

	return router.NewRouter(statsHandler), db
}

func insertStats(t *testing.T, db *sql.DB, stats model.Stats, cacheAge time.Duration) model.Stats {
	t.Helper()
	stats.CachedAt = time.Now().Add(-cacheAge).UTC().Truncate(time.Second)

	const query = `
		INSERT INTO stats (username, total_repos, total_stars, top_language, cached_at)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err := db.Exec(query,
		stats.Username, stats.TotalRepos, stats.TotalStars,
		stats.TopLanguage, stats.CachedAt.Format(time.RFC3339),
	)
	if err != nil {
		t.Fatal(err)
	}
	return stats
}

func readStatsByUsername(t *testing.T, db *sql.DB, username string) model.Stats {
	t.Helper()

	const query = `
		SELECT username, total_repos, total_stars, top_language, cached_at
		FROM stats WHERE username = ?;
	`
	var stats model.Stats
	var cachedAt string

	err := db.QueryRow(query, username).Scan(
		&stats.Username, &stats.TotalRepos, &stats.TotalStars,
		&stats.TopLanguage, &cachedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	stats.CachedAt, err = time.Parse(time.RFC3339, cachedAt)
	if err != nil {
		t.Fatal(err)
	}
	return stats
}

func requestStats(t *testing.T, app http.Handler, username string, status int) []byte {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, "/stats/"+username, nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != status {
		t.Fatalf("status = %d, want %d; body = %s",
			response.Code, status, response.Body.String())
	}

	return response.Body.Bytes()
}

func assertStatsEqual(t *testing.T, got, want model.Stats) {
	t.Helper()

	if got.Username != want.Username ||
		got.TotalRepos != want.TotalRepos ||
		got.TotalStars != want.TotalStars ||
		!reflect.DeepEqual(got.TopLanguage, want.TopLanguage) ||
		!got.CachedAt.Equal(want.CachedAt) {
		t.Fatalf("stats = %+v, want %+v", got, want)
	}
}

func TestGetStatsByUsername(t *testing.T) {
	t.Run("cache_miss_and_caches_stats", func(t *testing.T) {
		app, db := newTestStats(t)
		want := expectedGitHubStats

		body := requestStats(t, app, want.Username, http.StatusOK)
		var got model.Stats
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}

		want.CachedAt = got.CachedAt
		assertStatsEqual(t, got, want)

		stored := readStatsByUsername(t, db, want.Username)
		assertStatsEqual(t, stored, got)
	})

	t.Run("cache_hit", func(t *testing.T) {
		app, db := newTestStats(t)
		want := insertStats(t, db, cachedStatsFixture, 0)

		body := requestStats(t, app, want.Username, http.StatusOK)
		var got model.Stats
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		assertStatsEqual(t, got, want)

		stored := readStatsByUsername(t, db, want.Username)
		assertStatsEqual(t, stored, want)
	})

	t.Run("cache_miss_expired", func(t *testing.T) {
		app, db := newTestStats(t)
		insertStats(t, db, expiredStatsFixture, 2*time.Hour)
		want := expectedGitHubStats

		body := requestStats(t, app, want.Username, http.StatusOK)
		var got model.Stats
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}

		want.CachedAt = got.CachedAt
		assertStatsEqual(t, got, want)

		stored := readStatsByUsername(t, db, want.Username)
		assertStatsEqual(t, stored, got)
	})

	t.Run("null_top_language", func(t *testing.T) {
		app, db := newTestStats(t)
		want := insertStats(t, db, nullLanguageStatsFixture, 0)

		body := requestStats(t, app, want.Username, http.StatusOK)
		var got model.Stats
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		assertStatsEqual(t, got, want)

		var fields map[string]any
		if err := json.Unmarshal(body, &fields); err != nil {
			t.Fatal(err)
		}
		language, exists := fields["top_language"]
		if !exists || language != nil {
			t.Fatal("expected top_language to be present and null")
		}
	})

	t.Run("github_user_not_found_404", func(t *testing.T) {
		app, _ := newTestStats(t)
		body := requestStats(t, app, missingUsername, http.StatusNotFound)
		if got := strings.TrimSpace(string(body)); got != userNotFoundJSON {
			t.Fatalf("body = %s, want %s", got, userNotFoundJSON)
		}
	})

	t.Run("github_unavailable_503", func(t *testing.T) {
		t.Skip("Pendiente: simular un fallo de conexión con GitHub")
	})

	t.Run("github_rate_limit_503", func(t *testing.T) {
		t.Skip("Pendiente: simular un rate limit de GitHub")
	})

	t.Run("invalid_username_400", func(t *testing.T) {
		for _, testCase := range []struct {
			name     string
			username string
		}{
			{"empty_username", ""},
			{"invalid_characters", "invalid_user"},
			{"username_too_long", strings.Repeat("a", 40)},
			{"username_starting_with_hyphen", "-abdel"},
			{"username_ending_with_hyphen", "abdel-"},
		} {
			t.Run(testCase.name, func(t *testing.T) {
				app, _ := newTestStats(t)
				body := requestStats(t, app, testCase.username, http.StatusBadRequest)
				if got := strings.TrimSpace(string(body)); got != invalidUsernameJSON {
					t.Fatalf("body = %s, want %s", got, invalidUsernameJSON)
				}
			})
		}
	})
}
