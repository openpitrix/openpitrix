# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=openpitrix
TRAG.Gopkg:=openpitrix.io/openpitrix

TARG.Services:=api
TARG.Services+=app
TARG.Services+=repo
TARG.Services+=runtime
TARG.Services+=cluster

DOCKER_TAGS=latest

GO:=docker run --rm -it -v `pwd`:/go/src/$(TRAG.Gopkg) -w /go/src/$(TRAG.Gopkg) golang:1.9-alpine go
GO_WITH_MYSQL:=docker run --rm -it -v `pwd`:/go -w /go/src/$(TRAG.Gopkg) --link openpitrix-mysql:mysql -p 9527:9527 golang:1.9-alpine go

MYSQL_DATABASE:=$(TARG.Name)
MYSQL_ROOT_PASSWORD:=password

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all               to generate, test and release"
	@echo "  tools             to install depends tools"
	@echo "  init-vendor       to init vendor packages"
	@echo "  update-vendor     to update vendor packages"
	@echo "  generate          to generate restapi code"
	@echo "  build             to build the service"
	@echo "  test              to run go test ./..."
	@echo "  release           to build and release current version"
	@echo "  clean             to clean the temp files"

.PHONY: all
all: generate test release

.PHONY: init-vendor
init-vendor:
	@if [[ ! -f "$$(which govendor)" ]]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor init
	govendor add +external
	@echo "ok"

.PHONY: update-vendor
update-vendor:
	@if [[ ! -f "$$(which govendor)" ]]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor update +external
	govendor list
	@echo "ok"

.PHONY: tools
tools:
	docker pull golang:1.9-alpine
	docker pull chai2010/grpc-tools
	docker pull mysql
	@echo "ok"

.PHONY: generate
generate:
	cd ./api && make generate
	cd ./pkg/cmd/api && make
	@echo "ok"

.PHONY: mysql-start
mysql-start:
	@docker run --rm --name openpitrix-mysql -e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) -e MYSQL_DATABASE=$(MYSQL_DATABASE) -d mysql || docker start openpitrix-mysql
	@echo "ok"

.PHONY: mysql-stop
mysql-stop:
	@docker stop openpitrix-mysql
	@echo "ok"


.PHONY: build
build:
	$(GO) fmt ./...
	docker build -t $(TARG.Name) -f ./Dockerfile .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "ok"

.PHONY: run
run: mysql-start
	@for service in $(TARG.Services) ; do \
		docker run --rm -d $$service ; \
	done
	@echo "ok"

.PHONY: release
release:
	@echo "TODO"

.PHONY: test
test:
	$(GO) fmt ./...
	$(GO) test ./...
	@echo "ok"

.PHONY: clean
clean:
	@echo "ok"
