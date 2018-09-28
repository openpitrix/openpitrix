// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	pbempty "github.com/golang/protobuf/ptypes/empty"

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
	"openpitrix.io/openpitrix/pkg/util/senderutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func getClusterWrapper(ctx context.Context, clusterId string) (*models.ClusterWrapper, error) {
	clusterWrapper := new(models.ClusterWrapper)
	var cluster *models.Cluster
	var clusterCommons []*models.ClusterCommon
	var clusterNodes []*models.ClusterNode
	var clusterRoles []*models.ClusterRole
	var clusterLinks []*models.ClusterLink
	var clusterLoadbalancers []*models.ClusterLoadbalancer

	err := pi.Global().DB(ctx).
		Select(models.ClusterColumns...).
		From(constants.TableCluster).
		Where(db.Eq("cluster_id", clusterId)).
		LoadOne(&cluster)
	if err != nil {
		return nil, err
	}
	clusterWrapper.Cluster = cluster

	_, err = pi.Global().DB(ctx).
		Select(models.ClusterCommonColumns...).
		From(constants.TableClusterCommon).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterCommons)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterCommons = map[string]*models.ClusterCommon{}
	for _, clusterCommon := range clusterCommons {
		clusterWrapper.ClusterCommons[clusterCommon.Role] = clusterCommon
	}

	_, err = pi.Global().DB(ctx).
		Select(models.ClusterNodeColumns...).
		From(constants.TableClusterNode).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterNodes)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterNodesWithKeyPairs = map[string]*models.ClusterNodeWithKeyPairs{}
	for _, clusterNode := range clusterNodes {
		if stringutil.StringIn(clusterNode.Status, constants.DeletedStatuses) {
			continue
		}

		var nodeKeyPairs []*models.NodeKeyPair
		_, err = pi.Global().DB(ctx).
			Select(models.NodeKeyPairColumns...).
			From(constants.TableNodeKeyPair).
			Where(db.Eq("node_id", clusterNode.NodeId)).
			Load(&nodeKeyPairs)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		clusterNodeWithKeyPairs := &models.ClusterNodeWithKeyPairs{
			ClusterNode: clusterNode,
		}
		for _, nodeKeyPair := range nodeKeyPairs {
			clusterNodeWithKeyPairs.KeyPairId = append(clusterNodeWithKeyPairs.KeyPairId, nodeKeyPair.KeyPairId)
		}
		clusterWrapper.ClusterNodesWithKeyPairs[clusterNode.NodeId] = clusterNodeWithKeyPairs
	}

	_, err = pi.Global().DB(ctx).
		Select(models.ClusterRoleColumns...).
		From(constants.TableClusterRole).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterRoles)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterRoles = map[string]*models.ClusterRole{}
	for _, clusterRole := range clusterRoles {
		clusterWrapper.ClusterRoles[clusterRole.Role] = clusterRole
	}

	_, err = pi.Global().DB(ctx).
		Select(models.ClusterLinkColumns...).
		From(constants.TableClusterLink).
		Where(db.Eq("cluster_id", clusterId)).
		Load(&clusterLinks)
	if err != nil {
		return nil, err
	}

	clusterWrapper.ClusterLinks = map[string]*models.ClusterLink{}
	for _, clusterLink := range clusterLinks {
		clusterWrapper.ClusterLinks[clusterLink.Name] = clusterLink
	}

	_, err = pi.Global().DB(ctx).
		Select(models.ClusterLoadbalancerColumns...).
		From(constants.TableClusterLoadbalancer).
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

func getNodeKeyPairs(ctx context.Context, keyPairIds []string, nodeIds []string) ([]*models.NodeKeyPair, error) {
	var nodeKeyPairs []*models.NodeKeyPair
	for _, keyPairId := range keyPairIds {
		var singleNodeKeyPairs []*models.NodeKeyPair
		_, err := pi.Global().DB(ctx).
			Select(models.NodeKeyPairColumns...).
			From(constants.TableNodeKeyPair).
			Where(db.Eq("key_pair_id", keyPairId)).
			Where(db.Eq("node_id", nodeIds)).
			Load(&singleNodeKeyPairs)
		if err != nil {
			return nil, err
		}
		nodeKeyPairs = append(nodeKeyPairs, singleNodeKeyPairs...)
	}

	return nodeKeyPairs, nil
}

func updateTransitionStatus(ctx context.Context, cluster *models.Cluster) error {
	if cluster.TransitionStatus != "" {
		jobClient, err := jobclient.NewClient()
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
		jobs, err := jobClient.DescribeJobs(ctx, &pb.DescribeJobsRequest{
			ClusterId: pbutil.ToProtoString(cluster.ClusterId),
		})
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
		transitionStatus := ""
		for _, job := range jobs.JobSet {
			if !stringutil.StringIn(job.GetStatus().GetValue(), []string{constants.StatusSuccessful, constants.StatusFailed}) {
				transitionStatus = cluster.TransitionStatus
			}
		}
		cluster.TransitionStatus = transitionStatus
	}
	return nil
}

func (p *Server) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil || runtime.Owner != s.UserId {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, runtimeId)
	}

	if !plugins.IsVmbasedProviders(runtime.Provider) {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDescribeResourcesFailed)
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, runtime.Provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}

	return providerInterface.DescribeSubnets(ctx, req)
}

