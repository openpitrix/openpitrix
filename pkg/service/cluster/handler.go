// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"time"

	pb_empty "github.com/golang/protobuf/ptypes/empty"

	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func getCluster(clusterId, userId string) (*models.Cluster, error) {
	cluster := &models.Cluster{}
	err := pi.Global().Db.
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

func getClusterWrapper(clusterId string) (*models.ClusterWrapper, error) {
	clusterWrapper := new(models.ClusterWrapper)
	var cluster *models.Cluster
	var clusterCommons []*models.ClusterCommon
	var clusterNodes []*models.ClusterNode
	var clusterRoles []*models.ClusterRole
	var clusterLinks []*models.ClusterLink
	var clusterLoadbalancers []*models.ClusterLoadbalancer

	err := pi.Global().Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Where(db.Eq("cluster_id", clusterId)).
		LoadOne(&cluster)
	if err != nil {
		return nil, err
	}
	clusterWrapper.Cluster = cluster

	_, err = pi.Global().Db.
		Select(models.ClusterCommonColumns...).
		From(models.ClusterCommonTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterCommons)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterCommons = map[string]*models.ClusterCommon{}
	for _, clusterCommon := range clusterCommons {
		clusterWrapper.ClusterCommons[clusterCommon.Role] = clusterCommon
	}

	_, err = pi.Global().Db.
		Select(models.ClusterNodeColumns...).
		From(models.ClusterNodeTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterNodes)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterNodes = map[string]*models.ClusterNode{}
	for _, clusterNode := range clusterNodes {
		clusterWrapper.ClusterNodes[clusterNode.NodeId] = clusterNode
	}

	_, err = pi.Global().Db.
		Select(models.ClusterRoleColumns...).
		From(models.ClusterRoleTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterRoles)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterRoles = map[string]*models.ClusterRole{}
	for _, clusterRole := range clusterRoles {
		clusterWrapper.ClusterRoles[clusterRole.Role] = clusterRole
	}

	_, err = pi.Global().Db.
		Select(models.ClusterLinkColumns...).
		From(models.ClusterLinkTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterLinks)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterLinks = map[string]*models.ClusterLink{}
	for _, clusterLink := range clusterLinks {
		clusterWrapper.ClusterLinks[clusterLink.Name] = clusterLink
	}

	_, err = pi.Global().Db.
		Select(models.ClusterLoadbalancerColumns...).
		From(models.ClusterLoadbalancerTableName).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterLoadbalancers)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterLoadbalancers = map[string][]*models.ClusterLoadbalancer{}
	for _, clusterLoadbalancer := range clusterLoadbalancers {
		clusterWrapper.ClusterLoadbalancers[clusterLoadbalancer.Role] =
			append(clusterWrapper.ClusterLoadbalancers[clusterLoadbalancer.Role], clusterLoadbalancer)
	}

	return clusterWrapper, nil
}

func getClusterNode(nodeId, userId string) (*models.ClusterNode, error) {
	clusterNode := &models.ClusterNode{}
	err := pi.Global().Db.
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

func (p *Server) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, runtimeId)
	}

	providerInterface, err := plugins.GetProviderPlugin(runtime.Provider, nil)
	if err != nil {
		logger.Error("No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}

	return providerInterface.DescribeSubnets(ctx, req)
}

func (p *Server) CreateCluster(ctx context.Context, req *pb.CreateClusterRequest) (*pb.CreateClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	appId := req.GetAppId().GetValue()
	versionId := req.GetVersionId().GetValue()
	conf := req.GetConf().GetValue()
	clusterId := models.NewClusterId()
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, runtimeId)
	}

	providerInterface, err := plugins.GetProviderPlugin(runtime.Provider, nil)
	if err != nil {
		logger.Error("No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}
	clusterWrapper, err := providerInterface.ParseClusterConf(versionId, runtimeId, conf)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		}
		logger.Error("Parse cluster conf with versionId [%s] runtime [%s] failed: %+v", versionId, runtime, err)
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	clusterWrapper.Cluster.RuntimeId = runtimeId
	clusterWrapper.Cluster.Owner = s.UserId
	clusterWrapper.Cluster.ClusterId = clusterId
	clusterWrapper.Cluster.ClusterType = constants.NormalClusterType

	if reflectutil.In(runtime.Provider, constants.VmBaseProviders) {
		err = CheckVmBasedProvider(ctx, runtime, providerInterface, clusterWrapper)
		if err != nil {
			return nil, err
		}
	}

	err = RegisterClusterWrapper(clusterWrapper)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	directive := jsonutil.ToString(clusterWrapper)

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		appId,
		versionId,
		constants.ActionCreateCluster,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	res := &pb.CreateClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}
	return res, nil
}

