// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// start mysql in docker for test

package apps

import (
	"flag"
	"os"
	"testing"
)

var (
	fFalgDbName = flag.String("test-db-name", "root:password@tcp(mysql:3306)/openpitrix", "set mysql database name")
)

func TestMain(m *testing.M) {
	flag.Parse()

	rv := m.Run()
	os.Exit(rv)
}

func TestAppDatabase(t *testing.T) {
	// TODO
}
