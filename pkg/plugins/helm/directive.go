// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"context"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type JobDirective struct {
	Namespace   string
	RuntimeId   string
	Values      string
	ClusterName string
}

func decodeJobDirective(ctx context.Context, data string) (*JobDirective, error) {
	clusterWrapper, err := models.NewClusterWrapper(ctx, data)
	if err != nil {
		return nil, err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId

	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	namespace := runtime.Zone

	j := &JobDirective{
		Namespace:   namespace,
		RuntimeId:   runtimeId,
		Values:      clusterWrapper.Cluster.Env,
		ClusterName: clusterWrapper.Cluster.Name,
	}

	return j, nil
}

type TaskDirective struct {
	VersionId         string
	Namespace         string
	RuntimeId         string
	Values            string
	ClusterName       string
	RawClusterWrapper string
}

func encodeTaskDirective(v interface{}) string {
	return jsonutil.ToString(v)
}

func decodeTaskDirective(data string) (*TaskDirective, error) {
	var v TaskDirective
	err := jsonutil.Decode([]byte(data), &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
