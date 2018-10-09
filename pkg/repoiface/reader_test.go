// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func TestReader_GetIndex(t *testing.T) {
	var repo = &pb.Repo{
		RepoId:     pbutil.ToProtoString("test"),
		Url:        pbutil.ToProtoString("http://helm-chart-repo.pek3a.qingstor.com/svc-catalog-charts/"),
		Providers:  []string{constants.ProviderKubernetes},
		Type:       pbutil.ToProtoString("http"),
		Credential: pbutil.ToProtoString("{}"),
	}
	reader, err := NewReader(context.Background(), repo)
	require.NoError(t, err)
	index, err := reader.GetIndex(context.Background())
	require.NoError(t, err)
	require.NotEqual(t, 0, len(index.GetEntries()))
	t.Log(index.GetEntries())
}
