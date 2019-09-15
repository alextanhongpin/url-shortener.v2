package shorturl

import "github.com/alextanhongpin/url-shortener/database"

// There are no type safety for the sql queries - the query might fail when we
// execute the function. One way to check if the queries are "safe" is to
// prepare them when the application starts. This can reduce the syntax error,
// even though it does not help when the arguments are invalid (wrong order
// etc).

const (
	_ database.Stmt = iota

	// WithCode returns the long url for the given short url code, only if
	// the url has not expired yet.
	WithCode

	// Create creates a new entry of the short and long url with the expiry
	// date.
	Create

	// CheckExists returns true if the short url code already exists.
	CheckExists
)

var rawStmts = database.RawStmts{
	WithCode: `
		SELECT long_url 
		  FROM url
		 WHERE code = $1
	`,

	CheckExists: `
		SELECT EXISTS (SELECT 1 FROM url WHERE code = $1)
	`,

	Create: `
		INSERT INTO url (code, long_url, expire_at)
		VALUES ($1, $2, $3)
	`,
}
