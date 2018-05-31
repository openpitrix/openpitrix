// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package funcutil

import (
	"testing"

	. "github.com/chai2010/assert"
)

func TestCallerName(t *testing.T) {
	Assert(t, CallerName(0) == "openpitrix.io/openpitrix/pkg/util/funcutil.CallerName")
	Assert(t, CallerName(1) == "openpitrix.io/openpitrix/pkg/util/funcutil.TestCallerName")
}
