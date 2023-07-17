FROM golang:1.19 as builder
WORKDIR /go/src/github.com/skpr/fpm-metrics-adapter
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/metrics-adapter github.com/skpr/fpm-metrics-adapter/cmd/metrics-adapter

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/skpr/fpm-metrics-adapter/bin/metrics-adapter /usr/local/bin/metrics-adapter
CMD ["metrics-adapter"]
