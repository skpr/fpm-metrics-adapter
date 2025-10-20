package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func getClientset() *fake.Clientset {
	fakeClientset := fake.NewClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fail-ip",
				Namespace: "default",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Status: corev1.PodStatus{
				PodIP: "127.0.0.1",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "annotated-pod",
				Namespace: "default",
				Annotations: map[string]string{
					AnnotationProtocol: "https",
					AnnotationPort:     "443",
					AnnotationPath:     "/new-metrics",
				},
			},
			Status: corev1.PodStatus{
				PodIP: "127.0.0.70",
			},
		},
	)
	return fakeClientset
}

func TestGetConnFailIp(t *testing.T) {
	clientset := getClientset()

	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "fail-ip", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unable to find pod: %v", err)
	}

	_, err = getConn(pod)
	if err == nil {
		t.Fatalf("expected an error: %v", err)
	}
}

func TestGetConn(t *testing.T) {
	assertEndpoint(t, "test-pod", "http://127.0.0.1:80/metrics")
	assertEndpoint(t, "annotated-pod", "https://127.0.0.70:443/new-metrics")
}

func assertEndpoint(t *testing.T, name string, endpoint string) {
	clientset := getClientset()

	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unable to find pod: %v", err)
	}

	connectionString, err := getConn(pod)
	if err != nil {
		t.Fatalf("unable to find endpoint: %v", err)
	}

	if connectionString != endpoint {
		t.Fatalf("did not build expected endpoint for pod: %s, endpoint: %s", name, endpoint)
	}
}

func TestGetMetric(t *testing.T) {
	value := fpm.Status{
		IdleProcesses: 100,
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("unable to marshal value: %v", err)
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(jsonValue)
			if err != nil {
				t.Fatalf("unable to write response: %v", err)
			}
			return
		}
		http.NotFound(w, r)
	}))

	endpoint := fmt.Sprintf("%s/metrics", mockServer.URL)
	resp, err := getMetric(endpoint, fpm.MetricIdleProcesses)

	if err != nil {
		t.Fatalf("unable to scrape metrics: %v", err)
	}

	if resp != 100 {
		t.Fatalf("metrics scrape did not return 100")
	}
}
