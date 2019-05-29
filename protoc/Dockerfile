FROM alpine:latest AS builder
RUN apk add --no-cache build-base curl automake autoconf libtool git zlib-dev

ENV GRPC_VERSION=1.16.0 \
        GRPC_WEB_VERSION=1.0.4 \
        PROTOBUF_VERSION=3.6.1 \
        OUTDIR=/out

RUN mkdir -p /protobuf && \
        curl -L https://github.com/google/protobuf/archive/v${PROTOBUF_VERSION}.tar.gz | tar xvz --strip-components=1 -C /protobuf

RUN git clone --depth 1 --recursive -b v${GRPC_VERSION} https://github.com/grpc/grpc.git /grpc && \
        rm -rf grpc/third_party/protobuf && \
        ln -s /protobuf /grpc/third_party/protobuf

RUN mkdir -p /grpc-web && \
        curl -L https://github.com/grpc/grpc-web/archive/${GRPC_WEB_VERSION}.tar.gz | tar xvz --strip-components=1 -C /grpc-web

RUN cd /protobuf && \
        autoreconf -f -i -Wall,no-obsolete && \
        ./configure --prefix=/usr --enable-static=no && \
        make -j2 && make install

RUN cd grpc && \
        make -j2 plugins

RUN cd /protobuf && \
        make install DESTDIR=${OUTDIR}

RUN cd /grpc && \
        make install-plugins prefix=${OUTDIR}/usr

RUN cd /grpc-web/javascript/net/grpc/web && \
        make && \
        install protoc-gen-grpc-web ${OUTDIR}/usr/bin/

RUN find ${OUTDIR} -name "*.a" -delete -or -name "*.la" -delete

RUN apk add --no-cache go
ENV GOPATH=/go \
        PATH=/go/bin/:$PATH
RUN go get -u -v -ldflags '-w -s' \
        github.com/golang/protobuf/protoc-gen-go \
        && install -c ${GOPATH}/bin/protoc-gen-go ${OUTDIR}/usr/bin/


FROM alpine:latest
RUN apk add --no-cache libstdc++
COPY --from=builder /out/ /

RUN chmod a+x /usr/bin/protoc

ENTRYPOINT ["/usr/bin/protoc"]
