package auth

type Storer interface {
	GetUser(int64) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserForToken([]byte, Scope) (*User, error)
	DeleteUser(int64) error

	DeleteTokensForUser(int64, Scope) error
	InsertToken(*Token) error
	InsertUser(*User) error
}

type Service struct {
	storer Storer
}

func NewService(storer Storer) Service {
	return Service{
		storer: storer,
	}
}
