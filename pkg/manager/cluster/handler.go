// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jobClient "openpitrix.io/openpitrix/pkg/client/job"
	runtimeEnvClient "openpitrix.io/openpitrix/pkg/client/runtime_env"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) getRuntimeEnv(runtimeEnvId string) (*pb.RuntimeEnv, error) {
	runtimeEnvIds := []string{runtimeEnvId}
	response, err := runtimeEnvClient.DescribeRuntimeEnvs(&pb.DescribeRuntimeEnvsRequest{
		RuntimeEnvId: runtimeEnvIds,
	})
	if err != nil {
		logger.Errorf("Describe runtime env [%s] failed: %+v",
			strings.Join(runtimeEnvIds, ","), err)
		return nil, status.Errorf(codes.Internal, "Describe runtime env [%s] failed: %+v",
			strings.Join(runtimeEnvIds, ","), err)
	}

	if response.GetTotalCount().GetValue() == 0 {
		logger.Errorf("Runtime env [%s] not found", strings.Join(runtimeEnvIds, ","))
		return nil, status.Errorf(codes.PermissionDenied, "Runtime env [%s] not found",
			strings.Join(runtimeEnvIds, ","), err)
	}

	return response.RuntimeEnvSet[0], nil
}

func (p *Server) getRuntime(runtimeEnv *pb.RuntimeEnv) string {
	// TODO: need to parse runtime
	return runtimeEnv.GetLabels().GetValue()
}

func (p *Server) getCluster(clusterId, userId string) (*models.Cluster, error) {
	cluster := &models.Cluster{}
	err := p.Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Where(db.Eq("owner", userId)).
		LoadOne(&cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (p *Server) getClusterNode(nodeId, userId string) (*models.ClusterNode, error) {
	clusterNode := &models.ClusterNode{}
	err := p.Db.
		Select(models.ClusterNodeColumns...).
		From(models.ClusterNodeTableName).
		Where(db.Eq("node_id", nodeId)).
		Where(db.Eq("owner", userId)).
		LoadOne(&clusterNode)
	if err != nil {
		return nil, err
	}
	return clusterNode, nil
}

func (p *Server) CreateCluster(ctx context.Context, req *pb.CreateClusterRequest) (*pb.CreateClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	runtimeEnv, err := p.getRuntimeEnv(req.GetRuntimeEnvId().GetValue())
	if err != nil {
		return nil, err
	}

	runtime := p.getRuntime(runtimeEnv)
	runtimeInterface, err := plugins.GetRuntimePlugin(runtime)
	if err != nil {
		logger.Errorf("No such runtime [%s]. ", runtime)
		return nil, err
	}

	clusterId := models.NewClusterId()
	versionId := req.GetAppVersion().GetValue()
	conf := req.GetConf().GetValue()
	cluster, err := runtimeInterface.ParseClusterConf(versionId, conf)
	if err != nil {
		logger.Errorf("Parse cluster conf with versionId [%s] runtime [%s] failed. ", versionId, runtime)
		return nil, err
	}

	store := &Store{Pi: p.Pi}
	err = store.RegisterCluster(cluster)
	if err != nil {
		return nil, err
	}

	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusCreating).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusCreating, err)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		req.GetAppId().GetValue(),
		req.GetAppVersion().GetValue(),
		constants.ActionCreateCluster,
		"", // TODO: need to generate
		runtime,
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}
	res := &pb.CreateClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}
	return res, nil
}

func (p *Server) ModifyCluster(ctx context.Context, req *pb.ModifyClusterRequest) (*pb.ModifyClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	_, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get cluster [%s]", clusterId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "description", "status", "transition_status", "upgrade_status",
		"upgrade_time", "status_time")
	_, err = p.Db.
		Update(models.ClusterTableName).
		SetMap(attributes).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyCluster [%s]: %+v", clusterId, err)
	}

	res := &pb.ModifyClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
	}
	return res, nil
}

