# gin-nice-recovery

[Gin](https://gin-gonic.github.io/gin/) middleware to provide a nice user experience when recovering from a panic.

## Why?

The default `gin.Recovery()` middleware leaves the user looking a blank white page. This middleware calls the
specified handler, which can render a nice looking error page, return customized error JSON, or whatever is
required. It logs the same HTTP request information and stack trace as `gin.Recovery()`.

## Installation

```bash
$ go get github.com/ekyoung/gin-nice-recovery
```

## Usage

```go
package main

import (
    "github.com/gin-gonic/gin"

    "github.com/ekyoung/gin-nice-recovery"
)

func main() {
    router := gin.New()      // gin.Default() installs gin.Recovery() so use gin.New() instead
    router.Use(gin.Logger()) // Install the default logger, not required

    // Install nice.Recovery, passing the handler to call after recovery
    router.Use(nice.Recovery(recoveryHandler))

    // Load templates as usual
    router.LoadHTMLFiles("error.tmpl")

    // Define routes as usual
    router.GET("/", func(c *gin.Context) {
        panic("Doh!")
    })

    router.Run(":8080")
}

func recoveryHandler(c *gin.Context, err interface{}) {
    c.HTML(500, "error.tmpl", gin.H{
        "title": "Error",
        "err":   err,
    })
}
```

Complete example code can be found in the [example repository](https://github.com/ekyoung/gin-nice-recovery-example).