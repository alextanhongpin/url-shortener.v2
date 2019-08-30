package shorten

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
	repo      domain.Repository
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

	log.Println(req.Code)
	longURL, err := s.config.repo.GetByCode(req.Code)
	log.Println("got longurl", longURL)
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

	// To ensure that it is really unique, even for the same url,
	// use the unique user id to hash the url.
	e := domain.ShortURL(req)
	if e.Code != "" {
		exists, err := s.config.repo.CheckExists(e.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrAlreadyExists
		}
	} else {
		e.Code = s.config.shortener.Shorten(req.LongURL)
	}

	_, err := s.config.repo.Create(e)
	for err == ErrAlreadyExists {
		e.Code = s.config.shortener.Shorten(e.Code)
		_, err = s.config.repo.Create(e)
		log.Println("exists", err)
	}
	if err != nil {
		return nil, err
	}
	return &domain.PutResponse{Code: e.Code}, nil
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
