package main

import (
	"log"

	"github.com/ffss92/example/internal/data"
	"github.com/ffss92/example/internal/infra"
	"github.com/ffss92/example/internal/players"
)

type fakePlayerDb struct{}

// Its implicit, this works
func (f fakePlayerDb) InsertPlayer(player *data.Player) error {
	return nil
}

func main() {
	db, err := infra.ConnectSqlite("example.db")
	if err != nil {
		log.Fatal(err)
	}

	// Actual implementation
	store := data.NewStore(db)
	playerService := players.NewService(store)

	// Fake implementation for testing (Also works)
	// store := fakePlayerDb{}
	// playerService := players.NewService(store)

	playerService.Create(players.CreatePlayer{})
}
