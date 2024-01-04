package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/migrations"
)

type config struct {
	port   int
	dbPath string
}

// The main application struct.
type api struct {
	cfg    config
	logger *slog.Logger
	auth   auth.Service
}

func main() {
	// Config
	var cfg config

	flag.IntVar(&cfg.port, "port", 5000, "sets the server port")
	flag.StringVar(&cfg.dbPath, "db", "example.db", "sets the sqlite database path")
	flag.Parse()

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Database
	db, err := infra.ConnectSqlite(cfg.dbPath)
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
