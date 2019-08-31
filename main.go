package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/url-shortener/api/middleware"
	"github.com/alextanhongpin/url-shortener/infra"
	"github.com/alextanhongpin/url-shortener/pkg/health"
	"github.com/alextanhongpin/url-shortener/pkg/shorten"

	"github.com/go-chi/chi"
)

func main() {
	cfg := infra.NewConfig()

	// Container is responsible for starting/stopping all infrastructure.
	ctn := infra.NewContainer()

	// Make the type explicit.
	// Setup routes.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	// UseCase: Serve health endpoint.
	{
		ctl := health.NewController(cfg.Version)
		r.Mount("/health", ctl.Router())
	}

	// UseCase: Serve shortURL service endpoint.
	// We can also make this as a foogle (feature-toggle):
	// if (enableRoute)
	{
		repo := shorten.NewRepository(ctn.DB())
		svc := shorten.NewService(repo, shorten.NewShortener())
		ctl := shorten.NewController(svc)
		r.Mount("/v1", ctl.Router())
	}

	// Implements graceful shutdown.
	shutdown := grace.New(r, fmt.Sprint(cfg.Port))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		shutdown(ctx)
	}()

	// Listens to CTRL + C.
	<-grace.Signal()
}
