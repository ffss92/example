package main

import "github.com/go-chi/chi/v5"

// Sets up the api routes and middlewares.
func (a api) routes() *chi.Mux {
	r := chi.NewMux()

	// Middlewares
	r.Use(a.recoverer)
	r.Use(a.authenticate)

	// Routes
	// Auth
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", a.handleSignUp)
		r.Post("/sign-in", a.handleSignIn)
		r.Get("/me", a.requireAuth(a.handleCurrentUser))
	})

	// Posts
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", a.requireAuth(a.handleCreatePost))
		r.Get("/", a.handleListPosts)
		r.Get("/{postId}", a.handleGetPost)
		r.Put("/{postId}", a.requireAuth(a.handleUpdatePost))
		r.Delete("/{postId}", a.requireAuth(a.handleDeletePost))
		r.Post("/{postId}/likes", a.requireAuth(a.handleLikePost))
		r.Delete("/{postId}/likes", a.requireAuth(a.handleDislikePost))
		r.Post("/{postId}/comments", a.requireAuth(a.handleCreateComment))
		r.Get("/{postId}/comments", a.handleListComments)
		r.Delete("/{postId}/comments/{commentId}", a.requireAuth(a.handleDeleteComment))
		r.Put("/{postId}/comments/{commentId}", a.requireAuth(a.handleUpdateComment))
	})

	return r
}
