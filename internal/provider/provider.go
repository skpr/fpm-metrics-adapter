package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider/helpers"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

const (
	// MetricActiveProcesses provides te number of active processes.
	MetricActiveProcesses = "phpfpm_active_processes"

	// AnnotationProtocol is used for configuration which protocol is used for querying metrics.
	AnnotationProtocol = "fpm.skpr.io/protocol"
	// AnnotationPort is used for configuration which port is used for querying metrics.
	AnnotationPort = "fpm.skpr.io/port"
	// AnnotationPath is used for configuration which path is used for querying metrics.
	AnnotationPath = "fpm.skpr.io/path"

	// DefaultProtocol used when querying for metrics.
	DefaultProtocol = "http"
	// DefaultPort used when querying for metrics.
	DefaultPort = "80"
	// DefaultPath used when querying for metrics.
	DefaultPath = "/metrics"
)

// CustomMetricResource wraps provider.CustomMetricInfo in a struct which stores the Name and Namespace of the resource
// So that we can accurately store and retrieve the metric as if this were an actual metrics server.
type CustomMetricResource struct {
	provider.CustomMetricInfo
	types.NamespacedName
}

// Provider is a sample implementation of provider.MetricsProvider which stores a map of fake metrics
type Provider struct {
	client dynamic.Interface
	config *rest.Config
	mapper apimeta.RESTMapper
}

// New returns an instance of Provider, along with its restful.WebService that opens endpoints to post new fake metrics
func New(client dynamic.Interface, config *rest.Config, mapper apimeta.RESTMapper) provider.CustomMetricsProvider {
	return &Provider{
		client: client,
		config: config,
		mapper: mapper,
	}
}

// GetMetricByName returns a single metric by name.
func (p *Provider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	ref, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(p.config)
	if err != nil {
		return nil, err
	}

	metric, err := scrape(ctx, clientset, ref.Namespace, ref.Name, info.Metric)
	if err != nil {
		return nil, err
	}

	value := &custom_metrics.MetricValue{
		DescribedObject: ref,
		Metric: custom_metrics.MetricIdentifier{
			Name: info.Metric,
		},
		Timestamp: metav1.Time{Time: time.Now()},
		Value:     *resource.NewQuantity(metric, resource.DecimalExponent),
	}

	return value, nil
}

// GetMetricBySelector returns a set of metrics queried by selector.
// https://github.com/kubernetes-incubator/custom-metrics-apiserver/blob/master/test-adapter/provider/provider.go#L234
func (p *Provider) GetMetricBySelector(ctx context.Context, namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
	if err != nil {
		return nil, err
	}

	var items []custom_metrics.MetricValue

	for _, name := range names {
		n := types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}

		metric, err := p.GetMetricByName(ctx, n, info, metricSelector)
		if err != nil {
			log.Error(err)
			continue
		}

		items = append(items, *metric)
	}

	list := &custom_metrics.MetricValueList{
		Items: items,
	}

	return list, nil
}

// ListAllMetrics which this adapter exposes.
func (p *Provider) ListAllMetrics() []provider.CustomMetricInfo {
	return []provider.CustomMetricInfo{
		{
			GroupResource: schema.GroupResource{Group: "", Resource: "pods"},
			Metric:        MetricActiveProcesses,
			Namespaced:    true,
		},
	}
}

// Scrape the context of the PHP-FPM exporter.
func scrape(ctx context.Context, clientset *kubernetes.Clientset, namespace, name, metric string) (int64, error) {
	var status fpm.Status

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}

	endpoint, err := getConn(pod)
	if err != nil {
		return 0, err
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return 0, err
	}

	if metric == MetricActiveProcesses {
		return status.Processes.Active, nil
	}

	return 0, errors.New("not found")
}

// Helper function to get connection details from a Pod.
func getConn(pod *corev1.Pod) (string, error) {
	if pod.Status.PodIP == "" {
		return "", errors.New("not found: .Status.PodIP")
	}

	var (
		protocol = DefaultProtocol
		port     = DefaultPort
		path     = DefaultPath
	)

	if val, ok := pod.ObjectMeta.Annotations[AnnotationProtocol]; ok {
		protocol = val
	}

	if val, ok := pod.ObjectMeta.Annotations[AnnotationPort]; ok {
		port = val
	}

	if val, ok := pod.ObjectMeta.Annotations[AnnotationPath]; ok {
		path = val
	}

	return fmt.Sprintf("%s://%s:%s%s", protocol, pod.Status.PodIP, port, path), nil
}
