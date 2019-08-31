package shorten

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/alextanhongpin/url-shortener/domain"

	"github.com/leebenson/conform"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ErrAlreadyExists = errors.New("short url already exists")
	ErrDoesNotExists = errors.New("short url does not exist")
)

type Service struct {
	repo      domain.Repository
	shortener domain.Shortener
}

func NewService(repo domain.Repository, shortener domain.Shortener) *Service {
	return &Service{
		repo:      repo,
		shortener: shortener,
	}
}

func (s *Service) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	conform.Strings(&req)
	if err := domain.Validator.Struct(&req); err != nil {
		return nil, err
	}

	longURL, err := s.repo.GetByCode(req.Code)
	if err == sql.ErrNoRows {
		return nil, ErrDoesNotExists
	}
	if err != nil {
		return nil, err
	}
	return &domain.GetResponse{
		LongURL: longURL,
	}, nil
}

func (s *Service) Put(ctx context.Context, req domain.PutRequest) (*domain.PutResponse, error) {
	conform.Strings(&req)
	if err := domain.Validator.Struct(&req); err != nil {
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
			return nil, ErrAlreadyExists
		}
	} else {
		e.Code = s.shortener.Shorten(req.LongURL)
	}

	_, err := s.repo.Create(e)
	for err == ErrAlreadyExists {
		e.Code = s.shortener.Shorten(e.Code + fmt.Sprint(rand.Int()))
		_, err = s.repo.Create(e)
	}
	if err != nil {
		return nil, err
	}
	return &domain.PutResponse{Code: e.Code}, nil
}

func (s *Service) CheckExists(ctx context.Context, req domain.CheckExistsRequest) (*domain.CheckExistsResponse, error) {
	conform.Strings(&req)
	if err := domain.Validator.Struct(&req); err != nil {
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
