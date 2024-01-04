package main

import (
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/config"
	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/migrations"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	_ "github.com/joho/godotenv/autoload"
)

// The main application struct.
type api struct {
	cfg      config.Config
	logger   *slog.Logger
	auth     auth.Service
	validate *validator.Validate
	uni      *ut.UniversalTranslator
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

	// TODO: Add this to an internal package.
	// Validator + Translations
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Services
	authService := auth.NewService(store, validate)

	// API
	api := &api{
		cfg:      cfg,
		logger:   logger,
		auth:     authService,
		validate: validate,
		uni:      uni,
	}
	log.Fatal(api.serve())
}
