FROM alpine:3.21
ARG TARGETPLATFORM
RUN apk --no-cache add ca-certificates
COPY $TARGETPLATFORM/skpr-fpm-metrics-adapter-sidecar /usr/local/bin/skpr-fpm-metrics-adapter-sidecar
CMD ["skpr-fpm-metrics-adapter-sidecar"]
