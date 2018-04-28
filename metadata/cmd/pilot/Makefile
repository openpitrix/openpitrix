# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG:=openpitrix-pilot

KILL_APP:=ps aux | grep $(TARG) | grep -v grep | awk '{print $$2}' | xargs kill -9
SHOW_APP:=ps aux | grep $(TARG) | grep -v grep

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  gen-config   to generate app config"
	@echo "  run          to run app server"
	@echo "  start        to start app server"
	@echo "  restart      to restart app server"
	@echo "  ping         to ping app server"
	@echo "  stop         to stop app server"
	@echo "  clean        to stop app server and clean the temp files"
	@echo "  help         to show this help informatrion"


.PHONY: gen-config
gen-config: $(TARG).out
	./$(TARG).out gen-config > pilot-config.json

.PHONY: run
run:
	go build -o $(TARG).out
	./$(TARG).out serve

.PHONY: start
start: $(TARG).out
	nohup ./$(TARG).out serve > $(TARG).out.log 2>&1 &

.PHONY: restart
restart: $(TARG).out stop
	nohup ./$(TARG).out serve > $(TARG).out.log 2>&1 &

.PHONY: stop
stop:
	-$(KILL_APP)

.PHONY: ps
ps:
	-$(SHOW_APP)

.PHONY: ping
ping: $(TARG).out
	./$(TARG).out ping
	./$(TARG).out ping -endpoint-type=pilot
	./$(TARG).out ping -endpoint-type=frontgate
	./$(TARG).out ping -endpoint-type=frontgate -frontgate-id=openpitrix-frontgate-001
	./$(TARG).out ping -endpoint-type=drone
	./$(TARG).out ping -endpoint-type=drone -drone-host=localhost -drone-port=9112


.PHONY: confd-status
confd-status: $(TARG).out
	-./$(TARG).out confd-status -frontgate-id=openpitrix-frontgate-001 -drone-host=localhost

.PHONY: confd-start
confd-start: $(TARG).out
	-./$(TARG).out confd-start -drone-host=localhost


.PHONY: confd-stop
confd-stop: $(TARG).out
	-./$(TARG).out confd-stop -drone-host=localhost


.PHONY: clean
clean: stop
	-rm $(TARG).out
	-rm $(TARG).out.log


$(TARG).out:
	@go fmt ./...
	@go vet ./...
	go build -o $@

