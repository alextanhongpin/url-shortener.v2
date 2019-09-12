package database

import "github.com/kelseyhightower/envconfig"

// Config represents the database config.
type Config struct {
	Username          string `envconfig:"DB_USER"`
	Password          string `envconfig:"DB_PASS"`
	Database          string `envconfig:"DB_NAME"`
	Host              string `envconfig:"DB_HOST"`
	Port              string `envconfig:"DB_PORT"`
	EnableDBMigration bool   `envconfig:"ENABLE_DB_MIGRATION" default:"false"`
}

// NewConfig returns a new database config from the environment variables.
func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
