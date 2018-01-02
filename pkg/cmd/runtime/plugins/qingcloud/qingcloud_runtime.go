// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yunify/qingcloud-sdk-go/client"
	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/service"
	"github.com/yunify/qingcloud-sdk-go/utils"

	"openpitrix.io/openpitrix/pkg/cmd/runtime"
	"openpitrix.io/openpitrix/pkg/logger"
)

func init() {
	runtime.RegisterRuntime(new(QingcloudRuntime))
}

type QingcloudRuntime struct {
}

func (p *QingcloudRuntime) initService() (qingcloudService *service.QingCloudService, err error) {
	userConf, err := config.NewDefault()
	if err != nil {
		return
	}
	err = userConf.LoadUserConfig()
	if err != nil {
		return
	}
	qingcloudService, err = service.Init(userConf)
	if err != nil {
		return
	}
	return
}

func (p *QingcloudRuntime) initClusterService() (clusterService *service.ClusterService, err error) {
	qingcloudService, err := p.initService()
	if err != nil {
		logger.Errorf("Failed to init qingcloud api service: %v", err)
		return
	}
	clusterService, err = qingcloudService.Cluster(qingcloudService.Config.Zone)
	return
}

func (p *QingcloudRuntime) initJobService() (jobService *service.JobService, err error) {
	qingcloudService, err := p.initService()
	if err != nil {
		logger.Errorf("Failed to init qingcloud api service: %v", err)
		return
	}
	jobService, err = qingcloudService.Job(qingcloudService.Config.Zone)
	return
}

func (p *QingcloudRuntime) waitJobs(jobIds []string, timeout time.Duration, waitInterval time.Duration) error {
	jobService, err := p.initJobService()
	if err != nil {
		logger.Errorf("Failed to init job service when wait cluster jobs %s: %v", jobIds, err)
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(jobIds))
	done := make(chan error, 1)

	for _, jobId := range jobIds {
		go func() {
			defer wg.Done()
			done <- client.WaitJob(jobService, jobId, timeout, waitInterval)

			waitErr, ok := <-done
			if !ok {
				return
			}
			if waitErr != nil {
				logger.Errorf("Failed to wait cluster job [%s]: %v", jobId, waitErr)
				err = waitErr
				close(done)
			}
		}()
	}
	wg.Wait()
	return err
}

func (p *QingcloudRuntime) waitClustersStatus(clusterService *service.ClusterService, clusterIds string, status string, timeout time.Duration, waitInterval time.Duration) error {
	return utils.WaitForSpecificOrError(func() (bool, error) {
		clusters := strings.Split(clusterIds, ",")
		output, err := clusterService.DescribeClusters(&service.DescribeClustersInput{Clusters: service.StringSlice(clusters)})
		if err != nil {
			//network or api error
			return false, nil
		}

		if len(output.ClusterSet) != len(clusters) {
			return false, fmt.Errorf("Can not find clusters [%s]", clusterIds)
		}

		for _, cluster := range output.ClusterSet {
			if cluster.TransitionStatus == nil {
				logger.Errorf("Cluster [%s] transition_status is nil ", clusterIds)
				return false, nil
			}
			if service.StringValue(cluster.TransitionStatus) != "" {
				return false, nil
			}
			if service.StringValue(cluster.Status) != status {
				return false, fmt.Errorf("Cluster [%s] status [%s]", clusterIds, status)
			}
		}
		return true, nil

	}, timeout, waitInterval)
}

func (p *QingcloudRuntime) Name() string { return "qingcloud" }

func (p *QingcloudRuntime) Run(app string, args ...string) error {
	return errors.New("TODO")
}

func (p *QingcloudRuntime) DescribeVxnets() {
	// call QingCloud DescribeVxnets Api
	// var action = "DescribeVxnets"
}

func (p *QingcloudRuntime) GetGlobalUniqueId() string {
	var globalUniqueId = ""
	return globalUniqueId
}

func (p *QingcloudRuntime) GetClusterPrice(appId string, appVersion string, conf string) {
}

