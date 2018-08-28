// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	neturl "net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHttpInterface(t *testing.T) {
	var url = "https://kubernetes-charts.storage.googleapis.com/"
	u, err := neturl.Parse(url)
	require.NoError(t, err)
	httpInterface, err := NewHttpInterface(context.TODO(), u)
	require.NoError(t, err)
	body, err := httpInterface.ReadFile(context.Background(), "index.yaml")
	require.NoError(t, err)
	t.Log(len(body))
}
