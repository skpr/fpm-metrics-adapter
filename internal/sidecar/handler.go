package sidecar

import (
	"encoding/json"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/skpr/fpm-metrics-adapter/internal/fpm"
)

// Handler to return the latest status.
func (s *Server) handler(w http.ResponseWriter, _ *http.Request) {
	s.logger.Debug("Handling request")

	var jsonBytes []byte
	c, found := s.cache.Get("json")
	if found {
		s.logger.Debug("Cache hit")
		jsonBytes = c.([]byte)
	} else {
		s.logger.Debug("Cache miss")
		status, err := fpm.QueryStatus(s.config.Endpoint)
		if err != nil {
			s.logger.Error("Error querying fpm", "error", err.Error())
		}

		jsonBytes, err = json.Marshal(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.Error("failed to marshal status response", "error", err.Error())
			return
		}

		s.cache.Set("json", jsonBytes, cache.DefaultExpiration)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(jsonBytes)
	if err != nil {
		s.logger.Error("failed to write status response", "error", err.Error())
		return
	}

	s.logger.Debug("Request complete")
}
