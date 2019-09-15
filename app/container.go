package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/alextanhongpin/url-shortener/database"
	"gopkg.in/go-playground/validator.v9"

	"go.uber.org/zap"
)

type Container struct {
	db        *sql.DB
	log       *zap.Logger
	validator *validator.Validate
}

func (c *Container) Close() {
	if err := c.db.Close(); err != nil {
		log.Println(err)
	}
	if err := c.log.Sync(); err != nil {
		log.Println(err)
	}
}

func (c *Container) DB() *sql.DB {
	return c.db
}

func (c *Container) Validator() *validator.Validate {
	return c.validator
}

func NewContainer() *Container {
	return &Container{
		log:       initLogger(),
		db:        database.MustConn(database.NewConfig()),
		validator: validator.New(),
	}
}

func initLogger() *zap.Logger {
	log, _ := zap.NewProduction()
	host, _ := os.Hostname()
	log = log.With(zap.String("host", host))
	zap.ReplaceGlobals(log)
	return log
}
