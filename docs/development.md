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
CONTAINER ID        IMAGE                                      COMMAND                  CREATED             STATUS                   PORTS                                            NAMES
b0540d41b90f        openpitrix/dashboard                       "npm run prod:serve"     5 hours ago         Up 5 hours               0.0.0.0:8000->8000/tcp                           openpitrix-dashboard
162ee4e72a05        openpitrix                                 "api-gateway"            5 hours ago         Up 5 hours               0.0.0.0:9100->9100/tcp                           openpitrix-api-gateway
7fa3bda14643        openpitrix                                 "repo-indexer"           5 hours ago         Up 5 hours                                                                openpitrix-repo-indexer
f988950a170c        minio/minio:RELEASE.2018-09-25T21-34-43Z   "sh -c 'mkdir -p /da…"   5 hours ago         Up 5 hours (healthy)     0.0.0.0:19000->9000/tcp                          openpitrix-minio
7257c88b5048        openpitrix                                 "category-manager"       5 hours ago         Up 5 hours                                                                openpitrix-category-manager
893fea21a52a        openpitrix                                 "cluster-manager"        5 hours ago         Up 5 hours                                                                openpitrix-cluster-manager
fd670a665ca2        openpitrix                                 "account-service"        5 hours ago         Up 5 hours                                                                openpitrix-account-service
00ab68d59dc6        openpitrix                                 "job-manager"            5 hours ago         Up 5 hours                                                                openpitrix-job-manager
d4964a9c2f54        openpitrix                                 "app-manager"            5 hours ago         Up 5 hours                                                                openpitrix-app-manager
8c63b77a4af6        openpitrix                                 "task-manager"           5 hours ago         Up 5 hours                                                                openpitrix-task-manager
55e819d66118        openpitrix                                 "repo-manager"           5 hours ago         Up 5 hours                                                                openpitrix-repo-manager
79c7e40d4566        openpitrix                                 "runtime-manager"        5 hours ago         Up 5 hours                                                                openpitrix-runtime-manager
599d79142876        quay.io/coreos/etcd:v3.2.18                "etcd --data-dir /da…"   5 hours ago         Up 5 hours               2380/tcp, 0.0.0.0:12379->2379/tcp                openpitrix-etcd
0f01a2372bf3        openpitrix                                 "pilot -config=/opt/…"   5 hours ago         Up 5 hours               0.0.0.0:9110->9110/tcp, 0.0.0.0:9114->9114/tcp   openpitrix-pilot-service
03c66d4833d0        mysql:8.0.11                               "docker-entrypoint.s…"   5 hours ago         Up 5 hours               0.0.0.0:13306->3306/tcp                          openpitrix-db
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
$ ./deploy-k8s.sh -n openpitrix-system -b -d -o
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
openpitrix-api-gateway-deployment        1         1         1            1           2m
openpitrix-app-manager-deployment        1         1         1            1           2m
openpitrix-category-manager-deployment   1         1         1            1           2m
openpitrix-cluster-manager-deployment    1         1         1            1           1m
openpitrix-dashboard-deployment          1         1         1            1           1m
openpitrix-db-deployment                 1         1         1            1           57d
openpitrix-etcd-deployment               1         1         1            1           57d
openpitrix-account-service-deployment    1         1         1            1           1m
openpitrix-job-manager-deployment        1         1         1            1           1m
openpitrix-minio-deployment              1         1         1            1           5d
openpitrix-pilot-deployment              1         1         1            1           1m
openpitrix-repo-indexer-deployment       1         1         1            1           1m
openpitrix-repo-manager-deployment       1         1         1            1           1m
openpitrix-runtime-manager-deployment    1         1         1            1           1m
openpitrix-task-manager-deployment       1         1         1            1           1m

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
