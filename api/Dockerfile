FROM golang:1.9.2-alpine3.6 as builder

RUN apk add --no-cache \
    protobuf \
    git \
    curl \
    openssl

RUN cd /tmp;wget https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip;unzip protoc-3.5.0-linux-x86_64.zip;mv include /usr/local/include

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

FROM golang:1.9.2-alpine3.6
RUN apk add --no-cache protobuf

COPY --from=builder /usr/local/include /usr/local/include
COPY --from=builder /go/bin /go/bin