package sidecar

import (
	"encoding/json"
	"net/http"
)

// Handler to return the latest status.
func (s *Server) handler(w http.ResponseWriter, _ *http.Request) {
	s.logger.Debug("Handling request")

	jsonBytes, err := json.Marshal(s.status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Error("failed to marshal status response", "error", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		s.logger.Error("failed to write status response", "error", err.Error())
		return
	}

	s.logger.Debug("Request complete")
}
