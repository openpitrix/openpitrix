# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG:=openpitrix-frontgate

EXEC_APP:=nohup ./$(TARG).out > $(TARG).out.log 2>&1 &
KILL_APP:=ps aux | grep $(TARG) | grep -v grep | awk '{print $$2}' | xargs kill -9
SHOW_APP:=ps aux | grep $(TARG) | grep -v grep

EXEC_ETCD:=nohup etcd > etcd.out.log 2>&1 &
KILL_ETCD:=ps aux | grep etcd | grep -v grep | awk '{print $$2}' | xargs kill -9

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  gen-config   to generate app config"
	@echo "  etcd-start   to start etcd server"
	@echo "  etcd-stop    to stop etcd server"
	@echo "  run          to run app server"
	@echo "  start        to start app server"
	@echo "  restart      to restart app server"
	@echo "  ping         to ping app server"
	@echo "  confd-start  to start confd on drone server"
	@echo "  confd-stop   to stop confd on drone server"
	@echo "  stop         to stop app server"
	@echo "  clean        to stop app server and clean the temp files"
	@echo "  help         to show this help informatrion"


.PHONY: gen-config
gen-config: $(TARG).out
	./$(TARG).out gen-config > frontgate-config.json

etcd-start:
	$(EXEC_ETCD)
etcd-stop:
	-$(KILL_ETCD)

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

getv:
	-./$(TARG).out getv /abc

.PHONY: ping
ping: $(TARG).out
	-./$(TARG).out ping
	-./$(TARG).out ping -endpoint-type=frontgate
	-./$(TARG).out ping -endpoint-type=pilot
	-./$(TARG).out ping -endpoint-type=drone
	-./$(TARG).out ping -endpoint-type=drone -drone-host=localhost -drone-port=9112


.PHONY: confd-status
confd-status: $(TARG).out
	-./$(TARG).out confd-status -drone-host=localhost

.PHONY: confd-start
confd-start: $(TARG).out
	-./$(TARG).out confd-start -drone-host=localhost


.PHONY: confd-stop
confd-stop: $(TARG).out
	-./$(TARG).out confd-stop -drone-host=localhost


.PHONY: etcd-play
etcd-play:
	ETCDCTL_API=3 etcdctl get --prefix ""

	ETCDCTL_API=3 etcdctl put abc abc-value
	ETCDCTL_API=3 etcdctl put abc/aaa abc/aaa-value

	ETCDCTL_API=3 etcdctl get abc
	ETCDCTL_API=3 etcdctl get --prefix abc
	ETCDCTL_API=3 etcdctl get --prefix ""

	ETCDCTL_API=3 etcdctl put /127.0.0.1/msg localhost-msg
	ETCDCTL_API=3 etcdctl get --prefix ""
	ETCDCTL_API=3 etcdctl get /127.0.0.1/msg


.PHONY: clean
clean: stop
	-rm $(TARG).out
	-rm $(TARG).out.log

$(TARG).out:
	@go fmt ./...
	@go vet ./...
	go build -o $@


