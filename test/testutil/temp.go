// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package testutil

import (
	"path"

	"openpitrix.io/openpitrix/pkg/util/idutil"
)

const TmpPath = "/tmp/openpitrix-test"

func GetTmpDir() string {
	return path.Join(TmpPath, idutil.GetUuid(""))
}
