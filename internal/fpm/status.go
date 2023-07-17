package fpm

import (
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	fcgiclient "github.com/tomasen/fcgi_client"
)

var queryStatusRegexp = regexp.MustCompile(`(?m)^(.*):\s+(.*)$`)

// QueryStatus of the FPM worker pool.
func QueryStatus(address string) (Status, error) {
	var status Status

	env := map[string]string{
		"SCRIPT_FILENAME": "/status",
		"SCRIPT_NAME":     "/status",
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

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		return status, errors.New("")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return status, err
	}

	matches := queryStatusRegexp.FindAllStringSubmatch(string(body), -1)

	for _, match := range matches {
		key := match[1]

		value, err := strconv.Atoi(match[2])
		if err != nil {
			continue
		}

		if key == "total processes" {
			status.Processes.Total = int64(value)
		}

		if key == "active processes" {
			status.Processes.Active = int64(value)
		}
	}

	return status, nil
}
