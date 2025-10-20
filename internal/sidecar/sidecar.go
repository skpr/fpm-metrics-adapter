// Package sidecar for collecting metrics from FPM.
package sidecar

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// Server for collecting and returning
type Server struct {
	// Used for logging events.
	logger *slog.Logger
	// Configuration used by the HTTP server.
	config ServerConfig
	// Cache for flood control.
	cache *cache.Cache
}

// ServerConfig which is used by the HTTP server.
type ServerConfig struct {
	// Port that the server will service traffic on.
	Port string
	// Path which will return metrics responses for the metrics adapter.
	Path string
	// Endpoint for querying the latest FPM status information.
	Endpoint string
}

// NewServer for collecting and responding with the latest FPM status.
func NewServer(logger *slog.Logger, config ServerConfig) (*Server, error) {
	server := &Server{
		logger: logger,
		config: config,
		cache:  cache.New(1*time.Second, 10*time.Second),
	}

	return server, nil
}

// Run the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	router := http.NewServeMux()

	router.HandleFunc(s.config.Path, s.handler)

	srv := &http.Server{
		Addr:    s.config.Port,
		Handler: router,
	}

	s.logger.Info("Starting server")

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	s.logger.Info("Shutting down server")

	return nil
}
