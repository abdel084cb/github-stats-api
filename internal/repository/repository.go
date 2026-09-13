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
