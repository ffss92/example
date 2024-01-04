package main

import "github.com/go-chi/chi/v5"

// Sets up the api routes and middlewares.
func (a api) routes() *chi.Mux {
	r := chi.NewMux()

	// Middlewares
	r.Use(a.recoverer)
	r.Use(a.authenticate)

	// Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", a.handleSignUp)
		r.Post("/sign-in", a.handleSignIn)
		r.Get("/me", a.requireAuth(a.handleCurrentUser))
	})

	return r
}
