package main

import (
	"context"
	"time"

	"github.com/alextanhongpin/pkg/grace"

	"github.com/alextanhongpin/url-shortener/api/middleware"
	"github.com/alextanhongpin/url-shortener/app"
	"github.com/alextanhongpin/url-shortener/pkg/healthsvc"
	"github.com/alextanhongpin/url-shortener/pkg/shorturlsvc"

	"github.com/go-chi/chi"
)

func main() {
	cfg := app.NewConfig()

	// Container is responsible for starting/stopping all infrastructure.
	ctn := app.NewContainer()
	defer ctn.Close()

	// Make the type explicit.
	// Setup routes.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	// UseCase: Serve health endpoint.
	{
		ctl := healthsvc.NewController(cfg.Version)
		r.Mount("/health", ctl.Router())
	}

	// UseCase: Serve shortURL service endpoint.
	// We can also make this as a foogle (feature-toggle):
	// if (enableRoute)
	{
		repo := shorturlsvc.NewRepository(ctn.DB())
		svc := shorturlsvc.NewService()
		usecase := shorturlsvc.NewUseCase(repo, svc, ctn.Validator())
		ctl := shorturlsvc.NewController(usecase)
		r.Mount("/v1", ctl.Router())
	}

	// Implements graceful shutdown.
	shutdown := grace.New(r, cfg.Port)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		shutdown(ctx)
	}()

	// Listens to CTRL + C.
	<-grace.Signal()
}
