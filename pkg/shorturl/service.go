package shorturl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/alextanhongpin/url-shortener/domain"
	"gopkg.in/go-playground/validator.v9"

	"github.com/leebenson/conform"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Service struct {
	urls      domain.Repository
	shortener domain.Shortener
	before    func(interface{}) error
}

func NewService(urls domain.Repository, shortener domain.Shortener, validator *validator.Validate) *Service {
	return &Service{
		urls:      urls,
		shortener: shortener,
		before: func(req interface{}) error {
			// Trim the strings with the conform tag.
			conform.Strings(req)

			// Validate requests.
			return validator.Struct(req)
		},
	}
}

func (s *Service) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}
	longURL, err := s.urls.WithCode(req.Code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrDoesNotExists
	}
	if err != nil {
		return nil, err
	}
	return domain.NewGetResponse(longURL), nil
}

func (s *Service) Put(ctx context.Context, req domain.PutRequest) (*domain.PutResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}

	// To ensure that it is really unique, even for the same url,
	// use the unique user id to hash the url.
	e := domain.ShortURL(req)
	if e.Code != "" {
		exists, err := s.urls.CheckExists(e.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrAlreadyExists
		}
	} else {
		e.Code = s.shortener.Shorten(req.LongURL)
	}

	_, err := s.urls.Create(e)
	for errors.Is(err, domain.ErrAlreadyExists) {
		e.Code = s.shortener.Shorten(e.Code + fmt.Sprint(rand.Int()))
		_, err = s.urls.Create(e)
	}
	if err != nil {
		return nil, err
	}
	return domain.NewPutResponse(e.Code), nil
}

func (s *Service) CheckExists(ctx context.Context, req domain.CheckExistsRequest) (*domain.CheckExistsResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}
	exists, err := s.urls.CheckExists(req.Code)
	if err != nil {
		return nil, err
	}
	return domain.NewCheckExistsResponse(exists), nil
}