func (p *Server) ModifyCluster(ctx context.Context, req *pb.ModifyClusterRequest) (*pb.ModifyClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetCluster().GetClusterId().GetValue()
	_, err := getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	attributes := manager.BuildUpdateAttributes(req.Cluster, models.ClusterColumns...)
	logger.Debug("ModifyCluster got attributes: [%+v]", attributes)
	delete(attributes, "cluster_id")
	_, err = pi.Global().Db.
		Update(models.ClusterTableName).
		SetMap(attributes).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
	}

	for _, clusterNode := range req.ClusterNodeSet {
		nodeId := clusterNode.GetNodeId().GetValue()
		nodeAttributes := manager.BuildUpdateAttributes(clusterNode, models.ClusterNodeColumns...)
		delete(nodeAttributes, "cluster_id")
		delete(nodeAttributes, "node_id")
		_, err = pi.Global().Db.
			Update(models.ClusterNodeTableName).
			SetMap(nodeAttributes).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("node_id", nodeId)).
			Exec()
		if err != nil {
			logger.Error("ModifyCluster [%s] node [%s] failed. ", clusterId, nodeId)
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	for _, clusterRole := range req.ClusterRoleSet {
		role := clusterRole.GetRole().GetValue()
		roleAttributes := manager.BuildUpdateAttributes(clusterRole, models.ClusterRoleColumns...)
		delete(roleAttributes, "cluster_id")
		delete(roleAttributes, "role")
		_, err = pi.Global().Db.
			Update(models.ClusterRoleTableName).
			SetMap(roleAttributes).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			Exec()
		if err != nil {
			logger.Error("ModifyCluster [%s] role [%s] failed. ", clusterId, role)
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	for _, clusterCommon := range req.ClusterCommonSet {
		role := clusterCommon.GetRole().GetValue()
		commonAttributes := manager.BuildUpdateAttributes(clusterCommon, models.ClusterCommonColumns...)
		delete(commonAttributes, "cluster_id")
		delete(commonAttributes, "role")
		_, err = pi.Global().Db.
			Update(models.ClusterCommonTableName).
			SetMap(commonAttributes).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			Exec()
		if err != nil {
			logger.Error("ModifyCluster [%s] role [%s] common failed. ", clusterId, role)
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	for _, clusterLink := range req.ClusterLinkSet {
		name := clusterLink.GetName().GetValue()
		linkAttributes := manager.BuildUpdateAttributes(clusterLink, models.ClusterLinkColumns...)
		delete(linkAttributes, "cluster_id")
		delete(linkAttributes, "name")
		_, err = pi.Global().Db.
			Update(models.ClusterLinkTableName).
			SetMap(linkAttributes).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("name", name)).
			Exec()
		if err != nil {
			logger.Error("ModifyCluster [%s] name [%s] link failed. ", clusterId, name)
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	for _, clusterLoadbalancer := range req.ClusterLoadbalancerSet {
		role := clusterLoadbalancer.GetRole().GetValue()
		listenerId := clusterLoadbalancer.GetLoadbalancerListenerId().GetValue()
		loadbalancerAttributes := manager.BuildUpdateAttributes(clusterLoadbalancer, models.ClusterLoadbalancerColumns...)
		delete(loadbalancerAttributes, "cluster_id")
		delete(loadbalancerAttributes, "role")
		delete(loadbalancerAttributes, "loadbalancer_listener_id")
		_, err = pi.Global().Db.
			Update(models.ClusterLoadbalancerTableName).
			SetMap(loadbalancerAttributes).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			Where(db.Eq("loadbalancer_listener_id", listenerId)).
			Exec()
		if err != nil {
			logger.Error("ModifyCluster [%s] role [%s] loadbalancer listener id [%s] failed. ",
				clusterId, role, listenerId)
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	res := &pb.ModifyClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
	}
	return res, nil
}

func (p *Server) ModifyClusterNode(ctx context.Context, req *pb.ModifyClusterNodeRequest) (*pb.ModifyClusterNodeResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	nodeId := req.GetClusterNode().GetNodeId().GetValue()
	_, err := getClusterNode(nodeId, s.UserId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
	}

	attributes := manager.BuildUpdateAttributes(req.ClusterNode, models.ClusterNodeColumns...)
	_, err = pi.Global().Db.
		Update(models.ClusterNodeTableName).
		SetMap(attributes).
		Where(db.Eq("node_id", nodeId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourceFailed, nodeId)
	}

	res := &pb.ModifyClusterNodeResponse{
		NodeId: pbutil.ToProtoString(nodeId),
	}
	return res, nil
}

func (p *Server) AddTableClusterNodes(ctx context.Context, req *pb.AddTableClusterNodesRequest) (*pb_empty.Empty, error) {
	for _, clusterNode := range req.ClusterNodeSet {
		node := models.PbToClusterNode(clusterNode)
		err := RegisterClusterNode(node)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return nil, nil
}

func (p *Server) DeleteTableClusterNodes(ctx context.Context, req *pb.DeleteTableClusterNodesRequest) (*pb_empty.Empty, error) {
	for _, nodeId := range req.NodeId {
		_, err := pi.Global().Db.
			DeleteFrom(models.ClusterNodeTableName).
			Where(db.Eq("node_id", nodeId)).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return nil, nil
}

func (p *Server) DeleteClusters(ctx context.Context, req *pb.DeleteClustersRequest) (*pb.DeleteClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive, constants.StatusStopped, constants.StatusPending})
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorDeleteResourceFailed, clusterId)
		}

		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionDeleteClusters,
			directive,
			runtime.Provider,
			s.UserId,
		)

		jobId, err := jobclient.SendJob(newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourceFailed, clusterId)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.DeleteClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) UpgradeCluster(ctx context.Context, req *pb.UpgradeClusterRequest) (*pb.UpgradeClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	versionId := req.GetVersionId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusStopped})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorUpgradeResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		versionId,
		constants.ActionUpgradeCluster,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorUpgradeResourceFailed, clusterId)
	}

	return &pb.UpgradeClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) RollbackCluster(ctx context.Context, req *pb.RollbackClusterRequest) (*pb.RollbackClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorRollbackResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionRollbackCluster,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRollbackResourceFailed, clusterId)
	}

	return &pb.RollbackClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) ResizeCluster(ctx context.Context, req *pb.ResizeClusterRequest) (*pb.ResizeClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResizeResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionResizeCluster,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorResizeResourceFailed, clusterId)
	}

	return &pb.ResizeClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) AddClusterNodes(ctx context.Context, req *pb.AddClusterNodesRequest) (*pb.AddClusterNodesResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionAddClusterNodes,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}

	return &pb.AddClusterNodesResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) DeleteClusterNodes(ctx context.Context, req *pb.DeleteClusterNodesRequest) (*pb.DeleteClusterNodesResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorDeleteResourceNodeFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionDeleteClusterNodes,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourceNodeFailed, clusterId)
	}

	return &pb.DeleteClusterNodesResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) UpdateClusterEnv(ctx context.Context, req *pb.UpdateClusterEnvRequest) (*pb.UpdateClusterEnvResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorUpdateResourceEnvFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionUpdateClusterEnv,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorUpdateResourceEnvFailed, clusterId)
	}

	return &pb.UpdateClusterEnvResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) DescribeClusters(ctx context.Context, req *pb.DescribeClustersRequest) (*pb.DescribeClustersResponse, error) {
	var clusters []*models.Cluster
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterTableName))
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err := query.Load(&clusters)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var pbClusters []*pb.Cluster
	for _, cluster := range clusters {
		clusterId := cluster.ClusterId
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}
		pbClusters = append(pbClusters, models.ClusterWrapperToPb(clusterWrapper))
	}

	res := &pb.DescribeClustersResponse{
		ClusterSet: pbClusters,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) DescribeClusterNodes(ctx context.Context, req *pb.DescribeClusterNodesRequest) (*pb.DescribeClusterNodesResponse, error) {
	var clusterNodes []*models.ClusterNode
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().Db.
		Select(models.ClusterNodeColumns...).
		From(models.ClusterNodeTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterNodeTableName))
	_, err := query.Load(&clusterNodes)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var pbClusterNodes []*pb.ClusterNode
	for _, clusterNode := range clusterNodes {
		var clusterCommon *models.ClusterCommon
		var clusterRole *models.ClusterRole
		nodeId := clusterNode.NodeId
		role := clusterNode.Role
		clusterId := clusterNode.ClusterId
		err = pi.Global().Db.
			Select(models.ClusterCommonColumns...).
			From(models.ClusterCommonTableName).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			LoadOne(&clusterCommon)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
		}

		err = pi.Global().Db.
			Select(models.ClusterRoleColumns...).
			From(models.ClusterRoleTableName).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			LoadOne(&clusterRole)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
		}

		pbClusterNodes = append(pbClusterNodes,
			models.ClusterNodeWrapperToPb(clusterNode, clusterCommon, clusterRole))
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeClusterNodesResponse{
		ClusterNodeSet: pbClusterNodes,
		TotalCount:     count,
	}
	return res, nil
}

