package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// The process which will continually log the current status.
func (s *Server) logStatus(ctx context.Context) error {
	ticker := time.NewTicker(s.config.LogFrequency)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			jsonBytes, err := json.Marshal(s.status)
			if err != nil {
				s.logger.Error("failed to marshal status for logger", "error", err.Error())
			}

			fmt.Println(string(jsonBytes))
		}
	}
}
