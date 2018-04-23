// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
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

func (p *Server) CreateCluster(ctx context.Context, req *pb.CreateClusterRequest) (*pb.CreateClusterResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	// check image
	if utils.In(runtime.Provider, constants.VmBaseProviders) {
		_, err = pi.Global().GlobalConfig().GetRuntimeImageId(runtime.RuntimeUrl, runtime.Zone)
		if err != nil {
			return nil, err
		}
	}

	appId := req.GetAppId().GetValue()
	versionId := req.GetVersionId().GetValue()
	conf := req.GetConf().GetValue()

	clusterId := models.NewClusterId()

	providerInterface, err := plugins.GetProviderPlugin(runtime.Provider)
	if err != nil {
		logger.Errorf("No such provider [%s]. ", runtime.Provider)
		return nil, err
	}
	clusterWrapper, err := providerInterface.ParseClusterConf(versionId, conf)
	if err != nil {
		logger.Errorf("Parse cluster conf with versionId [%s] runtime [%s] failed. ", versionId, runtime)
		return nil, err
	}

	subnet, err := providerInterface.DescribeSubnet(runtimeId, clusterWrapper.Cluster.SubnetId)
	if err != nil {
		logger.Errorf("Describe subnet [%s] runtime [%s] failed. ", clusterWrapper.Cluster.SubnetId, runtime)
		return nil, err
	}
	vpcId := subnet.VpcId

	register := &Register{
		SubnetId: clusterWrapper.Cluster.SubnetId,
		VpcId:    vpcId,
		Runtime:  runtime,
		Owner:    s.UserId,
	}
	fg := &Frontgate{
		Runtime: runtime,
	}
	frontgate, err := fg.GetActiveFrontgate(vpcId, s.UserId, register)
	if err != nil {
		logger.Errorf("Get frontgate in vpc [%s] user [%s] failed. ", vpcId, s.UserId)
		return nil, err
	}

	register.ClusterId = clusterId
	register.FrontgateId = frontgate.ClusterId
	register.ClusterType = constants.NormalClusterType
	register.ClusterWrapper = clusterWrapper

	err = register.RegisterClusterWrapper()
	if err != nil {
		return nil, err
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

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

	clusterId := req.GetCluster().GetClusterId().GetValue()
	_, err := getCluster(clusterId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Get cluster [%s] failed", clusterId)
	}

	attributes := manager.BuildUpdateAttributes(req.Cluster, models.ClusterColumns...)
	logger.Debugf("ModifyCluster got attributes: [%+v]", attributes)
	delete(attributes, "cluster_id")
	_, err = pi.Global().Db.
		Update(models.ClusterTableName).
		SetMap(attributes).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] failed: %+v", clusterId, err)
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
			return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] node [%s] failed: %+v", clusterId, nodeId, err)
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
			return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] role [%s] failed: %+v", clusterId, role, err)
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
			return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] role [%s] common failed: %+v", clusterId, role, err)
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
			return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] name [%s] link failed: %+v", clusterId, name, err)
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
			return nil, status.Errorf(codes.Internal, "ModifyCluster [%s] role [%s] loadbalancer listener id [%s] failed: %+v",
				clusterId, role, listenerId, err)
		}
	}

	res := &pb.ModifyClusterResponse{
		ClusterId: utils.ToProtoString(clusterId),
	}
	return res, nil
}

func (p *Server) ModifyClusterNode(ctx context.Context, req *pb.ModifyClusterNodeRequest) (*pb.ModifyClusterNodeResponse, error) {
	s := sender.GetSenderFromContext(ctx)

	nodeId := req.GetClusterNode().GetNodeId().GetValue()
	_, err := getClusterNode(nodeId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Get cluster node [%s] failed", nodeId)
	}

	attributes := manager.BuildUpdateAttributes(req.ClusterNode, models.ClusterNodeColumns...)
	_, err = pi.Global().Db.
		Update(models.ClusterNodeTableName).
		SetMap(attributes).
		Where(db.Eq("node_id", nodeId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyClusterNode [%s] failed: %+v", nodeId, err)
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
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}

		directive, err := clusterWrapper.ToString()
		if err != nil {
			return nil, err
		}

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, err
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
	}

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionUpgradeCluster,
		directive,
		runtime.Provider,
		s.UserId,
	)

	jobId, err := jobclient.SendJob(newJob)
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
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
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
	}

	directive, err := clusterWrapper.ToString()
	if err != nil {
		return nil, err
	}

	runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, err
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
		return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
	}

	return &pb.UpdateClusterEnvResponse{
		ClusterId: utils.ToProtoString(clusterId),
		JobId:     utils.ToProtoString(jobId),
	}, nil
}

func (p *Server) DescribeClusters(ctx context.Context, req *pb.DescribeClustersRequest) (*pb.DescribeClustersResponse, error) {
	var clusters []*models.Cluster
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := pi.Global().Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterTableName))
	_, err := query.Load(&clusters)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeClusters failed: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeClusters failed: %+v", err)
	}

	var pbClusters []*pb.Cluster
	for _, cluster := range clusters {
		clusterId := cluster.ClusterId
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "DescribeClusters failed: %+v", err)
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
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := pi.Global().Db.
		Select(models.ClusterNodeColumns...).
		From(models.ClusterNodeTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.ClusterNodeTableName))
	_, err := query.Load(&clusterNodes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeClusterNodes: %+v", err)
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
			return nil, status.Errorf(codes.Internal, "DescribeClusterNodes [%s] common failed: %+v", nodeId, err)
		}

		err = pi.Global().Db.
			Select(models.ClusterRoleColumns...).
			From(models.ClusterRoleTableName).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			LoadOne(&clusterRole)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "DescribeClusterNodes [%s] role failed: %+v", nodeId, err)
		}

		pbClusterNodes = append(pbClusterNodes,
			models.ClusterNodeWrapperToPb(clusterNode, clusterCommon, clusterRole))
	}

	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeClusterNodes: %+v", err)
	}

	res := &pb.DescribeClusterNodesResponse{
		ClusterNodeSet: pbClusterNodes,
		TotalCount:     count,
	}
	return res, nil
}

func (p *Server) StopClusters(ctx context.Context, req *pb.StopClustersRequest) (*pb.StopClustersResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// TODO: check resource permission

	var jobIds []string
	for _, clusterId := range req.GetClusterId() {
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}

		directive, err := clusterWrapper.ToString()
		if err != nil {
			return nil, err
		}

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, err
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
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}

		directive, err := clusterWrapper.ToString()
		if err != nil {
			return nil, err
		}

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, err
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			logger.Errorf("Activate frontgate [%s] failed. ", clusterWrapper.Cluster.FrontgateId)
			return nil, err
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
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}

		directive, err := clusterWrapper.ToString()
		if err != nil {
			return nil, err
		}

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, err
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			logger.Errorf("Activate frontgate [%s] failed. ", clusterWrapper.Cluster.FrontgateId)
			return nil, err
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
		clusterWrapper, err := getClusterWrapper(clusterId)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Failed to get cluster [%s]", clusterId)
		}

		directive, err := clusterWrapper.ToString()
		if err != nil {
			return nil, err
		}

		runtime, err := runtimeclient.NewRuntime(clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, err
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
			return nil, status.Errorf(codes.Internal, "Send job [%s] failed: %+v", jobId, err)
		}
		jobIds = append(jobIds, jobId)
	}

	return &pb.CeaseClustersResponse{
		ClusterId: req.GetClusterId(),
		JobId:     jobIds,
	}, nil
}
