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
		Addr:         a.cfg.Addr(),
		Handler:      a.routes(),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
		ErrorLog:     slog.NewLogLogger(a.log.Handler(), slog.LevelError),
	}

	a.log.Info(fmt.Sprintf("listening on http://localhost:%d", a.cfg.Port))
	return srv.ListenAndServe()
}
