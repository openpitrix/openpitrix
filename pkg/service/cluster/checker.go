package cluster

import (
	"context"
	"fmt"

	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
)

func checkPermissionAndTransition(ctx context.Context, clusterId, userId string, status []string) (*models.Cluster, error) {
	cluster, err := getCluster(ctx, clusterId, userId)
	if err != nil {
		return nil, err
	}
	if cluster.TransitionStatus != "" {
		logger.Error(ctx, "Cluster [%s] is [%s], please try later", clusterId, cluster.TransitionStatus)
		return nil, fmt.Errorf("cluster [%s] is [%s], please try later", clusterId, cluster.TransitionStatus)
	}
	if status != nil && !reflectutil.In(cluster.Status, status) {
		logger.Error(ctx, "Cluster [%s] status is [%s] not in %s", clusterId, cluster.Status, status)
		return nil, fmt.Errorf("cluster [%s] status is [%s] not in %s", clusterId, cluster.Status, status)
	}
	return cluster, nil
}

func checkNodesPermissionAndTransition(ctx context.Context, nodeIds []string, userId string, status []string) ([]*models.ClusterNode, error) {
	clusterNodes, err := getClusterNodes(ctx, nodeIds, userId)
	if err != nil {
		return nil, err
	}
	for _, clusterNode := range clusterNodes {
		if clusterNode.TransitionStatus != "" {
			logger.Error(ctx, "Cluster node [%s] is [%s], please try later", clusterNode.NodeId, clusterNode.TransitionStatus)
			return nil, fmt.Errorf("cluster [%s] is [%s], please try later", clusterNode.NodeId, clusterNode.TransitionStatus)
		}
		if status != nil && !reflectutil.In(clusterNode.Status, status) {
			logger.Error(ctx, "Cluster [%s] status is [%s] not in %s", clusterNode.NodeId, clusterNode.Status, status)
			return nil, fmt.Errorf("cluster [%s] status is [%s] not in %s", clusterNode.NodeId, clusterNode.Status, status)
		}
	}

	return clusterNodes, nil
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
			Required("cluster_id", "role").
			Exec()
	case *pb.AddClusterNodesRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "role").
			Exec()
	case *pb.DeleteClusterNodesRequest:
		return manager.NewChecker(ctx, r).
			Required("cluster_id", "role", "node_id").
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
	}
	return nil
}

func CheckVmBasedProvider(ctx context.Context, runtime *runtimeclient.Runtime, providerInterface plugins.ProviderInterface,
	clusterWrapper *models.ClusterWrapper) error {

	// check pilot service
	pilotClient, err := pilotclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Connect to pilot service failed")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	_, err = pilotClient.PingPilot(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Error(ctx, "Pilot service is not running")
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// check image
	_, err = pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorValidateFailed)
	}

	// check subnet, vpc, eip
	subnetResponse, err := providerInterface.DescribeSubnets(ctx, &pb.DescribeSubnetsRequest{
		RuntimeId: pbutil.ToProtoString(runtime.RuntimeId),
		SubnetId:  []string{clusterWrapper.Cluster.SubnetId},
		Zone:      []string{clusterWrapper.Cluster.Zone},
	})
	if err != nil {
		logger.Error(ctx, "Describe subnet [%s] runtime [%s] failed. ", clusterWrapper.Cluster.SubnetId, runtime)
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
	err = providerInterface.CheckResource(ctx, clusterWrapper)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceQuotaNotEnough, err.Error())
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
