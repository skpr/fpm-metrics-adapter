#!/usr/bin/make -f

export CGO_ENABLED=0

IMAGE=skpr/fpm-adapter-metrics

# Builds the project.
build:
	docker build -f dockerfiles/metrics-adapter.dockerfile -t ${IMAGE}:metrics-adapter-latest .
	docker build -f dockerfiles/sidecar.dockerfile -t ${IMAGE}:sidecar-latest .

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

.PHONY: *