func (p *Server) ModifyClusterNode(ctx context.Context, req *pb.ModifyClusterNodeRequest) (*pb.ModifyClusterNodeResponse, error) {
	s := sender.GetSenderFromContext(ctx)

	nodeId := req.GetNodeId().GetValue()
	_, err := p.getClusterNode(nodeId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get cluster node [%s]", nodeId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "instance_id", "volume_id", "vxnet_id", "private_ip",
		"status", "transition_status", "health_status", "pub_key", "status_time")
	_, err = p.Db.
		Update(models.ClusterNodeTableName).
		SetMap(attributes).
		Where(db.Eq("node_id", nodeId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyClusterNode [%s]: %+v", nodeId, err)
	}

	res := &pb.ModifyClusterNodeResponse{
		NodeId: utils.ToProtoString(nodeId),
	}
	return res, nil
}

func (p *Server) DeleteClusters(ctx context.Context, req *pb.DeleteClustersRequest) (*pb.DeleteClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		cluster, err := p.getCluster(clusterId, s.UserId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}
		_, err = p.Db.
			Update(models.ClusterTableName).
			Set("transition_status", constants.StatusDeleting).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
				clusterId, constants.StatusDeleting, err)
		}

		runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
		if err != nil {
			return nil, err
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.AppVersion,
			constants.ActionDeleteClusters,
			"", // TODO: need to generate
			p.getRuntime(runtimeEnv),
			s.UserId,
		)

		jobId, err := jobClient.SendJob(newJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.DeleteClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) UpgradeCluster(ctx context.Context, req *pb.UpgradeClusterRequest) (*pb.UpgradeClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusUpgrading).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusUpgrading, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionUpgradeCluster,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.UpgradeClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) RollbackCluster(ctx context.Context, req *pb.RollbackClusterRequest) (*pb.RollbackClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusRollbacking).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusRollbacking, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionRollbackCluster,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.RollbackClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) ResizeCluster(ctx context.Context, req *pb.ResizeClusterRequest) (*pb.ResizeClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusResizing).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusResizing, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionResizeCluster,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.ResizeClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) AddClusterNodes(ctx context.Context, req *pb.AddClusterNodesRequest) (*pb.AddClusterNodesResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusScaling).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusScaling, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionAddClusterNodes,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.AddClusterNodesResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) DeleteClusterNodes(ctx context.Context, req *pb.DeleteClusterNodesRequest) (*pb.DeleteClusterNodesResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusScaling).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusScaling, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionDeleteClusterNodes,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.DeleteClusterNodesResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) UpdateClusterEnv(ctx context.Context, req *pb.UpdateClusterEnvRequest) (*pb.UpdateClusterEnvResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	clusterId := req.GetClusterId().GetValue()
	cluster, err := p.getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}
	_, err = p.Db.
		Update(models.ClusterTableName).
		Set("transition_status", constants.StatusUpdating).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
			clusterId, constants.StatusUpdating, err)
	}

	runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		cluster.AppId,
		cluster.AppVersion,
		constants.ActionUpdateClusterEnv,
		"", // TODO: need to generate
		p.getRuntime(runtimeEnv),
		s.UserId,
	)

	jobId, err := jobClient.SendJob(newJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.UpdateClusterEnvResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) DescribeClusters(ctx context.Context, req *pb.DescribeClustersRequest) (*pb.DescribeClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var clusters []*models.Cluster
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterTableName)).
		Where(db.Eq("owner", s.UserId))
	_, err := query.Load(&clusters)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeClusters: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeClusters: %+v", err)
	}

	res := &pb.DescribeClustersResponse{
		ClusterSet: models.ClustersToPbs(clusters),
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) DescribeClusterNodes(ctx context.Context, req *pb.DescribeClusterNodesRequest) (*pb.DescribeClusterNodesResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var clusterNodes []*models.ClusterNode
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.ClusterColumns...).
		From(models.ClusterNodeTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterNodeTableName)).
		Where(db.Eq("owner", s.UserId))
	_, err := query.Load(&clusterNodes)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeClusterNodes: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeClusterNodes: %+v", err)
	}

	res := &pb.DescribeClusterNodesResponse{
		ClusterNodeSet: models.ClusterNodesToPbs(clusterNodes),
		TotalCount:     count,
	}
	return res, nil
}

func (p *Server) StopClusters(ctx context.Context, req *pb.StopClustersRequest) (*pb.StopClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		cluster, err := p.getCluster(clusterId, s.UserId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}
		_, err = p.Db.
			Update(models.ClusterTableName).
			Set("transition_status", constants.StatusStopping).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
				clusterId, constants.StatusStopping, err)
		}

		runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
		if err != nil {
			return nil, err
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.AppVersion,
			constants.ActionStopClusters,
			"", // TODO: need to generate
			p.getRuntime(runtimeEnv),
			s.UserId,
		)

		jobId, err := jobClient.SendJob(newJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.StopClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) StartClusters(ctx context.Context, req *pb.StartClustersRequest) (*pb.StartClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		cluster, err := p.getCluster(clusterId, s.UserId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}
		_, err = p.Db.
			Update(models.ClusterTableName).
			Set("transition_status", constants.StatusStarting).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
				clusterId, constants.StatusStarting, err)
		}

		runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
		if err != nil {
			return nil, err
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.AppVersion,
			constants.ActionStartClusters,
			"", // TODO: need to generate
			p.getRuntime(runtimeEnv),
			s.UserId,
		)

		jobId, err := jobClient.SendJob(newJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.StartClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) RecoverClusters(ctx context.Context, req *pb.RecoverClustersRequest) (*pb.RecoverClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		cluster, err := p.getCluster(clusterId, s.UserId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}
		_, err = p.Db.
			Update(models.ClusterTableName).
			Set("transition_status", constants.StatusRecovering).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
				clusterId, constants.StatusRecovering, err)
		}

		runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
		if err != nil {
			return nil, err
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.AppVersion,
			constants.ActionRecoverClusters,
			"", // TODO: need to generate
			p.getRuntime(runtimeEnv),
			s.UserId,
		)

		jobId, err := jobClient.SendJob(newJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.RecoverClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) CeaseClusters(ctx context.Context, req *pb.CeaseClustersRequest) (*pb.CeaseClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		cluster, err := p.getCluster(clusterId, s.UserId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}
		_, err = p.Db.
			Update(models.ClusterTableName).
			Set("transition_status", constants.StatusCeasing).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Update cluster [%s] transition_status [%s] failed: %+v",
				clusterId, constants.StatusCeasing, err)
		}

		runtimeEnv, err := p.getRuntimeEnv(cluster.RuntimeEnvId)
		if err != nil {
			return nil, err
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.AppVersion,
			constants.ActionCeaseClusters,
			"", // TODO: need to generate
			p.getRuntime(runtimeEnv),
			s.UserId,
		)

		jobId, err := jobClient.SendJob(newJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.CeaseClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}
