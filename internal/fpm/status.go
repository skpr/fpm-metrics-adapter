// Package fpm for interacting with the FPM process.
package fpm

import (
	"encoding/json"
	"fmt"

	fcgiclient "github.com/tomasen/fcgi_client"
)

// QueryStatus of the FPM worker pool.
func QueryStatus(address string) (Status, error) {
	var status Status

	env := map[string]string{
		"SCRIPT_FILENAME": "/status",
		"SCRIPT_NAME":     "/status",
		"QUERY_STRING":    "json&full",
	}

	fcgi, err := fcgiclient.Dial("tcp", address)
	if err != nil {
		return status, err
	}
	defer fcgi.Close()

	resp, err := fcgi.Get(env)
	if err != nil {
		return status, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		return status, fmt.Errorf("status code was: %d", resp.StatusCode)
	}

	var response QueryResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return status, fmt.Errorf("failed to decode json: %w", err)
	}

	return Status(response), nil
}
