package infra

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func ConnectSqlite(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "example.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
