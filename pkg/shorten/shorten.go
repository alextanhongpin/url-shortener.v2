package shorten

import (
	"crypto/md5"
	"encoding/base64"
)

// Here we are implementing a shortener with md5, but we can implement other
// strategy, such as SHA, murmur, or maybe even base62 encoding.

// N represents the maximum length of the shortened url code.
const N = 6

// Shortener implements the Shortener interface.
type Shortener struct{}

// New returns a new Shortener implementation.
func New() *Shortener {
	return &Shortener{}
}

// Shorten takes a long url and return the short url code.
func (s *Shortener) Shorten(url string) string {
	return Shorten(url, N)
}

// Shorten shortens a string to length n.
func Shorten(str string, n int) string {
	h := md5.New()
	h.Write([]byte(str))
	code := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(code)[:n]
}
