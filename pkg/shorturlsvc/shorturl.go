package shorturlsvc

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// Here we are implementing a shortener with md5, but we can implement other
// strategy, such as SHA, murmur, or maybe even base62 encoding.

// N represents the maximum length of the shortened url code.
const N = 6

// Service implements the Shortener interface.
type Service struct{}

// NewService returns a new Shortener implementation.
func NewService() *Service {
	return &Service{}
}

// Shorten takes a long url and return the short url code.
func (s *Service) Shorten(url string) string {
	return Shorten(sha256.New(), url)[:N]
}

// Shorten shortens a string using SHA256.
func Shorten(h hash.Hash, str string) string {
	h.Write([]byte(str))
	code := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(code)
}
