package shorten

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/alextanhongpin/url-shortener/database"
	"github.com/alextanhongpin/url-shortener/domain"
)

const DuplicatePrimaryKeyViolation = "23505"

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

func (r *Repository) Create(entity domain.ShortURL) (bool, error) {
	res, err := r.stmts[Create].Exec(entity.Code, entity.LongURL, entity.ExpireAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// log.Printf("pqError: %#v", pqErr)
			if pqErr.Code == DuplicatePrimaryKeyViolation {
				return false, ErrAlreadyExists
			}
		}
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows == 1, err
}
