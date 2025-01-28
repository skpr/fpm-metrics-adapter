// Package sidecar for collecting metrics from FPM.
package sidecar

import (
	"context"
	"encoding/json"
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
	// Frequency for collecting FPM status information.
	Frequency time.Duration
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
		// We want to cancel the refresher if the server exits.
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
		// We want to shutdown the server if the refresh completes.
		defer cancel()

		s.logger.Info("Starting refresher")
		return s.refresh(ctx)
	})

	return group.Wait()
}

// Handler to return the latest status.
func (s *Server) handler(w http.ResponseWriter, _ *http.Request) {
	s.logger.Debug("Handling request")

	jsonBytes, err := json.Marshal(s.status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Error("failed to marshal status response", "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

	s.logger.Debug("Request complete")
}

// The process which will continually refresh the current status.
func (s *Server) refresh(ctx context.Context) error {
	ticker := time.NewTicker(s.config.Frequency)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.logger.Debug("Querying for FPM status")

			status, err := fpm.QueryStatus(s.config.Endpoint)
			if err != nil {
				s.logger.Error("failed to collect FPM status", "error", err.Error())
				continue
			}

			s.logger.Debug("Successfully queries latest FPM status", "response", status)

			s.status = status
		}
	}
}
