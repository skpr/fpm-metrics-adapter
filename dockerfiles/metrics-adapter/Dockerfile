FROM golang:1.25-alpine AS build
ADD . /go/src/github.com/skpr/fpm-metrics-adapter
WORKDIR /go/src/github.com/skpr/fpm-metrics-adapter
ENV CGO_ENABLED=0
RUN go build -o bin/skpr-fpm-metrics-adapter -ldflags='-extldflags "-static"' github.com/skpr/fpm-metrics-adapter/cmd/skpr-fpm-metrics-adapter

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/skpr/fpm-metrics-adapter/bin/skpr-fpm-metrics-adapter /usr/local/bin/skpr-fpm-metrics-adapter
CMD ["skpr-fpm-metrics-adapter"]
