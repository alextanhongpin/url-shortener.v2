package shorten

import (
	"database/sql"
	"time"

	"github.com/alextanhongpin/url-shortener/infra/database"
)

type Repository struct {
	stmts database.Stmts
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		stmts: database.Prepare(db, rawStmts),
	}
}

func (r *Repository) GetByCode(code string) (string, error) {
	var longURL string
	err := r.stmts[GetByCode].QueryRow(code).Scan(&longURL)
	return longURL, err
}

func (r *Repository) CheckExists(code string) (bool, error) {
	var exists bool
	err := r.stmts[CheckExists].QueryRow(code).Scan(&exists)
	return exists, err
}

func (r *Repository) Create(code, longURL string, expireAt time.Time) (bool, error) {
	res, err := r.stmts[Create].Exec(code, longURL, expireAt)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows == 1, err
}
