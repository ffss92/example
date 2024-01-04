package posts

import "errors"

var (
	ErrNotFound     = errors.New("post not found")
	ErrDuplicate    = errors.New("post already exists")
	ErrNotAllowed   = errors.New("user is not allowed to perform this operation")
	ErrAlreadyLiked = errors.New("user already like this post")
)
