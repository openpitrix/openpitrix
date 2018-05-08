# Developing for OpenPitrix

The [community repository](https://github.com/openpitrix) hosts all information about
building OpenPitrix from source, how to contribute code and documentation, who to contact about what, etc. If you find a requirement that this doc does not capture, or if you find other docs with references to requirements that are not simply links to this doc, please [submit an issue](https://github.com/openpitrix/openpitrix/issues/new).

----

## To start developing OpenPitrix

First of all, you should fork the project. Then follow one of the three options below to develop the project. Please note you should replace the official repo when using __go get__ or __git clone__ below with your own one.

### 1. You have a working [Docker Compose](https://docs.docker.com/compose/install) environment [recommend].
>You need to install [Docker](https://docs.docker.com/engine/installation/) first.

```shell
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ make compose-up
```

Exit docker runtime environment
```shell
$ make compose-down
```

### 2. You have a working [Docker](https://docs.docker.com/engine/installation/) environment.

```shell
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ docker network create -d bridge openpitrix-bridge
$ docker run --rm --name openpitrix-db -e MYSQL_ROOT_PASSWORD=password -v `pwd`/pkg/db/ddl:/docker-entrypoint-initdb.d  \
	--network openpitrix-bridge -p 13306:3306 -d mysql:5.7
$ docker run --rm --name openpitrix-etcd --network openpitrix-bridge -p 12379:2379 -d appcelerator/etcd
$ docker run --rm --name openpitrix-repo-db-ctrl -v `pwd`/pkg/db/schema/repo:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/repo -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-app-db-ctrl -v `pwd`/pkg/db/schema/app:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/app -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-runtime-db-ctrl -v `pwd`/pkg/db/schema/runtime:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/runtime -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-job-db-ctrl -v `pwd`/pkg/db/schema/job:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/cluster -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-task-db-ctrl -v `pwd`/pkg/db/schema/task:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/cluster -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-cluster-db-ctrl -v `pwd`/pkg/db/schema/cluster:/flyway/sql \
	--network openpitrix-bridge boxfuse/flyway:5.0.7-alpine \
	-url=jdbc:mysql://openpitrix-db/cluster -user=root -password=password -validateOnMigrate=false migrate
$ docker run --rm --name openpitrix-app-manager -e OPENPITRIX_MYSQL_DATABASE=app \
	--network openpitrix-bridge -d openpitrix app-manager
$ docker run --rm --name openpitrix-runtime-manager -e OPENPITRIX_MYSQL_DATABASE=runtime \
	--network openpitrix-bridge -d openpitrix runtime-manager
$ docker run --rm --name openpitrix-repo-indexer -e OPENPITRIX_MYSQL_DATABASE=repo \
	--network openpitrix-bridge -d openpitrix repo-indexer
$ docker run --rm --name openpitrix-repo-manager -e OPENPITRIX_MYSQL_DATABASE=repo \
	--network openpitrix-bridge -d openpitrix repo-manager
$ docker run --rm --name openpitrix-job-manager -e OPENPITRIX_MYSQL_DATABASE=cluster  \
	--network openpitrix-bridge -d openpitrix job-manager
$ docker run --rm --name openpitrix-task-manager -e OPENPITRIX_MYSQL_DATABASE=cluster \
	--network openpitrix-bridge -d openpitrix task-manager
$ docker run --rm --name openpitrix-cluster-manager -e OPENPITRIX_MYSQL_DATABASE=cluster \
	--network openpitrix-bridge -d openpitrix cluster-manager
$ docker run --rm --name openpitrix-pilot-service \
	--network openpitrix-bridge -d openpitrix pilot-service
$ docker run --rm --name openpitrix-api-gateway --network openpitrix-bridge -p 9100:9100 -d openpitrix api-gateway
```

Exit docker runtime environment
```shell
$ docker stop $(docker ps -f name=openpitrix -q)
```

### 3. You have a working [Go](prereqs.md#setting-up-go) environment.

- Install [protoc compiler](https://github.com/google/protobuf/releases/)
- Install protoc plugin:

```shell
$ go get github.com/golang/protobuf/protoc-gen-go
$ go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
$ go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
$ go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

- Get openpitrix source code and build service:

```shell
$ go get -d openpitrix.io/openpitrix
$ cd $GOPATH/src/openpitrix.io/openpitrix
$ make generate
$ GOBIN=`pwd`/bin go install ./cmd/...
```

- Install mysql server first. Then add the services name to the `/etc/hosts` file as follows.
>Note: If you install mysql server remotely then configure the server IP correspondingly. You may
need to create the database __openpitrix__ and change the user __root__ password to __password__ in advance. If the user __root__ password is different than the default one, then you need to specify the password in the command line when start OpenPitrix services.

```
127.0.0.1 openpitrix-api
127.0.0.1 openpitrix-repo
127.0.0.1 openpitrix-app
127.0.0.1 openpitrix-runtime
127.0.0.1 openpitrix-cluster
127.0.0.1 openpitrix-db
```

- Start OpenPitrix service:

```shell
$ ./bin/openpitrix-api &
$ ./bin/openpitrix-repo &
$ ./bin/openpitrix-app &
$ ./bin/openpitrix-runtime &
$ ./bin/openpitrix-cluster &
```

- Exit go runtime environment
```shell
$ ps aux | grep openpitrix- | grep -v grep | awk '{print $2}' | xargs kill -9
```

----

## Test OpenPitrix

Visit http://127.0.0.1:9100/swagger-ui in browser, and try it online, or test openpitrix api service via command line:

```shell
$ curl http://localhost:9100/v1/apps
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:9100/v1/apps/app-12345678
{"error":"App Id app-12345678 not exist","code":5}
$ curl http://localhost:9100/v1/appruntimes
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:9100/v1/clusters
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:9100/v1/repos
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
```

----

## DevOps

Please check [How to set up DevOps environment](devops.md).
