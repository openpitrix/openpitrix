# +-------------------------------------------------------------------------
# | Copyright (C) 2017 Yunify, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

TARG.Name:=openpitrix
TRAG.Gopkg:=openpitrix.io/openpitrix

TARG.Services:=api
TARG.Services+=app
TARG.Services+=repo
TARG.Services+=runtime

DOCKER_TAGS=latest

GO:=docker run --rm -it -v `pwd`:/go/src/$(TRAG.Gopkg) -w /go/src/$(TRAG.Gopkg) golang:1.9-alpine go
GO_WITH_MYSQL:=docker run --rm -it -v `pwd`:/go -w /go/src/$(TRAG.Gopkg) --link openpitrix-mysql:mysql -p 9527:9527 golang:1.9-alpine go
SWAGGER:=docker run --rm -it -v `pwd`:/go/src/$(TRAG.Gopkg) -w /go/src/$(TRAG.Gopkg) quay.io/goswagger/swagger

SWAGGER_SPEC_DIR:=./api
SWAGGER_SPEC_FILE:=./api/_all.json
SWAGGER_OUT_DIR:=./pkg/swagger

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
	docker pull quay.io/goswagger/swagger
	docker pull vidsyhq/multi-file-swagger-docker
	docker pull mysql
	@echo "ok"

.PHONY: generate
generate:
	cd ./api && make generate
	-mkdir -p $(SWAGGER_OUT_DIR)
	$(SWAGGER) generate server -f $(SWAGGER_SPEC_FILE) -t $(SWAGGER_OUT_DIR)
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
	@for service in $(TARG.Services) ; do \
		docker build -t $$service:$(DOCKER_TAGS) -f Dockerfile.$$service . ; \
	done
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
