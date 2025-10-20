FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY skpr-fpm-metrics-adapter /usr/local/bin/skpr-fpm-metrics-adapter
CMD ["skpr-fpm-metrics-adapter"]