func (p *QingcloudRuntime) CreateCluster(appConf string, shouldWait bool, args ...string) (clusterId string, err error) {
	logger.Infof("Creating cluster...")

	// call QingCloud CreateCluster Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when create cluster: %v", err)
		return
	}

	output, err := clusterService.CreateCluster(&service.CreateClusterInput{Conf: service.String(appConf)})
	if err != nil {
		logger.Errorf("Failed to create cluster: %v", err)
		return
	}

	clusterId = service.StringValue(output.ClusterID)

	if shouldWait {
		jobId := service.StringValue(output.JobID)
		jobIds := []string{jobId}
		err = p.waitJobs(jobIds, TIMEOUT_CREATE_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait create cluster jobs %s: %v", jobIds, err)
			return
		}
	}

	logger.Infof("Create cluster [%s] successful", clusterId)
	return
}

func (p *QingcloudRuntime) StopClusters(clusterIds string, shouldWait bool, args ...string) error {
	logger.Infof("Stoping cluster [%s]...", clusterIds)

	// call QingCloud StopClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when stop cluster: %v", err)
		return err
	}

	clusters := strings.Split(clusterIds, ",")
	force := 0
	if len(args) >= 1 {
		force, err = strconv.Atoi(args[0])
		if err != nil {
			logger.Errorf("Failed to stop cluster [%s] with force [%s]: %v", clusterIds, args[0], err)
			return err
		}
	}

	output, err := clusterService.StopClusters(&service.StopClustersInput{Clusters: service.StringSlice(clusters), Force: service.Int(force)})
	if err != nil {
		logger.Errorf("Failed to stop cluster [%s]: %v", clusterIds, err)
		return err
	}

	jobIdsMap := service.StringValueMap(output.JobIDs)

	if shouldWait {
		jobIds := []string{}
		for _, jobId := range jobIdsMap {
			jobIds = append(jobIds, jobId)
		}
		err = p.waitJobs(jobIds, TIMEOUT_STOP_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait stop cluster jobs %s: %v", jobIds, err)
			return err
		}
	}

	logger.Infof("Stop cluster [%s] successful", clusterIds)
	return err
}

func (p *QingcloudRuntime) StartClusters(clusterIds string, shouldWait bool, args ...string) error {
	logger.Infof("Starting cluster [%s]...", clusterIds)

	// call QingCloud StartClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when start cluster: %v", err)
		return err
	}

	clusters := strings.Split(clusterIds, ",")

	output, err := clusterService.StartClusters(&service.StartClustersInput{Clusters: service.StringSlice(clusters)})
	if err != nil {
		logger.Errorf("Failed to start cluster [%s]: %v", clusterIds, err)
		return err
	}

	jobIdsMap := service.StringValueMap(output.JobIDs)

	if shouldWait {
		jobIds := []string{}
		for _, jobId := range jobIdsMap {
			jobIds = append(jobIds, jobId)
		}
		err = p.waitJobs(jobIds, TIMEOUT_START_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait start cluster jobs %s: %v", jobIds, err)
			return err
		}
	}

	logger.Infof("Start cluster [%s] successful", clusterIds)
	return err
}

func (p *QingcloudRuntime) DeleteClusters(clusterIds string, shouldWait bool, args ...string) error {
	logger.Infof("Deleting cluster [%s]...", clusterIds)

	// call QingCloud DeleteClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when delete cluster: %v", err)
		return err
	}

	clusters := strings.Split(clusterIds, ",")

	force := 0
	if len(args) >= 1 {
		force, err = strconv.Atoi(args[0])
		if err != nil {
			logger.Errorf("Failed to delete cluster [%s] with force [%s]: %v", clusterIds, args[0], err)
			return err
		}
	}

	output, err := clusterService.DeleteClusters(&service.DeleteClustersInput{Clusters: service.StringSlice(clusters), Force: service.Int(force)})
	if err != nil {
		logger.Errorf("Failed to delete cluster [%s]: %v", clusterIds, err)
		return err
	}

	jobIdsMap := service.StringValueMap(output.JobIDs)

	if shouldWait {
		jobIds := []string{}
		for _, jobId := range jobIdsMap {
			jobIds = append(jobIds, jobId)
		}
		err = p.waitJobs(jobIds, TIMEOUT_DELETE_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait delete cluster jobs %s: %v", jobIds, err)
			return err
		}
	}

	logger.Infof("Delete cluster [%s] successful", clusterIds)
	return err
}

