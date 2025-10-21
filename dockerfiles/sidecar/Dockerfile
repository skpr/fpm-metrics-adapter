FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY skpr-fpm-metrics-adapter-sidecar /usr/local/bin/skpr-fpm-metrics-adapter-sidecar

CMD ["skpr-fpm-metrics-adapter-sidecar"]
