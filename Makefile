# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=openpitrix
TRAG.Gopkg:=openpitrix.io/openpitrix
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

DOCKER_TAGS=latest
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` openpitrix/openpitrix-builder
GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_FILES:=./cmd ./test ./pkg
DB_TEST:=OP_DB_UNIT_TEST=1 OPENPITRIX_MYSQL_HOST=127.0.0.1 OPENPITRIX_MYSQL_PORT=13306
ETCD_TEST:=OP_ETCD_UNIT_TEST=1 OPENPITRIX_ETCD_ENDPOINTS=127.0.0.1:12379
define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -E "^(test|cmd|pkg)/.+\.go"))
endef
define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef

COMPOSE_APP_SERVICES=openpitrix-runtime-manager openpitrix-app-manager openpitrix-repo-indexer openpitrix-api-gateway openpitrix-repo-manager openpitrix-job-manager openpitrix-task-manager openpitrix-cluster-manager openpitrix-pilot-service
COMPOSE_DB_CTRL=openpitrix-app-db-ctrl openpitrix-repo-db-ctrl openpitrix-runtime-db-ctrl openpitrix-job-db-ctrl openpitrix-task-db-ctrl openpitrix-cluster-db-ctrl
CMD?=...
comma:= ,
empty:=
space:= $(empty) $(empty)
CMDS=$(subst $(comma),$(space),$(CMD))

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

.PHONY: buildflyway
buildflyway:
	docker build -t $(TARG.Name):flyway -f ./Dockerfile.flyway .

.PHONY: build
build: fmt buildflyway
	mkdir -p ./tmp/bin
	$(call get_build_flags)
	$(RUN_IN_DOCKER) time go install -v -ldflags '$(BUILD_FLAG)' $(foreach cmd,$(CMDS),$(TRAG.Gopkg)/cmd/$(cmd))
	$(RUN_IN_DOCKER) time go install -v -ldflags '$(BUILD_FLAG)' $(TRAG.Gopkg)/metadata/cmd/...
	docker build -t $(TARG.Name) -t $(TARG.Name):metadata -f ./Dockerfile.dev ./tmp/bin
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-update
compose-update: build compose-up compose-migrate-db
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
	docker-compose exec openpitrix-db bash -c "cat /docker-entrypoint-initdb.d/*.sql | mysql -uroot -ppassword"
	docker-compose up $(COMPOSE_DB_CTRL)

compose-update-%:
	CMD=$* make build
	docker-compose up -d --no-deps openpitrix-$*
	@echo "compose-update done"

.PHONY: compose-up
compose-up:
	docker-compose up -d openpitrix-db && sleep 20 && docker-compose up -d
	@echo "compose-up done"

.PHONY: compose-up-app
compose-up-app:
	docker-compose -f docker-compose-app.yml up -d openpitrix-db && sleep 20 && docker-compose -f docker-compose-app.yml up -d
	@echo "compose-up app service done"

.PHONY: compose-down
compose-down:
	docker-compose down
	@echo "compose-down done"

.PHONY: release
release:
	@echo "TODO"

.PHONY: test
test:
	make unit-test
	make e2e-test
	@echo "test done"


.PHONY: e2e-test
e2e-test:
	go test -v -a -tags="integration" ./test/...
	@echo "e2e-test done"

.PHONY: ci-test
ci-test:
	make fmt-check
	# build with production Dockerfile, not dev version
	docker build -t $(TARG.Name) -f ./Dockerfile .
	docker build -t $(TARG.Name):metadata -f ./Dockerfile.metadata .
	make compose-up
	sleep 20
	make unit-test
	make e2e-test
	@echo "ci-test done"

.PHONY: clean
clean:
	-make -C ./pkg/version clean
	@echo "ok"

.PHONY: unit-test
unit-test:
	$(DB_TEST) $(ETCD_TEST) go test -v -a -tags="etcd db" ./...
	@echo "unit-test done"
