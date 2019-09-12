package shorturl_test

import (
	"crypto/sha256"
	"testing"

	"github.com/alextanhongpin/url-shortener/pkg/shorten"
)

func testShortenString(t *testing.T, str string) {
	shortener := shorten.New()
	expected := shorten.Shorten(sha256.New(), str)[:6]
	actual := shortener.Shorten(str)

	if expected != actual {
		t.Errorf("expected %s, actual %s", expected, actual)
	}
}

func TestShortener(t *testing.T) {
	t.Run("shortens a valid url", func(t *testing.T) {
		testShortenString(t, "https://www.google.com")
	})

	t.Run("shortens an empty string", func(t *testing.T) {
		testShortenString(t, "")
	})
}
