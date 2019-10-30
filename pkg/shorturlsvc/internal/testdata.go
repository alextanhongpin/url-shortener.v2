package internal

import (
	"context"
	"errors"

	"github.com/alextanhongpin/url-shortener/domain/shorturl"
)

func NewMockUseCase() *mockUseCase {
	return &mockUseCase{Err: errors.New("not implemented")}
}

type mockUseCase struct {
	GetResponse         *shorturl.GetResponse
	PutResponse         *shorturl.PutResponse
	CheckExistsResponse *shorturl.CheckExistsResponse
	Err                 error
}

func (m *mockUseCase) Get(context.Context, shorturl.GetRequest) (*shorturl.GetResponse, error) {
	return m.GetResponse, m.Err
}

func (m *mockUseCase) Put(context.Context, shorturl.PutRequest) (*shorturl.PutResponse, error) {
	return m.PutResponse, m.Err
}

func (m *mockUseCase) CheckExists(context.Context, shorturl.CheckExistsRequest) (*shorturl.CheckExistsResponse, error) {
	return m.CheckExistsResponse, m.Err
}
