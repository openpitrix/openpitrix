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
$ ./bin/api
```

Test openpitrix/api service:

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
$ docker run --rm --name openpitrix-mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=openpitrix -d mysql
$ docker run --rm -it --link openpitrix-mysql:openpitrix-mysql -p 8080:8080 openpitrix api
```

Test openpitrix/api service:

```
$ curl http://localhost:8080/v1/apps
{"items":null}
$ curl http://localhost:8080/v1/apps/app-12345678
{"code":"500","message":"sql: no rows in result set"}
```
