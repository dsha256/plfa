package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	version = "1.0.0"

	idleTimeout  = 1 * time.Minute
	readTimeout  = 10 * time.Second
	writeTimeout = 30 * time.Second

	gracefulShutdownTimeout = 5 * time.Second
)

type Server struct {
	logger *jsonlog.Logger
	repo   repository.AggregateRepository
}

func NewServer(logger *jsonlog.Logger, repo repository.AggregateRepository) *Server {
	return &Server{logger: logger, repo: repo}
}

func (s *Server) Serve(port string) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.routes(),
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		qs := <-quit
		s.logger.PrintInfo("shutting down server", map[string]string{"signal": qs.String()})

		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
	}()

	s.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
	})

	// Calling Shutdown() on the server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	s.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
