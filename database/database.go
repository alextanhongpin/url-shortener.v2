package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
)

// MustConn initialize a connection with the database.
func MustConn(cfg Config) *sql.DB {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Host,
		cfg.Port,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		if err := db.Ping(); err != nil {
			log.Println("retrying db connection in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		panic(err)
	}

	if cfg.EnableDBMigration {
		log.Println("[db]: migration started")
		err = runMigration(db)
		if err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
		log.Println("[db]: migration completed with:", err)
	} else {
		log.Println("[db]: migration disabled")
	}

	return db
}
