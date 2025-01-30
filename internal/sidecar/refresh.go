package sidecar

import (
	"context"
	"time"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// The process which will continually refresh the current status.
func (s *Server) refreshStatus(ctx context.Context) error {
	ticker := time.NewTicker(s.config.EndpointPoll)

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
