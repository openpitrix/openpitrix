// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/transport"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Provider struct {
	ctx context.Context
}

func NewProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx,
	}
}

func (p *Provider) getChart(versionId string) (*chart.Chart, error) {
	ctx := clientutil.SetSystemUserToContext(p.ctx)
	appClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, err
	}

	req := pb.GetAppVersionPackageRequest{
		VersionId: pbutil.ToProtoString(versionId),
	}

	resp, err := appClient.GetAppVersionPackage(ctx, &req)
	if err != nil {
		return nil, err
	}

	pkg := resp.GetPackage()
	r := bytes.NewReader(pkg)

	c, err := chartutil.LoadArchive(r)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p *Provider) ParseClusterConf(versionId, runtimeId, conf string) (*models.ClusterWrapper, error) {
	c, err := p.getChart(versionId)
	if err != nil {
		logger.Error(p.ctx, "Load helm chart from app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	parser := Parser{
		ctx:       p.ctx,
		Chart:     c,
		Conf:      conf,
		VersionId: versionId,
		RuntimeId: runtimeId,
	}
	clusterWrapper, err := parser.Parse()
	if err != nil {
		logger.Error(p.ctx, "Parse app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	return clusterWrapper, nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	jobDirective, err := decodeJobDirective(p.ctx, job.Directive)
	if err != nil {
		return nil, err
	}

	switch job.JobAction {
	case constants.ActionCreateCluster:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionCreateCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionUpgradeCluster:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionUpgradeCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionUpdateClusterEnv:
		td := TaskDirective{
			VersionId:         job.VersionId,
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			Values:            jobDirective.Values,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionUpgradeCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionRollbackCluster:
		td := TaskDirective{
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionRollbackCluster, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionDeleteClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionDeleteClusters, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	case constants.ActionCeaseClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", constants.ProviderKubernetes, constants.ActionCeaseClusters, tdj, job.Owner, false)
		tl := models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}

		return &tl, nil
	default:
		return nil, fmt.Errorf("the job action [%s] is not supported", job.JobAction)
	}
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	helmHandler := GetHelmHandler(p.ctx, taskDirective.RuntimeId)

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		c, err := p.getChart(taskDirective.VersionId)
		if err != nil {
			return err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return err
		}

		logger.Debug(p.ctx, "Install helm release with name [%+v], namespace [%+v], values [%s]", taskDirective.ClusterName, taskDirective.Namespace, rawVals)

		err = helmHandler.InstallReleaseFromChart(c, taskDirective.Namespace, rawVals, taskDirective.ClusterName)
		if err != nil {
			return err
		}
	case constants.ActionUpgradeCluster:
		c, err := p.getChart(taskDirective.VersionId)
		if err != nil {
			return err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return err
		}

		logger.Debug(p.ctx, "Update helm release [%+v] with values [%s]", taskDirective.ClusterName, rawVals)

		err = helmHandler.UpdateReleaseFromChart(taskDirective.ClusterName, c, rawVals)
		if err != nil {
			return err
		}
	case constants.ActionRollbackCluster:
		err = helmHandler.RollbackRelease(taskDirective.ClusterName)
		if err != nil {
			return err
		}
	case constants.ActionDeleteClusters:
		err = helmHandler.DeleteRelease(taskDirective.ClusterName, false)
		if err != nil {
			return err
		}
	case constants.ActionCeaseClusters:
		err = helmHandler.DeleteRelease(taskDirective.ClusterName, true)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("the task action [%s] is not supported", task.TaskAction)
	}

	return nil
}

func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	helmHandler := GetHelmHandler(p.ctx, taskDirective.RuntimeId)

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		switch task.TaskAction {
		case constants.ActionCreateCluster:
			fallthrough
		case constants.ActionUpgradeCluster:
			fallthrough
		case constants.ActionRollbackCluster:
			resp, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				return true, err
			}

			switch resp.Info.Status.Code {
			case release.Status_FAILED:
				logger.Debug(p.ctx, "Helm release gone to failed")
				return true, fmt.Errorf("release failed")
			case release.Status_DEPLOYED:
				clusterWrapper, err := models.NewClusterWrapper(p.ctx, taskDirective.RawClusterWrapper)
				if err != nil {
					return true, err
				}

				kubeHandler := GetKubeHandler(p.ctx, taskDirective.RuntimeId)
				err = kubeHandler.WaitPodsRunning(taskDirective.RuntimeId, taskDirective.Namespace, clusterWrapper.ClusterRoles, timeout, waitInterval)
				if err != nil {
					return true, err
				}

				return true, nil
			}
		case constants.ActionDeleteClusters:
			resp, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				return true, err
			}

			if resp.Info.Status.Code == release.Status_DELETED {
				return true, nil
			}
		case constants.ActionCeaseClusters:
			_, err := helmHandler.ReleaseStatus(taskDirective.ClusterName)
			if err != nil {
				if _, ok := err.(transport.ConnectionError); ok {
					return false, nil
				}
				return true, nil
			}
		}
		return false, nil
	}, timeout, waitInterval)

	return err
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	return nil, nil
}

