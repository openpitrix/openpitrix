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

TARG:=apphub

GO:=docker run --rm -it -v $(shell go env GOPATH):/go -w /go/src/$(TARG) golang:1.9-alpine go
SWAGGER:=docker run --rm -it -v $(shell go env GOPATH):/go -w /go/src/$(TARG) quay.io/goswagger/swagger

SWAGGER_SPEC_FILE:=./src/api/swagger-spec/_all.json
SWAGGER_OUT_DIR:=./src/api/swagger

default: generate test

validate:
	$(SWAGGER) validate $(SWAGGER_SPEC_FILE)

init-vendor:
	govendor init
	govendor add +external

update-vendor:
	govendor update +external
	govendor list

tools:
	go get github.com/kardianos/govendor
	docker pull golang:1.9-alpine
	docker pull quay.io/goswagger/swagger
	docker pull vidsyhq/multi-file-swagger-docker

generate:
	cd ./src/api/swagger-spec && make
	-mkdir -p $(SWAGGER_OUT_DIR)
	$(SWAGGER) generate server -f $(SWAGGER_SPEC_FILE) -t $(SWAGGER_OUT_DIR)

fmt:
	$(GO) fmt ./...

test:
	$(GO) fmt ./...
	$(GO) test ./...

clean:
