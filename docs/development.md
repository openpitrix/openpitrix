# Developing for OpenPitrix

The [community repository](https://github.com/openpitrix) hosts all information about
building OpenPitrix from source, how to contribute code and documentation, who to contact about what, etc. If you find a requirement that this doc does not capture, or if you find other docs with references to requirements that are not simply links to this doc, please [submit an issue](https://github.com/openpitrix/openpitrix/issues/new).

----

## To start developing OpenPitrix

First of all, you should fork the project. Then follow one of the options below to develop the project. Please note you should replace the official repo when using __go get__ or __git clone__ below with your own one.

### 1. You have a working [Docker Compose](https://docs.docker.com/compose/install) environment [recommend].

> You need to install [Docker](https://docs.docker.com/engine/installation/) first.

#### How to deploy
```shell
$ git clone https://github.com/openpitrix/openpitrix
$ cd openpitrix
$ make build
$ make compose-up
```

#### Verifying the deploy
```shell
$ docker ps
CONTAINER ID        IMAGE                         COMMAND                  CREATED             STATUS              PORTS                               NAMES
c5fbe681d670        openpitrix/dashboard          "npm run prod:serve"     26 minutes ago      Up 26 minutes       0.0.0.0:8000->8000/tcp              openpitrix-dashboard
69d849bf8429        openpitrix                    "api-gateway"            26 minutes ago      Up 26 minutes       0.0.0.0:9100->9100/tcp              openpitrix-api-gateway
0ee764b83fd6        openpitrix                    "repo-indexer"           26 minutes ago      Up 26 minutes                                           openpitrix-repo-indexer
e0de82867bf2        openpitrix                    "task-manager"           26 minutes ago      Up 26 minutes                                           openpitrix-task-manager
a7d3bfd02cc9        openpitrix                    "app-manager"            26 minutes ago      Up 26 minutes                                           openpitrix-app-manager
4104d75d3ec9        openpitrix                    "cluster-manager"        26 minutes ago      Up 26 minutes                                           openpitrix-cluster-manager
f10d5e45cf0a        openpitrix                    "job-manager"            26 minutes ago      Up 26 minutes                                           openpitrix-job-manager
c67c9d43d762        openpitrix                    "category-manager"       26 minutes ago      Up 26 minutes                                           openpitrix-category-manager
a7e346810492        openpitrix                    "repo-manager"           26 minutes ago      Up 26 minutes                                           openpitrix-repo-manager
77ad45001c06        openpitrix                    "runtime-manager"        26 minutes ago      Up 26 minutes                                           openpitrix-runtime-manager
7e11808e4c46        quay.io/coreos/etcd:v3.2.18   "etcd -listen-client…"   26 minutes ago      Up 26 minutes       2380/tcp, 0.0.0.0:12379->2379/tcp   openpitrix-etcd
f00be9e66ba1        openpitrix                    "pilot -config=/opt/…"   27 minutes ago      Up 27 minutes       0.0.0.0:9110->9110/tcp              openpitrix-pilot-service
2d312d6eb408        mysql:8.0.11                  "docker-entrypoint.s…"   31 minutes ago      Up 31 minutes       0.0.0.0:13306->3306/tcp             openpitrix-db
```

#### How to upgrade

```shell
$ make compose-update
```

#### How to clean environment
```shell
$ make compose-down
```

### 2. You have a working [Kubernetes](https://kubernetes.io/docs/setup/) environment.

> You need to install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) first.

#### How to deploy

Clone source code from master branch: 
```shell
$ git clone https://github.com/openpitrix/openpitrix.git
$ cd openpitrix/deploy/kubernetes/scripts
$ ./deploy-k8s.sh -n openpitrix-system -b -d
```

Or go to the [OpenPitrix release](https://github.com/openpitrix/openpitrix/releases) page to download the deploy package. You can also run the following command to download and extract the latest release deploy package automatically:
```shell
$ curl -L https://git.io/GetOpenPitrix | sh -
$ cd openpitrix-${OPENPITRIX_VERSION}-kubernetes/kubernetes/scripts
$ ./deploy-k8s.sh -n openpitrix-system -b -d
```

If the dashboard is required:
```shell
$ ./deploy-k8s.sh -n openpitrix-system -s
```

If create clusters in vm-based runtime is required, the metadata model is needed:
```shell
$ ./deploy-k8s.sh -n openpitrix-system -m
```

#### Verifying the deploy
```shell
$ kubectl get deployment -n openpitrix-system
NAME                                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
openpitrix-api-gateway-deployment        1         1         1            1           3h
openpitrix-app-manager-deployment        1         1         1            1           3h
openpitrix-category-manager-deployment   1         1         1            1           3h
openpitrix-cluster-manager-deployment    1         1         1            1           3h
openpitrix-db-deployment                 1         1         1            1           3h
openpitrix-etcd-deployment               1         1         1            1           3h
openpitrix-job-manager-deployment        1         1         1            1           3h
openpitrix-repo-indexer-deployment       1         1         1            1           3h
openpitrix-repo-manager-deployment       1         1         1            1           3h
openpitrix-runtime-manager-deployment    1         1         1            1           3h
openpitrix-task-manager-deployment       1         1         1            1           3h
```

#### How to upgrade

```shell
$ ./deploy-k8s.sh -n openpitrix-system -b -d
```

If the dashboard is required:
```shell
$ ./deploy-k8s.sh -n openpitrix-system -s
```

If create clusters in vm-based runtime is required, the metadata model is needed:
```shell
$ ./deploy-k8s.sh -n openpitrix-system -m
```

#### How to clean environment
```shell
$ ./clean -n openpitrix-system
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
