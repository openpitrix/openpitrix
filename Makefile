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

TARG:=openpitrix

GO:=docker run --rm -it -v $(shell go env GOPATH):/go -w /go/src/$(TARG) golang:1.9-alpine go
GO_WITH_MYSQL:=docker run --rm -it -v $(shell go env GOPATH):/go -w /go/src/$(TARG) --link openpitrix-mysql:mysql -p 9527:9527 golang:1.9-alpine go
SWAGGER:=docker run --rm -it -v $(shell go env GOPATH):/go -w /go/src/$(TARG) quay.io/goswagger/swagger

SWAGGER_SPEC_FILE:=./src/api/swagger-spec/_all.json
SWAGGER_OUT_DIR:=./src/api/swagger

MYSQL_DATABASE:=$(TARG)
MYSQL_ROOT_PASSWORD:=password

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

all: generate test release

init-vendor:
	govendor init
	govendor add +external
	@echo "ok"

update-vendor:
	govendor update +external
	govendor list
	@echo "ok"

tools:
	go get github.com/kardianos/govendor
	docker pull golang:1.9-alpine
	docker pull quay.io/goswagger/swagger
	docker pull vidsyhq/multi-file-swagger-docker
	docker pull mysql
	@echo "ok"

generate:
	cd ./src/api/swagger-spec && make generate
	-mkdir -p $(SWAGGER_OUT_DIR)
	$(SWAGGER) generate server -f $(SWAGGER_SPEC_FILE) -t $(SWAGGER_OUT_DIR)
	@echo "ok"

mysql-start:
	docker run --rm --name openpitrix-mysql -e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) -e MYSQL_DATABASE=$(MYSQL_DATABASE) -d mysql || docker start openpitrix-mysql
	@echo "ok"

mysql-stop:
	docker stop openpitrix-mysql
	@echo "ok"

run: mysql-start
	$(GO_WITH_MYSQL) run ./src/cmd/openpitrix-server/main.go

build:
	@echo "TODO"

release:
	@echo "TODO"

test:
	$(GO) fmt ./...
	$(GO) test ./...
	@echo "ok"

clean:
	@echo "ok"
