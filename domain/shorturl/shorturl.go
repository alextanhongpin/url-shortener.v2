package shorturl

import (
	"errors"

	"github.com/alextanhongpin/url-shortener/domain/datetime"
)

var (
	ErrAlreadyExists = errors.New("short url already exists")
	ErrDoesNotExists = errors.New("short url does not exist")
)

// URL is the entity.
type URL struct {
	Code      string            `json:"code"`
	LongURL   string            `json:"long_url"`
	CreatedAt datetime.DateTime `json:"created_at,omitempty"`
	UpdatedAt datetime.DateTime `json:"updated_at,omitempty"`
}

func New(code, url string) URL {
	return URL{
		Code:    code,
		LongURL: url,
	}
}
