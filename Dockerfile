FROM golang:1.20-alpine3.17 as build

WORKDIR /go/src/github.com/skpr/fpm-metrics-adapter
COPY . /go/src/github.com/skpr/fpm-metrics-adapter

RUN apk add gcc musl-dev

RUN go build -ldflags "-linkmode external -extldflags -static" -o fpm-metrics-adapter cmd/metrics-adapter/main.go

FROM alpine:3.17

COPY --from=build /go/src/github.com/skpr/fpm-metrics-adapter/fpm-metrics-adapter /usr/local/bin/fpm-metrics-adapter

CMD ["/usr/local/bin/fpm-metrics-adapter"]
