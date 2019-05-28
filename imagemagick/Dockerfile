FROM alpine:latest

RUN apk --update --no-cache add imagemagick

WORKDIR /workspace
ENTRYPOINT ["/usr/bin/convert"]
