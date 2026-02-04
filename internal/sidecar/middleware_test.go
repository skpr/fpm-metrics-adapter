package sidecar

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

type FpmCountClient struct {
	count int
	throw bool
}

func (client *FpmCountClient) QueryStatus() (fpm.Status, error) {
	client.count++
	if client.throw {
		return fpm.Status{}, fmt.Errorf("error")
	}
	return fpm.Status{
		ActiveProcesses: int64(5 * client.count),
	}, nil
}

// TestMetricsRefreshFloodControl tests that the metrics middleware will not
// refresh metrics more than once per second.
func TestMetricsRefreshFloodControl(t *testing.T) {
	client := &FpmCountClient{}

	config := ServerConfig{}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	server, err := NewServer(logger, config, client)
	if err != nil {
		t.Fatal(err)
	}

	triggerMetricsMiddleware(server)
	assert.Equal(t, 1, client.count)
	assert.Equal(t, float64(5), testutil.ToFloat64(server.metrics.ActiveProcesses))
	// Within a second, client call should be skipped.
	triggerMetricsMiddleware(server)
	assert.Equal(t, 1, client.count)
	assert.Equal(t, float64(5), testutil.ToFloat64(server.metrics.ActiveProcesses))

	server.metrics.LastUpdate = time.Now().Add(-5 * time.Second)

	triggerMetricsMiddleware(server)
	assert.Equal(t, 2, client.count)
	assert.Equal(t, float64(10), testutil.ToFloat64(server.metrics.ActiveProcesses))
	// Within a second, client call should be skipped.
	triggerMetricsMiddleware(server)
	assert.Equal(t, 2, client.count)
	assert.Equal(t, float64(10), testutil.ToFloat64(server.metrics.ActiveProcesses))
}

// TestMetricsRefreshQueryStatusError tests that the metrics middleware will
// return cached result if query status throws an error.
func TestMetricsRefreshQueryStatusError(t *testing.T) {
	client := &FpmCountClient{}

	config := ServerConfig{}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	server, err := NewServer(logger, config, client)
	if err != nil {
		t.Fatal(err)
	}

	triggerMetricsMiddleware(server)
	assert.Equal(t, 1, client.count)
	assert.Equal(t, float64(5), testutil.ToFloat64(server.metrics.ActiveProcesses))

	server.metrics.LastUpdate = time.Now().Add(-time.Second)
	client.throw = true

	triggerMetricsMiddleware(server)
	assert.Equal(t, 2, client.count)
	// Value cached in event of query status throwing error.
	assert.Equal(t, float64(5), testutil.ToFloat64(server.metrics.ActiveProcesses))
}

// triggerMetricsMiddleware makes a http request to trigger middleware.
func triggerMetricsMiddleware(server *Server) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := server.RefreshMetricsMiddleware(next)
	handler.ServeHTTP(rec, req)
}
