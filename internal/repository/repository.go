package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/abdelbassat/github-stats-api/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateDatabase(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS stats (
		username TEXT PRIMARY KEY,
		total_repos INTEGER NOT NULL,
		total_stars INTEGER NOT NULL,
		top_language TEXT,
		cached_at TEXT NOT NULL
	);`

	_, err := db.Exec(query)

	return err
}

func (r *Repository) GetStatsByUsername(ctx context.Context, username string) (model.Stats, error) {
	const query = `
		SELECT * FROM stats WHERE username = ?;
	`

	row := r.db.QueryRowContext(ctx, query, username)

	if row.Err() != nil {
		return model.Stats{}, row.Err()
	}

	var result model.Stats
	var cachedAtText string
	var topLanguage sql.NullString

	err := row.Scan(
		&result.Username,
		&result.TotalRepos,
		&result.TotalStars,
		&topLanguage,
		&cachedAtText,
	)
	if err != nil {
		return model.Stats{}, err
	}

	result.CachedAt, err = time.Parse(time.RFC3339, cachedAtText)
	if err != nil {
		return model.Stats{}, err
	}

	if topLanguage.Valid {
		language := topLanguage.String
		result.TopLanguage = &language
	}

	return result, nil
}

func (r *Repository) CacheStatsByUsername(ctx context.Context, stats model.Stats) error {
	const query = `
		INSERT INTO stats (
			username,
			total_repos,
			total_stars,
			top_language,
			cached_at
		)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
			total_repos = excluded.total_repos,
			total_stars = excluded.total_stars,
			top_language = excluded.top_language,
			cached_at = excluded.cached_at;
	`

	cachedAtText := stats.CachedAt.UTC().Format(time.RFC3339)

	var topLanguage any
	if stats.TopLanguage != nil {
		topLanguage = *stats.TopLanguage
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		stats.Username,
		stats.TotalRepos,
		stats.TotalStars,
		topLanguage,
		cachedAtText,
	)
	if err != nil {
		return err
	}

	return nil
}
