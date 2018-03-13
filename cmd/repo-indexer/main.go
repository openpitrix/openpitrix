// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix repo manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/manager/repo_indexer"
)

func main() {
	cfg := config.LoadConf()
	repo_indexer.Serve(cfg)
}
