package auth

type Storer interface {
	// Users
	InsertUser(*User) error
	GetUser(int64) (*User, error)
	GetUserByUsername(string) (*User, error)
	GetUserForToken([]byte, Scope) (*User, error)
	DeleteUser(int64) error

	// Tokens
	DeleteTokensForUser(int64, Scope) error
	InsertToken(*Token) error
}

type Service struct {
	storer Storer
}

func NewService(storer Storer) Service {
	return Service{
		storer: storer,
	}
}