func (p *Server) DeleteNodeKeyPairs(ctx context.Context, req *pb.DeleteNodeKeyPairsRequest) (*pb.DeleteNodeKeyPairsResponse, error) {
	nodeKeyPairs := req.NodeKeyPair
	for _, nodeKeyPair := range nodeKeyPairs {
		_, err := CheckClusterNodePermission(ctx, nodeKeyPair.GetNodeId().GetValue())
		if err != nil {
			return nil, err
		}
		_, err = CheckKeyPairPermission(ctx, nodeKeyPair.GetKeyPairId().GetValue())
		if err != nil {
			return nil, err
		}
	}
	for _, nodeKeyPair := range nodeKeyPairs {
		_, err := pi.Global().DB(ctx).
			DeleteFrom(constants.TableNodeKeyPair).
			Where(db.Eq("key_pair_id", nodeKeyPair.GetKeyPairId().GetValue())).
			Where(db.Eq("node_id", nodeKeyPair.GetNodeId().GetValue())).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDetachKeyPairsFailed)
		}
	}
	return &pb.DeleteNodeKeyPairsResponse{}, nil
}

func (p *Server) AddNodeKeyPairs(ctx context.Context, req *pb.AddNodeKeyPairsRequest) (*pb.AddNodeKeyPairsResponse, error) {
	nodeKeyPairs := req.NodeKeyPair
	for _, nodeKeyPair := range nodeKeyPairs {
		_, err := CheckClusterNodePermission(ctx, nodeKeyPair.GetNodeId().GetValue())
		if err != nil {
			return nil, err
		}
		_, err = CheckKeyPairPermission(ctx, nodeKeyPair.GetKeyPairId().GetValue())
		if err != nil {
			return nil, err
		}
	}
	for _, nodeKeyPair := range nodeKeyPairs {
		nodeKeyPair := &models.NodeKeyPair{
			NodeId:    nodeKeyPair.GetNodeId().GetValue(),
			KeyPairId: nodeKeyPair.GetKeyPairId().GetValue(),
		}
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableNodeKeyPair).
			Record(nodeKeyPair).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorAttachKeyPairsFailed)
		}
	}

	return &pb.AddNodeKeyPairsResponse{}, nil
}

func (p *Server) CreateKeyPair(ctx context.Context, req *pb.CreateKeyPairRequest) (*pb.CreateKeyPairResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	owner := s.UserId
	name := req.GetName().GetValue()
	description := req.GetDescription().GetValue()
	pubKey := req.GetPubKey().GetValue()
	now := time.Now()
	newKeyPair := &models.KeyPair{
		KeyPairId:   models.NewKeyPairId(),
		Name:        name,
		Description: description,
		Owner:       owner,
		PubKey:      pubKey,
		CreateTime:  now,
		StatusTime:  now,
	}

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableKeyPair).
		Record(newKeyPair).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	res := &pb.CreateKeyPairResponse{
		KeyPairId: pbutil.ToProtoString(newKeyPair.KeyPairId),
	}
	return res, nil
}

func (p *Server) DescribeKeyPairs(ctx context.Context, req *pb.DescribeKeyPairsRequest) (*pb.DescribeKeyPairsResponse, error) {
	var keyPairs []*models.KeyPair
	var keyPairWithNodes []*models.KeyPairWithNodes
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	query := pi.Global().DB(ctx).
		Select(models.KeyPairColumns...).
		From(constants.TableKeyPair).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableKeyPair))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&keyPairs)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	for _, keyPair := range keyPairs {
		var nodeKeyPairs []*models.NodeKeyPair
		query = pi.Global().DB(ctx).
			Select(models.NodeKeyPairColumns...).
			From(constants.TableNodeKeyPair).
			Where(db.Eq("key_pair_id", keyPair.KeyPairId))
		_, err := query.Load(&nodeKeyPairs)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		keyPairWithNodesItem := &models.KeyPairWithNodes{
			KeyPair: keyPair,
		}

		for _, nodeKeyPair := range nodeKeyPairs {
			keyPairWithNodesItem.NodeId = append(keyPairWithNodesItem.NodeId, nodeKeyPair.NodeId)
		}
		keyPairWithNodes = append(keyPairWithNodes, keyPairWithNodesItem)
	}

	keyPairSet := models.KeyPairNodesToPbs(keyPairWithNodes)

	res := &pb.DescribeKeyPairsResponse{
		KeyPairSet: keyPairSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) DeleteKeyPairs(ctx context.Context, req *pb.DeleteKeyPairsRequest) (*pb.DeleteKeyPairsResponse, error) {
	keyPairIds := req.KeyPairId
	keyPairs, err := CheckKeyPairsPermission(ctx, keyPairIds)
	if err != nil {
		return nil, err
	}

	var attachedKeyPairs []*models.KeyPair
	_, err = pi.Global().DB(ctx).
		Select(models.NodeKeyPairColumns...).
		From(constants.TableNodeKeyPair).
		Where(db.Eq("key_pair_id", keyPairIds)).
		Load(&attachedKeyPairs)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}
	if len(attachedKeyPairs) > 0 {
		var attachedKeyPairIds []string
		for _, attachedKeyPair := range attachedKeyPairs {
			attachedKeyPairIds = append(attachedKeyPairIds, attachedKeyPair.KeyPairId)
		}
		err = fmt.Errorf("key pairs [%s] are still attached", strings.Join(attachedKeyPairIds, ","))
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDeleteResourceFailed, strings.Join(attachedKeyPairIds, ","))
	}

	var deleteKeyPairIds []string
	for _, keyPair := range keyPairs {
		deleteKeyPairIds = append(deleteKeyPairIds, keyPair.KeyPairId)
	}

	_, err = pi.Global().DB(ctx).
		DeleteFrom(constants.TableKeyPair).
		Where(db.Eq("key_pair_id", deleteKeyPairIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed)
	}

	res := &pb.DeleteKeyPairsResponse{
		KeyPairId: deleteKeyPairIds,
	}
	return res, nil
}

