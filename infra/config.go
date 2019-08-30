package infra

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application config.
type Config struct {
	Version string `envconfig:"VERSION" default:"v1"`
	Port    int    `envconfig:"PORT" default:"8080"`
}

// NewConfig returns the application config.
func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
