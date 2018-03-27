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
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -E "^(test|cmd|pkg)/.+\.go"))
endef

COMPOSE_APP_SERVICES=openpitrix-runtime-env-manager openpitrix-app-manager openpitrix-repo-indexer openpitrix-api-gateway openpitrix-repo-manager openpitrix-job-manager openpitrix-task-manager openpitrix-cluster-manager
COMPOSE_DB_CTRL=openpitrix-app-db-ctrl openpitrix-repo-db-ctrl openpitrix-runtime-db-ctrl openpitrix-job-db-ctrl openpitrix-task-db-ctrl openpitrix-cluster-db-ctrl

.PHONY: all
all: generate build

.PHONY: help
help:
# TODO: update help info to last version
	@echo "TODO"

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

.PHONY: update-builder
update-builder:
	docker pull openpitrix/openpitrix-builder
	@echo "update-builder done"

.PHONY: generate-in-local
generate-in-local:
	cd ./api && make generate
	cd ./pkg/cmd/api && make
	go generate ./pkg/version/

.PHONY: generate
generate:
	$(RUN_IN_DOCKER) make generate-in-local
	@echo "generate done"

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

.PHONY: fmt-check
fmt-check: fmt-all
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)

.PHONY: build
build: fmt
	docker build -t $(TARG.Name) -f ./Dockerfile .
	@docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-update
compose-update: build compose-up
	@echo "compose-update done"

.PHONY: compose-update-service-without-deps
compose-update-service-without-deps: build
	docker-compose up -d --no-dep $(COMPOSE_APP_SERVICES)
	@echo "compose-update-service-without-deps done"

.PHONY: compose-logs-f
compose-logs-f:
	docker-compose logs -f $(COMPOSE_APP_SERVICES)

.PHONY: compose-migrate-db
compose-migrate-db:
	docker-compose up $(COMPOSE_DB_CTRL)

compose-update-%:
	docker-compose up -d --no-deps $*
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
	@make ci-unit-test
	@echo "ci-test done"

.PHONY: clean
clean:
	-make -C ./pkg/version clean
	@echo "ok"

.PHONY: test-db-up
test-db-up:
	docker-compose -p openpitrix-test-db -f ./docker-compose.test.yml up -d openpitrix-db && \
	sleep 20 && \
	docker-compose -p openpitrix-test-db -f ./docker-compose.test.yml up -d

.PHONY: test-db-down
test-db-down:
	docker-compose -p openpitrix-test-db down

.PHONY: ci-unit-test
ci-unit-test: test-db-up
	cd ./pkg/manager/runtime_env/ && OPTESTCONFIG_DBTEST=true OPTESTCONFIG_DB_DATABASE="runtime" go test -v ./...
	make test-db-down
