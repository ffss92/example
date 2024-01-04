package auth_test

import (
	"log"
	"os"
	"testing"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/migrations"
)

var (
	seedUser  *auth.User
	seedToken *auth.Token
)

var testStore auth.Storer

func newTestService() auth.Service {
	return auth.NewService(testStore)
}

func TestMain(m *testing.M) {
	db, err := infra.ConnectSqlite(":memory:")
	if err != nil {
		log.Fatalf("failed to connect to in mem db: %s", err)
	}
	defer db.Close()

	if err := migrations.Up(db, "sqlite"); err != nil {
		log.Fatalf("failed to migrate the db: %s", err)
	}

	testStore = data.NewStore(db)
	service := newTestService()

	user, err := service.CreateUser(auth.CreateUserParams{
		Username: "test",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("failed to seed user: %s", err)
	}
	seedUser = user

	token, err := service.Authenticate(auth.CredentialsParam{
		Username: "test",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("failed to seed token: %s", err)
	}
	seedToken = token

	code := m.Run()
	os.Exit(code)
}
