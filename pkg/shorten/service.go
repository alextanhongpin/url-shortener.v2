package shorten

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/alextanhongpin/url-shortener/domain"

	"github.com/leebenson/conform"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Service struct {
	repo      domain.Repository
	shortener domain.Shortener
	before    func(interface{}) error
}

func NewService(repo domain.Repository, shortener domain.Shortener) *Service {
	return &Service{
		repo:      repo,
		shortener: shortener,
		before: func(req interface{}) error {
			// Trim the strings with the conform tag.
			conform.Strings(req)

			// Validate requests.
			return domain.Validator.Struct(req)
		},
	}
}

func (s *Service) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}

	longURL, err := s.repo.GetByCode(req.Code)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDoesNotExists
	}
	if err != nil {
		return nil, err
	}
	return &domain.GetResponse{
		LongURL: longURL,
	}, nil
}

func (s *Service) Put(ctx context.Context, req domain.PutRequest) (*domain.PutResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}

	// To ensure that it is really unique, even for the same url,
	// use the unique user id to hash the url.
	e := domain.ShortURL(req)
	if e.Code != "" {
		exists, err := s.repo.CheckExists(e.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrAlreadyExists
		}
	} else {
		e.Code = s.shortener.Shorten(req.LongURL)
	}

	_, err := s.repo.Create(e)
	for err == domain.ErrAlreadyExists {
		e.Code = s.shortener.Shorten(e.Code + fmt.Sprint(rand.Int()))
		_, err = s.repo.Create(e)
	}
	if err != nil {
		return nil, err
	}
	return &domain.PutResponse{Code: e.Code}, nil
}

func (s *Service) CheckExists(ctx context.Context, req domain.CheckExistsRequest) (*domain.CheckExistsResponse, error) {
	if err := s.before(&req); err != nil {
		return nil, err
	}
	exists, err := s.repo.CheckExists(req.Code)
	if err != nil {
		return nil, err
	}
	return &domain.CheckExistsResponse{
		Exist: exists,
	}, nil
}
