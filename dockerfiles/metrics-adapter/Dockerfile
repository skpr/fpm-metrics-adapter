FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/skpr-fpm-metrics-adapter /usr/local/bin/skpr-fpm-metrics-adapter

CMD ["skpr-fpm-metrics-adapter"]
