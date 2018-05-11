package cluster

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/logger"
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
