// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"
	"io/ioutil"
)

const Tmpl = `
// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

const InitialGlobalConfig = %s
`

func main() {
	yamlContent, err := ioutil.ReadFile("./global_config.init.yaml")
	if err != nil {
		panic(err)
	}
	data := fmt.Sprintf(Tmpl, "`"+string(yamlContent)+"`")
	err = ioutil.WriteFile("../../pkg/config/init_global_config.go", []byte(data), 0666)
	if err != nil {
		panic(err)
	}
}
