# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.9.2-alpine3.6

RUN apk add --no-cache \
	git \
	curl \
	protobuf \
	hugo \
	caddy

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get github.com/go-swagger/go-swagger/cmd/swagger
