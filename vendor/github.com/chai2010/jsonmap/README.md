# Map for Json/Struct/Key-Value-Database

[![Build Status](https://travis-ci.org/chai2010/jsonmap.svg)](https://travis-ci.org/chai2010/jsonmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/chai2010/jsonmap)](https://goreportcard.com/report/github.com/chai2010/jsonmap)
[![GoDoc](https://godoc.org/github.com/chai2010/jsonmap?status.svg)](https://godoc.org/github.com/chai2010/jsonmap)

## Install

1. `go get github.com/chai2010/jsonmap`
2. `go run hello.go`

## Example

```go
package main

import (
	"fmt"
	"sort"

	"github.com/chai2010/jsonmap"
)

func main() {
	var jsonMap = jsonmap.JsonMap{
		"a": map[string]interface{}{
			"sub-a": "value-sub-a",
		},
		"b": map[string]interface{}{
			"sub-b": "value-sub-b",
		},
		"c": 123,
		"d": 3.14,
		"e": true,

		"x": map[string]interface{}{
			"a": map[string]interface{}{
				"sub-a": "value-sub-a",
			},
			"b": map[string]interface{}{
				"sub-b": "value-sub-b",
			},
			"c": 123,
			"d": 3.14,
			"e": true,

			"z": map[string]interface{}{
				"zz": "value-zz",
			},
		},
	}

	m := jsonMap.ToMapString("/")

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, m[k])
	}

	// Output:
	// /a/sub-a value-sub-a
	// /b/sub-b value-sub-b
	// /c 123
	// /d 3.14
	// /e true
	// /x/a/sub-a value-sub-a
	// /x/b/sub-b value-sub-b
	// /x/c 123
	// /x/d 3.14
	// /x/e true
	// /x/z/zz value-zz
}
```

## BUGS

Report bugs to <chaishushan@gmail.com>.

Thanks!
