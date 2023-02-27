package server

import (
	"context"
	"proj/config"
	"proj/internal/handler"

	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler *handler.Handler) chan error {
	serverErr := make(chan error)
	s.httpServer = &http.Server{
		Addr:    cfg.RunAddress,
		Handler: handler.Router,
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	return serverErr
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
