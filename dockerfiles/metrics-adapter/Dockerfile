FROM alpine:3.21
ARG TARGETPLATFORM
RUN apk --no-cache add ca-certificates
COPY $TARGETPLATFORM/skpr-fpm-metrics-adapter /usr/local/bin/skpr-fpm-metrics-adapter
CMD ["skpr-fpm-metrics-adapter"]
