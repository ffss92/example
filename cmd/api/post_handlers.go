package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ffss92/example/internal/pagination"
	"github.com/ffss92/example/internal/posts"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (a api) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())

	var input posts.CreatePostParams
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}

	input.User = user
	post, err := a.posts.CreatePost(input)
	if err != nil {
		var ve validator.ValidationErrors
		switch {
		case errors.As(err, &ve):
			a.validationError(w, r, ve)
		case errors.Is(err, posts.ErrDuplicate):
			a.conflictError(w, r, err)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleListPosts(w http.ResponseWriter, r *http.Request) {
	pag := pagination.FromRequest(r)
	posts, err := a.posts.ListPosts(pag)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleGetPost(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	post, err := a.posts.GetPost(postId)
	if err != nil {
		switch {
		case errors.Is(err, posts.ErrDuplicate):
			a.notFoundError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	if err := a.posts.DeletePost(user, postId); err != nil {
		switch {
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		case errors.Is(err, posts.ErrNotAllowed):
			a.forbiddenError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a api) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	var input posts.UpdatePostParams
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}

	input.User = user
	input.PostID = postId
	if err := a.posts.UpdatePost(input); err != nil {
		switch {
		case errors.Is(err, posts.ErrNotAllowed):
			a.forbiddenError(w, r)
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		case errors.Is(err, posts.ErrDuplicate):
			a.conflictError(w, r, err)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a api) handleLikePost(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	if err := a.posts.LikePost(user.ID, postId); err != nil {
		switch {
		case errors.Is(err, posts.ErrAlreadyLiked):
			a.clientError(w, r, err)
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, m{"message": "post liked"}); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleDislikePost(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	if err := a.posts.DislikePost(user.ID, postId); err != nil {
		switch {
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a api) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	var input posts.CreateCommentParams
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}
	input.User = user
	input.PostID = postId

	comment, err := a.posts.CreateComment(input)
	if err != nil {
		var ve validator.ValidationErrors
		switch {
		case errors.As(err, &ve):
			a.validationError(w, r, ve)
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, comment); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleListComments(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	p := pagination.FromRequest(r)
	comments, err := a.posts.ListComments(postId, p)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, comments); err != nil {
		a.serverError(w, r, err)
	}
}

func (a api) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())

	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	commentId, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	err = a.posts.DeleteComment(posts.DeletePostParams{
		User:      user,
		PostID:    postId,
		CommentID: commentId,
	})
	if err != nil {
		switch {
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		case errors.Is(err, posts.ErrNotAllowed):
			a.forbiddenError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a api) handleUpdateComment(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())

	postId, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	commentId, err := strconv.ParseInt(chi.URLParam(r, "commentId"), 10, 64)
	if err != nil || postId < 1 {
		a.notFoundError(w, r)
		return
	}

	var input posts.UpdateCommentParams
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}

	input.CommentID = commentId
	input.PostID = postId
	input.User = user

	if err := a.posts.UpdateComment(input); err != nil {
		var ve validator.ValidationErrors
		switch {
		case errors.As(err, &ve):
			a.validationError(w, r, ve)
		case errors.Is(err, posts.ErrNotFound):
			a.notFoundError(w, r)
		default:
			a.serverError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