func (p *Server) AttachKeyPairs(ctx context.Context, req *pb.AttachKeyPairsRequest) (*pb.AttachKeyPairsResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	nodeIds := req.GetNodeId()
	clusterNodes, err := CheckClusterNodesPermission(ctx, nodeIds)
	if err != nil {
		return nil, err
	}
	err = checkNodesPermissionAndTransition(ctx, clusterNodes, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
	}

	keyPairIds := req.GetKeyPairId()
	keyPairs, err := CheckKeyPairsPermission(ctx, keyPairIds)
	if err != nil {
		return nil, err
	}

	existNodeKeyPairs, err := getNodeKeyPairs(ctx, keyPairIds, nodeIds)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
	}
	if len(existNodeKeyPairs) > 0 {
		err = fmt.Errorf("keypair [%s] has already been attached to [%s]", existNodeKeyPairs[0].KeyPairId, existNodeKeyPairs[0].NodeId)
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
	}

	clusterNodeIds := make(map[string][]string)
	clusterNodeMap := make(map[string]*models.ClusterNode)
	for _, clusterNode := range clusterNodes {
		_, isExist := clusterNodeIds[clusterNode.ClusterId]
		if isExist {
			clusterNodeIds[clusterNode.ClusterId] = append(clusterNodeIds[clusterNode.ClusterId], clusterNode.NodeId)
		} else {
			clusterNodeIds[clusterNode.ClusterId] = []string{clusterNode.NodeId}
		}
		clusterNodeMap[clusterNode.NodeId] = clusterNode
	}

	keyPairMap := make(map[string]*models.KeyPair)
	for _, keyPair := range keyPairs {
		keyPairMap[keyPair.KeyPairId] = keyPair
	}

	var jobIds []string
	for clusterId, nodeIds := range clusterNodeIds {
		cluster, err := CheckClusterPermission(ctx, clusterId)
		if err != nil {
			return nil, err
		}
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive, constants.StatusPending})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
		}
		runtime, err := runtimeclient.NewRuntime(ctx, cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
		}

		if !plugins.IsVmbasedProviders(runtime.Provider) {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAttachKeyPairsFailed)
		}

		var nodeKeyPairDetails models.NodeKeyPairDetails
		for _, nodeId := range nodeIds {
			for _, keyPairId := range keyPairIds {
				nodeKeyPairDetail := models.NodeKeyPairDetail{
					NodeKeyPair: &models.NodeKeyPair{
						KeyPairId: keyPairId,
						NodeId:    nodeId,
					},
					ClusterNode: clusterNodeMap[nodeId],
					KeyPair:     keyPairMap[keyPairId],
				}
				nodeKeyPairDetails = append(nodeKeyPairDetails, nodeKeyPairDetail)
			}
		}

		directive := jsonutil.ToString(nodeKeyPairDetails)

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.VersionId,
			constants.ActionAttachKeyPairs,
			directive,
			runtime.Provider,
			s.UserId,
			cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorAttachKeyPairsFailed)
		}
		jobIds = append(jobIds, jobId)
	}

	res := &pb.AttachKeyPairsResponse{
		JobId: jobIds,
	}
	return res, nil
}

