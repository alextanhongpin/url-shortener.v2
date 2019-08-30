package domain

import (
	"context"
	"time"
)

// Shortener represents the operations for the url shortener.
type Shortener interface {
	Shorten(longURL string) (code string)
}

// URLRepository represents the query layer for the database.
type URLRepository interface {
	GetByCode(code string) (longURL string, err error)
	Create(code, longURL string, expireAt time.Time) (bool, error)
	CheckExists(code string) (bool, error)
}

// URLService represents the service layer for url operations.
type URLService interface {
	Get(context.Context, GetRequest) (*GetResponse, error)
	Put(context.Context, PutRequest) (*PutResponse, error)
	CheckExists(context.Context, CheckExistsRequest) (*CheckExistsResponse, error)
}

type (
	// GetRequest is the request body.
	GetRequest struct {
		Code string `json:"code" validate:"required,max=6" conform:"trim"`
	}

	// GetResponses is the response body.
	GetResponse struct {
		LongURL string `json:"long_url"`
	}
)
type (
	// PutRequest is the request body.
	PutRequest struct {
		Code     string    `json:"code" validate:"max=6" conform:"trim"`
		LongURL  string    `json:"long_url" validate:"required,url" conform:"trim"`
		ExpireAt time.Time `json:"expire_at"`
	}

	// PutResponse is the response body.
	PutResponse struct {
		Code string `json:"code"`
	}
)

type (
	// CheckExistsRequest is the request body.
	CheckExistsRequest struct {
		Code string `json:"code" validate:"required,max=6"`
	}

	// CheckExistsResponse is the response body.
	CheckExistsResponse struct {
		Exist bool `json:"exist"`
	}
)
