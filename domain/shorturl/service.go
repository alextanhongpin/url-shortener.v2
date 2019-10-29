package shorturl

// Service represents the operations for the url shortener.
type Service interface {
	Shorten(longURL string) (code string)
}