func (p *Server) DetachKeyPairs(ctx context.Context, req *pb.DetachKeyPairsRequest) (*pb.DetachKeyPairsResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	nodeIds := req.GetNodeId()
	clusterNodes, err := CheckClusterNodesPermission(ctx, nodeIds)
	if err != nil {
		return nil, err
	}
	err = checkNodesPermissionAndTransition(ctx, clusterNodes, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
	}

	keyPairIds := req.GetKeyPairId()
	keyPairs, err := CheckKeyPairsPermission(ctx, keyPairIds)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
	}

	existNodeKeyPairs, err := getNodeKeyPairs(ctx, keyPairIds, nodeIds)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
	}
	if len(existNodeKeyPairs) < len(keyPairIds)*len(nodeIds) {
		err = fmt.Errorf("keypair has not been attached to node")
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
	}

	clusterNodeIds := make(map[string][]string)
	clusterNodeMap := make(map[string]*models.ClusterNode)
	for _, clusterNode := range clusterNodes {
		_, isExist := clusterNodeIds[clusterNode.ClusterId]
		if isExist {
			clusterNodeIds[clusterNode.ClusterId] = append(clusterNodeIds[clusterNode.ClusterId], clusterNode.NodeId)
		} else {
			clusterNodeIds[clusterNode.ClusterId] = []string{clusterNode.NodeId}
		}
		clusterNodeMap[clusterNode.NodeId] = clusterNode
	}

	keyPairMap := make(map[string]*models.KeyPair)
	for _, keyPair := range keyPairs {
		keyPairMap[keyPair.KeyPairId] = keyPair
	}

	var jobIds []string
	for clusterId, nodeIds := range clusterNodeIds {
		cluster, err := CheckClusterPermission(ctx, clusterId)
		if err != nil {
			return nil, err
		}
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive, constants.StatusPending})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
		}
		runtime, err := runtimeclient.NewRuntime(ctx, cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
		}

		if !plugins.IsVmbasedProviders(runtime.Provider) {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDetachKeyPairsFailed)
		}

		var nodeKeyPairDetails models.NodeKeyPairDetails
		for _, nodeId := range nodeIds {
			for _, keyPairId := range keyPairIds {
				nodeKeyPairDetail := models.NodeKeyPairDetail{
					NodeKeyPair: &models.NodeKeyPair{
						KeyPairId: keyPairId,
						NodeId:    nodeId,
					},
					ClusterNode: clusterNodeMap[nodeId],
					KeyPair:     keyPairMap[keyPairId],
				}
				nodeKeyPairDetails = append(nodeKeyPairDetails, nodeKeyPairDetail)
			}
		}

		directive := jsonutil.ToString(nodeKeyPairDetails)

		newJob := models.NewJob(
			constants.PlaceHolder,
			clusterId,
			cluster.AppId,
			cluster.VersionId,
			constants.ActionDetachKeyPairs,
			directive,
			runtime.Provider,
			s.UserId,
			cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDetachKeyPairsFailed)
		}
		jobIds = append(jobIds, jobId)
	}

	res := &pb.DetachKeyPairsResponse{
		JobId: jobIds,
	}
	return res, nil
}

func (p *Server) CreateCluster(ctx context.Context, req *pb.CreateClusterRequest) (*pb.CreateClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	appId := req.GetAppId().GetValue()
	versionId := req.GetVersionId().GetValue()
	conf := req.GetConf().GetValue()
	clusterId := models.NewClusterId()
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil || runtime.Owner != s.UserId {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAccessDenied, runtimeId)
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, runtime.Provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}
	clusterWrapper := new(models.ClusterWrapper)
	err = providerInterface.ParseClusterConf(ctx, versionId, runtimeId, conf, clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Parse cluster conf with versionId [%s] runtime [%s] failed: %+v", versionId, runtimeId, err)
		if gerr.IsGRPCError(err) {
			return nil, err
		}
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	if clusterWrapper.Cluster.Zone == "" {
		clusterWrapper.Cluster.Zone = runtime.Zone
	}
	clusterWrapper.Cluster.RuntimeId = runtimeId
	clusterWrapper.Cluster.Owner = s.UserId
	clusterWrapper.Cluster.ClusterId = clusterId
	clusterWrapper.Cluster.ClusterType = constants.NormalClusterType

	if plugins.IsVmbasedProviders(runtime.Provider) {
		err = CheckVmBasedProvider(ctx, runtime, providerInterface, clusterWrapper)
		if err != nil {
			return nil, err
		}
	} else {
		err = providerInterface.CheckResource(ctx, clusterWrapper)
		if err != nil {
			return nil, err
		}
	}

	err = RegisterClusterWrapper(ctx, clusterWrapper)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
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
		runtimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	res := &pb.CreateClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}
	return res, nil
}