func (p *Provider) CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	helmHandler := GetHelmHandler(p.ctx, clusterWrapper.Cluster.RuntimeId)

	err := helmHandler.CheckClusterNameIsUnique(clusterWrapper.Cluster.Name)
	if err != nil {
		logger.Error(p.ctx, "Cluster name [%s] already existed in runtime [%s]: %+v",
			clusterWrapper.Cluster.Name, clusterWrapper.Cluster.RuntimeId, err)
		return err
	}
	return nil
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	return nil, nil
}

func (p *Provider) updateClusterEnv(job *models.Job) error {
	clusterWrapper, err := models.NewClusterWrapper(p.ctx, job.Directive)
	if err != nil {
		return err
	}

	ctx := clientutil.SetSystemUserToContext(p.ctx)
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return err
	}

	var clusterRoles []*models.ClusterRole
	for _, clusterRole := range clusterWrapper.ClusterRoles {
		clusterRoles = append(clusterRoles, clusterRole)
	}

	modifyClusterRequest := &pb.ModifyClusterRequest{
		Cluster: &pb.Cluster{
			ClusterId:   pbutil.ToProtoString(job.ClusterId),
			Description: pbutil.ToProtoString(clusterWrapper.Cluster.Description),
		},
		ClusterRoleSet: models.ClusterRolesToPbs(clusterRoles),
	}
	_, err = clusterClient.ModifyCluster(ctx, modifyClusterRequest)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) updateClusterNodes(job *models.Job) error {
	clusterWrapper, err := models.NewClusterWrapper(p.ctx, job.Directive)
	if err != nil {
		return err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId

	runtime, err := runtimeclient.NewRuntime(p.ctx, runtimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone

	kubeHandler := GetKubeHandler(p.ctx, runtimeId)
	pbClusterNodes, err := kubeHandler.GetKubePodsAsClusterNodes(namespace, job.ClusterId, job.Owner, clusterWrapper.ClusterRoles)
	if err != nil {
		logger.Error(p.ctx, "Get kubernetes pods failed, %+v", err)
		return err
	}

	ctx := clientutil.SetSystemUserToContext(p.ctx)
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return err
	}

	// get all old node ids
	describeNodesRequest := &pb.DescribeClusterNodesRequest{
		ClusterId: pbutil.ToProtoString(job.ClusterId),
	}
	describeNodesResponse, err := clusterClient.DescribeClusterNodes(ctx, describeNodesRequest)
	if err != nil {
		logger.Error(p.ctx, "Get old nodes failed, %+v", err)
		return err
	}
	var nodeIds []string
	for _, clusterNode := range describeNodesResponse.ClusterNodeSet {
		nodeIds = append(nodeIds, clusterNode.GetNodeId().GetValue())
	}

	if len(nodeIds) != 0 {
		// delete old nodes from table
		deleteNodesRequest := &pb.DeleteTableClusterNodesRequest{
			NodeId: nodeIds,
		}
		_, err = clusterClient.DeleteTableClusterNodes(ctx, deleteNodesRequest)
		if err != nil {
			logger.Error(p.ctx, "Delete old nodes failed, %+v", err)
		}
	}

	if len(pbClusterNodes) != 0 {
		// add new nodes into table
		addNodesRequest := &pb.AddTableClusterNodesRequest{
			ClusterNodeSet: pbClusterNodes,
		}
		_, err = clusterClient.AddTableClusterNodes(ctx, addNodesRequest)
		if err != nil {
			logger.Error(p.ctx, "Add new nodes failed, %+v", err)
		}
	}
	return nil
}

func (p *Provider) UpdateClusterStatus(job *models.Job) error {
	switch job.JobAction {
	case constants.ActionCreateCluster:
		fallthrough
	case constants.ActionUpgradeCluster:
		fallthrough
	case constants.ActionRollbackCluster:
		err := p.updateClusterNodes(job)
		if err != nil {
			logger.Error(p.ctx, "Update cluster nodes failed, %+v", err)
			return err
		}
	case constants.ActionUpdateClusterEnv:
		err := p.updateClusterNodes(job)
		if err != nil {
			logger.Error(p.ctx, "Update cluster nodes failed, %+v", err)
			return err
		}

		err = p.updateClusterEnv(job)
		if err != nil {
			logger.Error(p.ctx, "Update cluster env failed, %+v", err)
			return err
		}
	}

	return nil
}

func (p *Provider) ValidateCredential(url, credential, zone string) error {
	kubeHandler := GetKubeHandler(p.ctx, "")
	return kubeHandler.ValidateCredential(credential, zone)
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) ([]string, error) {
	kubeHandler := GetKubeHandler(p.ctx, "")
	return kubeHandler.DescribeRuntimeProviderZones(credential)
}
