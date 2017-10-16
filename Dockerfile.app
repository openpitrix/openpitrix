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

FROM golang:1.9-alpine
WORKDIR /go/src/openpitrix.io/openpitrix/ 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-server ./cmd/app

FROM alpine:latest
COPY --from=0 /go/src/openpitrix.io/openpitrix/app-server /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app-server"]
