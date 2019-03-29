# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

COMPOSE_DB_CTRL=openpitrix-db-init openpitrix-im-db-init openpitrix-am-db-init openpitrix-notification-db-init openpitrix-app-db-ctrl openpitrix-repo-db-ctrl openpitrix-runtime-db-ctrl openpitrix-job-db-ctrl openpitrix-task-db-ctrl openpitrix-cluster-db-ctrl openpitrix-iam-db-ctrl openpitrix-attachment-db-ctrl openpitrix-isv-db-ctrl openpitrix-notification-db-ctrl openpitrix-im-db-ctrl openpitrix-am-db-ctrl

default: pilot-ca compose-up

.PHONY: pilot-ca
pilot-ca:
	cd pilot && make

.PHONY: clean
clean: compose-clean

.PHONY: compose-migrate-db
compose-migrate-db:
	until docker-compose exec openpitrix-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	docker-compose up $(COMPOSE_DB_CTRL)

.PHONY: compose-up
compose-up: pull
	docker-compose up -d openpitrix-db
	make compose-migrate-db
	docker-compose up --remove-orphans -d
	@echo "compose-up done"

.PHONY: pull
pull:
	docker-compose pull

.PHONY: compose-clean
compose-clean:
	docker-compose down

.PHONY: compose-put-global-config
compose-put-global-config:
	@test -s ./global_config.yaml || { echo "[./global_config.yaml] not exist"; exit 1; }
	cat ./global_config.yaml | docker run -i --rm openpitrix/openpitrix:latest opctl validate_global_config
	cat ./global_config.yaml | docker-compose exec -T openpitrix-etcd /bin/sh -c "export ETCDCTL_API=3 && etcdctl put openpitrix/global_config"

.PHONY: compose-get-global-config
compose-get-global-config:
	docker-compose exec -T openpitrix-etcd /bin/sh -c "export ETCDCTL_API=3 && etcdctl get openpitrix/global_config --print-value-only" > deploy/config/global_config.yaml
