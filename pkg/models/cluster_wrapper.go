// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type ClusterWrapper struct {
	ctx                      context.Context
	Cluster                  *Cluster
	ClusterNodesWithKeyPairs map[string]*ClusterNodeWithKeyPairs // key=nodeId
	ClusterCommons           map[string]*ClusterCommon           // key=role
	ClusterLinks             map[string]*ClusterLink             // key=name
	ClusterRoles             map[string]*ClusterRole             // key=role
	ClusterLoadbalancers     map[string][]*ClusterLoadbalancer   // key=role
}

func NewClusterWrapper(ctx context.Context, data string) (*ClusterWrapper, error) {
	clusterWrapper := &ClusterWrapper{
		ctx: ctx,
	}
	err := jsonutil.Decode([]byte(data), clusterWrapper)
	if err != nil {
		logger.Error(ctx, "Decode [%s] into cluster wrapper failed: %+v", data, err)
	}
	return clusterWrapper, err
}

func ClusterWrapperToPb(clusterWrapper *ClusterWrapper) *pb.Cluster {

	pbCluster := ClusterToPb(clusterWrapper.Cluster)

	var clusterCommons []*ClusterCommon
	var clusterNodesWithKeyPairs []*ClusterNodeWithKeyPairs
	var clusterRoles []*ClusterRole
	var clusterLinks []*ClusterLink
	var clusterLoadbalancers []*ClusterLoadbalancer

	for _, clusterCommon := range clusterWrapper.ClusterCommons {
		clusterCommons = append(clusterCommons, clusterCommon)
	}
	pbCluster.ClusterCommonSet = ClusterCommonsToPbs(clusterCommons)

	for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
		clusterNodesWithKeyPairs = append(clusterNodesWithKeyPairs, clusterNode)
	}
	pbCluster.ClusterNodeSet = ClusterNodesWithKeyPairsToPbs(clusterNodesWithKeyPairs)

	for _, clusterRole := range clusterWrapper.ClusterRoles {
		clusterRoles = append(clusterRoles, clusterRole)
	}
	pbCluster.ClusterRoleSet = ClusterRolesToPbs(clusterRoles)

	for _, clusterLink := range clusterWrapper.ClusterLinks {
		clusterLinks = append(clusterLinks, clusterLink)
	}
	pbCluster.ClusterLinkSet = ClusterLinksToPbs(clusterLinks)

	for _, clusterLoadbalancer := range clusterWrapper.ClusterLoadbalancers {
		clusterLoadbalancers = append(clusterLoadbalancers, clusterLoadbalancer...)
	}
	pbCluster.ClusterLoadbalancerSet = ClusterLoadbalancersToPbs(clusterLoadbalancers)

	return pbCluster
}

func ClusterNodeWrapperToPb(clusterNode *ClusterNodeWithKeyPairs, clusterCommon *ClusterCommon,
	clusterRole *ClusterRole) *pb.ClusterNode {

	pbClusterNode := ClusterNodeWithKeyPairsToPb(clusterNode)
	pbClusterNode.ClusterCommon = ClusterCommonToPb(clusterCommon)
	pbClusterNode.ClusterRole = ClusterRoleToPb(clusterRole)

	return pbClusterNode
}

