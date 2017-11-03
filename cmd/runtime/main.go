// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix runtime server
package main

import (
	"openpitrix.io/openpitrix/pkg/cmd/runtime"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {
	runtime.Main(config.MustLoadUserConfig())
}
