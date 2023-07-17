FROM alpine:3.17

RUN apk --no-cache add ca-certificates
COPY skpr sidecar /usr/local/bin/

CMD ["sidecar"]