package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/klog"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/apiserver/metrics"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"

	customprovider "github.com/skpr/fpm-metrics-adapter/internal/provider"
)

// Adapter for custom metrics.
type Adapter struct {
	basecmd.AdapterBase

	// Message is printed on successful startup
	Message string
}

// Helper function to instantiate the custom metrics provider.
func (a *Adapter) makeProviderOrDie() (provider.CustomMetricsProvider, error) {
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

	return customprovider.New(client, config, mapper), nil
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := &Adapter{}
	cmd.Name = "skpr-fpm-metrics-adapter"

	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	logs.AddFlags(cmd.Flags())
	if err := cmd.Flags().Parse(os.Args); err != nil {
		panic(err)
	}

	testProvider, err := cmd.makeProviderOrDie()
	if err != nil {
		panic(err)
	}

	cmd.WithCustomMetrics(testProvider)

	if err := metrics.RegisterMetrics(legacyregistry.Register); err != nil {
		panic(err)
	}

	klog.Infof(cmd.Message)

	go func() {
		// Open port for POSTing fake metrics
		server := &http.Server{
			Addr:              ":8080",
			ReadHeaderTimeout: 3 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	if err := cmd.Run(wait.NeverStop); err != nil {
		panic(err)
	}
}
