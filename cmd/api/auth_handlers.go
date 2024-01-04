package main

import (
	"errors"
	"net/http"

	"github.com/ffss92/example/internal/auth"
)

// Registers a user.
func (a api) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var input auth.CreateUserParams
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}

	user, err := a.auth.CreateUser(input)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrDuplicateEmail):
			a.conflictError(w, r, err)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusCreated, user); err != nil {
		a.serverError(w, r, err)
	}
}

// Logs in a user.
func (a api) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var input auth.CredentialsParam
	if err := readJSON(r, &input); err != nil {
		a.clientError(w, r, err)
		return
	}

	token, err := a.auth.Authenticate(input)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			a.invalidCredsError(w, r)
		default:
			a.serverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, token); err != nil {
		a.serverError(w, r, err)
	}
}

// Returns the current authenticated user.
func (a api) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := a.mustGetUser(r.Context())

	if err := writeJSON(w, http.StatusOK, user); err != nil {
		a.serverError(w, r, err)
	}
}
