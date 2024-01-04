package data

import (
	"database/sql"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/posts"
)

// We make sure at compile time that Store implements the correct interfaces without allocating.
var _ auth.Storer = (*Store)(nil)
var _ posts.Storer = (*Store)(nil)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
	}
}
