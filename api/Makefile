# Copyright 2018 The OpenPitrix Authors. All rights reserved.
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
	mkdir -p ../pkg/pb
	rm -rf ../pkg/pb/*
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/cmd/api/spec ${PROTO_FILES}
	cd ../pkg/cmd/api/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' > ./api.swagger.json
	cd ../pkg/cmd/api/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm
	# generate client code
	rm -rf ../test/models/* ../test/client/*
	cd ../test/ && swagger generate client -f ../pkg/cmd/api/spec/api.swagger.json -A openpitrix
	# generate frontgate/drone code
	cd frontgate && make generate
	cd drone && make generate
	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg ../test
	@echo "ok"

tools:
	docker pull openpitrix/openpitrix:builder
	@echo "ok"

clean:
	@echo "ok"
