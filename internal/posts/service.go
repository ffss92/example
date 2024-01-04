package posts

type Storer interface {
	InsertPost(*Post) error
	GetPost(int64) (*Post, error)
	ListPosts(limit, offset int) ([]*Post, error)
	DeletePost(int64) error
	UpdatePost(*Post) error
	CountPosts() (int, error)

	InsertPostLike(userId int64, postId int64) error
	DeletePostLike(userId int64, postId int64) error
}

type Service struct {
	storer Storer
}

func NewService(storer Storer) Service {
	return Service{
		storer: storer,
	}
}
