# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=openpitrix
TRAG.Gopkg:=openpitrix.io/openpitrix
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

DOCKER_TAGS=latest
BUILDER_IMAGE=openpitrix/openpitrix-builder:release-v0.2.3
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` $(BUILDER_IMAGE)
GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_RACE:=go build -race
GO_VET:=go vet
GO_FILES:=./cmd ./test ./pkg
GO_PATH_FILES:=./cmd/... ./test/... ./pkg/...
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

COMPOSE_APP_SERVICES=openpitrix-runtime-manager openpitrix-app-manager openpitrix-category-manager openpitrix-repo-indexer openpitrix-api-gateway openpitrix-repo-manager openpitrix-job-manager openpitrix-task-manager openpitrix-cluster-manager openpitrix-market-manager openpitrix-pilot-service openpitrix-iam-service
COMPOSE_DB_CTRL=openpitrix-app-db-ctrl openpitrix-repo-db-ctrl openpitrix-runtime-db-ctrl openpitrix-job-db-ctrl openpitrix-task-db-ctrl openpitrix-cluster-db-ctrl openpitrix-iam-db-ctrl openpitrix-market-db-ctrl
CMD?=...
WITH_METADATA?=yes
WITH_K8S=no
comma:= ,
empty:=
space:= $(empty) $(empty)
CMDS=$(subst $(comma),$(space),$(CMD))

.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_%-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: generate build ## Run generate and build

.PHONY: init-vendor
init-vendor: ## Initialize vendor and add dependence
	@if [ ! -f "$$(which govendor)" ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor init
	govendor add +external
	@echo "init-vendor done"

.PHONY: update-vendor
update-vendor: ## Update dependence
	@if [ ! -f "$$(which govendor)" ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor update +external
	govendor list
	@echo "update-vendor done"

.PHONY: update-builder
update-builder: ## Pull openpitrix-builder image
	docker pull $(BUILDER_IMAGE)
	@echo "update-builder done"

.PHONY: generate-in-local
generate-in-local: ## Generate code from protobuf file in local
	cd ./api && make generate
	cd ./pkg/apigateway && make

.PHONY: generate
generate: generate-global-config ## Generate code from protobuf file in docker
	$(RUN_IN_DOCKER) make generate-in-local
	@echo "generate done"

.PHONY: generate-global-config
generate-global-config: ## Generate global config
	$(RUN_IN_DOCKER) go generate openpitrix.io/openpitrix/deploy/config

.PHONY: fmt-all
fmt-all: ## Format all code
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

.PHONY: check
check: ## go vet and race
	$(GO_RACE) $(GO_PATH_FILES)
	$(GO_VET) $(GO_PATH_FILES)

.PHONY: fmt
fmt: ## Format changed files
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		$(RUN_IN_DOCKER) $(GO_FMT) ${DIFF_FILES}, \
		$(info cannot find modified files from git) \
	)
	@echo "fmt done"

.PHONY: fmt-check
fmt-check: fmt-all ## Check whether all files be formatted
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)

.PHONY: build-flyway
build-flyway: ## Build custom flyway image
	docker build -t $(TARG.Name):flyway -f ./pkg/db/Dockerfile ./pkg/db/

.PHONY: build
build: fmt build-flyway ## Build all openpitrix images
	mkdir -p ./tmp/bin
	$(call get_build_flags)
	$(RUN_IN_DOCKER) time go install -tags netgo -v -ldflags '$(BUILD_FLAG)' $(foreach cmd,$(CMDS),$(TRAG.Gopkg)/cmd/$(cmd))
ifneq ($(WITH_METADATA),no)
	$(RUN_IN_DOCKER) time go install -tags netgo -v -ldflags '$(BUILD_FLAG)' $(TRAG.Gopkg)/metadata/cmd/...
endif
	docker build -t $(TARG.Name) -t $(TARG.Name):metadata -f ./Dockerfile.dev ./tmp/bin
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-update
compose-update: build compose-up compose-migrate-db ## Update service in docker compose
	@echo "compose-update done"

.PHONY: compose-update-service-without-deps
compose-update-service-without-deps: build ## Update service in docker compose without dependence
	docker-compose up -d --no-dep $(COMPOSE_APP_SERVICES)
	@echo "compose-update-service-without-deps done"

.PHONY: compose-logs-f
compose-logs-f: ## Follow openpitrix log in docker compose
	docker-compose logs --tail 5 -f $(COMPOSE_APP_SERVICES)

.PHONY: compose-migrate-db
compose-migrate-db: ## Migrate db in docker compose
	until docker-compose exec openpitrix-db bash -c "cat /docker-entrypoint-initdb.d/*.sql | mysql -uroot -ppassword"; do echo "ddl waiting for mysql"; sleep 2; done;
	docker-compose up $(COMPOSE_DB_CTRL)

compose-update-%: ## Update "openpitrix-%" service in docker compose
	CMD=$* make build
	docker-compose up -d --no-deps openpitrix-$*
	@echo "compose-update done"

.PHONY: compose-put-global-config
compose-put-global-config: ## Put global config in docker compose
	@test -s deploy/config/global_config.yaml || { echo "[deploy/config/global_config.yaml] not exist"; exit 1; }
	cat deploy/config/global_config.yaml | docker run -i --rm openpitrix opctl validate_global_config
	cat deploy/config/global_config.yaml | docker-compose exec -T openpitrix-etcd /bin/sh -c "export ETCDCTL_API=3 && etcdctl put openpitrix/global_config"

.PHONY: generate-certs
generate-certs: ## Generate tls certificates
	cd ./deploy/kubernetes/tls-config && make

.PHONY: compose-up
compose-up: generate-certs ## Launch openpitrix in docker compose
	docker-compose up -d openpitrix-db
	until docker-compose exec openpitrix-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	make compose-migrate-db
	docker-compose up -d
	@echo "compose-up done"

.PHONY: compose-down
compose-down: ## Shutdown docker compose
	docker-compose down
	@echo "compose-down done"

release-%: ## Release version
	@if [ "`echo "$*" | grep -E "^openpitrix-v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	mkdir deploy/$*-kubernetes; \
	cp -r deploy/config deploy/kubernetes deploy/$*-kubernetes/; \
	cd deploy/ && tar -czvf $*-kubernetes.tar.gz $*-kubernetes; \
	cd ../; \
	mkdir deploy/$*-docker-compose; \
	cp -r deploy/docker-compose/. deploy/$*-docker-compose; \
	cp -r deploy/config/global_config.init.yaml deploy/$*-docker-compose/global_config.yaml; \
	cd deploy/ && tar -czvf $*-docker-compose.tar.gz $*-docker-compose; \
	fi

