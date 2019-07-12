package log

import (
	"github.com/prometheus/common/log"
	"github.com/rs/xid"
)

const (
	// KeyRequest for identifying a whole request (set of logs).
	KeyRequest = "request"
	// KeyFunction for identifying which function this log belongs to.
	KeyFunction = "function"
)

// New logger for server interactions.
func New(name string) log.Logger {
	return log.With(KeyRequest, xid.New().String()).With(KeyFunction, name)
}