func (p *Server) ModifyCluster(ctx context.Context, req *pb.ModifyClusterRequest) (*pb.ModifyClusterResponse, error) {
	clusterId := req.GetCluster().GetClusterId().GetValue()
	_, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req.Cluster, models.ClusterColumns...)
	logger.Debug(ctx, "ModifyCluster got attributes: [%+v]", attributes)
	delete(attributes, "cluster_id")
	if len(attributes) != 0 {
		_, err = pi.Global().DB(ctx).
			Update(constants.TableCluster).
			SetMap(attributes).
			Where(db.Eq("cluster_id", clusterId)).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
		}
	}

	for _, clusterNode := range req.ClusterNodeSet {
		nodeId := clusterNode.GetNodeId().GetValue()
		nodeAttributes := manager.BuildUpdateAttributes(clusterNode, models.ClusterNodeColumns...)
		delete(nodeAttributes, "cluster_id")
		delete(nodeAttributes, "node_id")
		if len(nodeAttributes) != 0 {
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterNode).
				SetMap(nodeAttributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("node_id", nodeId)).
				Exec()
			if err != nil {
				logger.Error(ctx, "ModifyCluster [%s] node [%s] failed. ", clusterId, nodeId)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	for _, clusterRole := range req.ClusterRoleSet {
		role := clusterRole.GetRole().GetValue()
		roleAttributes := manager.BuildUpdateAttributes(clusterRole, models.ClusterRoleColumns...)
		delete(roleAttributes, "cluster_id")
		delete(roleAttributes, "role")
		if len(roleAttributes) != 0 {
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterRole).
				SetMap(roleAttributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("role", role)).
				Exec()
			if err != nil {
				logger.Error(ctx, "ModifyCluster [%s] role [%s] failed. ", clusterId, role)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	for _, clusterCommon := range req.ClusterCommonSet {
		role := clusterCommon.GetRole().GetValue()
		commonAttributes := manager.BuildUpdateAttributes(clusterCommon, models.ClusterCommonColumns...)
		delete(commonAttributes, "cluster_id")
		delete(commonAttributes, "role")
		if len(commonAttributes) != 0 {
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterCommon).
				SetMap(commonAttributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("role", role)).
				Exec()
			if err != nil {
				logger.Error(ctx, "ModifyCluster [%s] role [%s] common failed. ", clusterId, role)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	for _, clusterLink := range req.ClusterLinkSet {
		name := clusterLink.GetName().GetValue()
		linkAttributes := manager.BuildUpdateAttributes(clusterLink, models.ClusterLinkColumns...)
		delete(linkAttributes, "cluster_id")
		delete(linkAttributes, "name")
		if len(linkAttributes) != 0 {
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterLink).
				SetMap(linkAttributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("name", name)).
				Exec()
			if err != nil {
				logger.Error(ctx, "ModifyCluster [%s] name [%s] link failed. ", clusterId, name)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	for _, clusterLoadbalancer := range req.ClusterLoadbalancerSet {
		role := clusterLoadbalancer.GetRole().GetValue()
		listenerId := clusterLoadbalancer.GetLoadbalancerListenerId().GetValue()
		loadbalancerAttributes := manager.BuildUpdateAttributes(clusterLoadbalancer, models.ClusterLoadbalancerColumns...)
		delete(loadbalancerAttributes, "cluster_id")
		delete(loadbalancerAttributes, "role")
		delete(loadbalancerAttributes, "loadbalancer_listener_id")
		if len(loadbalancerAttributes) != 0 {
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterLoadbalancer).
				SetMap(loadbalancerAttributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("role", role)).
				Where(db.Eq("loadbalancer_listener_id", listenerId)).
				Exec()
			if err != nil {
				logger.Error(ctx, "ModifyCluster [%s] role [%s] loadbalancer listener id [%s] failed. ",
					clusterId, role, listenerId)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	res := &pb.ModifyClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
	}
	return res, nil
}

func (p *Server) ModifyClusterNode(ctx context.Context, req *pb.ModifyClusterNodeRequest) (*pb.ModifyClusterNodeResponse, error) {
	nodeId := req.GetClusterNode().GetNodeId().GetValue()
	_, err := CheckClusterNodePermission(ctx, nodeId)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req.ClusterNode, models.ClusterNodeColumns...)
	_, err = pi.Global().DB(ctx).
		Update(constants.TableClusterNode).
		SetMap(attributes).
		Where(db.Eq("node_id", nodeId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, nodeId)
	}

	res := &pb.ModifyClusterNodeResponse{
		NodeId: pbutil.ToProtoString(nodeId),
	}
	return res, nil
}

func (p *Server) ModifyClusterAttributes(ctx context.Context, req *pb.ModifyClusterAttributesRequest) (*pb.ModifyClusterAttributesResponse, error) {
	clusterId := req.GetClusterId().GetValue()
	_, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req, models.ClusterColumns...)
	_, err = pi.Global().DB(ctx).
		Update(constants.TableCluster).
		SetMap(attributes).
		Where(db.Eq("cluster_id", clusterId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
	}

	res := &pb.ModifyClusterAttributesResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
	}
	return res, nil
}

func (p *Server) ModifyClusterNodeAttributes(ctx context.Context, req *pb.ModifyClusterNodeAttributesRequest) (*pb.ModifyClusterNodeAttributesResponse, error) {
	nodeId := req.GetNodeId().GetValue()
	_, err := CheckClusterNodePermission(ctx, nodeId)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req, models.ClusterNodeColumns...)
	_, err = pi.Global().DB(ctx).
		Update(constants.TableClusterNode).
		SetMap(attributes).
		Where(db.Eq("node_id", nodeId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, nodeId)
	}

	res := &pb.ModifyClusterNodeAttributesResponse{
		NodeId: pbutil.ToProtoString(nodeId),
	}
	return res, nil
}

func (p *Server) AddTableClusterNodes(ctx context.Context, req *pb.AddTableClusterNodesRequest) (*pbempty.Empty, error) {
	for _, clusterNode := range req.ClusterNodeSet {
		node := models.PbToClusterNode(clusterNode)
		err := RegisterClusterNode(ctx, node.ClusterNode)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pbempty.Empty{}, nil
}

func (p *Server) DeleteTableClusterNodes(ctx context.Context, req *pb.DeleteTableClusterNodesRequest) (*pbempty.Empty, error) {
	for _, nodeId := range req.NodeId {
		_, err := pi.Global().DB(ctx).
			DeleteFrom(constants.TableClusterNode).
			Where(db.Eq("node_id", nodeId)).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pbempty.Empty{}, nil
}

func (p *Server) DeleteClusters(ctx context.Context, req *pb.DeleteClustersRequest) (*pb.DeleteClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	clusterIds := req.GetClusterId()
	clusters, err := CheckClustersPermission(ctx, clusterIds)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, cluster := range clusters {
		err = updateTransitionStatus(ctx, cluster)
		if err != nil {
			return nil, err
		}
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive, constants.StatusStopped, constants.StatusPending})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDeleteResourceFailed, cluster.ClusterId)
		}

		clusterWrapper, err := getClusterWrapper(ctx, cluster.ClusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, cluster.ClusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}
		newJob := models.NewJob(
			constants.PlaceHolder,
			cluster.ClusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionDeleteClusters,
			directive,
			runtime.Provider,
			s.UserId,
			clusterWrapper.Cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, cluster.ClusterId)
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
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	versionId := req.GetVersionId().GetValue()
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusStopped})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorUpgradeResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
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
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpgradeResourceFailed, clusterId)
	}

	return &pb.UpgradeClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) RollbackCluster(ctx context.Context, req *pb.RollbackClusterRequest) (*pb.RollbackClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRollbackResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	directive := jsonutil.ToString(clusterWrapper)

	runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
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
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRollbackResourceFailed, clusterId)
	}

	return &pb.RollbackClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) ResizeCluster(ctx context.Context, req *pb.ResizeClusterRequest) (*pb.ResizeClusterResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResizeResourceFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	if clusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType || !plugins.IsVmbasedProviders(runtime.Provider) {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}

	var roleResizeResources models.RoleResizeResources
	for _, pbRoleResource := range req.RoleResource {
		roleResource := models.PbToRoleResource(pbRoleResource)
		clusterRole, isExist := clusterWrapper.ClusterRoles[roleResource.Role]
		if !isExist {
			err = fmt.Errorf("role [%s] not found", roleResource.Role)
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceRoleNotFound, clusterId, roleResource.Role)
		}

		if isSame, roleResizeResource := roleResource.IsSame(clusterRole); !isSame && roleResizeResource != nil {
			roleResizeResources = append(roleResizeResources, roleResizeResource)
			attributes := map[string]interface{}{
				"cpu":           clusterRole.Cpu,
				"memory":        clusterRole.Memory,
				"gpu":           clusterRole.Gpu,
				"instance_size": clusterRole.InstanceSize,
				"storage_size":  clusterRole.StorageSize,
			}
			_, err = pi.Global().DB(ctx).
				Update(constants.TableClusterRole).
				SetMap(attributes).
				Where(db.Eq("cluster_id", clusterId)).
				Where(db.Eq("role", roleResource.Role)).
				Exec()
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, clusterId)
			}
		}
	}

	if len(roleResizeResources) == 0 {
		err = fmt.Errorf("cluster [%s] is already the resource type", clusterId)
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorResizeResourceFailed, clusterId)
	}

	directive := jsonutil.ToString(roleResizeResources)

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionResizeCluster,
		directive,
		runtime.Provider,
		s.UserId,
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorResizeResourceFailed, clusterId)
	}

	return &pb.ResizeClusterResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) AddClusterNodes(ctx context.Context, req *pb.AddClusterNodesRequest) (*pb.AddClusterNodesResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}
	role := req.GetRole().GetValue()
	count := int(req.GetNodeCount().GetValue())
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	owner := clusterWrapper.Cluster.Owner

	runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	if clusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType || !plugins.IsVmbasedProviders(runtime.Provider) {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}

	var roleNodes []*models.ClusterNodeWithKeyPairs
	for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
		if clusterNode.Role == role {
			roleNodes = append(roleNodes, clusterNode)
		}
	}

	for i := 1; i <= count; i++ {
		clusterWrapper.ClusterNodesWithKeyPairs[string(i)] = &models.ClusterNodeWithKeyPairs{
			ClusterNode: &models.ClusterNode{
				Status: constants.StatusPending,
				Role:   role,
			},
		}
	}

	conf := ""
	if len(req.AdvancedParam) > 0 {
		conf = req.AdvancedParam[0]
	}
	providerInterface, err := plugins.GetProviderPlugin(ctx, runtime.Provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}
	err = providerInterface.ParseClusterConf(ctx, clusterWrapper.Cluster.VersionId, runtime.RuntimeId, conf, clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Parse cluster conf with versionId [%s] runtime [%s] failed: %+v", clusterWrapper.Cluster.VersionId, runtime.RuntimeId, err)
		if gerr.IsGRPCError(err) {
			return nil, err
		}
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	// register new role
	if len(roleNodes) == 0 {
		if len(req.AdvancedParam) == 0 {
			err = fmt.Errorf("conf parameter is needed when role [%s] node does not exist", role)
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorAddResourceNodeFailed, clusterId)
		}

		clusterWrapper.ClusterRoles[role].ClusterId = clusterId
		err = RegisterClusterRole(ctx, clusterWrapper.ClusterRoles[role])
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorAddResourceNodeFailed)
		}
	}

	// register new nodes
	for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
		if clusterNode.Status == constants.StatusPending {
			clusterNode.ClusterNode.Owner = owner
			err = RegisterClusterNode(ctx, clusterNode.ClusterNode)
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorAddResourceNodeFailed)
			}
		}
	}

	// reload clusterWrapper from db
	clusterWrapper, err = getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}
	directive := jsonutil.ToString(clusterWrapper)
	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionAddClusterNodes,
		directive,
		runtime.Provider,
		s.UserId,
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}

	return &pb.AddClusterNodesResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) DeleteClusterNodes(ctx context.Context, req *pb.DeleteClusterNodesRequest) (*pb.DeleteClusterNodesResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDeleteResourceNodeFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}

	runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	if clusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType || !plugins.IsVmbasedProviders(runtime.Provider) {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorAddResourceNodeFailed, clusterId)
	}

	// TODO: check
	nodeIds := req.GetNodeId()
	for _, nodeId := range nodeIds {
		clusterNode, isExist := clusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if !isExist || clusterNode.Status != constants.StatusActive {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
		}
		clusterNode.Status = constants.StatusDeleting
	}

	directive := jsonutil.ToString(clusterWrapper)

	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionDeleteClusterNodes,
		directive,
		runtime.Provider,
		s.UserId,
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceNodeFailed, clusterId)
	}

	return &pb.DeleteClusterNodesResponse{
		ClusterId: pbutil.ToProtoString(clusterId),
		JobId:     pbutil.ToProtoString(jobId),
	}, nil
}

