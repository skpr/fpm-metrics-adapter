#!/usr/bin/make -f

export CGO_ENABLED=0

default: lint test

# Run all lint checking with exit codes for CI.
lint:
	revive -config revive.toml -set_exit_status ./cmd/... ./internal/...

# Run tests with coverage reporting.
test:
	go test -cover ./...
