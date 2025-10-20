FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY skpr-fpm-metrics-adapter-sidecar /usr/local/bin/skpr-fpm-metrics-adapter-sidecar
CMD ["skpr-fpm-metrics-adapter-sidecar"]
