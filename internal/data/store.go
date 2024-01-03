package data

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	// Ignore errors for brevity
	db.Exec(`CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY, name TEXT NOT NULL);`)

	return Store{
		db: db,
	}
}
