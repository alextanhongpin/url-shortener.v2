package database

import (
	"database/sql"
	"fmt"
)

// Stmt represents a unique id for the stmt.
type Stmt uint

// RawStmts holds the unprepared statements.
type RawStmts map[Stmt]string

// Stmts holds the prepared statements.c
type Stmts map[Stmt]*sql.Stmt

// MustPrepare takes the raw statements and returns the prepared statements.
func (rawStmts RawStmts) MustPrepare(db *sql.DB) Stmts {
	stmts := make(Stmts)
	var err error
	for key, value := range rawStmts {
		stmts[key], err = db.Prepare(value)
		if err != nil {
			stmtErr := fmt.Errorf("prepareError: %w with %s", err, value)
			panic(stmtErr)
		}
	}
	return stmts
}
