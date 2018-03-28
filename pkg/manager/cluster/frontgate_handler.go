// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"encoding/json"

	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

func (f *Frontgate) parseConf(subnetId, conf string) (string, error) {
	decodeConf := make(map[string]interface{})
	err := json.Unmarshal([]byte(conf), &decodeConf)
	if err != nil {
		return "", err
	}
	decodeConf["version_id"] = constants.FrontgateVersionId
	decodeConf["subnet"] = subnetId
	resConf, err := json.Marshal(&decodeConf)
	if err != nil {
		return "", err
	}
	return string(resConf), nil
}

func (f *Frontgate) getConf(subnetId string) (string, error) {
	conf := constants.FrontgateDefaultConf
	if f.GlobalConfig().Cluster.FrontgateConf != "" {
		conf = f.GlobalConfig().Cluster.FrontgateConf
	}
	return f.parseConf(subnetId, conf)
}

func (f *Frontgate) CreateCluster(register *Register) (string, error) {
	clusterId := models.NewClusterId()

	conf, err := f.getConf(register.SubnetId)
	if err != nil {
		logger.Errorf("Get frontgate cluster conf failed. ")
		return clusterId, err
	}
	clusterWrapper, err := f.Runtime.RuntimeInterface.ParseClusterConf(constants.FrontgateVersionId, conf)
	if err != nil {
		logger.Errorf("Parse frontgate cluster conf failed. ")
		return clusterId, err
	}

	register.ClusterId = clusterId
	register.FrontgateId = ""
	register.ClusterType = constants.FrontgateClusterType
	register.ClusterWrapper = clusterWrapper

	err = register.RegisterClusterWrapper()
	if err != nil {
		return clusterId, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionCreateCluster,
		"", // TODO: need to generate
		f.Runtime.Runtime,
		register.Owner,
	)

	_, err = jobclient.SendJob(newJob)
	return clusterId, err
}

func (f *Frontgate) StartCluster(frontgate *models.Cluster) error {
	newJob := models.NewJob(
		constants.PlaceHolder,
		frontgate.ClusterId,
		frontgate.AppId,
		frontgate.VersionId,
		constants.ActionRecoverClusters,
		"", // TODO: need to generate
		f.Runtime.Runtime,
		frontgate.Owner,
	)

	_, err := jobclient.SendJob(newJob)
	return err
}

func (f *Frontgate) RecoverCluster(frontgate *models.Cluster) error {
	newJob := models.NewJob(
		constants.PlaceHolder,
		frontgate.ClusterId,
		frontgate.AppId,
		frontgate.VersionId,
		constants.ActionRecoverClusters,
		"", // TODO: need to generate
		f.Runtime.Runtime,
		frontgate.Owner,
	)

	_, err := jobclient.SendJob(newJob)
	return err
}
