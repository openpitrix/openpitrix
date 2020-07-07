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
	var url = "http://helm-chart-repo.pek3a.qingstor.com/svc-catalog-charts/"
	u, err := neturl.Parse(url)
	require.NoError(t, err)
	httpInterface, err := NewHttpInterface(context.TODO(), u)
	require.NoError(t, err)
	body, err := httpInterface.ReadFile(context.Background(), "index.yaml")
	require.NoError(t, err)
	t.Log(len(body))

	url = "https://helm.elastic.co/"
	u, err = neturl.Parse(url)
	require.NoError(t, err)
	httpInterface, err = NewHttpInterface(context.TODO(), u)
	require.NoError(t, err)
	body, err = httpInterface.ReadFile(context.Background(), "https://helm.elastic.co/helm/metricbeat/metricbeat-7.3.2.tgz")
	require.NoError(t, err)
	t.Log(len(body))

	// httpbin guide: https://httpbin.org/#/Auth/get_basic_auth__user___passwd_
	url = "https://1:2@httpbin.org/basic-auth/username-1/password-2"
	u, err = neturl.Parse(url)
	require.NoError(t, err)
	httpInterface, err = NewHttpInterface(context.TODO(), u)
	require.NoError(t, err)
	body, err = httpInterface.ReadFile(context.Background(), "https://httpbin.org/basic-auth/1/2")
	require.NoError(t, err)
	t.Log(len(body))
}
