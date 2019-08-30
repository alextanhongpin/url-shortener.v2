package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
)

// Stmt represents a unique id for the stmt.
type Stmt uint

// RawStmts holds the unprepared statements.
type RawStmts map[Stmt]string

// Stmts holds the prepared statements.c
type Stmts map[Stmt]*sql.Stmt

type Config struct {
	Username string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASS"`
	Database string `envconfig:"DB_NAME"`
	Host     string `envconfig:"DB_HOST"`
	Port     string `envconfig:"DB_PORT"`
}

func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}

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
		err := db.Ping()
		if err != nil {
			log.Println("retrying db connection in 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}
	log.Println("[db]: migration started")
	err = runMigration(db)
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	log.Println("[db]: migration completed with:", err)
	return db
}

// Prepare takes the raw statements and returns the prepared statements.
func Prepare(db *sql.DB, rawStmts RawStmts) Stmts {
	stmts := make(Stmts)
	var err error
	for key, value := range rawStmts {
		stmts[key], err = db.Prepare(value)
		if err != nil {
			panic(err)
		}
	}
	return stmts
}