func (p *Server) StopClusters(ctx context.Context, req *pb.StopClustersRequest) (*pb.StopClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusActive})
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorStopResourceFailed, clusterId)
		}
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionStopClusters,
			directive,
			runtime.Provider,
			s.UserId,
		)

		jobId, err := jobclient.SendJob(newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorStopResourceFailed, clusterId)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.StopClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) StartClusters(ctx context.Context, req *pb.StartClustersRequest) (*pb.StartClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusStopped})
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorStartResourceFailed, clusterId)
		}
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorStartResourceFailed, clusterId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionStartClusters,
			directive,
			runtime.Provider,
			s.UserId,
		)

		jobId, err := jobclient.SendJob(newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorStartResourceFailed, clusterId)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.StartClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) RecoverClusters(ctx context.Context, req *pb.RecoverClustersRequest) (*pb.RecoverClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusDeleted})
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorRecoverResourceFailed, clusterId)
		}
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRecoverResourceFailed, clusterId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionRecoverClusters,
			directive,
			runtime.Provider,
			s.UserId,
		)

		jobId, err := jobclient.SendJob(newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRecoverResourceFailed, clusterId)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.RecoverClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

func (p *Server) CeaseClusters(ctx context.Context, req *pb.CeaseClustersRequest) (*pb.CeaseClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		err := checkPermissionAndTransition(clusterId, s.UserId, []string{constants.StatusDeleted})
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorCeaseResourceFailed, clusterId)
		}
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionCeaseClusters,
			directive,
			runtime.Provider,
			s.UserId,
		)

		jobId, err := jobclient.SendJob(newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCeaseResourceFailed, clusterId)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.CeaseClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}

