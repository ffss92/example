package players

import "github.com/ffss92/example/internal/data"

type Storer interface {
	InsertPlayer(*data.Player) error
}

type Service struct {
	storer Storer
}

// Inject the storter interface to the Service
func NewService(storer Storer) Service {
	return Service{
		storer: storer,
	}
}

type CreatePlayer struct {
	Name string `json:"name"`
}

func (s Service) Create(input CreatePlayer) (*data.Player, error) {
	// Some business login
	player := &data.Player{
		Name: input.Name,
	}

	if err := s.storer.InsertPlayer(player); err != nil {
		return nil, err
	}

	return player, nil
}
