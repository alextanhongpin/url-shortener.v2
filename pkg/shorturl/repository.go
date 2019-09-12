package shorturl

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/alextanhongpin/url-shortener/database"
	"github.com/alextanhongpin/url-shortener/domain"
)

const DuplicatePrimaryKeyViolation = "23505"

type Repository struct {
	stmts database.Stmts
	// db sq.StatementBuilderType
}

func NewRepository(db *sql.DB) *Repository {
	// dbCache := sq.NewStmtCacheProxy(db)
	return &Repository{
		// db: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(dbCache),
		stmts: rawStmts.MustPrepare(db),
	}
}

func (r *Repository) WithCode(code string) (string, error) {
	var longURL string
	// err := r.db.Select("long_url").From("url").Where(sq.Eq{"code": code}).Scan(&longURL)
	err := r.stmts[WithCode].QueryRow(code).Scan(&longURL)
	return longURL, err

}

func (r *Repository) CheckExists(code string) (bool, error) {
	var exists bool
	err := r.stmts[CheckExists].QueryRow(code).Scan(&exists)
	// err := r.db.Select("EXISTS").FromSelect(r.db.Select("1").From("url").Where(sq.Eq{"code": code}), "").Scan(&exists)
	return exists, err
}

func (r *Repository) Create(entity domain.ShortURL) (bool, error) {
	// res, err := r.db.Insert("url").Columns("code", "long_url", "expire_at").Values(entity.Code, entity.LongURL, entity.ExpireAt).Exec()
	res, err := r.stmts[Create].Exec(entity.Code, entity.LongURL, entity.ExpireAt)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == DuplicatePrimaryKeyViolation {
			return false, domain.ErrAlreadyExists
		}
	}
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows == 1, err
}
