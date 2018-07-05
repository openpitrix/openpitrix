# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

build:
	cd .. && docker build -t openpitrix/openpitrix:metadata -f ./Dockerfile.metadata .
	docker images openpitrix/openpitrix:metadata

pull:
	docker pull openpitrix/openpitrix:metadata

frontgate:
	docker run --rm -it openpitrix/openpitrix:metadata \
		frontgate -config=/opt/openpitrix/conf/frontgate.json serve

drone:
	docker run --rm -it openpitrix/openpitrix:metadata \
		drone -config=/opt/openpitrix/conf/drone.json serve

run-help:
	docker run --rm -it openpitrix/openpitrix:metadata frontgate -h
	docker run --rm -it openpitrix/openpitrix:metadata drone -h

clean:
	-docker image prune -f 1>/dev/null 2>&1
