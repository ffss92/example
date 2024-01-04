package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/config"
	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/internal/posts"
	"github.com/ffss92/example/internal/validate"
	"github.com/ffss92/example/migrations"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	_ "github.com/joho/godotenv/autoload"
)

// The main application struct.
type api struct {
	cfg   config.Config
	log   *slog.Logger
	uni   *ut.UniversalTranslator
	auth  auth.Service
	posts posts.Service
}

func main() {
	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Dump(os.Stdout)

	// Database
	db, err := infra.ConnectSqlite(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := data.NewStore(db)

	// Migrations
	if err := migrations.Up(db, "sqlite"); err != nil {
		log.Fatal(err)
	}

	// Translations (user friendly validation errors)
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate.Validator(), trans)

	// Services
	auth := auth.NewService(store)
	posts := posts.NewService(store)

	// API
	api := &api{
		cfg:   cfg,
		log:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		uni:   uni,
		auth:  auth,
		posts: posts,
	}
	log.Fatal(api.serve())
}