func (p *QingcloudRuntime) RecoverClusters(clusterIds string, shouldWait bool, args ...string) error {
	logger.Infof("Recovering cluster [%s]...", clusterIds)

	// call QingCloud RecoverClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when recover cluster: %v", err)
		return err
	}

	clusters := strings.Split(clusterIds, ",")

	_, err = clusterService.RecoverClusters(&service.RecoverClustersInput{Resources: service.StringSlice(clusters)})
	if err != nil {
		logger.Errorf("Failed to recover cluster [%s]: %v", clusterIds, err)
		return err
	}

	if shouldWait {
		err = p.waitClustersStatus(clusterService, clusterIds, STATUS_ACTIVE, TIMEOUT_RECOVER_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait cluster [%s] status [%s]: %v", clusterIds, STATUS_ACTIVE, err)
			return err
		}
	}

	logger.Infof("Recover cluster [%s] successful", clusterIds)
	return err
}

func (p *QingcloudRuntime) CeaseClusters(clusterIds string, shouldWait bool, args ...string) error {
	logger.Infof("Ceasing cluster [%s]...", clusterIds)

	// call QingCloud CeaseClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when cease cluster: %v", err)
		return err
	}

	clusters := strings.Split(clusterIds, ",")

	output, err := clusterService.CeaseClusters(&service.CeaseClustersInput{Clusters: service.StringSlice(clusters)})
	if err != nil {
		logger.Errorf("Failed to cease cluster [%s]: %v", clusterIds, err)
		return err
	}

	jobIdsMap := service.StringValueMap(output.JobIDs)

	if shouldWait {
		jobIds := []string{}
		for _, jobId := range jobIdsMap {
			jobIds = append(jobIds, jobId)
		}
		err = p.waitJobs(jobIds, TIMEOUT_CEASE_CLUSTER, WAIT_INTERVAL)
		if err != nil {
			logger.Errorf("Failed to wait cease cluster jobs %s: %v", jobIds, err)
			return err
		}
	}

	logger.Infof("Cease cluster [%s] successful", clusterIds)
	return err
}

func (p *QingcloudRuntime) DescribeClusters(clusterIds string, args ...string) (output *service.DescribeClustersOutput, err error) {
	// call QingCloud DescribeClusters Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when cease cluster: %v", err)
		return
	}

	clusters := strings.Split(clusterIds, ",")

	output, err = clusterService.DescribeClusters(&service.DescribeClustersInput{Clusters: service.StringSlice(clusters)})
	if err != nil {
		logger.Errorf("Failed to describe cluster [%s]: %v", clusterIds, err)
		return
	}

	logger.Printf("count: %+v", service.IntValue(output.TotalCount))

	return
}

func (p *QingcloudRuntime) DescribeClusterNodes(clusterId string, args ...string) (output *service.DescribeClusterNodesOutput, err error) {
	// call QingCloud DescribeClusterNodes Api
	clusterService, err := p.initClusterService()
	if err != nil {
		logger.Errorf("Failed to init cluster service when describe cluster nodes: %v", err)
		return
	}

	output, err = clusterService.DescribeClusterNodes(&service.DescribeClusterNodesInput{Cluster: service.String(clusterId)})
	if err != nil {
		logger.Errorf("Failed to describe cluster [%s] nodes: %v", clusterId, err)
		return
	}

	logger.Printf("count: %+v", service.IntValue(output.TotalCount))

	return
}