bin-release-%: ## Bin release version
	@if [ "`echo "$*" | grep -E "^openpitrix-v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	mkdir deploy/$*-bin; \
	docker cp openpitrix-api-gateway:/usr/local/bin/op deploy/$*-bin; \
	docker cp openpitrix-api-gateway:/usr/local/bin/opctl deploy/$*-bin; \
	docker cp openpitrix-api-gateway:/usr/local/bin/frontgate deploy/$*-bin; \
	docker cp openpitrix-api-gateway:/usr/local/bin/drone deploy/$*-bin; \
	cd deploy/ && tar -czvf $*-bin.tar.gz $*-bin; \
	fi

.PHONY: test
test: ## Run all tests
	make unit-test
	make e2e-test
	@echo "test done"


.PHONY: e2e-test
e2e-test: ## Run integration tests
	cd ./test/init/ && sh init_config.sh
	go test -v -a -tags="integration" ./test/...
ifeq ($(WITH_K8S),yes)
	go test -v -a -timeout 0 -tags="k8s" ./test/...
endif
	@echo "e2e-test done"

.PHONY: clean
clean: ## Clean generated version file
	-make -C ./pkg/version clean
	cd ./deploy/kubernetes/tls-config && make clean
	@echo "ok"

.PHONY: unit-test
unit-test: ## Run unit tests
	$(DB_TEST) $(ETCD_TEST) go test -v -a -tags="etcd db" ./...
	@echo "unit-test done"

build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t openpitrix/openpitrix:latest .; \
	docker build -t openpitrix/openpitrix:metadata -f ./Dockerfile.metadata .; \
	docker build -t openpitrix/openpitrix:flyway -f ./pkg/db/Dockerfile ./pkg/db/;\
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t openpitrix/openpitrix:$* .; \
	docker build -t openpitrix/openpitrix:metadata-$* -f ./Dockerfile.metadata .; \
	docker build -t openpitrix/openpitrix:flyway-$* -f ./pkg/db/Dockerfile ./pkg/db/; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push openpitrix/openpitrix:latest; \
	docker push openpitrix/openpitrix:metadata; \
	docker push openpitrix/openpitrix:flyway; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push openpitrix/openpitrix:$*; \
	docker push openpitrix/openpitrix:metadata-$*; \
	docker push openpitrix/openpitrix:flyway-$*; \
	fi