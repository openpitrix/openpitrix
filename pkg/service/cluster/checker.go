package cluster

import (
	"context"
	"fmt"
	"strings"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
)

func checkPermissionAndTransition(clusterId, userId string, status []string) error {
	cluster, err := getCluster(clusterId, userId)
	if err != nil {
		return err
	}
	if cluster.TransitionStatus != "" {
		logger.Error("Cluster [%s] is [%s], please try later", clusterId, cluster.TransitionStatus)
		return fmt.Errorf("cluster [%s] is [%s], please try later", clusterId, cluster.TransitionStatus)
	}
	if status != nil && !reflectutil.In(cluster.Status, status) {
		logger.Error("Cluster [%s] status is [%s] not in %s", clusterId, cluster.Status, status)
		return fmt.Errorf("cluster [%s] status is [%s] not in %s", clusterId, cluster.Status, status)
	}
	return nil
}

func isActionSupported(clusterId, role, action string) bool {
	clusterWrapper, err := getClusterWrapper(clusterId)
	if err != nil {
		return false
	}
	clusterCommon, exist := clusterWrapper.ClusterCommons[role]
	if !exist {
		logger.Error("Cluster [%s] has no role [%s]", clusterId, role)
		return false
	}
	advanceActions := clusterCommon.AdvancedActions
	if advanceActions == "" {
		return false
	}
	actions := strings.Split(advanceActions, ",")
	if reflectutil.In(action, actions) {
		return true
	} else {
		return false
	}
}

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
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
	// check image
	_, _, err := pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorValidateFailed)
	}

	// check subnet, vpc, eip
	subnetResponse, err := providerInterface.DescribeSubnets(ctx, &pb.DescribeSubnetsRequest{
		RuntimeId: pbutil.ToProtoString(runtime.RuntimeId),
		SubnetId:  []string{clusterWrapper.Cluster.SubnetId},
	})
	if err != nil {
		logger.Error("Describe subnet [%s] runtime [%s] failed. ", clusterWrapper.Cluster.SubnetId, runtime)
		return gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, clusterWrapper.Cluster.SubnetId)
	}
	vpcId := ""
	if subnetResponse != nil && len(subnetResponse.SubnetSet) == 1 {
		vpcId = subnetResponse.SubnetSet[0].GetVpcId().GetValue()
	}
	if vpcId == "" {
		err = fmt.Errorf("subnet [%s] not found or vpc not bind eip", clusterWrapper.Cluster.SubnetId)
		return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorSubnetNotFound, clusterWrapper.Cluster.SubnetId)
	}
	clusterWrapper.Cluster.VpcId = vpcId

	// check resource quota
	err = providerInterface.CheckResourceQuotas(ctx, clusterWrapper)
	if err != nil {
		return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceQuotaNotEnough, err.Error())
	}

	fg := &Frontgate{
		Runtime: runtime,
	}
	frontgate, err := fg.GetActiveFrontgate(clusterWrapper)
	if err != nil {
		logger.Error("Get frontgate in vpc [%s] user [%s] failed. ", clusterWrapper.Cluster.VpcId, clusterWrapper.Cluster.Owner)
		return err
	}

	clusterWrapper.Cluster.FrontgateId = frontgate.ClusterId

	return nil
}
