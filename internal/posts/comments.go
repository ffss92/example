package posts

import (
	"time"

	"github.com/ffss92/example/internal/auth"
	"github.com/ffss92/example/internal/pagination"
	"github.com/ffss92/example/internal/validate"
)

type Comment struct {
	ID        int64     `json:"id"`
	Comment   string    `json:"comment"`
	Author    string    `json:"author"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentParams struct {
	User    *auth.User `json:"-" validate:"required"`
	PostID  int64      `json:"-" validate:"required"`
	Comment string     `json:"comment" validate:"required,max=60"`
}

func (s Service) CreateComment(params CreateCommentParams) (*Comment, error) {
	if err := validate.Struct(params); err != nil {
		return nil, err
	}

	post, err := s.storer.GetPost(params.PostID)
	if err != nil {
		return nil, err
	}

	comment := &Comment{
		Comment: params.Comment,
		PostID:  post.ID,
		Author:  params.User.Username,
		UserID:  params.User.ID,
	}
	if err := s.storer.InsertPostComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s Service) ListComments(postId int64, p pagination.Pagination) (*pagination.Paginated[*Comment], error) {
	count, err := s.storer.CountPostComments(postId)
	if err != nil {
		return nil, err
	}

	comments, err := s.storer.ListPostComments(postId, p.Limit(), p.Offset())
	if err != nil {
		return nil, err
	}

	paginated := pagination.NewPaginated(p.Page(), p.Limit(), count, comments)
	return &paginated, nil
}

type DeletePostParams struct {
	User      *auth.User
	PostID    int64
	CommentID int64
}

// Post owners can delete comments.
func (s Service) DeleteComment(params DeletePostParams) error {
	post, err := s.storer.GetPost(params.PostID)
	if err != nil {
		return err
	}

	comment, err := s.storer.GetComment(params.PostID, params.CommentID)
	if err != nil {
		return err
	}

	isOwner := post.UserID == params.User.ID
	isAuthor := comment.UserID == params.User.ID
	if !isAuthor || !isOwner {
		return ErrNotAllowed
	}

	if err := s.storer.DeletePostComment(comment.ID); err != nil {
		return err
	}

	return nil
}

type UpdateCommentParams struct {
	User      *auth.User `json:"-" validate:"required"`
	PostID    int64      `json:"-" validate:"required"`
	CommentID int64      `json:"-" validate:"required"`
	Comment   string     `json:"comment" validate:"required,max=60"`
}

func (s Service) UpdateComment(params UpdateCommentParams) error {
	if err := validate.Struct(params); err != nil {
		return err
	}

	comment, err := s.storer.GetComment(params.PostID, params.CommentID)
	if err != nil {
		return err
	}

	isOwner := params.User.ID == comment.UserID
	if !isOwner {
		return ErrNotAllowed
	}

	comment.Comment = params.Comment
	if err := s.storer.UpdatePostComment(comment); err != nil {
		return err
	}
	return nil
}
