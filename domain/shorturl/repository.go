package shorturl

// Repository represents the query layer for the database.
type Repository interface {
	// CRUD.
	Create(URL) (string, error)

	// Scopes.
	WithCode(code string) (longURL string, err error)
	CheckExists(code string) (bool, error)
}
