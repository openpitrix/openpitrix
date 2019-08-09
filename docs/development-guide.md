# Development Guide

This document walks you through how to get started developing OpenPitrix and development workflow.

## Preparing the environment

### Go

OpenPitrix development is written in [Go](http://golang.org/). If you don't have a Go development environment, please [set one up](http://golang.org/doc/code.html), GO requires `>= 1.12`.

> Tips: 
> - Ensure your GOPATH and PATH have been configured in accordance with the Go
environment instructions. 
> - It's recommended to install [macOS GNU tools](https://www.topbug.net/blog/2013/04/14/install-and-use-gnu-command-line-tools-in-mac-os-x) for Mac OS.


### Dependency management

OpenPitrix uses `dep` to manage dependencies in the `vendor/` tree, execute following command to install [dep](https://github.com/golang/dep).

```go
go get -u github.com/golang/dep/cmd/dep
```

### Test

In the development process, it is recommended to install an single-node [all-in-one](https://github.com/openpitrix/openpitrix#all-in-one) environment (Kubernetes-based) for quick testing.

> Tip: It also supports to use Docker for Desktop ships with Kubernetes as the test environment.

## Development Workflow

![](https://pek3b.qingstor.com/kubesphere-docs/png/20190708115631.png)

### 1 Fork in the cloud

1. Visit https://github.com/openpitrix/openpitrix
2. Click `Fork` button to establish a cloud-based fork.

### 2 Clone fork to local storage

Per Go's [workspace instructions][https://golang.org/doc/code.html#Workspaces], place OpenPitrix' code on your `GOPATH` using the following cloning procedure.

1. Define a local working directory:

```bash
$ export working_dir=$GOPATH/src/openpitrix.io
$ export user={your github profile name}
```

2. Create your clone locally:

```bash
$ mkdir -p $working_dir
$ cd $working_dir
$ git clone https://github.com/$user/openpitrix.git
$ cd $working_dir/openpitrix
$ git remote add upstream https://github.com/openpitrix/openpitrix.git

# Never push to upstream master
$ git remote set-url --push upstream no_push

# Confirm that your remotes make sense:
$ git remote -v
```

### 3 Keep your branch in sync

```bash
git fetch upstream
git checkout master
git rebase upstream/master
```

### 4 Add new features or fix issues

Branch from it:

```bash
$ git checkout -b myfeature
```

Then edit code on the myfeature branch.

**Run and test**

```bash
$ make build
$ make compose-up
```

Run `make help` for additional information on these make targets.

### 5 Development in new branch

**Sync with upstream**

After the test is completed, suggest you to keep your local in sync with upstream which can avoid conflicts.

```
# Rebase your the master branch of your local repo.
$ git checkout master
$ git rebase upstream/master

# Then make your development branch in sync with master branch
git checkout new_feature
git rebase -i master
```
**Commit local changes**

```bash
$ git add <file>
$ git commit -s -m "add your description"
```

### 6 Push to your folk

When ready to review (or just to establish an offsite backup or your work), push your branch to your fork on github.com:

```
$ git push -f ${your_remote_name} myfeature
```

### 7 Create a PR

- Visit your fork at https://github.com/$user/openpitrix
- Click the` Compare & Pull Request` button next to your myfeature branch.
- Check out the [pull request process](pull-request.md) for more details and advice.


## CI/CD

OpenPitrix uses [Travis CI](https://travis-ci.org/) as a CI/CD tool.

The components of OpenPitrix under `/cmd` folder need to be compiled and build include following:

After your PR is merged，Travis CI will compile the entire project and build the image, and push the image `openpitrix/[component-name]:latest` to Dockerhub (e.g. `openpitrix/openpitrix-api-gateway:latest`)

## Code conventions

Please reference [Code conventions](https://github.com/kubernetes/community/blob/master/contributors/guide/coding-conventions.md) and follow with the rules.

**Note:**

> - All new packages and most new significant functionality must come with unit tests
> - Comment your code in English, see [Go's commenting conventions
](http://blog.golang.org/godoc-documenting-go-code)