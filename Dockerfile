# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.13-alpine as builder

RUN apk add --no-cache git curl openssl

WORKDIR /go/src/openpitrix.io/openpitrix/
COPY . .

RUN mkdir -p /openpitrix_bin
RUN go generate openpitrix.io/openpitrix/pkg/version && \
	CGO_ENABLED=0 GOBIN=/openpitrix_bin go install -ldflags '-w -s' -v -tags netgo openpitrix.io/openpitrix/cmd/hyperpitrix

FROM alpine:3.7

RUN apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing gops
RUN apk add --no-cache curl wget

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /openpitrix_bin/* /usr/local/bin/

RUN apk add --update ca-certificates && \
    update-ca-certificates && \
    adduser -D -g openpitrix -u 1002 openpitrix && \
    chown -R openpitrix:openpitrix /usr/local/bin/

USER openpitrix

CMD ["sh"]