package sidecar

import (
	"net/http"
	"time"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Handler wraps promhttp.Handler to fetch data.
func (s *Server) RefreshMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Flood control for requests to fpm.
		if time.Now().After(s.metrics.LastUpdate.Add(1 * time.Second)) {
			s.logger.Debug("collecting FPM status")

			s.metrics.LastUpdate = time.Now()

			status, err := fpm.QueryStatus(s.config.Endpoint)
			if err != nil {
				s.logger.Error("failed to collect FPM status", "error", err.Error())
			}

			s.metrics.ListenQueue.Set(float64(status.ListenQueue))
			s.metrics.ListenQueueLen.Set(float64(status.ListenQueueLen))
			s.metrics.IdleProcesses.Set(float64(status.IdleProcesses))
			s.metrics.ActiveProcesses.Set(float64(status.ActiveProcesses))
			s.metrics.TotalProcesses.Set(float64(status.TotalProcesses))
			s.metrics.MaxActiveProcesses.Set(float64(status.MaxActiveProcesses))
		}

		next.ServeHTTP(w, r)
	})
}
