package shorturlsvc

import (
	"database/sql"

	"github.com/alextanhongpin/url-shortener/database"
	"github.com/alextanhongpin/url-shortener/domain/shorturl"
)

const DuplicatePrimaryKeyViolation = "23505"

type Repository struct {
	stmts database.Stmts
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		stmts: rawStmts.MustPrepare(db),
	}
}

func (r *Repository) WithCode(code string) (string, error) {
	var longURL string
	err := r.stmts[WithCode].QueryRow(code).Scan(&longURL)
	return longURL, err

}

func (r *Repository) CheckExists(code string) (bool, error) {
	var exists bool
	err := r.stmts[CheckExists].QueryRow(code).Scan(&exists)
	return exists, err
}

func (r *Repository) Create(entity shorturl.URL) (string, error) {
	var id string
	err := r.stmts[Create].QueryRow(entity.Code, entity.LongURL).Scan(&id)
	return id, err
}
