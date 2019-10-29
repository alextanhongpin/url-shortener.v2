package shorturl

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// UseCase represents the business logic layer for url operations.
type UseCase interface {
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
		Code string `json:"code" validate:"max=6" conform:"trim"`
		URL  string `json:"url" validate:"required,url" conform:"trim"`
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
	enc.AddString("url", req.URL)
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
