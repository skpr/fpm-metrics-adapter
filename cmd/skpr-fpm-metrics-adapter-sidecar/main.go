// Package skpr-fpm-metrics-adapter-sidecar for the sidecar entrypoint.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/christgf/env"
	"github.com/spf13/cobra"

	"github.com/skpr/fpm-metrics-adapter/internal/sidecar"
)

var (
	cmdLong = `
		Run the sidecar which collects profiles and prints them to stdout.`

	cmdExample = `
		# Run the sidecar with the defaults.
		skpr-metrics-adapter-sidecar

		# Enable debug logs.
		export SKPR_FPM_METRICS_ADAPTER_LOG_LEVEL=debug
		skpr-metrics-adapter-sidecar`
)

// Options for this sidecar application.
type Options struct {
	ServerConfig sidecar.ServerConfig
	LogLevel     string
	StatusLogger sidecar.LogStatus
}

func main() {
	o := Options{}

	cmd := &cobra.Command{
		Use:     "skpr-metrics-adapter-sidecar",
		Short:   "Run the FPM metrics adapter sidecar",
		Long:    cmdLong,
		Example: cmdExample,
		RunE: func(cmd *cobra.Command, _ []string) error {
			lvl := new(slog.LevelVar)

			switch o.LogLevel {
			case "info":
				lvl.Set(slog.LevelInfo)
			case "debug":
				lvl.Set(slog.LevelDebug)
			case "warn":
				lvl.Set(slog.LevelWarn)
			case "error":
				lvl.Set(slog.LevelError)
			default:
				lvl.Set(slog.LevelInfo)
			}

			logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: lvl,
			}))

			logger.Info("Booting sidecar")

			server, err := sidecar.NewServer(logger, o.ServerConfig, o.StatusLogger)
			if err != nil {
				return fmt.Errorf("failed to start server: %w", err)
			}

			err = server.Run(cmd.Context())
			if err != nil {
				return fmt.Errorf("server failed: %w", err)
			}

			logger.Info("Sidecar finished")

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&o.LogLevel, "log-level", env.String("SKPR_FPM_METRICS_ADAPTER_LOG_LEVEL", "info"), "Set the logging level")
	cmd.PersistentFlags().StringVar(&o.ServerConfig.Port, "port", env.String("SKPR_FPM_METRICS_ADAPTER_PORT", ":80"), "Port which our metrics endpoint will be served on")
	cmd.PersistentFlags().StringVar(&o.ServerConfig.Path, "path", env.String("SKPR_FPM_METRICS_ADAPTER_PATH", "/metrics"), "Path which our metrics endpoint will be served on")
	cmd.PersistentFlags().StringVar(&o.ServerConfig.Endpoint, "endpoint", env.String("SKPR_FPM_METRICS_ADAPTER_ENDPOINT", "127.0.0.1:9000"), "Endpoint which we will poll for FPM status information")
	cmd.PersistentFlags().DurationVar(&o.ServerConfig.EndpointPoll, "endpoint-poll", env.Duration("SKPR_FPM_METRICS_ADAPTER_ENDPOINT_POLL", 10*time.Second), "How frequently to poll the endpoint for status information")
	cmd.PersistentFlags().BoolVar(&o.StatusLogger.Enabled, "log-enabled", env.Bool("SKPR_FPM_METRICS_ADAPTER_LOG_STATUS_ENABLED", false), "If the status logger is enabled")
	cmd.PersistentFlags().DurationVar(&o.StatusLogger.Frequency, "log-frequency", env.Duration("SKPR_FPM_METRICS_ADAPTER_LOG_STATUS_FREQUENCY", 30*time.Second), "How frequently to log status information for external systems")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
