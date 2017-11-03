# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PROTO_FILES=$(wildcard ./*.proto)

PROTOC:=docker run --rm -it -v $(HOME):$(HOME) -w `pwd` chai2010/grpc-tools protoc

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  tools             to install depends tools"
	@echo "  help              to show this help informatrion"

generate:
	mkdir -p ../pkg/service.pb
	$(PROTOC) --go_out=plugins=grpc:../pkg/service.pb ${PROTO_FILES}
	$(PROTOC) --grpc-gateway_out=logtostderr=true:../pkg/service.pb ${PROTO_FILES}
	$(PROTOC) --swagger_out=logtostderr=true:../pkg/service.pb ${PROTO_FILES}
	$(PROTOC) --swagger_out=logtostderr=true:../pkg/cmd/api/spec ${PROTO_FILES}
	@echo "ok"

tools:
	docker pull chai2010/grpc-tools
	@echo "ok"