func (p *Server) UpdateClusterEnv(ctx context.Context, req *pb.UpdateClusterEnvRequest) (*pb.UpdateClusterEnvResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterId := req.GetClusterId().GetValue()
	cluster, err := CheckClusterPermission(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	conf := req.GetEnv().GetValue()
	err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorUpdateResourceEnvFailed, clusterId)
	}
	clusterWrapper, err := getClusterWrapper(ctx, clusterId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
	}
	versionId := clusterWrapper.Cluster.VersionId
	runtimeId := clusterWrapper.Cluster.RuntimeId
	clusterName := clusterWrapper.Cluster.Name

	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, runtime.Provider)
	if err != nil {
		logger.Error(ctx, "No such provider [%s]. ", runtime.Provider)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
	}
	err = providerInterface.ParseClusterConf(ctx, versionId, runtimeId, conf, clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Parse cluster conf with versionId [%s] runtime [%s] conf [%s] failed: %+v",
			versionId, runtime, conf, err)
		if gerr.IsGRPCError(err) {
			return nil, err
		}
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	clusterWrapper.Cluster.ClusterId = clusterId
	clusterWrapper.Cluster.Name = clusterName
	directive := jsonutil.ToString(clusterWrapper)
	newJob := models.NewJob(
		constants.PlaceHolder,
		clusterId,
		clusterWrapper.Cluster.AppId,
		clusterWrapper.Cluster.VersionId,
		constants.ActionUpdateClusterEnv,
		directive,
		runtime.Provider,
		s.UserId,
		clusterWrapper.Cluster.RuntimeId,
	)

	jobId, err := jobclient.SendJob(ctx, newJob)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceEnvFailed, clusterId)
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

	query := pi.Global().DB(ctx).
		Select(models.ClusterColumns...).
		From(constants.TableCluster).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableCluster))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&clusters)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var pbClusters []*pb.Cluster
	for _, cluster := range clusters {
		clusterId := cluster.ClusterId
		clusterWrapper, err := getClusterWrapper(ctx, clusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterId)
		}

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		if runtime.Provider == constants.ProviderKubernetes {
			providerInterface, err := plugins.GetProviderPlugin(ctx, runtime.Provider)
			if err != nil {
				logger.Error(ctx, "No such provider [%s]. ", runtime.Provider)
				return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtime.Provider)
			}

			err = providerInterface.DescribeClusterDetails(ctx, clusterWrapper)
			if err != nil {
				logger.Warn(ctx, "Describe cluster details failed: %+v", err)
			}
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

	query := pi.Global().DB(ctx).
		Select(models.ClusterNodeColumns...).
		From(constants.TableClusterNode).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableClusterNode))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&clusterNodes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var pbClusterNodes []*pb.ClusterNode
	for _, clusterNode := range clusterNodes {
		var clusterCommon *models.ClusterCommon
		var clusterRole *models.ClusterRole
		nodeId := clusterNode.NodeId
		role := clusterNode.Role
		clusterId := clusterNode.ClusterId
		err = pi.Global().DB(ctx).
			Select(models.ClusterCommonColumns...).
			From(constants.TableClusterCommon).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			LoadOne(&clusterCommon)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
		}

		err = pi.Global().DB(ctx).
			Select(models.ClusterRoleColumns...).
			From(constants.TableClusterRole).
			Where(db.Eq("cluster_id", clusterId)).
			Where(db.Eq("role", role)).
			LoadOne(&clusterRole)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, nodeId)
		}

		var nodeKeyPairs []*models.NodeKeyPair
		_, err = pi.Global().DB(ctx).
			Select(models.NodeKeyPairColumns...).
			From(constants.TableNodeKeyPair).
			Where(db.Eq("node_id", clusterNode.NodeId)).
			Load(&nodeKeyPairs)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		clusterNodeWithKeyPairs := &models.ClusterNodeWithKeyPairs{
			ClusterNode: clusterNode,
		}

		for _, nodeKeyPair := range nodeKeyPairs {
			clusterNodeWithKeyPairs.KeyPairId = append(clusterNodeWithKeyPairs.KeyPairId, nodeKeyPair.KeyPairId)
		}

		pbClusterNodes = append(pbClusterNodes,
			models.ClusterNodeWrapperToPb(clusterNodeWithKeyPairs, clusterCommon, clusterRole))
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeClusterNodesResponse{
		ClusterNodeSet: pbClusterNodes,
		TotalCount:     count,
	}
	return res, nil
}

