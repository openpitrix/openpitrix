// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix job manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/service/job"
)

func main() {
	cfg := config.GetConf()
	job.Serve(cfg)
}
