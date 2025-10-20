package sidecar

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Handler wraps promhttp.Handler to fetch data.
func (s *Server) Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		// Flood control for requests to fpm.
		if time.Now().After(s.lastUpdate.Add(1 * time.Second)) {
			s.logger.Debug("collecting FPM status")

			s.lastUpdate = time.Now()

			status, err := fpm.QueryStatus(s.config.Endpoint)
			if err != nil {
				s.logger.Error("failed to collect FPM status", "error", err.Error())
			}

			listenQueue.Set(float64(status.ListenQueue))
			listenQueueLen.Set(float64(status.ListenQueueLen))
			idleProcesses.Set(float64(status.IdleProcesses))
			activeProcesses.Set(float64(status.ActiveProcesses))
			totalProcesses.Set(float64(status.TotalProcesses))
			maxActiveProcesses.Set(float64(status.MaxActiveProcesses))
		}

		handler := promhttp.Handler()
		handler.ServeHTTP(writer, req)
	}
}
