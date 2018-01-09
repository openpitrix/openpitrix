# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM openpitrix/openpitrix-builder as builder

WORKDIR /go/src/openpitrix.io/openpitrix/
COPY . .

RUN go generate openpitrix.io/openpitrix/pkg/version && \
	go install  openpitrix.io/openpitrix/cmd/...     && \
	mv /go/bin/openpitrix-api     /go/bin/api        && \
	mv /go/bin/openpitrix-app     /go/bin/app        && \
	mv /go/bin/openpitrix-cluster /go/bin/cluster    && \
	mv /go/bin/openpitrix-repo    /go/bin/repo       && \
	mv /go/bin/openpitrix-runtime /go/bin/runtime


FROM alpine:3.6

# Glog
ENV OPENPITRIX_CONFIG_GLOG_LOGTOSTDERR=
ENV OPENPITRIX_CONFIG_GLOG_ALSOLOGTOSTDERR=
ENV OPENPITRIX_CONFIG_GLOG_STDERRTHRESHOLD=
ENV OPENPITRIX_CONFIG_GLOG_LOGDIR=

ENV OPENPITRIX_CONFIG_GLOG_LOGBACKTRACEAT=
ENV OPENPITRIX_CONFIG_GLOG_V=
ENV OPENPITRIX_CONFIG_GLOG_VMODULE=

ENV OPENPITRIX_CONFIG_GLOG_COPYSTANDARDLOGTO=

# database
ENV OPENPITRIX_CONFIG_DB_TYPE=
ENV OPENPITRIX_CONFIG_DB_DBNAME=
ENV OPENPITRIX_CONFIG_DB_ENCODING=
ENV OPENPITRIX_CONFIG_DB_ENGINE=
ENV OPENPITRIX_CONFIG_DB_HOST=
ENV OPENPITRIX_CONFIG_DB_PORT=
ENV OPENPITRIX_CONFIG_DB_ROOTNAME=
ENV OPENPITRIX_CONFIG_DB_ROOTPASSWORD=
ENV OPENPITRIX_CONFIG_DB_USERNAME=
ENV OPENPITRIX_CONFIG_DB_USERPASSWORD=

# api service
ENV OPENPITRIX_CONFIG_API_HOST=
ENV OPENPITRIX_CONFIG_API_PORT=

# app service
ENV OPENPITRIX_CONFIG_APP_HOST=
ENV OPENPITRIX_CONFIG_APP_PORT=

# runtime service
ENV OPENPITRIX_CONFIG_RUNTIME_HOST=
ENV OPENPITRIX_CONFIG_RUNTIME_PORT=

# cluster service
ENV OPENPITRIX_CONFIG_CLUSTER_HOST=
ENV OPENPITRIX_CONFIG_CLUSTER_PORT=

# repo service
ENV OPENPITRIX_CONFIG_REPO_HOST=
ENV OPENPITRIX_CONFIG_REPO_PORT=

COPY --from=builder /go/bin/* /usr/local/bin/

CMD ["sh"]
