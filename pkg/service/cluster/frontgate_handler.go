// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
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

func (f *Frontgate) getConf(subnetId string) (string, error) {
	conf := constants.FrontgateDefaultConf
	if pi.Global().GlobalConfig().Cluster.FrontgateConf != "" {
		conf = pi.Global().GlobalConfig().Cluster.FrontgateConf
	}
	return f.parseConf(subnetId, conf)
}

func (f *Frontgate) CreateCluster(clusterWrapper *models.ClusterWrapper) (string, error) {
	clusterId := models.NewClusterId()

	conf, err := f.getConf(clusterWrapper.Cluster.SubnetId)
	if err != nil {
		logger.Error("Get frontgate cluster conf failed. ")
		return clusterId, err
	}
	providerInterface, err := plugins.GetProviderPlugin(f.Runtime.Provider, nil)
	if err != nil {
		logger.Error("No such provider [%s]. ", f.Runtime.Provider)
		return clusterId, err
	}
	frontgateWrapper, err := providerInterface.ParseClusterConf(constants.FrontgateVersionId, clusterWrapper.Cluster.RuntimeId, conf)
	if err != nil {
		logger.Error("Parse frontgate cluster conf failed. ")
		return clusterId, err
	}

	frontgateWrapper.Cluster.ClusterId = clusterId
	frontgateWrapper.Cluster.SubnetId = clusterWrapper.Cluster.SubnetId
	frontgateWrapper.Cluster.VpcId = clusterWrapper.Cluster.VpcId
	frontgateWrapper.Cluster.Owner = clusterWrapper.Cluster.Owner
	frontgateWrapper.Cluster.ClusterType = constants.FrontgateClusterType
	frontgateWrapper.Cluster.FrontgateId = ""
	frontgateWrapper.Cluster.RuntimeId = f.Runtime.RuntimeId

	err = RegisterClusterWrapper(frontgateWrapper)
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
		f.Runtime.Provider,
		frontgateWrapper.Cluster.Owner,
	)

	_, err = jobclient.SendJob(newJob)
	return clusterId, err
}

func (f *Frontgate) StartCluster(frontgate *models.Cluster) error {
	clusterWrapper, err := getClusterWrapper(frontgate.ClusterId)
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
		f.Runtime.Provider,
		frontgate.Owner,
	)

	_, err = jobclient.SendJob(newJob)
	return err
}

func (f *Frontgate) RecoverCluster(frontgate *models.Cluster) error {
	clusterWrapper, err := getClusterWrapper(frontgate.ClusterId)
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
		f.Runtime.Provider,
		frontgate.Owner,
	)

	_, err = jobclient.SendJob(newJob)
	return err
}
