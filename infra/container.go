package infra

import (
	"database/sql"

	"github.com/alextanhongpin/url-shortener/infra/database"
)

type Container struct {
	db *sql.DB
}

func (c *Container) Close() {
	c.db.Close()
}

func (c *Container) DB() *sql.DB {
	return c.db
}

func NewContainer() *Container {
	return &Container{
		db: database.MustConn(database.NewConfig()),
	}
}
