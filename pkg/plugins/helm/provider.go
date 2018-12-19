// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/transport"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Provider struct{}

func NewProvider() *Provider {
	return new(Provider)
}

func (p *Provider) getChart(ctx context.Context, versionId string) (*chart.Chart, error) {
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

func (p *Provider) ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error {
	c, err := p.getChart(ctx, versionId)
	if err != nil {
		logger.Error(ctx, "Load helm chart from app version [%s] failed: %+v", versionId, err)
		return err
	}

	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return err
	}
	namespace := runtime.Zone

	parser := Parser{
		ctx:       ctx,
		Chart:     c,
		Conf:      conf,
		VersionId: versionId,
		RuntimeId: runtimeId,
		Namespace: namespace,
	}
	err = parser.Parse(clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Parse app version [%s] failed: %+v", versionId, err)
		return err
	}

	return nil
}

func (p *Provider) SplitJobIntoTasks(ctx context.Context, job *models.Job) (*models.TaskLayer, error) {
	jobDirective, err := decodeJobDirective(ctx, job.Directive)
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

func (p *Provider) HandleSubtask(ctx context.Context, task *models.Task) error {
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	helmHandler := GetHelmHandler(ctx, taskDirective.RuntimeId)

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		c, err := p.getChart(ctx, taskDirective.VersionId)
		if err != nil {
			return err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return err
		}

		logger.Debug(ctx, "Install helm release with name [%+v], namespace [%+v], values [%s]", taskDirective.ClusterName, taskDirective.Namespace, rawVals)

		err = helmHandler.InstallReleaseFromChart(c, taskDirective.Namespace, rawVals, taskDirective.ClusterName)
		if err != nil {
			return err
		}
	case constants.ActionUpgradeCluster:
		c, err := p.getChart(ctx, taskDirective.VersionId)
		if err != nil {
			return err
		}

		rawVals, err := ConvertJsonToYaml([]byte(taskDirective.Values))
		if err != nil {
			return err
		}

		logger.Debug(ctx, "Update helm release [%+v] with values [%s]", taskDirective.ClusterName, rawVals)

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

func (p *Provider) WaitSubtask(ctx context.Context, task *models.Task) error {
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return err
	}

	helmHandler := GetHelmHandler(ctx, taskDirective.RuntimeId)

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
				logger.Debug(ctx, "Helm release gone to failed")
				return true, fmt.Errorf("release failed")
			case release.Status_DEPLOYED:
				clusterWrapper, err := models.NewClusterWrapper(ctx, taskDirective.RawClusterWrapper)
				if err != nil {
					return true, err
				}

				kubeHandler := GetKubeHandler(ctx, taskDirective.RuntimeId)
				err = kubeHandler.WaitWorkloadReady(taskDirective.RuntimeId, taskDirective.Namespace, clusterWrapper.ClusterRoles, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
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
				if strings.Contains(err.Error(), "not found") {
					logger.Warn(nil, "Waiting on a helm release not existed, %+v", err)
					return true, nil
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
				if strings.Contains(err.Error(), "not found") {
					logger.Warn(nil, "Waiting on a helm release not existed, %+v", err)
					return true, nil
				}
				return true, nil
			}
		}
		return false, nil
	}, task.GetTimeout(constants.WaitHelmTaskTimeout), constants.WaitTaskInterval)

	return err
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	return nil, nil
}

func (p *Provider) CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	helmHandler := GetHelmHandler(ctx, clusterWrapper.Cluster.RuntimeId)

	err := helmHandler.CheckClusterNameIsUnique(clusterWrapper.Cluster.Name)
	if err != nil {
		logger.Error(ctx, "Cluster name [%s] already existed in runtime [%s]: %+v",
			clusterWrapper.Cluster.Name, clusterWrapper.Cluster.RuntimeId, err)
		return err
	}
	return nil
}

func (p *Provider) DescribeVpc(ctx context.Context, runtimeId, vpcId string) (*models.Vpc, error) {
	return nil, nil
}

func (p *Provider) ValidateRuntime(ctx context.Context, runtimeId, zone string, runtimeCredential *models.RuntimeCredential, needCreate bool) error {
	kubeHandler := GetKubeHandler(ctx, runtimeId)
	return kubeHandler.ValidateRuntime(zone, runtimeCredential, needCreate)
}

func (p *Provider) DescribeRuntimeProviderZones(ctx context.Context, runtimeCredential *models.RuntimeCredential) ([]string, error) {
	kubeHandler := GetKubeHandler(ctx, "")
	return kubeHandler.DescribeRuntimeProviderZones(runtimeCredential)
}

func (p *Provider) DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) (*models.ClusterWrapper, error) {
	kubeHandler := GetKubeHandler(ctx, cluster.Cluster.RuntimeId)
	return cluster, kubeHandler.DescribeClusterDetails(cluster)
}
