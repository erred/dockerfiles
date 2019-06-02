FROM golang:alpine AS build

WORKDIR /app
COPY . .

RUN apk add --update --no-cache ca-certificates
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o=cf-dns-update

FROM scratch

COPY --from=build /app/cf-dns-update /bin
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/bin/cf-dns-update"]
