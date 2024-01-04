package posts

import (
	"time"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/pagination"
	"github.com/ffss92/example/internal/validate"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Likes     int64     `json:"likes"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID int64 `json:"-"`
}

type CreatePostParams struct {
	Title   string `json:"title" validate:"required,min=2,max=50"`
	Content string `json:"content" validate:"required,max=250"`
}

func (s Service) CreatePost(user *auth.User, params CreatePostParams) (*Post, error) {
	if err := validate.Struct(params); err != nil {
		return nil, err
	}

	post := &Post{
		Title:   params.Title,
		Content: params.Content,
		UserID:  user.ID,
		Author:  user.Username,
		Likes:   0,
	}

	if err := s.storer.InsertPost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s Service) GetPost(id int64) (*Post, error) {
	post, err := s.storer.GetPost(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s Service) ListPosts(p pagination.Pagination) (*pagination.Paginated[*Post], error) {
	count, err := s.storer.CountPosts()
	if err != nil {
		return nil, err
	}

	posts, err := s.storer.ListPosts(p.Limit(), p.Offset())
	if err != nil {
		return nil, err
	}

	paginated := pagination.NewPaginated(p.Page(), p.Limit(), count, posts)
	return &paginated, nil
}

func (s Service) DeletePost(user *auth.User, postId int64) error {
	post, err := s.storer.GetPost(postId)
	if err != nil {
		return err
	}

	if post.UserID != user.ID {
		return ErrNotAllowed
	}

	if err := s.storer.DeletePost(postId); err != nil {
		return err
	}

	return nil
}

type UpdatePostParams struct {
	Title   string `json:"title" validate:"required,min=2,max=50"`
	Content string `json:"content" validate:"required,max=250"`
}

func (s Service) UpdatePost(user *auth.User, postId int64, params UpdatePostParams) error {
	if err := validate.Struct(params); err != nil {
		return err
	}

	post, err := s.storer.GetPost(postId)
	if err != nil {
		return err
	}
	if post.UserID != user.ID {
		return ErrNotAllowed
	}

	post.Title = params.Title
	post.Content = params.Content
	if err := s.storer.UpdatePost(post); err != nil {
		return err
	}

	return nil
}
