FROM golang:alpine

RUN apk add --no-cache git curl protobuf

RUN curl -L -o /usr/bin/protoc-gen-grpc-web \
        https://github.com/grpc/grpc-web/releases/download/1.0.4/protoc-gen-grpc-web-1.0.4-linux-x86_64 && \
        chmod +x /usr/bin/protoc-gen-grpc-web

RUN GOBIN=/usr/bin go get github.com/golang/protobuf/protoc-gen-go && rm -rf /root/.cache /go

ENTRYPOINT ["/usr/bin/protoc"]
