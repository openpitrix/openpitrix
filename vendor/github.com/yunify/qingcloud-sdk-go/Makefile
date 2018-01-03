SHELL := /bin/bash

.PHONY: all check vet lint update generate build unit test release clean

PREFIX=qingcloud-sdk-go
VERSION=$(shell cat version.go | grep "Version\ =" | sed -e s/^.*\ //g | sed -e s/\"//g)
DIRS_TO_CHECK=$(shell ls -d */ | grep -vE "vendor|test")
PKGS_TO_CHECK=$(shell go list ./... | grep -vE "vendor|test")
PKGS_TO_RELEASE=$(shell go list ./... | grep -vE "/vendor/|/test")
FILES_TO_RELEASE=$(shell find . -name "*.go" | grep -vE "/vendor/|/test|.*_test.go")
FILES_TO_RELEASE_WITH_VENDOR=$(shell find . -name "*.go" | grep -vE "/test|.*_test.go")
LINT_IGNORE_DOC="service\/.*\.go:.+(comment on exported|should have comment or be unexported)"
LINT_IGNORE_CONFLICT="service\/.*\.go:.+(type name will be used as)"

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  all               to check, build, test and release this SDK"
	@echo "  check             to vet and lint the SDK"
	@echo "  generate          to generate service code"
	@echo "  build             to build the SDK"
	@echo "  unit              to run all sort of unit tests except runtime"
	@echo "  unit-test         to run unit test"
	@echo "  unit-benchmark    to run unit test with benchmark"
	@echo "  unit-coverage     to run unit test with coverage"
	@echo "  unit-race         to run unit test with race"
	@echo "  unit-runtime      to run test with go1.5, go1.6, go 1.7 in docker"
	@echo "  test              to run service test"
	@echo "  release           to build and release current version"
	@echo "  release-source    to pack the source code"
	@echo "  clean             to clean the coverage files"

all: check build unit release

check: vet lint

vet:
	@echo "go tool vet, skipping vendor packages"
	@go tool vet -all ${DIRS_TO_CHECK}
	@echo "ok"

lint:
	@echo "golint, skipping vendor packages"
	@lint=$$(for pkg in ${PKGS_TO_CHECK}; do golint $${pkg}; done); \
	 lint=$$(echo "$${lint}" | grep -vE -e ${LINT_IGNORE_DOC} -e ${LINT_IGNORE_CONFLICT}); \
	 if [[ -n $${lint} ]]; then echo "$${lint}"; exit 1; fi
	@echo "ok"

generate: snips ../qingcloud-api-specs/package.json
	./snips \
		-f=../qingcloud-api-specs/2013-08-30/swagger/api_v2.0.json \
		-t=./template \
		-o=./service
	go fmt ./service/...
	@echo "ok"

snips:
	curl -L https://github.com/yunify/snips/releases/download/v0.2.16/snips-v0.2.16-${shell go env GOOS}_amd64.tar.gz | tar zx

../qingcloud-api-specs/package.json:
	-go get -d github.com/yunify/qingcloud-api-specs
	file %@

build:
	@echo "build the SDK"
	GOOS=linux GOARCH=amd64 go build ${PKGS_TO_CHECK}
	GOOS=darwin GOARCH=amd64 go build ${PKGS_TO_CHECK}
	GOOS=windows GOARCH=amd64 go build ${PKGS_TO_CHECK}
	@echo "ok"


unit: unit-test unit-benchmark unit-coverage unit-race

unit-test:
	@echo "run unit test"
	go test -v ${PKGS_TO_CHECK}
	@echo "ok"

unit-benchmark:
	@echo "run unit test with benchmark"
	go test -v -bench=. ${PKGS_TO_CHECK}
	@echo "ok"

unit-coverage:
	@echo "run unit test with coverage"
	for pkg in ${PKGS_TO_CHECK}; do \
		output="coverage$${pkg#github.com/yunify/qingcloud-sdk-go}"; \
		mkdir -p $${output}; \
		go test -v -cover -coverprofile="$${output}/profile.out" $${pkg}; \
		if [[ -e "$${output}/profile.out" ]]; then \
			go tool cover -html="$${output}/profile.out" -o "$${output}/profile.html"; \
		fi; \
	done
	@echo "ok"

unit-race:
	@echo "run unit test with race"
	go test -v -race -cpu=1,2,4 ${PKGS_TO_CHECK}
	@echo "ok"

unit-runtime: unit-runtime-go-1.5 unit-runtime-go-1.6 unit-runtime-go-1.7

export define DOCKERFILE_GO_1_7
FROM golang:1.7

ADD . /go/src/github.com/yunify/qingcloud-sdk-go
WORKDIR /go/src/github.com/yunify/qingcloud-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.7:
	@echo "run test in go 1.7"
	echo "$${DOCKERFILE_GO_1_7}" > "dockerfile_go_1.7"
	docker build -f "./dockerfile_go_1.7" -t "${PREFIX}:go-1.7" .
	rm -f "./dockerfile_go_1.7"
	docker run --name "${PREFIX}-go-1.7-unit" -t "${PREFIX}:go-1.7"
	docker rm "${PREFIX}-go-1.7-unit"
	docker rmi "${PREFIX}:go-1.7"
	@echo "ok"

export define DOCKERFILE_GO_1_6
FROM golang:1.6

ADD . /go/src/github.com/yunify/qingcloud-sdk-go
WORKDIR /go/src/github.com/yunify/qingcloud-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.6:
	@echo "run test in go 1.6"
	echo "$${DOCKERFILE_GO_1_6}" > "dockerfile_go_1.6"
	docker build -f "./dockerfile_go_1.6" -t "${PREFIX}:go-1.6" .
	rm -f "./dockerfile_go_1.6"
	docker run --name "${PREFIX}-go-1.6-unit" -t "${PREFIX}:go-1.6"
	docker rm "${PREFIX}-go-1.6-unit"
	docker rmi "${PREFIX}:go-1.6"
	@echo "ok"

export define DOCKERFILE_GO_1_5
FROM golang:1.5
ENV GO15VENDOREXPERIMENT="1"

ADD . /go/src/github.com/yunify/qingcloud-sdk-go
WORKDIR /go/src/github.com/yunify/qingcloud-sdk-go

CMD ["make", "build", "unit"]
endef

unit-runtime-go-1.5:
	@echo "run test in go 1.5"
	echo "$${DOCKERFILE_GO_1_5}" > "dockerfile_go_1.5"
	docker build -f "./dockerfile_go_1.5" -t "${PREFIX}:go-1.5" .
	rm -f "./dockerfile_go_1.5"
	docker run --name "${PREFIX}-go-1.5-unit" -t "${PREFIX}:go-1.5"
	docker rm "${PREFIX}-go-1.5-unit"
	docker rmi "${PREFIX}:go-1.5"
	@echo "ok"

test:
	pushd "./test"; go test ; popd
	@echo "ok"

release: release-source release-source-with-vendor

release-source:
	@echo "pack the source code"
	mkdir -p "release"
	@zip -FS "release/${PREFIX}-source-v${VERSION}.zip" ${FILES_TO_RELEASE}
	@echo "ok"

release-source-with-vendor:
	@echo "pack the source code"
	mkdir -p "release"
	@zip -FS "release/${PREFIX}-source-with-vendor-v${VERSION}.zip" ${FILES_TO_RELEASE_WITH_VENDOR}
	@echo "ok"

clean:
	rm -rf $${PWD}/coverage
	@echo "ok"
