package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/config"
	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/migrations"
	_ "github.com/joho/godotenv/autoload"
)

// The main application struct.
type api struct {
	cfg    config.Config
	logger *slog.Logger
	auth   auth.Service
}

func main() {
	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Dump(os.Stdout)

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Database
	db, err := infra.ConnectSqlite(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Store
	store := data.NewStore(db)

	// Migrations
	if err := migrations.Up(db, "sqlite"); err != nil {
		log.Fatal(err)
	}

	// Services
	authService := auth.NewService(store)

	// API
	api := &api{
		cfg:    cfg,
		logger: logger,
		auth:   authService,
	}
	log.Fatal(api.serve())
}
