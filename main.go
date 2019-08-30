package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alextanhongpin/pkg/grace"
	"github.com/alextanhongpin/url-shortener/domain"
	"github.com/alextanhongpin/url-shortener/infra"
	"github.com/alextanhongpin/url-shortener/pkg/shorten"

	"github.com/go-chi/chi"
)

func main() {
	// Container is responsible for starting/stopping all infrastructure.
	ctn := infra.NewContainer()

	// Make the type explicit.
	var svc domain.URLService
	{
		cfg := shorten.DefaultConfig(ctn.DB())
		svc = shorten.NewService(cfg)
	}

	// Setup routes.
	r := chi.NewRouter()
	r.Mount("/", shorten.NewController(svc).Router())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})

	// Implements graceful shutdown.
	shutdown := grace.New(r, "8080")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		shutdown(ctx)
	}()

	// Listens to CTRL + C.
	<-grace.Signal()
}