func (p *Server) StopClusters(ctx context.Context, req *pb.StopClustersRequest) (*pb.StopClustersResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	clusterIds := req.GetClusterId()
	clusters, err := CheckClustersPermission(ctx, clusterIds)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, cluster := range clusters {
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusActive})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorStopResourceFailed, cluster.ClusterId)
		}
		clusterWrapper, err := getClusterWrapper(ctx, cluster.ClusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, cluster.ClusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		if !plugins.IsVmbasedProviders(runtime.Provider) {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorStopResourceFailed, cluster.ClusterId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			cluster.ClusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionStopClusters,
			directive,
			runtime.Provider,
			s.UserId,
			clusterWrapper.Cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorStopResourceFailed, cluster.ClusterId)
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

	clusterIds := req.GetClusterId()
	clusters, err := CheckClustersPermission(ctx, clusterIds)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, cluster := range clusters {
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusStopped})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorStartResourceFailed, cluster.ClusterId)
		}
		clusterWrapper, err := getClusterWrapper(ctx, cluster.ClusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, cluster.ClusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		if !plugins.IsVmbasedProviders(runtime.Provider) {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorStartResourceFailed, cluster.ClusterId)
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(ctx, clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorStartResourceFailed, cluster.ClusterId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			cluster.ClusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionStartClusters,
			directive,
			runtime.Provider,
			s.UserId,
			clusterWrapper.Cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorStartResourceFailed, cluster.ClusterId)
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

	clusterIds := req.GetClusterId()
	clusters, err := CheckClustersPermission(ctx, clusterIds)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, cluster := range clusters {
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusDeleted})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorRecoverResourceFailed, cluster.ClusterId)
		}
		clusterWrapper, err := getClusterWrapper(ctx, cluster.ClusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, cluster.ClusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		fg := &Frontgate{
			Runtime: runtime,
		}
		err = fg.ActivateFrontgate(ctx, clusterWrapper.Cluster.FrontgateId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRecoverResourceFailed, cluster.ClusterId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			cluster.ClusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionRecoverClusters,
			directive,
			runtime.Provider,
			s.UserId,
			clusterWrapper.Cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRecoverResourceFailed, cluster.ClusterId)
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

	clusterIds := req.GetClusterId()
	clusters, err := CheckClustersPermission(ctx, clusterIds)
	if err != nil {
		return nil, err
	}

	var jobIds []string
	for _, cluster := range clusters {
		err = checkPermissionAndTransition(ctx, cluster, []string{constants.StatusDeleted})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCeaseResourceFailed, cluster.ClusterId)
		}
		clusterWrapper, err := getClusterWrapper(ctx, cluster.ClusterId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, cluster.ClusterId)
		}

		directive := jsonutil.ToString(clusterWrapper)

		runtime, err := runtimeclient.NewRuntime(ctx, clusterWrapper.Cluster.RuntimeId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.RuntimeId)
		}

		newJob := models.NewJob(
			constants.PlaceHolder,
			cluster.ClusterId,
			clusterWrapper.Cluster.AppId,
			clusterWrapper.Cluster.VersionId,
			constants.ActionCeaseClusters,
			directive,
			runtime.Provider,
			s.UserId,
			clusterWrapper.Cluster.RuntimeId,
		)

		jobId, err := jobclient.SendJob(ctx, newJob)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCeaseResourceFailed, cluster.ClusterId)
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
type appStatistic struct {
	AppId string `db:"app_id"`
	Count uint32 `db:"COUNT(cluster_id)"`
}

