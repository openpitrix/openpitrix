// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix app server
package main

import (
	"openpitrix.io/openpitrix/pkg/config-v2"
	"openpitrix.io/openpitrix/pkg/handlers/apps"
)

func main() {
	cfg := config.MustLoadUserConfig()
	apps.ListenAndServeAppsServer(cfg)
}
