# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.13-alpine as builder

RUN apk add --no-cache git curl openssl

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/chai2010/protorpc/protoc-gen-stdrpc
RUN go get golang.org/x/tools/cmd/goimports

# swagger-0.13.0
# RUN go get github.com/go-swagger/go-swagger/cmd/swagger
RUN mkdir -p /swagger && cd /swagger \
	&& wget https://github.com/go-swagger/go-swagger/releases/download/0.13.0/swagger_linux_amd64 \
	&& chmod +x swagger_linux_amd64 && mv swagger_linux_amd64 /go/bin/swagger

RUN wget https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v1.16.0/protoc-gen-grpc-gateway-v1.16.0-linux-x86_64 \
    && chmod -x protoc-gen-grpc-gateway-v1.16.0-linux-x86_64 && mv protoc-gen-grpc-gateway-v1.16.0-linux-x86_64 /go/bin/protoc-gen-grpc-gateway

RUN wget https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v1.16.0/protoc-gen-swagger-v1.16.0-linux-x86_64 \
    && chmod -x protoc-gen-swagger-v1.16.0-linux-x86_64 && mv protoc-gen-swagger-v1.16.0-linux-x86_64 /go/bin/protoc-gen-swagger

# the protoc can't run on alpine,
# we only need the protobuf's stdarnd library in the `/protoc/include`.
RUN mkdir -p /protoc && cd /protoc \
	&& wget https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip \
	&& unzip protoc-3.5.0-linux-x86_64.zip

FROM golang:1.13-alpine

RUN apk add --no-cache git protobuf make curl openssl jq rsync upx build-base

COPY --from=builder /protoc/include /usr/local/include
COPY --from=builder /go/bin /go/bin
