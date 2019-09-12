package shorturl

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// Here we are implementing a shortener with md5, but we can implement other
// strategy, such as SHA, murmur, or maybe even base62 encoding.

// N represents the maximum length of the shortened url code.
const N = 6

// Shortener implements the Shortener interface.
type Shortener struct{}

// NewShortener returns a new Shortener implementation.
func NewShortener() *Shortener {
	return &Shortener{}
}

// Shorten takes a long url and return the short url code.
func (s *Shortener) Shorten(url string) string {
	return Shorten(sha256.New(), url)[:N]
}

// ShortenSha256 shortens a string using SHA256.
func Shorten(h hash.Hash, str string) string {
	h.Write([]byte(str))
	code := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(code)
}
