package main

import (
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
func (a *Adapter) makeProviderOrDie() provider.CustomMetricsProvider {
	config, err := a.ClientConfig()
	if err != nil {
		klog.Fatalf("unable to construct client config: %v", err)
	}

	client, err := a.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		klog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return customprovider.New(client, config, mapper)
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := &Adapter{}
	cmd.Name = "skpr-fpm-metrics-adapter"

	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	logs.AddFlags(cmd.Flags())
	if err := cmd.Flags().Parse(os.Args); err != nil {
		klog.Fatalf("unable to parse flags: %v", err)
	}

	testProvider := cmd.makeProviderOrDie()
	cmd.WithCustomMetrics(testProvider)

	if err := metrics.RegisterMetrics(legacyregistry.Register); err != nil {
		klog.Fatal("unable to register metrics: %v", err)
	}

	klog.Infof(cmd.Message)

	go func() {
		// Open port for POSTing fake metrics
		server := &http.Server{
			Addr:              ":8080",
			ReadHeaderTimeout: 3 * time.Second,
		}
		klog.Fatal(server.ListenAndServe())
	}()

	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}