func (p *Server) GetClusterStatistics(ctx context.Context, req *pb.GetClusterStatisticsRequest) (*pb.GetClusterStatisticsResponse, error) {
	res := &pb.GetClusterStatisticsResponse{
		LastTwoWeekCreated: make(map[string]uint32),
		TopTenRuntimes:     make(map[string]uint32),
		TopTenApps:         make(map[string]uint32),
	}
	clusterCount, err := pi.Global().DB(ctx).
		Select(constants.ColumnClusterId).
		From(constants.TableCluster).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Count()
	if err != nil {
		logger.Error(ctx, "Failed to get cluster count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.ClusterCount = clusterCount

	err = pi.Global().DB(ctx).
		Select("COUNT(DISTINCT runtime_id)").
		From(constants.TableCluster).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		LoadOne(&res.RuntimeCount)
	if err != nil {
		logger.Error(ctx, "Failed to get runtime count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var cs []*clusterStatistic
	_, err = pi.Global().DB(ctx).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(cluster_id)").
		From(constants.TableCluster).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(constants.ColumnCreateTime, time2week)).
		Limit(14).Load(&cs)

	if err != nil {
		logger.Error(ctx, "Failed to get cluster statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range cs {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var rs []*runtimeStatistic
	_, err = pi.Global().DB(ctx).
		Select("runtime_id", "COUNT(cluster_id)").
		From(constants.TableCluster).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		GroupBy(constants.ColumnRuntimeId).
		OrderDir("COUNT(cluster_id)", false).
		Limit(10).Load(&rs)

	if err != nil {
		logger.Error(ctx, "Failed to get runtime statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.TopTenRuntimes[a.RuntimeId] = a.Count
	}

	var as []*appStatistic
	_, err = pi.Global().DB(ctx).
		Select("app_id", "COUNT(cluster_id)").
		From(constants.TableCluster).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(db.Neq(constants.ColumnAppId, []string{"", constants.FrontgateAppId})).
		GroupBy(constants.ColumnAppId).
		OrderDir("COUNT(cluster_id)", false).
		Limit(10).Load(&as)

	if err != nil {
		logger.Error(ctx, "Failed to get app statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range as {
		res.TopTenApps[a.AppId] = a.Count
	}

	return res, nil
}
