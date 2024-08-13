package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dsn)
	return db, err
}
