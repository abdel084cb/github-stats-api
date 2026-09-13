package repository

import (
	"database/sql"
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
