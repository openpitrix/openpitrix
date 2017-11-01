# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:alpine as builder

WORKDIR /go/src/openpitrix.io/openpitrix/
COPY . .
RUN go install ./cmd/...

FROM alpine

COPY --from=builder /go/bin/api /usr/local/bin/
COPY --from=builder /go/bin/app /usr/local/bin/
COPY --from=builder /go/bin/repo /usr/local/bin/
COPY --from=builder /go/bin/runtime /usr/local/bin/

ENV OPENPITRIX_DATABASE_DBNAME=openpitrix
ENV OPENPITRIX_DATABASE_ENCODING=utf8
ENV OPENPITRIX_DATABASE_ENGINE=InnoDB
ENV OPENPITRIX_DATABASE_HOST=openpitrix-mysql
ENV OPENPITRIX_DATABASE_PORT=3306
ENV OPENPITRIX_DATABASE_ROOTPASSWORD=password
ENV OPENPITRIX_DATABASE_TYPE=mysql
ENV OPENPITRIX_HOST=0.0.0.0
ENV OPENPITRIX_LOGLEVEL=warn
ENV OPENPITRIX_PORT=8080

CMD ["sh"]
