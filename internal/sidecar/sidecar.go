// Package sidecar for collecting metrics from FPM.
package sidecar

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Server for collecting and returning
type Server struct {
	// Used for logging events.
	logger *slog.Logger
	// The most recently collected FPM status.
	status fpm.Status
	// Configuration used by the HTTP server.
	config ServerConfig
}

// ServerConfig which is used by the HTTP server.
type ServerConfig struct {
	// Port that the server will service traffic on.
	Port string
	// Path which will return metrics responses for the metrics adapter.
	Path string
	// Endpoint for querying the latest FPM status information.
	Endpoint string
	// EndpointPoll frequency for collecting FPM status information.
	EndpointPoll time.Duration
	// LogFrequency for how often a status is logged for external systems.
	LogFrequency time.Duration
}

// NewServer for collecting and responding with the latest FPM status.
func NewServer(logger *slog.Logger, config ServerConfig) (*Server, error) {
	server := &Server{
		logger: logger,
		config: config,
	}

	return server, nil
}

// Run the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	group, ctx := errgroup.WithContext(ctx)

	router := http.NewServeMux()

	router.HandleFunc(s.config.Path, s.handler)

	srv := &http.Server{
		Addr:    s.config.Port,
		Handler: router,
	}

	// Run the HTTP server.
	group.Go(func() error {
		// We want to shutdown all other tasks if this logger exits.
		defer cancel()

		s.logger.Info("Starting server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	// Run the HTTP server.
	group.Go(func() error {
		<-ctx.Done()
		s.logger.Info("Shutting down server")
		return srv.Shutdown(ctx)
	})

	// Query FPM periodically for statistics.
	group.Go(func() error {
		// We want to shutdown all other tasks if this logger exits.
		defer cancel()

		s.logger.Info("Starting refresher")
		return s.refreshStatus(ctx)
	})

	// Logger for emmit metrics as a log event for external systems.
	group.Go(func() error {
		// We want to shutdown all other tasks if this logger exits.
		defer cancel()

		s.logger.Info("Starting logger")
		return s.logStatus(ctx)
	})

	return group.Wait()
}
