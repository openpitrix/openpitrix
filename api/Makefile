# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis
GRPC_GATEWAY_PATH=third_party/grpc-ecosystem/grpc-gateway/

#PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH) -I$(VENDOR_PATH)
PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH)

GOPATH:=$(shell go env GOPATH)
PWD:=$(shell pwd)

METADATA_TYPES_PBFILES=$(sort $(wildcard metadata/types/*.proto))
METADATA_METADATA_PBFILES=$(sort $(wildcard metadata/metadata/*.proto))
METADATA_DRONE_PBFILES=$(sort $(wildcard metadata/drone/*.proto))
METADATA_FRONTGATE_PBFILES=$(sort $(wildcard metadata/frontgate/*.proto))
METADATA_PILOT_PBFILES=$(sort $(wildcard metadata/pilot/*.proto))

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  tools             to install depends tools"
	@echo "  help              to show this help informatrion"

generate_opctl:
	go run ../cmd/opctl/opctl-gen-yaml/main.go -f ../pkg/apigateway/spec/api.swagger.json > ../cmd/opctl/config.yaml
	go run ../cmd/opctl/opctl-gen-cmd/main.go -f ../cmd/opctl/config.yaml > ../cmd/opctl/all_cmd.go
	goimports -l -w -e -local=openpitrix ../cmd/opctl/all_cmd.go

generate: $(PROTO_FILES) Makefile
	mkdir -p ../pkg/pb
	rm -rf ../pkg/pb/*

	make metadata_types
	make metadata_drone
	make metadata_frontgate
	make metadata_pilot

	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/apigateway/spec ${PROTO_FILES}
	cd ../pkg/apigateway/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' | jq -f ./preprocess.jq > ./api.swagger.json
	cd ../pkg/apigateway/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm
#	# generate client code
	make generate_opctl
#
	rm -rf ../test/models/* ../test/client/*
	cd ../test/ && swagger generate client -f ../pkg/apigateway/spec/api.swagger.json -A openpitrix
#	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg ../test
#	# format generated files
	goimports -l -w -e -local=openpitrix ../pkg/pb
	@echo "ok"

metadata_types: ${METADATA_TYPES_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/types
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/types/*
	protoc $(PROTOC_FLAGS) --go_out=$(GOPATH)/src ${METADATA_TYPES_PBFILES}
	ls $(PWD)/../build/replace
	go run $(PWD)/../build/replace/main.go -old=,omitempty -new= -dir=$(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/types
	@echo "ok"

metadata_drone: ${METADATA_DRONE_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/drone
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/drone/*
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${METADATA_DRONE_PBFILES}
	@echo "ok"

metadata_frontgate: ${METADATA_FRONTGATE_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/frontgate
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/frontgate/*
	protoc $(PROTOC_FLAGS) --stdrpc_out=$(GOPATH)/src ${METADATA_FRONTGATE_PBFILES}
	@echo "ok"

metadata_pilot: ${METADATA_PILOT_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/pilot
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/metadata/pilot/*
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${METADATA_PILOT_PBFILES}
	@echo "ok"

tools:
	docker pull openpitrix/openpitrix:builder
	@echo "ok"

clean:
	@echo "ok"
