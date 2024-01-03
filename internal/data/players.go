package data

type Player struct {
	ID   int64
	Name string
}

func (s Store) InsertPlayer(p *Player) error {
	query := `INSERT INTO players (name) VALUES ($1) RETURNING id`

	err := s.db.QueryRow(query, p.Name).Scan(
		&p.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
