# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:alpine as builder

RUN apk add --no-cache git
RUN go get openpitrix.io/openpitrix/cmd/repo

FROM alpine

COPY --from=builder /go/bin/repo /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/repo"]
CMD []
