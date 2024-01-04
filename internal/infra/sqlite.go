package infra

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func ConnectSqlite(dbPath string) (*sql.DB, error) {
	pragmas := "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	db, err := sql.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}
