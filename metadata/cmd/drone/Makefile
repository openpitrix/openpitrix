# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG:=openpitrix-drone

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
	./$(TARG).out gen-config > drone-config.json

.PHONY: run
run:
	go build -o $(TARG).out
	./$(TARG).out serve

.PHONY: start
start: $(TARG).out
	nohup ./$(TARG).out > $(TARG).out.log 2>&1 &

.PHONY: restart
restart: $(TARG).out stop
	nohup ./$(TARG).out > $(TARG).out.log 2>&1 &

.PHONY: stop
stop:
	-$(KILL_APP)

.PHONY: ps
ps:
	-$(SHOW_APP)

.PHONY: ping
ping: $(TARG).out
	-./$(TARG).out ping

.PHONY: list
list: $(TARG).out
	-./$(TARG).out list

.PHONY: confd-status
confd-status: $(TARG).out
	-./$(TARG).out confd-status

.PHONY: confd-start
confd-start: $(TARG).out
	-./$(TARG).out confd-start


.PHONY: confd-stop
confd-stop: $(TARG).out
	-./$(TARG).out confd-stop


etcd-play:
	ETCDCTL_API=3 etcdctl get /127.0.0.1/msg

	-./$(TARG).out getv /127.0.0.1/msg
	-./$(TARG).out getv /self/msg
	-./$(TARG).out getv /msg


.PHONY: clean
clean: stop
	-rm $(TARG).out
	-rm $(TARG).out.log


$(TARG).out:
	@go fmt ./...
	@go vet ./...
	go build -o $@

