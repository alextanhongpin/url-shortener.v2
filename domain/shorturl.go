package domain

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	ErrAlreadyExists = errors.New("short url already exists")
	ErrDoesNotExists = errors.New("short url does not exist")
)

// Shortener represents the operations for the url shortener.
type Shortener interface {
	Shorten(longURL string) (code string)
}

// ShortURL is the entity.
type ShortURL struct {
	Code     string    `json:"code"`
	LongURL  string    `json:"long_url"`
	ExpireAt time.Time `json:"expire_at"`
}

// Repository represents the query layer for the database.
type Repository interface {
	// CRUD.
	Create(ShortURL) (bool, error)
	// All(limit, offset int) ([]ShortURL, error)
	// Find(id int) (ShortURL, error)
	// Update(ShortURL) (bool, error)
	// Delete(id int) (bool, error)

	// Scopes.
	WithCode(code string) (longURL string, err error)
	CheckExists(code string) (bool, error)
}

// Service represents the service layer for url operations.
type Service interface {
	Get(context.Context, GetRequest) (*GetResponse, error)
	Put(context.Context, PutRequest) (*PutResponse, error)
	CheckExists(context.Context, CheckExistsRequest) (*CheckExistsResponse, error)
}

type (
	// GetRequest is the request body.
	GetRequest struct {
		Code string `json:"code" validate:"required,max=6" conform:"trim"`
	}

	// GetResponse is the response body.
	GetResponse struct {
		LongURL string `json:"long_url"`
	}
)

func NewGetResponse(longURL string) *GetResponse {
	return &GetResponse{longURL}
}

func (req *GetRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("code", req.Code)
	return nil
}

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

func NewPutResponse(code string) *PutResponse {
	return &PutResponse{code}
}

func (req *PutRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("code", req.Code)
	enc.AddString("long_url", req.LongURL)
	enc.AddTime("expire_at", req.ExpireAt)
	return nil
}

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

func NewCheckExistsResponse(exist bool) *CheckExistsResponse {
	return &CheckExistsResponse{exist}
}

func (req *CheckExistsRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("code", req.Code)
	return nil
}
