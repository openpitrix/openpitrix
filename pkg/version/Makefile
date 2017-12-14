# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

GO:=docker run --rm -it -v `pwd`:/go/src/$(TRAG.Gopkg) -w /go/src/$(TRAG.Gopkg) openpitrix/openpitrix:builder go

default:
	$(GO) generate

clean:
	-rm z_update_version.go
