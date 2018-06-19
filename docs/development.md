# Developing for OpenPitrix

The [community repository](https://github.com/openpitrix) hosts all information about
building OpenPitrix from source, how to contribute code and documentation, who to contact about what, etc. If you find a requirement that this doc does not capture, or if you find other docs with references to requirements that are not simply links to this doc, please [submit an issue](https://github.com/openpitrix/openpitrix/issues/new).

----

## To start developing OpenPitrix

First of all, you should fork the project. Then follow one of the options below to develop the project. Please note you should replace the official repo when using __go get__ or __git clone__ below with your own one.

### 1. You have a working [Docker Compose](https://docs.docker.com/compose/install) environment [recommend].

> You need to install [Docker](https://docs.docker.com/engine/installation/) first.

```shell
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ make compose-up
```

#### How to upgrade

```shell
$ make compose-update
```

#### Exit docker runtime environment
```shell
$ make compose-down
```

### 2. You have a working [Kubernetes](https://kubernetes.io/docs/setup/) environment.

> You need to install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) first.

```shell
$ git clone https://github.com/openpitrix/openpitrix.git
$ cd openpitrix/devops/scripts
$ ./deploy-k8s.sh -v v0.1.0 -b -d
```

#### How to upgrade

```shell
$ ./deploy-k8s.sh -v latest -m dbctrl
$ kubectl delete -f ../kubernetes/openpitrix
$ kubectl apply -f ../kubernetes/openpitrix
```

#### Exit kubernetes runtime environment
```shell
$ ./clean
```

----

## Test OpenPitrix

Visit http://127.0.0.1:9100/swagger-ui in browser, and try it online, or test openpitrix api service via command line:

```shell
$ curl http://localhost:9100/v1/repos
{}
$ curl -XPOST -d '{"name": "repo1", "description": "repo1", "type": "https", "url": "https://kubernetes-charts.storage.googleapis.com", "credential": "{}", "visibility": "public", "providers": ["kubernetes"]}' http://localhost:9100/v1/repos
{"repo":{"repo_id":"repo-69z7YN1r2mWl","name":"repo1","description":"repo1","type":"https","url":"https://kubernetes-charts.storage.googleapis.com","credential":"{}","visibility":"public","owner":"system","providers":["kubernetes"],"status":"active","create_time":"2018-05-25T03:40:55.280010221Z","status_time":"2018-05-25T03:40:55.280010431Z"}}
$ curl http://localhost:9100/v1/repos
{"total_count":1,"repo_set":[{"repo_id":"repo-69z7YN1r2mWl","name":"repo1","description":"repo1","type":"https","url":"https://kubernetes-charts.storage.googleapis.com","credential":"{}","visibility":"public","owner":"system","providers":["kubernetes"],"status":"active","create_time":"2018-05-25T03:40:55Z","status_time":"2018-05-25T03:40:55Z"}]}
$ curl http://localhost:9100/v1/apps
{"total_count":131,"app_set":[{"app_id":"app-3wK0YkoXZKLr","name":"sugarcrm","repo_id":"repo-69z7YN1r2mWl","description":"DEPRECATED SugarCRM enables businesses to create extraordinary customer relationships with the most innovative and affordable CRM solution in the market.","status":"active","home":"http://www.sugarcrm.com/","icon":"https://bitnami.com/assets/stacks/sugarcrm/img/sugarcrm-stack-110x117.png","screenshots":"","maintainers":"","keywords":"sugarcrm,CRM","sources":"https://github.com/bitnami/bitnami-docker-sugarcrm","readme":"","chart_name":"sugarcrm","owner":"system","create_time":"2018-05-25T03:42:06Z","status_time":"2018-05-25T03:42:06Z","latest_app_version":{"version_id":"appv-x17NoPGlOrJB","app_id":"app-3wK0YkoXZKLr","owner":"system","name":"1.0.7 [6.5.26]","description":"DEPRECATED SugarCRM enables businesses to create extraordinary customer relationships with the most innovative and affordable CRM solution in the market.","package_name":"https://kubernetes-charts.storage.googleapis.com/sugarcrm-1.0.7.tgz","status":"active","create_time":"2018-05-25T03:42:06Z","status_time":"2018-05-25T03:42:06Z","sequence":16}}, ...]}
$ curl http://localhost:9100/v1/repo_events
{"total_count":1,"repo_event_set":[{"repo_event_id":"repoe-5L1EA4Oqwx18","repo_id":"repo-69z7YN1r2mWl","owner":"system","status":"successful","result":"","create_time":"2018-05-25T03:40:56Z","status_time":"2018-05-25T03:40:56Z"}]}
```

----

## DevOps

Please check [How to set up DevOps environment](devops.md).
