package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alextanhongpin/pkg/grace"

	"github.com/alextanhongpin/url-shortener/domain"
	"github.com/alextanhongpin/url-shortener/infra"
	"github.com/alextanhongpin/url-shortener/pkg/shorten"

	"github.com/go-chi/chi"
)

var DeployedAt = time.Now()

func main() {
	cfg := infra.NewConfig()

	// Container is responsible for starting/stopping all infrastructure.
	ctn := infra.NewContainer()

	// Make the type explicit.
	var svc domain.Service
	{
		cfg := shorten.DefaultConfig(ctn.DB())
		svc = shorten.NewService(cfg)
	}

	// Setup routes.
	r := chi.NewRouter()
	r.Mount("/", shorten.NewController(svc).Router())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		type Health struct {
			DeployedAt time.Time `json:"deployed_at"`
			Uptime     string    `json:"uptime"`
			Version    string    `json:"version"`
		}
		json.NewEncoder(w).Encode(Health{
			DeployedAt: DeployedAt,
			Uptime:     time.Since(DeployedAt).String(),
			Version:    cfg.Version,
		})
	})

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