func PbToClusterWrapper(pbCluster *pb.Cluster) *ClusterWrapper {
	clusterWrapper := new(ClusterWrapper)
	clusterWrapper.Cluster = PbToCluster(pbCluster)

	clusterWrapper.ClusterCommons = make(map[string]*ClusterCommon)
	for _, pbClusterCommon := range pbCluster.ClusterCommonSet {
		clusterWrapper.ClusterCommons[pbClusterCommon.GetRole().GetValue()] = PbToClusterCommon(pbClusterCommon)
	}

	clusterWrapper.ClusterNodesWithKeyPairs = make(map[string]*ClusterNodeWithKeyPairs)
	for _, pbClusterNode := range pbCluster.ClusterNodeSet {
		clusterWrapper.ClusterNodesWithKeyPairs[pbClusterNode.GetNodeId().GetValue()] = PbToClusterNode(pbClusterNode)
	}

	clusterWrapper.ClusterRoles = make(map[string]*ClusterRole)
	for _, pbClusterRole := range pbCluster.ClusterRoleSet {
		clusterWrapper.ClusterRoles[pbClusterRole.GetRole().GetValue()] = PbToClusterRole(pbClusterRole)
	}

	clusterWrapper.ClusterLinks = make(map[string]*ClusterLink)
	for _, pbClusterLink := range pbCluster.ClusterLinkSet {
		clusterWrapper.ClusterLinks[pbClusterLink.GetName().GetValue()] = PbToClusterLink(pbClusterLink)
	}

	clusterWrapper.ClusterLoadbalancers = make(map[string][]*ClusterLoadbalancer)
	for _, pbClusterLoadbalancer := range pbCluster.ClusterLoadbalancerSet {
		clusterWrapper.ClusterLoadbalancers[pbClusterLoadbalancer.GetRole().GetValue()] =
			append(clusterWrapper.ClusterLoadbalancers[pbClusterLoadbalancer.GetRole().GetValue()],
				PbToClusterLoadbalancer(pbClusterLoadbalancer))
	}
	return clusterWrapper
}

func (c *ClusterWrapper) GetCommonAttribute(role, attributeName string) interface{} {
	if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
		role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
	}

	clusterCommon, exist := c.ClusterCommons[role]
	if !exist {
		logger.Error(c.ctx, "No such role [%s] in cluster [%s]. ",
			role, c.Cluster.ClusterId)
		return nil
	}

	return clusterCommon.GetAttribute(attributeName)
}

/*
endpoints is in the following format:
{
  "client_port": {
	  "port": 2181,
	  "protocol": "tcp"
  },
  "reserved_ips": {
	"write_vip":{

	},
	 "read_vip":{

	}
  }
}
where client_port is a developer-defined name. Port either is an integer or a reference
to an env variable such as env.<port> or env.<role>.<port>. It may have multiple endpoints defined.
*/
func (c *ClusterWrapper) GetEndpoints() (map[string]map[string]interface{}, error) {
	if c.Cluster.Endpoints != "" {
		endpoints := make(map[string]map[string]interface{})
		err := jsonutil.Decode([]byte(c.Cluster.Endpoints), &endpoints)
		if err != nil {
			logger.Error(c.ctx, "Unmarshal cluster [%s] endpoints failed: %+v", c.Cluster.ClusterId, err)
			return nil, err
		}
		for _, service := range endpoints {
			port, exist := service["port"]
			if !exist {
				continue
			} else {
				switch v := port.(type) {
				case string:
					portInfo := strings.Split(v, ".")
					var param string
					var cRole *ClusterRole
					if len(portInfo) >= 2 {
						if portInfo[0] == "env" {
							// no role associated with, choose the first node
							param = strings.Join(portInfo[1:], ".")
							for _, clusterRole := range c.ClusterRoles {
								cRole = clusterRole
								break
							}
						} else {
							// the first part of the port should be role name
							role := portInfo[0]
							param = strings.Join(portInfo[2:], ".")
							cRole = c.ClusterRoles[role]
						}
					} else {
						logger.Error(c.ctx, "Link [%s] in endpoints must be in env.x or <role name>.env.x for the cluster [%s]",
							port, c.Cluster.ClusterId)
						return nil, fmt.Errorf("Cluster [%s] endpoints link error. ", c.Cluster.ClusterId)
					}
					if cRole == nil {
						logger.Error(c.ctx, "Can't find the node of the cluster [%s] for the endpoints", c.Cluster.ClusterId)
						return nil, fmt.Errorf("Cluster [%s] endpoints parse failed. ", c.Cluster.ClusterId)
					}
					env := make(map[string]interface{})
					err = jsonutil.Decode([]byte(cRole.Env), &env)
					if err != nil {
						logger.Error(c.ctx, "Unmarshal cluster [%s] env failed: %+v", c.Cluster.ClusterId, err)
						return nil, err
					}
					value, exist := env[param]
					if exist {
						service["port"] = value
					}

				default:
					continue
				}
			}
		}
		return endpoints, nil
	} else {
		return nil, nil
	}
}
