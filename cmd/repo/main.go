// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix repo server
package main

import (
	"openpitrix.io/openpitrix/pkg/cmd/repo"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {
	repo.Main(config.MustLoadUserConfig())
}
