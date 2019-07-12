#!/usr/bin/make -f

export CGO_ENABLED=0

IMAGE=skpr/fpm-adapter-metrics
VERSION=$(shell git describe --tags --always)

# Builds the project.
build:
	docker build -f dockerfiles/metrics-adapter.dockerfile -t ${IMAGE}:metrics-adapter-${VERSION} .
	docker build -f dockerfiles/sidecar.dockerfile -t ${IMAGE}:sidecar-${VERSION} .

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

release:
	docker push ${IMAGE}:metrics-adapter-${VERSION}
	docker push ${IMAGE}:sidecar-${VERSION}

.PHONY: *