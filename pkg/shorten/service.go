package shorten

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alextanhongpin/url-shortener/domain"
	"github.com/leebenson/conform"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ErrAlreadyExists = errors.New("short url already exists")
	ErrDoesNotExists = errors.New("short url does not exist")
)

// Config represents the configuration of the Service.
type Config struct {
	validator *validator.Validate
	shortener domain.Shortener
	repo      domain.URLRepository
}

func DefaultConfig(db *sql.DB) Config {
	return Config{
		validator: validator.New(),
		shortener: New(),
		repo:      NewRepository(db),
	}
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	if config.validator == nil {
		config.validator = validator.New()
	}
	return &Service{
		config: config,
	}
}

func (s *Service) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	conform.Strings(&req)
	if err := s.config.validator.Struct(&req); err != nil {
		return nil, err
	}

	longURL, err := s.config.repo.GetByCode(req.Code)
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
	if err := s.config.validator.Struct(&req); err != nil {
		return nil, err
	}

	code := req.Code
	if code != "" {
		exists, err := s.config.repo.CheckExists(code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrAlreadyExists
		}
	} else {
		code = s.config.shortener.Shorten(req.LongURL)
	}

	_, err := s.config.repo.Create(code, req.LongURL, req.ExpireAt)
	for err == ErrAlreadyExists {
		code = s.config.shortener.Shorten(code)
		_, err = s.config.repo.Create(code, req.LongURL, req.ExpireAt)
	}
	if err != nil {
		return nil, err
	}
	return &domain.PutResponse{Code: code}, nil
}

func (s *Service) CheckExists(ctx context.Context, req domain.CheckExistsRequest) (*domain.CheckExistsResponse, error) {
	conform.Strings(&req)
	if err := s.config.validator.Struct(&req); err != nil {
		return nil, err
	}
	exists, err := s.config.repo.CheckExists(req.Code)
	if err != nil {
		return nil, err
	}
	return &domain.CheckExistsResponse{
		Exist: exists,
	}, nil
}
