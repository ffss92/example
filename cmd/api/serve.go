package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Starts listening
func (a api) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.port),
		Handler:      a.routes(),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		ErrorLog:     slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
	}

	a.logger.Info(fmt.Sprintf("listening on http://localhost:%d", a.cfg.port))
	return srv.ListenAndServe()
}
