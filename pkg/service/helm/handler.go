package helm

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func getChartAndAppId(ctx context.Context, versionId string) (*chart.Chart, string, error) {
	var err error

	client, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, "", err
	}

	req := &pb.GetAppVersionPackageRequest{
		VersionId: pbutil.ToProtoString(versionId),
	}

	resp, err := client.GetAppVersionPackage(ctx, req)
	if err != nil {
		return nil, "", err
	}

	pkg := resp.GetPackage()
	r := bytes.NewReader(pkg)
	chrt, err := loader.LoadArchive(r)
	if err != nil {
		return nil, "", err
	}

	return chrt, resp.GetAppId().GetValue(), nil
}

func (p *Server) ParseClusterConf(ctx context.Context, req *pb.ParseClusterConfRequest) (*pb.ParseClusterConfResponse, error) {
	versionId := req.GetVersionId().GetValue()
	runtimeId := req.GetRuntimeId().GetValue()
	conf := req.GetConf().GetValue()
	cluster := models.PbToClusterWrapper(req.GetCluster())

	c, appId, err := getChartAndAppId(ctx, versionId)
	if err != nil {
		logger.Error(ctx, "Load helm chart from app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
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
	err = parser.Parse(cluster, appId)
	if err != nil {
		logger.Error(ctx, "Parse app version [%s] failed: %+v", versionId, err)
		return nil, err
	}

	return &pb.ParseClusterConfResponse{
		Cluster: models.ClusterWrapperToPb(cluster),
	}, err
}

func (p *Server) SplitJobIntoTasks(ctx context.Context, req *pb.SplitJobIntoTasksRequest) (*pb.SplitJobIntoTasksResponse, error) {
	job := models.PbToJob(req.GetJob())
	jobDirective, err := decodeJobDirective(ctx, job.Directive)
	if err != nil {
		return nil, err
	}

	tl := new(models.TaskLayer)

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

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionCreateCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
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

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionUpgradeCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
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

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionUpgradeCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionRollbackCluster:
		td := TaskDirective{
			Namespace:         jobDirective.Namespace,
			RuntimeId:         jobDirective.RuntimeId,
			ClusterName:       jobDirective.ClusterName,
			RawClusterWrapper: job.Directive,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionRollbackCluster, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionDeleteClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionDeleteClusters, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	case constants.ActionCeaseClusters:
		td := TaskDirective{
			RuntimeId:   jobDirective.RuntimeId,
			ClusterName: jobDirective.ClusterName,
		}
		tdj := encodeTaskDirective(td)

		task := models.NewTask(constants.PlaceHolder, job.JobId, "", jobDirective.RuntimeId, constants.ActionCeaseClusters, tdj, sender.OwnerPath(job.OwnerPath), false)
		tl = &models.TaskLayer{
			Tasks: []*models.Task{task},
			Child: nil,
		}
	default:
		return nil, fmt.Errorf("the job action [%s] is not supported", job.JobAction)
	}
	return &pb.SplitJobIntoTasksResponse{
		TaskLayer: models.TaskLayerToPb(tl),
	}, nil
}

func (s *Server) HandleSubtask(ctx context.Context, req *pb.HandleSubtaskRequest) (*pb.HandleSubtaskResponse, error) {
	task := models.PbToTask(req.GetTask())
	directive, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return nil, err
	}

	proxy := NewProxy(ctx, directive.RuntimeId)
	//todo HELM_DRIVER
	cfg, _ := proxy.GetHelmConfig("configmap", directive.Namespace, DefaultCredentialGetter)

	switch task.TaskAction {
	case constants.ActionCreateCluster:
		c, _, err := getChartAndAppId(ctx, directive.VersionId)
		if err != nil {
			return nil, err
		}

		rawVals := make(map[string]interface{})
		err = jsonutil.Decode([]byte(directive.Values), &rawVals)
		if err != nil {
			return nil, err
		}
		logger.Debug(ctx, "Install helm release with name [%+v], namespace [%+v], values [%s]", directive.ClusterName, directive.Namespace, rawVals)

		err = proxy.InstallReleaseFromChart(cfg, c, rawVals, directive.ClusterName)
		if err != nil {
			return nil, err
		}
	case constants.ActionUpgradeCluster:
		c, _, err := getChartAndAppId(ctx, directive.VersionId)
		if err != nil {
			return nil, err
		}
		rawVals := make(map[string]interface{})
		err = jsonutil.Decode([]byte(directive.Values), &rawVals)
		if err != nil {
			return nil, err
		}

		logger.Debug(ctx, "Update helm release [%+v] with values [%s]", directive.ClusterName, rawVals)
		err = proxy.UpdateReleaseFromChart(cfg, directive.ClusterName, c, rawVals)
		if err != nil {
			return nil, err
		}

	case constants.ActionRollbackCluster:
		err = proxy.RollbackRelease(cfg, directive.ClusterName)
		if err != nil {
			return nil, err
		}
	case constants.ActionDeleteClusters:
		err = proxy.DeleteRelease(cfg, directive.ClusterName, false)
		if err != nil {
			return nil, err
		}
	case constants.ActionCeaseClusters:
		// todo about namespaces
		//cfg, err = proxy.GetHelmConfig("configmap", DefaultCredentialGetter)
		err = proxy.DeleteRelease(cfg, directive.ClusterName, true)
		if err != nil {
			logger.Debug(ctx, "Cease helm release [%+v] error: [%s]", directive.ClusterName, err.Error())
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the task action [%s] is not supported", task.TaskAction)
	}
	return &pb.HandleSubtaskResponse{
		Task: models.TaskToPb(task),
	}, nil
}

func (s *Server) WaitSubtask(ctx context.Context, req *pb.WaitSubtaskRequest) (*pb.WaitSubtaskResponse, error) {
	task := models.PbToTask(req.GetTask())
	taskDirective, err := decodeTaskDirective(task.Directive)
	if err != nil {
		return nil, err
	}

	proxy := NewProxy(ctx, taskDirective.RuntimeId)
	//todo HELMDRIVER
	cfg, err := proxy.GetHelmConfig("configmap", taskDirective.Namespace, DefaultCredentialGetter)
	if err != nil {
		logger.Debug(ctx, "get helm action config error [%s]", err.Error())
		return nil, err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		switch task.TaskAction {
		case constants.ActionCreateCluster:
			fallthrough
		case constants.ActionUpgradeCluster:
			fallthrough
		case constants.ActionRollbackCluster:
			status, err := proxy.ReleaseStatus(cfg, taskDirective.ClusterName)
			if err != nil {
				if isConnectionError(err) {
					return false, nil
				}
				return true, err
			}

			switch status {
			case release.StatusFailed:
				logger.Debug(ctx, "Helm release gone to failed")
				return true, fmt.Errorf("release failed")
			case release.StatusDeployed:
				clusterWrapper, err := models.NewClusterWrapper(ctx, taskDirective.RawClusterWrapper)
				if err != nil {
					return true, err
				}
				err = proxy.WaitWorkloadReady(
					taskDirective.RuntimeId,
					taskDirective.Namespace,
					clusterWrapper.ClusterRoles,
					task.GetTimeout(constants.WaitTaskTimeout),
					constants.WaitTaskInterval,
				)
				if err != nil {
					return true, err
				}

				return true, nil
			}
		case constants.ActionDeleteClusters:
			status, err := proxy.ReleaseStatus(cfg, taskDirective.ClusterName)
			if err != nil {
				if isConnectionError(err) {
					return false, nil
				}
				if strings.Contains(err.Error(), "not found") {
					logger.Warn(nil, "Waiting on a helm release not existed, %+v", err)
					return true, nil
				}
				return true, err
			}

			if status == release.StatusUninstalled {
				return true, nil
			}
		case constants.ActionCeaseClusters:
			_, err := proxy.ReleaseStatus(cfg, taskDirective.ClusterName)
			if err != nil {
				if isConnectionError(err) {
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

	if err != nil {
		return nil, err
	} else {
		return &pb.WaitSubtaskResponse{
			Task: models.TaskToPb(task),
		}, nil
	}
}

func (p *Server) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	return nil, fmt.Errorf("the action DescribeSubnets is not supported")
}

func (p *Server) CheckResource(ctx context.Context, req *pb.CheckResourceRequest) (*pb.CheckResourceResponse, error) {
	cluster := models.PbToClusterWrapper(req.GetCluster())
	proxy := NewProxy(ctx, cluster.Cluster.RuntimeId)

	err := proxy.CheckClusterNameIsUnique(cluster.Cluster.Name, cluster.Cluster.Zone)
	if err != nil {
		logger.Error(ctx, "Cluster name [%s] already existed in runtime [%s]: %+v",
			cluster.Cluster.Name, cluster.Cluster.RuntimeId, err)
		return &pb.CheckResourceResponse{
			Ok: pbutil.ToProtoBool(false),
		}, err
	} else {
		return &pb.CheckResourceResponse{
			Ok: pbutil.ToProtoBool(true),
		}, nil
	}
}

func (p *Server) DescribeVpc(ctx context.Context, req *pb.DescribeVpcRequest) (*pb.DescribeVpcResponse, error) {
	return nil, fmt.Errorf("the action DescribeSubnets is not supported")
}

func (p *Server) DescribeClusterDetails(ctx context.Context, req *pb.DescribeClusterDetailsRequest) (*pb.DescribeClusterDetailsResponse, error) {
	cluster := models.PbToClusterWrapper(req.GetCluster())
	proxy := NewProxy(ctx, cluster.Cluster.RuntimeId)
	err := proxy.DescribeClusterDetails(cluster)
	cluster.Cluster.AdditionalInfo = jsonutil.ToString(proxy.WorkloadInfo)
	return &pb.DescribeClusterDetailsResponse{
		Cluster: models.ClusterWrapperToPb(cluster),
	}, err
}

func (p *Server) ValidateRuntime(ctx context.Context, req *pb.ValidateRuntimeRequest) (*pb.ValidateRuntimeResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	zone := req.GetZone().GetValue()
	needCreate := req.GetNeedCreate().GetValue()
	runtimeCredential := models.PbToRuntimeCredential(req.GetRuntimeCredential())

	proxy := NewProxy(ctx, runtimeId)
	err := proxy.ValidateRuntime(zone, runtimeCredential, needCreate)
	if err != nil {
		return &pb.ValidateRuntimeResponse{
			Ok: pbutil.ToProtoBool(false),
		}, err
	} else {
		return &pb.ValidateRuntimeResponse{
			Ok: pbutil.ToProtoBool(true),
		}, nil
	}
}

func (p *Server) DescribeZones(ctx context.Context, req *pb.DescribeZonesRequest) (*pb.DescribeZonesResponse, error) {
	runtimeCredential := models.PbToRuntimeCredential(req.GetRuntimeCredential())
	proxy := NewProxy(ctx, "")
	zones, err := proxy.DescribeRuntimeProviderZones(runtimeCredential)
	return &pb.DescribeZonesResponse{
		Zones: zones,
	}, err
}
