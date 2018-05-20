// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package reporeader

import (
	neturl "net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHttpReader(t *testing.T) {
	var url = "https://kubernetes-charts.storage.googleapis.com/"
	u, err := neturl.Parse(url)
	require.NoError(t, err)
	httpReader := NewHttpReader(u)
	body, err := httpReader.GetIndexYaml()
	require.NoError(t, err)
	t.Log(len(body))
}
