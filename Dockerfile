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

CMD ["sh"]
