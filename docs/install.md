# To start developing OpenPitrix

The [community repository](https://github.com/openpitrix) hosts all information about
building OpenPitrix from source, how to contribute code
and documentation, who to contact about what, etc.

## Develop OpenPitrix

If you want to build OpenPitrix right away there are three options:

##### You have a working [Go environment].

```
$ go get -d openpitrix.io/openpitrix
$ cd $GOPATH/src/openpitrix.io/openpitrix
$ GOBIN=`pwd`/bin go install ./cmd/...
$ docker run --rm --name openpitrix-mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=openpitrix -p 3306:3306 -d mysql:5.6
$ ./bin/api &
$ ./bin/repo &
$ ./bin/app &
$ ./bin/runtime &
$ ./bin/cluster &
```

##### You have a working [Docker environment].

```
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ docker run --rm --name openpitrix-mysql \
    -e MYSQL_ROOT_PASSWORD=password \
    -e MYSQL_DATABASE=openpitrix \
    -p 3306:3306 -d mysql:5.6
$ docker run --rm --name openpitrix-app \
    --link openpitrix-mysql:openpitrix-mysql \
    -d openpitrix app
$ docker run --rm --name openpitrix-runtime \
    --link openpitrix-mysql:openpitrix-mysql \
    -d openpitrix runtime
$ docker run --rm --name openpitrix-cluster \
    --link openpitrix-mysql:openpitrix-mysql \
    -d openpitrix cluster
$ docker run --rm --name openpitrix-repo \
    --link openpitrix-mysql:openpitrix-mysql \
    -d openpitrix repo
$ docker run --rm --name openpitrix-api \
    --link openpitrix-app:openpitrix-app \
    --link openpitrix-runtime:openpitrix-runtime \
    --link openpitrix-cluster:openpitrix-cluster \
    --link openpitrix-repo:openpitrix-repo \
    -p 8080:8080 -d openpitrix api
```

Exit docker runtime environment
```
$ docker kill $(docker ps -f name=openpitrix -q -a)
```

##### You have a working [Docker-Compose environment].

```
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ docker-compose up -d
```

Exit docker runtime environment
```
$ docker-compose down
```

## Test OpenPitrix

Visit http://127.0.0.1:8080/swagger-ui in browser, and try it online.

Or test openpitrix/api service in command line:

```
$ curl http://localhost:8080/v1/apps
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:8080/v1/apps/app-12345678
{"error":"App Id app-12345678 not exist","code":5}
$ curl http://localhost:8080/v1/appruntimes
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:8080/v1/clusters
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
$ curl http://localhost:8080/v1/repos
{"total_items":0,"total_pages":0,"page_size":10,"current_page":1}
```