type clusterStatistic struct {
	Date  string `db:"DATE_FORMAT(create_time, '%Y-%m-%d')"`
	Count uint32 `db:"COUNT(cluster_id)"`
}
type runtimeStatistic struct {
	RuntimeId string `db:"runtime_id"`
	Count     uint32 `db:"COUNT(cluster_id)"`
}

func (p *Server) GetClusterStatistics(ctx context.Context, req *pb.GetClusterStatisticsRequest) (*pb.GetClusterStatisticsResponse, error) {
	res := &pb.GetClusterStatisticsResponse{
		LastTwoWeekCreated: make(map[string]uint32),
		TopTenRuntimes:     make(map[string]uint32),
	}
	clusterCount, err := pi.Global().Db.Select(models.ColumnClusterId).From(models.ClusterTableName).Count()
	if err != nil {
		logger.Error("Failed to get cluster count, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.ClusterCount = clusterCount

	err = pi.Global().Db.
		Select("COUNT(DISTINCT runtime_id)").
		From(models.ClusterTableName).
		LoadOne(&res.RuntimeCount)
	if err != nil {
		logger.Error("Failed to get runtime count, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var cs []*clusterStatistic
	_, err = pi.Global().Db.
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(cluster_id)").
		From(models.ClusterTableName).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(models.ColumnCreateTime, time2week)).
		Limit(14).Load(&cs)

	if err != nil {
		logger.Error("Failed to get cluster statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range cs {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var rs []*runtimeStatistic
	_, err = pi.Global().Db.
		Select("runtime_id", "COUNT(cluster_id)").
		From(models.ClusterTableName).
		GroupBy(models.ColumnRuntimeId).
		OrderDir("COUNT(cluster_id)", false).
		Limit(10).Load(&rs)

	if err != nil {
		logger.Error("Failed to get runtime statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.TopTenRuntimes[a.RuntimeId] = a.Count
	}

	return res, nil
}
