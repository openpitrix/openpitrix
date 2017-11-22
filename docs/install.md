## To start developing OpenPitrix

The [community repository](https://github.com/openpitrix) hosts all information about
building OpenPitrix from source, how to contribute code
and documentation, who to contact about what, etc.

If you want to build OpenPitrix right away there are two options:

##### You have a working [Go environment].

```
$ go get -d openpitrix.io/openpitrix
$ cd $GOPATH/src/openpitrix.io/openpitrix
$ GOBIN=`pwd`/bin go install ./cmd/...
$ docker run --rm --name openpitrix-mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=openpitrix -p 3306:3306 -d mysql
$ ./bin/api &
$ ./bin/repo &
$ ./bin/app &
$ ./bin/runtime &
$ ./bin/cluster &
```

Visit http://127.0.0.1:8080/swagger-ui in browser, the try in online.

Or test openpitrix/api service in command line:

```
$ curl http://localhost:8080/v1/apps
{"items":null}
$ curl http://localhost:8080/v1/apps/app-12345678
{"code":"500","message":"sql: no rows in result set"}
```

##### You have a working [Docker environment].

```
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ docker run --rm --name openpitrix-mysql \
    -e MYSQL_ROOT_PASSWORD=password \
    -e MYSQL_DATABASE=openpitrix \
    -p 3306:3306 -d mysql
$ docker run --rm --name openpitrix-app \
    --link openpitrix-mysql:openpitrix-mysql \
    -p 8081:8081 -d openpitrix app
$ docker run --rm --name openpitrix-runtime \
    --link openpitrix-mysql:openpitrix-mysql \
    -p 8082:8082 -d openpitrix runtime
$ docker run --rm --name openpitrix-cluster \
    --link openpitrix-mysql:openpitrix-mysql \
    -p 8083:8083 -d openpitrix cluster
$ docker run --rm --name openpitrix-repo \
    --link openpitrix-mysql:openpitrix-mysql \
    -p 8084:8084 -d openpitrix repo
$ docker run --rm --name openpitrix-api \
    --link openpitrix-app:openpitrix-app \
    --link openpitrix-runtime:openpitrix-runtime \
    --link openpitrix-cluster:openpitrix-cluster \
    --link openpitrix-repo:openpitrix-repo \
    -p 8080:8080 -d openpitrix api
```

Visit http://127.0.0.1:8080/swagger-ui in browser, the try in online.

Or test openpitrix/api service in command line:

```
$ curl http://localhost:8080/v1/apps
{"items":null}
$ curl http://localhost:8080/v1/apps/app-12345678
{"code":"500","message":"sql: no rows in result set"}
```

##### You have a working [Docker-Compose environment].

```
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ docker-compose up -d
$ curl http://localhost:8080/v1/apps
$ curl http://localhost:8080/v1/appruntimes
$ curl http://localhost:8080/v1/clusters
$ curl http://localhost:8080/v1/repos
$ docker-compose down
```
