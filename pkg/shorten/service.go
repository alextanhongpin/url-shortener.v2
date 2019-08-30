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
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ErrAlreadyExists = errors.New("short url already exists")
	ErrDoesNotExists = errors.New("short url does not exist")
)

// Config represents the cfguration of the Service.
type Config struct {
	validator *validator.Validate
	shortener domain.Shortener
	repo      domain.Repository
}

func DefaultConfig(db *sql.DB) Config {
	return Config{
		validator: validator.New(),
		shortener: NewShortener(),
		repo:      NewRepository(db),
	}
}

type Service struct {
	cfg Config
}

func NewService(cfg Config) *Service {
	if cfg.validator == nil {
		cfg.validator = validator.New()
	}
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) Get(ctx context.Context, req domain.GetRequest) (*domain.GetResponse, error) {
	conform.Strings(&req)
	if err := s.cfg.validator.Struct(&req); err != nil {
		return nil, err
	}

	longURL, err := s.cfg.repo.GetByCode(req.Code)
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
	if err := s.cfg.validator.Struct(&req); err != nil {
		return nil, err
	}

	// To ensure that it is really unique, even for the same url,
	// use the unique user id to hash the url.
	e := domain.ShortURL(req)
	if e.Code != "" {
		exists, err := s.cfg.repo.CheckExists(e.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrAlreadyExists
		}
	} else {
		e.Code = s.cfg.shortener.Shorten(req.LongURL)
	}

	_, err := s.cfg.repo.Create(e)
	for err == ErrAlreadyExists {
		e.Code = s.cfg.shortener.Shorten(e.Code + fmt.Sprint(rand.Int()))
		_, err = s.cfg.repo.Create(e)
	}
	if err != nil {
		return nil, err
	}
	return &domain.PutResponse{Code: e.Code}, nil
}

func (s *Service) CheckExists(ctx context.Context, req domain.CheckExistsRequest) (*domain.CheckExistsResponse, error) {
	conform.Strings(&req)
	if err := s.cfg.validator.Struct(&req); err != nil {
		return nil, err
	}

	exists, err := s.cfg.repo.CheckExists(req.Code)
	if err != nil {
		return nil, err
	}
	return &domain.CheckExistsResponse{
		Exist: exists,
	}, nil
}
