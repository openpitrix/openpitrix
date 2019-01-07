# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM openpitrix/openpitrix-builder as builder

WORKDIR /go/src/openpitrix.io/openpitrix/
COPY . .

RUN apk add --no-cache util-linux

RUN mkdir -p /metadata_bin
RUN go generate openpitrix.io/openpitrix/pkg/version && \
	CGO_ENABLED=0 GOOS=linux GOBIN=/metadata_bin go install -ldflags '-w -s' -tags netgo openpitrix.io/openpitrix/metadata/cmd/...

RUN find /metadata_bin -type f -exec upx {} \;


FROM golang:1.10.1-alpine3.7 as metad-builder

RUN apk add --no-cache upx git

WORKDIR /go/src/openpitrix.io/
RUN git clone https://github.com/openpitrix/metad.git && cd metad

RUN mkdir -p /metad_bin
RUN go generate openpitrix.io/metad/pkg/version && \
CGO_ENABLED=0 GOOS=linux GOBIN=/metad_bin go install -ldflags '-w -s' -tags netgo openpitrix.io/metad

RUN find /metad_bin -type f -exec upx {} \;


FROM alpine:3.7
RUN apk add --update ca-certificates && update-ca-certificates

COPY --from=builder /usr/bin/nsenter        /usr/bin/

COPY --from=builder /metadata_bin/drone     /usr/local/bin/
COPY --from=builder /metadata_bin/frontgate /usr/local/bin/

COPY --from=metad-builder /metad_bin/* /usr/local/bin/

RUN apk add --no-cache supervisor
COPY metadata/supervisord/supervisord.conf /etc/
COPY metadata/supervisord/start-supervisord.sh /usr/local/bin/
COPY metadata/supervisord/frontgate/start-frontgate.sh /usr/local/bin/
COPY metadata/supervisord/drone/start-drone.sh /usr/local/bin/
COPY metadata/supervisord/metad/start-metad.sh /usr/local/bin/

RUN mkdir -p /etc/supervisor.d
COPY metadata/supervisord/frontgate/frontgate.ini /etc/supervisor.d
COPY metadata/supervisord/drone/drone.ini /etc/supervisor.d
COPY metadata/supervisord/metad/metad.ini /etc/supervisor.d

ENTRYPOINT ["start-supervisord.sh"]
