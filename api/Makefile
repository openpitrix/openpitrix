# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis
GRPC_GATEWAY_PATH=../vendor/github.com/grpc-ecosystem/grpc-gateway/
VENDOR_PATH=../vendor/

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH) -I$(VENDOR_PATH)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  tools             to install depends tools"
	@echo "  help              to show this help informatrion"

generate: $(PROTO_FILES) Makefile
	mkdir -p ../pkg/service.pb
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/service.pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --govalidators_out=:../pkg/service.pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true:../pkg/service.pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true:../pkg/service.pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true:../pkg/cmd/api/spec ${PROTO_FILES}

	# generate client code
	cat ../pkg/cmd/api/spec/*.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' > ../test/api.swagger.json
	cd ../test/ && swagger generate client -f api.swagger.json -A openpitrix
	rm ../test/api.swagger.json

	@echo "ok"

tools:
	docker pull openpitrix/openpitrix:builder
	@echo "ok"

clean:
	@echo "ok"
