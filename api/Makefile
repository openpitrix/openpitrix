# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis
GRPC_GATEWAY_PATH=../vendor/github.com/grpc-ecosystem/grpc-gateway/
VENDOR_PATH=../vendor/

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(GRPC_GATEWAY_PATH) -I$(PROTOC_INC_PATH) -I$(VENDOR_PATH)

GOPATH:=$(shell go env GOPATH)
PWD:=$(shell pwd)

OPRNPITRIX_TYPES_PBFILES=$(sort $(wildcard openpitrix/types/*.proto))
OPRNPITRIX_METADATA_PBFILES=$(sort $(wildcard openpitrix/metadata/*.proto))
OPRNPITRIX_DRONE_PBFILES=$(sort $(wildcard openpitrix/drone/*.proto))
OPRNPITRIX_FRONTGATE_PBFILES=$(sort $(wildcard openpitrix/frontgate/*.proto))
OPRNPITRIX_PILOT_PBFILES=$(sort $(wildcard openpitrix/pilot/*.proto))

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  generate          to generate grpc/grpc-gateway/swagger code"
	@echo "  tools             to install depends tools"
	@echo "  help              to show this help informatrion"

generate: $(PROTO_FILES) Makefile
	mkdir -p ../pkg/pb
	rm -rf ../pkg/pb/*

	make openpitrix_types
	make openpitrix_drone
	make openpitrix_frontgate
	make openpitrix_pilot

	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PROTO_FILES}
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/cmd/api/spec ${PROTO_FILES}
	cd ../pkg/cmd/api/spec/ && cat ./*.swagger.json 0.swagger.json | jq --slurp 'reduce .[] as $$item ({}; . * $$item)' | jq -f ./preprocess.jq > ./api.swagger.json
	cd ../pkg/cmd/api/spec/ && find . | grep .swagger.json | grep -v "api" | xargs rm
	# generate client code
	rm -rf ../test/models/* ../test/client/*
	cd ../test/ && swagger generate client -f ../pkg/cmd/api/spec/api.swagger.json -A openpitrix
	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg ../test
	@echo "ok"

openpitrix_types: ${OPRNPITRIX_TYPES_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/types
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/types/*
	protoc $(PROTOC_FLAGS) --go_out=$(GOPATH)/src ${OPRNPITRIX_TYPES_PBFILES}
	ls $(PWD)/../build/replace
	go run $(PWD)/../build/replace/main.go -old=,omitempty -new= -dir=$(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/types
	@echo "ok"

openpitrix_drone: ${OPRNPITRIX_DRONE_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/drone
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/drone/*
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${OPRNPITRIX_DRONE_PBFILES}
	@echo "ok"

openpitrix_frontgate: ${OPRNPITRIX_FRONTGATE_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/frontgate
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/frontgate/*
	protoc $(PROTOC_FLAGS) --stdrpc_out=$(GOPATH)/src ${OPRNPITRIX_FRONTGATE_PBFILES}
	@echo "ok"

openpitrix_pilot: ${OPRNPITRIX_PILOT_PBFILES}
	mkdir -p $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/pilot
	rm -rf   $(GOPATH)/src/openpitrix.io/openpitrix/pkg/pb/pilot/*
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${OPRNPITRIX_PILOT_PBFILES}
	@echo "ok"


tools:
	docker pull openpitrix/openpitrix:builder
	@echo "ok"

clean:
	@echo "ok"
