// Package sidecar for collecting metrics from FPM.
package sidecar

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
	"log/slog"
	"net/http"
	"time"
)

// Server for collecting and returning
type Server struct {
	// Used for logging events.
	logger *slog.Logger
	// Configuration used by the HTTP server.
	config ServerConfig
	// LastUpdate timestamp
	lastUpdate time.Time
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

var (
	listenQueue = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricListenQueue,
		Help: "The number of items in the listen queue.",
	})
	listenQueueLen = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricListenQueueLen,
		Help: "The total size of the listen queue.",
	})
	idleProcesses = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricIdleProcesses,
		Help: "The number of idle fpm processes.",
	})
	activeProcesses = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricActiveProcesses,
		Help: "The number of active fpm processes.",
	})
	totalProcesses = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricTotalProcesses,
		Help: "The total number of processes available in fpm.",
	})
	maxActiveProcesses = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: fpm.MetricMaxActiveProcesses,
		Help: "The maximum number of active processes since the FPM master process was started.",
	})
)

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
	prometheus.MustRegister(listenQueue)
	prometheus.MustRegister(listenQueueLen)
	prometheus.MustRegister(idleProcesses)
	prometheus.MustRegister(activeProcesses)
	prometheus.MustRegister(totalProcesses)
	prometheus.MustRegister(maxActiveProcesses)

	s.logger.Info("Starting server")

	http.Handle(s.config.Path, s.Handler())
	err := http.ListenAndServe(s.config.Port, nil)
	if err != nil {
		return err
	}

	return nil
}
