package data

import (
	"database/sql"

	"github.com/ffss92/example/internal/auth"
)

// We make sure at compile time that Store implements users.Storer without allocating.
var _ auth.Storer = (*Store)(nil)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
	}
}
