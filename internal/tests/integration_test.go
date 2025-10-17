//go:build integration

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/christgf/env"
	"io"
	"net/http"
	"testing"
)

type ExpectedMetricsResponse struct {
	PhpfpmProcessManager     string `json:"phpfpm_process_manager"`
	PhpfpmListenQueue        int    `json:"phpfpm_listen_queue"`
	PhpfpmListenQueueLen     int    `json:"phpfpm_listen_queue_len"`
	PhpfpmIdleProcesses      int    `json:"phpfpm_idle_processes"`
	PhpfpmActiveProcesses    int    `json:"phpfpm_active_processes"`
	PhpfpmTotalProcesses     int    `json:"phpfpm_total_processes"`
	PhpfpmMaxActiveProcesses int    `json:"phpfpm_max_active_processes"`
}

func TestMetrics(t *testing.T) {
	resp, err := getUrl("/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := &ExpectedMetricsResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		t.Fatal(err)
	}

	if result.PhpfpmProcessManager != "dynamic" {
		t.Errorf("Expected %s, got %s", "dynamic", result.PhpfpmProcessManager)
	}
	if result.PhpfpmActiveProcesses != 1 {
		t.Errorf("Expected %d, got %d", 1, result.PhpfpmActiveProcesses)
	}
	if result.PhpfpmIdleProcesses != 1 {
		t.Errorf("Expected %d, got %d", 1, result.PhpfpmIdleProcesses)
	}
}

func TestUnknownUrl(t *testing.T) {
	resp, err := getUrl("/unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound, got %v", resp.StatusCode)
	}
}

func getUrl(path string) (resp *http.Response, err error) {
	server := env.String("TEST_SERVER_URL", "http://127.0.0.1:8080")

	// You could also test a path that should fail
	return http.Get(fmt.Sprintf("%s%s", server, path))
}
