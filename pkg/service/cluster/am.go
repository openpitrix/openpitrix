// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"

	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
	"openpitrix.io/openpitrix/pkg/util/tlsutil"
)

func checkPermissionAndTransition(ctx context.Context, cluster *models.Cluster, status []string) error {
	if cluster.TransitionStatus != "" {
		logger.Error(ctx, "Cluster [%s] is [%s], please try later", cluster.ClusterId, cluster.TransitionStatus)
		return fmt.Errorf("cluster [%s] is [%s], please try later", cluster.ClusterId, cluster.TransitionStatus)
	}
	if status != nil && !reflectutil.In(cluster.Status, status) {
		logger.Error(ctx, "Cluster [%s] status is [%s] not in %s", cluster.ClusterId, cluster.Status, status)
		return fmt.Errorf("cluster [%s] status is [%s] not in %s", cluster.ClusterId, cluster.Status, status)
	}
	return nil
}

func checkNodesPermissionAndTransition(ctx context.Context, clusterNodes []*models.ClusterNode, status []string) error {
	for _, clusterNode := range clusterNodes {
		if clusterNode.TransitionStatus != "" {
			logger.Error(ctx, "Cluster node [%s] is [%s], please try later", clusterNode.NodeId, clusterNode.TransitionStatus)
			return fmt.Errorf("cluster [%s] is [%s], please try later", clusterNode.NodeId, clusterNode.TransitionStatus)
		}
		if status != nil && !reflectutil.In(clusterNode.Status, status) {
			logger.Error(ctx, "Cluster [%s] status is [%s] not in %s", clusterNode.NodeId, clusterNode.Status, status)
			return fmt.Errorf("cluster [%s] status is [%s] not in %s", clusterNode.NodeId, clusterNode.Status, status)
		}
	}
	return nil
}

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateKeyPairRequest:
		return manager.NewChecker(ctx, r).
			Required("key").
			Exec()
	case *pb.DeleteKeyPairsRequest:
		return manager.NewChecker(ctx, r).
			Required("key_pair_id").
			Exec()
	case *pb.DescribeSubnetsRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_id").
			Exec()
	case *pb.CreateClusterRequest:
		return manager.NewChecker(ctx, r).
			Required("app_id", "version_id", "runtime_id", "conf").
			Exec()
	case *pb.DeleteClustersRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.UpgradeClusterRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "version_id").
			Exec()
	case *pb.RollbackClusterRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.ResizeClusterRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.AddClusterNodesRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "node_count").
			Exec()
	case *pb.DeleteClusterNodesRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "node_id").
			Exec()
	case *pb.UpdateClusterEnvRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "env").
			Exec()
	case *pb.StopClustersRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.StartClustersRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.RecoverClustersRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.CeaseClustersRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id").
			Exec()
	case *pb.GetClusterStatisticsRequest:
		return manager.NewChecker(ctx, r).
			Exec()
	}
	return nil
}

func CheckVmBasedProvider(ctx context.Context, runtime *models.RuntimeDetails, providerClient pb.RuntimeProviderManagerClient,
	clusterWrapper *models.ClusterWrapper) error {

	// check pilot service
	pilotClient, err := pilotclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Connect to pilot service failed")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	tlsConfig, err := pilotClient.GetPilotClientTLSConfig(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Error(ctx, "Get pilot client tls config failed")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	tlsPilotClientConfig, err := tlsutil.NewClientTLSConfigFromString(
		tlsConfig.ClientCrtData,
		tlsConfig.ClientKeyData,
		tlsConfig.CaCrtData,
		tlsConfig.PilotServerName,
	)

	pilotTLSClient, err := pilotclient.NewTLSClient(tlsPilotClientConfig)
	if err != nil {
		logger.Error(ctx, "Connect to pilot service tls port failed")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	_, err = pilotTLSClient.PingPilot(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Error(ctx, "Pilot service is not running or the endpoint is wrong")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// check image
	_, err = pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorValidateFailed)
	}

	// check subnet, vpc, eip
	subnetResponse, err := providerClient.DescribeSubnets(ctx, &pb.DescribeSubnetsRequest{
		RuntimeId: pbutil.ToProtoString(runtime.RuntimeId),
		SubnetId:  []string{clusterWrapper.Cluster.SubnetId},
		Zone:      []string{clusterWrapper.Cluster.Zone},
	})
	if err != nil {
		logger.Error(ctx, "Describe subnet [%s] runtime [%s] failed. ", clusterWrapper.Cluster.SubnetId, runtime.RuntimeId)
		return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.SubnetId)
	}
	vpcId := ""
	if subnetResponse != nil && len(subnetResponse.SubnetSet) == 1 {
		vpcId = subnetResponse.SubnetSet[0].GetVpcId().GetValue()
	}
	if vpcId == "" {
		err = fmt.Errorf("subnet [%s] not found or vpc not bind eip", clusterWrapper.Cluster.SubnetId)
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorSubnetNotFound, clusterWrapper.Cluster.SubnetId)
	}
	clusterWrapper.Cluster.VpcId = vpcId

	// check resource
	response, err := providerClient.CheckResource(ctx, &pb.CheckResourceRequest{
		RuntimeId: pbutil.ToProtoString(runtime.RuntimeId),
		Cluster:   models.ClusterWrapperToPb(clusterWrapper),
	})
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceQuotaNotEnough, err.Error())
	}
	if !response.Ok.GetValue() {
		return fmt.Errorf("response is not ok")
	}

	fg := &Frontgate{
		Runtime: runtime,
	}
	frontgate, err := fg.GetActiveFrontgate(ctx, clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Get frontgate in vpc [%s] user [%s] failed. ", clusterWrapper.Cluster.VpcId, clusterWrapper.Cluster.Owner)
		return err
	}

	clusterWrapper.Cluster.FrontgateId = frontgate.ClusterId

	return nil
}
