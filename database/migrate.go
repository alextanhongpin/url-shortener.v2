package database

import (
	"database/sql"

	"github.com/alextanhongpin/url-shortener/database/migrations"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
)

func runMigration(db *sql.DB) error {
	src, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		}))
	if err != nil {
		return err
	}
	tgt, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("go-bindata", src, "postgres", tgt)
	if err != nil {
		return err
	}
	return m.Up()
}
