package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationsFS embed.FS

// Applies migrations to the database.
func Up(db *sql.DB, dialect string) error {
	goose.SetBaseFS(migrationsFS)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}
	return goose.Up(db, ".")

}
