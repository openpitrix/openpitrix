# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default:
	docker build -t openpitrix/openpitrix-builder .
	@echo "ok"

pull:
	docker pull openpitrix/openpitrix-builder
	@echo "ok"

run:
	docker run --rm -it -v `pwd`:/root openpitrix/openpitrix-builder

clean:
	@echo "ok"
