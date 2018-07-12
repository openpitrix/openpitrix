# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default:
	cd ../../pilot/tls-config && make

	cp ../../pilot/tls-config/openpitrix-ca.crt .

	cp ../../pilot/tls-config/pilot-client.key .
	cp ../../pilot/tls-config/pilot-client.crt .

clean:
	-rm *.key *.crt *.csr *.srl

