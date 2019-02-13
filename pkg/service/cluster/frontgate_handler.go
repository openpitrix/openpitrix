// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"

	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (f *Frontgate) parseConf(subnetId, conf string) (string, error) {
	decodeConf := make(map[string]interface{})
	err := jsonutil.Decode([]byte(conf), &decodeConf)
	if err != nil {
		return "", err
	}
	decodeConf["version_id"] = constants.FrontgateVersionId
	decodeConf["subnet"] = subnetId
	resConf := jsonutil.ToString(&decodeConf)
	return resConf, nil
}

func (f *Frontgate) getConf(ctx context.Context, subnetId, runtimeUrl, runtimeZone string) (string, error) {
	conf := constants.FrontgateDefaultConf
	if pi.Global().GlobalConfig().Cluster.FrontgateConf != "" {
		conf = pi.Global().GlobalConfig().Cluster.FrontgateConf
	}

	imageConfig, err := pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtimeUrl, runtimeZone)
	if err != nil {
		return "", err
	}
	if imageConfig.FrontgateConf != "" {
		conf = imageConfig.FrontgateConf
	}

	return f.parseConf(subnetId, conf)
}

func (f *Frontgate) CreateCluster(ctx context.Context, clusterWrapper *models.ClusterWrapper) (string, error) {
	clusterId := models.NewClusterId()

	conf, err := f.getConf(ctx, clusterWrapper.Cluster.SubnetId, f.Runtime.RuntimeUrl, f.Runtime.Zone)
	if err != nil {
		logger.Error(ctx, "Get frontgate cluster conf failed. ")
		return clusterId, err
	}
	frontgateWrapper := new(models.ClusterWrapper)
	providerClient, err := providerclient.NewRuntimeProviderManagerClient()
	if err != nil {
		return clusterId, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	response, err := providerClient.ParseClusterConf(ctx, &pb.ParseClusterConfRequest{
		RuntimeId: pbutil.ToProtoString(clusterWrapper.Cluster.RuntimeId),
		VersionId: pbutil.ToProtoString(constants.FrontgateVersionId),
		Conf:      pbutil.ToProtoString(conf),
		Cluster:   models.ClusterWrapperToPb(frontgateWrapper),
	})
	if err != nil {
		logger.Error(ctx, "Parse frontgate cluster conf failed.")
		return clusterId, err
	}

	frontgateWrapper = models.PbToClusterWrapper(response.Cluster)
	frontgateWrapper.Cluster.Zone = clusterWrapper.Cluster.Zone
	frontgateWrapper.Cluster.Debug = clusterWrapper.Cluster.Debug
	frontgateWrapper.Cluster.ClusterId = clusterId
	frontgateWrapper.Cluster.SubnetId = clusterWrapper.Cluster.SubnetId
	frontgateWrapper.Cluster.VpcId = clusterWrapper.Cluster.VpcId
	frontgateWrapper.Cluster.Owner = clusterWrapper.Cluster.Owner
	frontgateWrapper.Cluster.OwnerPath = clusterWrapper.Cluster.OwnerPath
	frontgateWrapper.Cluster.ClusterType = constants.FrontgateClusterType
	frontgateWrapper.Cluster.FrontgateId = ""
	frontgateWrapper.Cluster.RuntimeId = f.Runtime.RuntimeId

	err = RegisterClusterWrapper(ctx, frontgateWrapper)
	if err != nil {
		return clusterId, err
	}

	directive := jsonutil.ToString(frontgateWrapper)
	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		frontgateWrapper.Cluster.AppId,
		frontgateWrapper.Cluster.VersionId,
		constants.ActionCreateCluster,
		directive,
		f.Runtime.Runtime.Provider,
		frontgateWrapper.Cluster.OwnerPath,
		frontgateWrapper.Cluster.RuntimeId,
	)

	_, err = jobclient.SendJob(ctx, newJob)
	return clusterId, err
}

func (f *Frontgate) StartCluster(ctx context.Context, frontgate *models.Cluster) error {
	clusterWrapper, err := getClusterWrapper(ctx, frontgate.ClusterId)
	if err != nil {
		return err
	}

	directive := jsonutil.ToString(clusterWrapper)
	newJob := models.NewJob(
		constants.PlaceHolder,
		frontgate.ClusterId,
		frontgate.AppId,
		frontgate.VersionId,
		constants.ActionStartClusters,
		directive,
		f.Runtime.Runtime.Provider,
		frontgate.OwnerPath,
		frontgate.RuntimeId,
	)

	_, err = jobclient.SendJob(ctx, newJob)
	return err
}

func (f *Frontgate) RecoverCluster(ctx context.Context, frontgate *models.Cluster) error {
	clusterWrapper, err := getClusterWrapper(ctx, frontgate.ClusterId)
	if err != nil {
		return err
	}

	directive := jsonutil.ToString(clusterWrapper)
	newJob := models.NewJob(
		constants.PlaceHolder,
		frontgate.ClusterId,
		frontgate.AppId,
		frontgate.VersionId,
		constants.ActionRecoverClusters,
		directive,
		f.Runtime.Runtime.Provider,
		frontgate.OwnerPath,
		frontgate.RuntimeId,
	)

	_, err = jobclient.SendJob(ctx, newJob)
	return err
}
