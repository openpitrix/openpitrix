# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default: pilot-ca compose-up

.PHONY: pilot-ca
pilot-ca:
	cd pilot && make

.PHONY: clean
clean: compose-clean

.PHONY: compose-up
compose-up: pull
	docker-compose up -d openpitrix-db openpitrix-etcd
	until docker-compose exec openpitrix-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	docker-compose up openpitrix-db-init
	docker-compose up openpitrix-app-db-ctrl openpitrix-repo-db-ctrl openpitrix-runtime-db-ctrl openpitrix-job-db-ctrl openpitrix-task-db-ctrl openpitrix-cluster-db-ctrl openpitrix-iam-db-ctrl
	docker-compose up -d

.PHONY: pull
pull:
	docker pull openpitrix/openpitrix:flyway
	docker pull openpitrix/openpitrix:latest
	docker-compose pull

.PHONY: compose-clean
compose-clean:
	docker-compose down

.PHONY: compose-put-global-config
compose-put-global-config:
	@test -s ./global_config.yaml || { echo "[./global_config.yaml] not exist"; exit 1; }
	cat ./global_config.yaml | docker run -i --rm openpitrix/openpitrix:latest opctl validate_global_config
	cat ./global_config.yaml | docker-compose exec -T openpitrix-etcd /bin/sh -c "export ETCDCTL_API=3 && etcdctl put openpitrix/global_config"
