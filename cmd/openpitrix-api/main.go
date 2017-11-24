// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// install httpie:
// - macOS: brew install httpie
// - Ubuntu: apt-get install httpie
// - CentOS: yum install httpie

// http get :8080/v1/apps

// curl http://localhost:8080/v1/apps

// openpitrix server
package main

import (
	"openpitrix.io/openpitrix/pkg/cmd/api"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {
	api.Main(config.MustLoadUserConfig())
}
