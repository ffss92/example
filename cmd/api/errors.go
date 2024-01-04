package main

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	authScheme = "Bearer"
)

type errorResponse struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Writes a standard error message to the client.
func (a api) writeError(w http.ResponseWriter, r *http.Request, status int, res errorResponse) {
	if err := writeJSON(w, status, res); err != nil {
		a.logger.Error("failed to write error response", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Logs the error and writes a 500 error response to the client.
func (a api) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error("unexpected error in handler", slog.String("err", err.Error()))
	a.writeError(w, r, http.StatusInternalServerError, errorResponse{
		Message: "internal server error",
	})
}

// Writes a 400 error response to the client.
func (a api) clientError(w http.ResponseWriter, r *http.Request, err error) {
	a.writeError(w, r, http.StatusBadRequest, errorResponse{
		Message: err.Error(),
	})
}

// Writes a 401 response for an invalid auth token.
func (a api) invalidTokenError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", authScheme)
	a.writeError(w, r, http.StatusUnauthorized, errorResponse{
		Message: "token is expired or invalid",
	})
}

// Writes a 401 response for an unauthed user.
func (a api) authRequiredError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", authScheme)
	a.writeError(w, r, http.StatusUnauthorized, errorResponse{
		Message: "you must be authenticated to access this resource",
	})
}

// Writes a 409 response for a conflict error.
func (a api) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	a.writeError(w, r, http.StatusConflict, errorResponse{
		Message: err.Error(),
	})
}

// Writes a 401 response to the client for invalid credentials.
func (a api) invalidCredsError(w http.ResponseWriter, r *http.Request) {
	a.writeError(w, r, http.StatusUnauthorized, errorResponse{
		Message: "invalid user credentials",
	})
}

// Writes a 422 response to the client with the validation details.
func (a api) validationError(w http.ResponseWriter, r *http.Request, ve validator.ValidationErrors) {
	details := make(map[string]string, len(ve))

	for _, e := range ve {
		trans, _ := a.uni.GetTranslator("en")
		details[e.Field()] = e.Translate(trans)
	}

	a.writeError(w, r, http.StatusUnprocessableEntity, errorResponse{
		Message: "validation failed",
		Details: details,
	})
}
