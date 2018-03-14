# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM openpitrix/openpitrix-builder as builder

WORKDIR /go/src/openpitrix.io/openpitrix/
COPY . .

RUN go generate openpitrix.io/openpitrix/pkg/version && \
	go install  openpitrix.io/openpitrix/cmd/...

FROM alpine:3.6
RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /go/bin/* /usr/local/bin/
COPY ./pkg/db/schema /schema
COPY ./pkg/db/ddl /ddl


CMD ["sh"]
