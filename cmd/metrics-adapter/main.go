package main

import (
	"flag"
	"os"

	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog"

	customprovider "github.com/skpr/fpm-metrics-adapter/pkg/provider"
)

// Adapter for setting up the custom metrics server.
type Adapter struct {
	basecmd.AdapterBase

	// Message is printed on succesful startup
	Message string
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := &Adapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags
	cmd.Flags().Parse(os.Args)

	p := cmd.makeProviderOrDie()
	cmd.WithCustomMetrics(p)

	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
}

// Helper function to setup the custom metrics provider.
func (a *Adapter) makeProviderOrDie() provider.CustomMetricsProvider {
	config, err := a.ClientConfig()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
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
