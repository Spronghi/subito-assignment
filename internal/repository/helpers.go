package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewSQLiteDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
