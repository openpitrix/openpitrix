# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:alpine as builder

RUN apk add --no-cache git
RUN go get openpitrix.io/openpitrix/cmd/runtime

FROM alpine

COPY --from=builder /go/bin/runtime /usr/local/bin/

# database
ENV OPENPITRIX_DATABASE_DBNAME=openpitrix
ENV OPENPITRIX_DATABASE_ENCODING=utf8
ENV OPENPITRIX_DATABASE_ENGINE=InnoDB
ENV OPENPITRIX_DATABASE_HOST=openpitrix-mysql
ENV OPENPITRIX_DATABASE_PORT=3306
ENV OPENPITRIX_DATABASE_ROOTPASSWORD=password
ENV OPENPITRIX_DATABASE_TYPE=mysql

# api service
ENV OPENPITRIX_APISERVICE_HOST=openpitrix-api
ENV OPENPITRIX_APISERVICE_PORT=8080

# repo service
ENV OPENPITRIX_REPOSERVICE_HOST=openpitrix-repo
ENV OPENPITRIX_REPOSERVICE_PORT=8081

# app service
ENV OPENPITRIX_APPSERVICE_HOST=openpitrix-app
ENV OPENPITRIX_APPSERVICE_PORT=8082

# app runtime service
ENV OPENPITRIX_APPRUNTIMESERVICE_HOST=openpitrix-runtime
ENV OPENPITRIX_APPRUNTIMESERVICE_PORT=8083

# cluster service
ENV OPENPITRIX_CLUSTERSERVICE_HOST=openpitrix-cluster
ENV OPENPITRIX_CLUSTERSERVICE_PORT=8084

ENTRYPOINT ["/usr/local/bin/runtime"]
CMD []
