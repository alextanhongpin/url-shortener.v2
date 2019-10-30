package shorturlsvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alextanhongpin/url-shortener/domain/shorturl"
	"github.com/lib/pq"

	"github.com/leebenson/conform"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ErrUsedShortID = errors.New("short id is already used")
)

type UseCase struct {
	urls    shorturl.Repository
	service shorturl.Service
	before  func(interface{}) error
}

func NewUseCase(urls shorturl.Repository, service shorturl.Service, validator *validator.Validate) *UseCase {
	return &UseCase{
		urls:    urls,
		service: service,
		before: func(req interface{}) error {
			// Trim the strings with the conform tag.
			conform.Strings(req)

			// Validate requests.
			return validator.Struct(req)
		},
	}
}

func (u *UseCase) Get(ctx context.Context, req shorturl.GetRequest) (*shorturl.GetResponse, error) {
	if err := u.before(&req); err != nil {
		return nil, err
	}
	longURL, err := u.urls.WithCode(req.Code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, shorturl.ErrDoesNotExists
	}
	if err != nil {
		return nil, err
	}
	return shorturl.NewGetResponse(longURL), nil
}

func (u *UseCase) Put(ctx context.Context, req shorturl.PutRequest) (*shorturl.PutResponse, error) {
	if err := u.before(&req); err != nil {
		return nil, err
	}

	// To ensure that it is really unique, even for the same url,
	// use the unique user id to hash the url.

	// If user provides an custom short url...
	if req.Code == "" {
		req.Code = u.service.Shorten(req.URL)
	}

	surl := shorturl.New(req.Code, req.URL)
	id, err := u.urls.Create(surl)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == DuplicatePrimaryKeyViolation {
			return nil, ErrUsedShortID
		}
	}
	if err != nil {
		return nil, err
	}

	return shorturl.NewPutResponse(id), nil
}

func (u *UseCase) CheckExists(ctx context.Context, req shorturl.CheckExistsRequest) (*shorturl.CheckExistsResponse, error) {
	if err := u.before(&req); err != nil {
		return nil, err
	}
	exists, err := u.urls.CheckExists(req.Code)
	if err != nil {
		return nil, err
	}
	return shorturl.NewCheckExistsResponse(exists), nil
}
