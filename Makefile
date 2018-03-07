# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=openpitrix
TRAG.Gopkg:=openpitrix.io/openpitrix

DOCKER_TAGS=latest
RUN_IN_DOCKER:=docker run --rm -it -v `pwd`:/go/src/$(TRAG.Gopkg) -w /go/src/$(TRAG.Gopkg) -e USER_ID=`id -u` -e GROUP_ID=`id -g` openpitrix/openpitrix-builder
GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_FILES:=./cmd ./test ./pkg
define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only | grep -E "^(test|cmd|pkg)/.+\.go"))
endef

MYSQL_DATABASE:=$(TARG.Name)
MYSQL_ROOT_PASSWORD:=password

.PHONY: all
all: generate build

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all               to generate, test and release"
	@echo "  start             to start services (port:9100)"
	@echo "  stop              to stop services"
	@echo "  tools             to install depends tools"
	@echo "  init-vendor       to init vendor packages"
	@echo "  update-vendor     to update vendor packages"
	@echo "  generate          to generate restapi code"
	@echo "  build             to build the services"
	@echo "  build-in-docker   to build the services using docker without installing any dependence"
	@echo "  test              to run go test ./..."
	@echo "  release           to build and release current version"
	@echo "  clean             to clean the temp files"

.PHONY: init-vendor
init-vendor:
	@if [[ ! -f "$$(which govendor)" ]]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor init
	govendor add +external
	@echo "init-vendor done"

.PHONY: update-vendor
update-vendor:
	@if [[ ! -f "$$(which govendor)" ]]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor update +external
	govendor list
	@echo "update-vendor done"

#.PHONY: tools
#tools:
#	docker pull openpitrix/openpitrix:builder
#	docker pull mysql:5.6
#	@echo "ok"

.PHONY: generate-in-local
generate-in-local:
	cd ./api && make generate
	cd ./pkg/cmd/api && make
	go generate ./pkg/version/

.PHONY: generate
generate:
	$(RUN_IN_DOCKER) make generate-in-local
	@echo "generate done"

#.PHONY: mysql-start
#mysql-start:
#	@docker run --rm --name openpitrix-db -e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) -e MYSQL_DATABASE=$(MYSQL_DATABASE) -p 3306:3306 -d mysql:5.6 || docker start openpitrix-db
#	@echo "ok"
#
#.PHONY: mysql-stop
#mysql-stop:
#	@docker stop openpitrix-db
#	@echo "ok"

.PHONY: fmt-all
fmt-all:
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

.PHONY: fmt
fmt:
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		$(RUN_IN_DOCKER) $(GO_FMT) ${DIFF_FILES}, \
		$(info cannot find modified files from git) \
	)
	@echo "fmt done"

.PHONY: build
build: fmt
	docker build -t $(TARG.Name) -f ./Dockerfile .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-update
compose-update: build compose-up
	@echo "compose-update done"

.PHONY: compose-up
compose-up:
	docker-compose up -d openpitrix-db && sleep 20 && docker-compose up -d 
	@echo "compose-up done"

.PHONY: compose-down
compose-down:
	docker-compose down
	@echo "compose-down done"

.PHONY: release
release:
	@echo "TODO"

.PHONY: test
test:
	go test ./...
	@echo "test done"


.PHONY: e2e-test
e2e-test:
	cd test && go test -v
	@echo "e2e-test done"

.PHONY: ci-test
ci-test: compose-update
	sleep 20
	@make e2e-test
	@echo "ci-test done"

.PHONY: clean
clean:
	-make -C ./pkg/version clean
	@echo "ok"
