// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/utils/iptool"
)

func MakeDroneId(suffix string) string {
	return fmt.Sprintf("drone@%s/%s", iptool.GetLocalIP(), strings.TrimSpace(suffix))
}
