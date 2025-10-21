// Package sidecar for collecting metrics from FPM.
package sidecar

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Server for collecting and returning
type Server struct {
	// Used for logging events.
	logger *slog.Logger
	// Configuration used by the HTTP server.
	config ServerConfig
	// Metrics for the server
	metrics Metrics
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

type Metrics struct {
	// The last time the FPM status was updated.
	LastUpdate time.Time
	// Prometheus metrics.
	ListenQueue        prometheus.Gauge
	ListenQueueLen     prometheus.Gauge
	IdleProcesses      prometheus.Gauge
	ActiveProcesses    prometheus.Gauge
	TotalProcesses     prometheus.Gauge
	MaxActiveProcesses prometheus.Gauge
}

// NewServer for collecting and responding with the latest FPM status.
func NewServer(logger *slog.Logger, config ServerConfig) (*Server, error) {
	server := &Server{
		logger: logger,
		config: config,
		metrics: Metrics{
			ListenQueue: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricListenQueue,
				Help: "The number of items in the listen queue.",
			}),
			ListenQueueLen: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricListenQueueLen,
				Help: "The total size of the listen queue.",
			}),
			IdleProcesses: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricIdleProcesses,
				Help: "The number of idle fpm processes.",
			}),
			ActiveProcesses: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricActiveProcesses,
				Help: "The number of active fpm processes.",
			}),
			TotalProcesses: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricTotalProcesses,
				Help: "The total number of processes available in fpm.",
			}),
			MaxActiveProcesses: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: fpm.MetricMaxActiveProcesses,
				Help: "The maximum number of active processes since the FPM master process was started.",
			}),
		},
	}

	return server, nil
}

// Run the HTTP server.
func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Registering metrics")

	var errs []error

	metrics := []prometheus.Collector{
		s.metrics.ListenQueue,
		s.metrics.ListenQueueLen,
		s.metrics.IdleProcesses,
		s.metrics.ActiveProcesses,
		s.metrics.TotalProcesses,
		s.metrics.MaxActiveProcesses,
	}

	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to register metrics: %w", errors.Join(errs...))
	}

	mux := http.NewServeMux()
	mux.Handle(s.config.Path, s.RefreshMetricsMiddleware(promhttp.Handler()))

	s.logger.Info("Starting server")

	return http.ListenAndServe(s.config.Port, mux)
}
