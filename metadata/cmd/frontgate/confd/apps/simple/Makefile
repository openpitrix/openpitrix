# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG:=openpitrix-simple-app

KILL_APP:=ps aux | grep $(TARG) | grep -v grep | awk '{print $$2}' | xargs kill -9
SHOW_APP:=ps aux | grep $(TARG) | grep -v grep

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check-config to check app config"
	@echo "  check-server to check app server"
	@echo "  run          to run app server"
	@echo "  start        to start app server"
	@echo "  restart      to restart app server"
	@echo "  stop         to stop app server"
	@echo "  clean        to stop app server and clean the temp files"
	@echo "  help         to show this help informatrion"


.PHONY: check-config
check-config: $(TARG).out
	./$(TARG).out -check-config

.PHONY: check-server
check-server: $(TARG).out
	./$(TARG).out -check-server

.PHONY: run
run: $(TARG).out
	./$(TARG).out

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


.PHONY: etcd-put-key
etcd-put-key:
	ETCDCTL_API=3 etcdctl put /listen_port 3999
	ETCDCTL_API=3 etcdctl put /key key-value

.PHONY: clean
clean: stop
	-rm $(TARG).out
	-rm $(TARG).out.log


$(TARG).out: main.go
	@go fmt ./...
	@go vet ./...
	go build -o $@
