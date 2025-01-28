package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/christgf/env"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/metrics/legacyregistry"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/apiserver/metrics"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"

	customprovider "github.com/skpr/fpm-metrics-adapter/internal/provider"
)

var (
	// Name for identifying which adapter is providing metrics.
	adapterName = "skpr-fpm-metrics-adapter"

	cmdLong = `
		Run the metrics adapter with collects FPM status information.`

	cmdExample = `
		# Run the adapter with the defaults.
		skpr-fpm-metrics-adapter

		# Run the adapter with a longer cache expiration.
		export SKPR_FPM_METRICS_ADAPTER_CACHE_EXPIRATION=120s
		skpr-fpm-metrics-adapter`
)

// Adapter for custom metrics.
type Adapter struct {
	basecmd.AdapterBase
}

// Helper function to instantiate the custom metrics provider.
func (a *Adapter) getProvider(logger *slog.Logger, cacheExpiration time.Duration) (provider.CustomMetricsProvider, error) {
	config, err := a.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to construct client config: %w", err)
	}

	client, err := a.DynamicClient()
	if err != nil {
		return nil, fmt.Errorf("unable to construct dynamic client: %w", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		return nil, fmt.Errorf("unable to construct discovery REST mapper: %w", err)
	}

	return customprovider.New(logger, client, config, mapper, cacheExpiration), nil
}

// Options for this sidecar application.
type Options struct {
	CacheExpiration time.Duration
	LogLevel        string
}

func main() {
	o := Options{}

	cmd := &cobra.Command{
		Use:     adapterName,
		Short:   "Run the Kubernetes metrics adapter",
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

			logger.Info("Starting metrics adapter")

			adapter := &Adapter{}
			adapter.Name = adapterName

			logger.Info("Getting provider")

			provider, err := adapter.getProvider(logger, o.CacheExpiration)
			if err != nil {
				return fmt.Errorf("failed to get provider: %w", err)
			}

			logger.Info("Registering metrics")

			adapter.WithCustomMetrics(provider)

			if err := metrics.RegisterMetrics(legacyregistry.Register); err != nil {
				return fmt.Errorf("failed to register metrics: %w", err)
			}

			logger.Info("Running adapter")

			if err := adapter.Run(wait.NeverStop); err != nil {
				return fmt.Errorf("failed to run adapter: %w", err)
			}

			logger.Info("Metrics adapter finished")

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&o.LogLevel, "log-level", env.String("SKPR_FPM_METRICS_ADAPTER_LOG_LEVEL", "info"), "Set the logging level")
	cmd.PersistentFlags().DurationVar(&o.CacheExpiration, "cache-expiration", env.Duration("SKPR_FPM_METRICS_ADAPTER_CACHE_EXPIRATION", 10*time.Second), "How long to keep cached metrics")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
