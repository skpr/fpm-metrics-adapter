FROM alpine:3.17

RUN apk --no-cache add ca-certificates
COPY metrics-adapter /usr/local/bin/

CMD ["metrics-adapter"]
