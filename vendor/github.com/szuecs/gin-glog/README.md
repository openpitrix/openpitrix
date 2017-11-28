# Gin-Glog

[Gin](https://github.com/gin-gonic/gin) middleware for Logging with
[glog](https://github.com/golang/glog). It is meant as drop in
replacement for the default logger used in Gin.

[![Build Status](https://travis-ci.org/szuecs/gin-glog.svg?branch=master)](https://travis-ci.org/szuecs/gin-glog)
[![Coverage Status](https://coveralls.io/repos/szuecs/gin-glog/badge.svg?branch=master&service=github)](https://coveralls.io/github/szuecs/gin-glog?branch=master)
[![Go Report Card](https://goreportcard.com/badge/szuecs/gin-glog)](https://goreportcard.com/report/szuecs/gin-glog)
[![GoDoc](https://godoc.org/github.com/szuecs/gin-glog?status.svg)](https://godoc.org/github.com/szuecs/gin-glog)

## Project Context and Features

When it comes to choosing a Go framework, there's a lot of confusion
about what to use. The scene is very fragmented, and detailed
comparisons of different frameworks are still somewhat rare. Meantime,
how to handle dependencies and structure projects are big topics in
the Go community. We've liked using Gin for its speed,
accessibility, and usefulness in developing microservice
architectures. In creating Gin-Glog, we wanted to take fuller
advantage of [Gin](https://github.com/gin-gonic/gin)'s capabilities
and help other devs do likewise.

Gin-Glog replaces the default logger of [Gin](https://github.com/gin-gonic/gin) to use
[Glog](https://github.com/golang/glog).

## How Glog is different compared to other loggers

Glog is an efficient pure Go implementation of leveled logs. If you
use Glog, you do not use blocking calls for writing logs. A goroutine
in the background will flush queued loglines into appropriate
logfiles. It provides logrotation, maintains symlinks to current files
and creates flags to your command line interface.

## Requirements

Gin-Glog uses the following [Go](https://golang.org/) packages as
dependencies:

- github.com/gin-gonic/gin
- github.com/golang/glog

## Installation

Assuming you've installed Go and Gin, run this:

    go get github.com/szuecs/gin-glog

## Usage
### Example

Start your webapp to log to STDERR and /tmp:

    % ./webapp -log_dir="/tmp" -alsologtostderr -v=2

```go
package main

import (
    "flag"
    "time"

    "github.com/golang/glog"
    "github.com/szuecs/gin-glog"
    "github.com/gin-gonic/gin"
)

func main() {
    flag.Parse()
    router := gin.New()
    router.Use(ginglog.Logger(3 * time.Second))
    //..
    router.Use(gin.Recovery())

    glog.Warning("warning")
    glog.Error("err")
    glog.Info("info")
    glog.V(2).Infoln("This line will be printed if you use -v=N with N >= 2.")

    router.Run(":8080")
}
```

STDERR output of the example above. Lines with prefix "[Gin-Debug]"
are hardcoded output of Gin and will not appear in your logfile:

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
     - using env:   export GIN_MODE=release
     - using code:  gin.SetMode(gin.ReleaseMode)

    W0306 16:37:12.836001     367 main.go:18] warning
    E0306 16:37:12.836335     367 main.go:19] err
    I0306 16:37:12.836402     367 main.go:20] info
    I0306 16:26:33.901278   32538 main.go:19] This line will be printed if you use -v=N with N >= 2.
    [GIN-debug] Listening and serving HTTP on :8080


## Synopsis

Glog will add the following flags to your binary:

    -alsologtostderr
          log to standard error as well as files
    -log_backtrace_at value
          when logging hits line file:N, emit a stack trace (default :0)
    -log_dir string
          If non-empty, write log files in this directory
    -logtostderr
          log to standard error instead of files
    -stderrthreshold value
          logs at or above this threshold go to stderr
    -v value
          log level for V logs
    -vmodule value
          comma-separated list of pattern=N settings for file-filtered logging


## Contributing/TODO

We welcome contributions from the community; just submit a pull
request. To help you get started, here are some items that we'd love
help with:

- Remove hardcoded logs in [Gin](https://github.com/gin-gonic/gin)
- the code base

Please use github issues as starting point for contributions, new
ideas or bugreports.

## Contact

* E-Mail: team-techmonkeys@zalando.de
* IRC on freenode: #zalando-techmonkeys

## Contributors

Thanks to:

- &lt;your name&gt;

## License

See [LICENSE](LICENSE) file.
