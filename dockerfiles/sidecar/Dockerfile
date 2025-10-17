FROM golang:1.25-alpine AS build
ADD . /go/src/github.com/skpr/fpm-metrics-adapter
WORKDIR /go/src/github.com/skpr/fpm-metrics-adapter
ENV CGO_ENABLED=0
RUN go build -o bin/skpr-fpm-metrics-adapter-sidecar -ldflags='-extldflags "-static"' github.com/skpr/fpm-metrics-adapter/cmd/skpr-fpm-metrics-adapter-sidecar

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/skpr/fpm-metrics-adapter/bin/skpr-fpm-metrics-adapter-sidecar /usr/local/bin/skpr-fpm-metrics-adapter-sidecar
CMD ["skpr-fpm-metrics-adapter-sidecar"]
