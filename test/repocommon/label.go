// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repocommon

import (
	"net/url"

	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func GenerateLabels() string {
	v := url.Values{}
	v.Add("key1", idutil.GetUuid(""))
	v.Add("key2", idutil.GetUuid(""))
	v.Add("key3", idutil.GetUuid(""))
	v.Add("key4", idutil.GetUuid(""))
	v.Add("key5", idutil.GetUuid(""))
	return v.Encode()
}
