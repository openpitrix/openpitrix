// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package service

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

var _ fmt.State
var _ time.Time

type ClusterService struct {
	Config     *config.Config
	Properties *ClusterServiceProperties
}

type ClusterServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Cluster(zone string) (*ClusterService, error) {
	properties := &ClusterServiceProperties{
		Zone: &zone,
	}

	return &ClusterService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/add_cluster_nodes.html
func (s *ClusterService) AddClusterNodes(i *AddClusterNodesInput) (*AddClusterNodesOutput, error) {
	if i == nil {
		i = &AddClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddClusterNodes",
		RequestMethod: "GET",
	}

	x := &AddClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type AddClusterNodesInput struct {
	Cluster      *string   `json:"cluster" name:"cluster" location:"params"`       // Required
	NodeCount    *int      `json:"node_count" name:"node_count" location:"params"` // Required
	NodeName     *string   `json:"node_name" name:"node_name" location:"params"`
	NodeRole     *string   `json:"node_role" name:"node_role" location:"params"`
	PrivateIPs   []*string `json:"private_ips" name:"private_ips" location:"params"`
	ResourceConf *string   `json:"resource_conf" name:"resource_conf" location:"params"`
}

func (v *AddClusterNodesInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "AddClusterNodesInput",
		}
	}

	if v.NodeCount == nil {
		return errors.ParameterRequiredError{
			ParameterName: "NodeCount",
			ParentName:    "AddClusterNodesInput",
		}
	}

	return nil
}

type AddClusterNodesOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	ClusterID  *string   `json:"cluster_id" name:"cluster_id" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	NewNodeIDs []*string `json:"new_node_ids" name:"new_node_ids" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/associate_eip_to_cluster_node.html
func (s *ClusterService) AssociateEIPToClusterNode(i *AssociateEIPToClusterNodeInput) (*AssociateEIPToClusterNodeOutput, error) {
	if i == nil {
		i = &AssociateEIPToClusterNodeInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AssociateEipToClusterNode",
		RequestMethod: "GET",
	}

	x := &AssociateEIPToClusterNodeOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type AssociateEIPToClusterNodeInput struct {
	ClusterNode *string `json:"cluster_node" name:"cluster_node" location:"params"` // Required
	EIP         *string `json:"eip" name:"eip" location:"params"`                   // Required
	NIC         *string `json:"nic" name:"nic" location:"params"`
}

func (v *AssociateEIPToClusterNodeInput) Validate() error {

	if v.ClusterNode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ClusterNode",
			ParentName:    "AssociateEIPToClusterNodeInput",
		}
	}

	if v.EIP == nil {
		return errors.ParameterRequiredError{
			ParameterName: "EIP",
			ParentName:    "AssociateEIPToClusterNodeInput",
		}
	}

	return nil
}

type AssociateEIPToClusterNodeOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/cease_clusters.html
func (s *ClusterService) CeaseClusters(i *CeaseClustersInput) (*CeaseClustersOutput, error) {
	if i == nil {
		i = &CeaseClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CeaseClusters",
		RequestMethod: "GET",
	}

	x := &CeaseClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type CeaseClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"` // Required
}

func (v *CeaseClustersInput) Validate() error {

	if len(v.Clusters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Clusters",
			ParentName:    "CeaseClustersInput",
		}
	}

	return nil
}

type CeaseClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/change_cluster_vxnet.html
func (s *ClusterService) ChangeClusterVxNet(i *ChangeClusterVxNetInput) (*ChangeClusterVxNetOutput, error) {
	if i == nil {
		i = &ChangeClusterVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeClusterVxnet",
		RequestMethod: "GET",
	}

	x := &ChangeClusterVxNetOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ChangeClusterVxNetInput struct {
	Cluster    *string     `json:"cluster" name:"cluster" location:"params"` // Required
	PrivateIPs interface{} `json:"private_ips" name:"private_ips" location:"params"`
	Roles      []*string   `json:"roles" name:"roles" location:"params"`
	VxNet      *string     `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *ChangeClusterVxNetInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "ChangeClusterVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "ChangeClusterVxNetInput",
		}
	}

	return nil
}

type ChangeClusterVxNetOutput struct {
	Message   *string `json:"message" name:"message"`
	Action    *string `json:"action" name:"action" location:"elements"`
	ClusterID *string `json:"cluster_id" name:"cluster_id" location:"elements"`
	JobID     *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int    `json:"ret_code" name:"ret_code" location:"elements"`
	VxNetID   *string `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/create_cluster.html
func (s *ClusterService) CreateCluster(i *CreateClusterInput) (*CreateClusterOutput, error) {
	if i == nil {
		i = &CreateClusterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateCluster",
		RequestMethod: "GET",
	}

	x := &CreateClusterOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type CreateClusterInput struct {
	Conf *string `json:"conf" name:"conf" location:"params"` // Required
}

func (v *CreateClusterInput) Validate() error {

	if v.Conf == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Conf",
			ParentName:    "CreateClusterInput",
		}
	}

	return nil
}

type CreateClusterOutput struct {
	Message     *string   `json:"message" name:"message"`
	Action      *string   `json:"action" name:"action" location:"elements"`
	AppID       *string   `json:"app_id" name:"app_id" location:"elements"`
	AppVersion  *string   `json:"app_version" name:"app_version" location:"elements"`
	ClusterID   *string   `json:"cluster_id" name:"cluster_id" location:"elements"`
	ClusterName *string   `json:"cluster_name" name:"cluster_name" location:"elements"`
	JobID       *string   `json:"job_id" name:"job_id" location:"elements"`
	NodeIDs     []*string `json:"node_ids" name:"node_ids" location:"elements"`
	RetCode     *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNetID     *string   `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/create_cluster_from_snapshot.html
func (s *ClusterService) CreateClusterFromSnapshot(i *CreateClusterFromSnapshotInput) (*CreateClusterFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateClusterFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateClusterFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateClusterFromSnapshotOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type CreateClusterFromSnapshotInput struct {
	Conf       *string `json:"conf" name:"conf" location:"params"`               // Required
	SnapshotID *string `json:"snapshot_id" name:"snapshot_id" location:"params"` // Required
}

func (v *CreateClusterFromSnapshotInput) Validate() error {

	if v.Conf == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Conf",
			ParentName:    "CreateClusterFromSnapshotInput",
		}
	}

	if v.SnapshotID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SnapshotID",
			ParentName:    "CreateClusterFromSnapshotInput",
		}
	}

	return nil
}

type CreateClusterFromSnapshotOutput struct {
	Message     *string   `json:"message" name:"message"`
	Action      *string   `json:"action" name:"action" location:"elements"`
	AppID       *string   `json:"app_id" name:"app_id" location:"elements"`
	AppVersion  *string   `json:"app_version" name:"app_version" location:"elements"`
	ClusterID   *string   `json:"cluster_id" name:"cluster_id" location:"elements"`
	ClusterName *string   `json:"cluster_name" name:"cluster_name" location:"elements"`
	JobID       *string   `json:"job_id" name:"job_id" location:"elements"`
	NodeIDs     []*string `json:"node_ids" name:"node_ids" location:"elements"`
	RetCode     *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNetID     *string   `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/delete_cluster_nodes.html
func (s *ClusterService) DeleteClusterNodes(i *DeleteClusterNodesInput) (*DeleteClusterNodesOutput, error) {
	if i == nil {
		i = &DeleteClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteClusterNodes",
		RequestMethod: "GET",
	}

	x := &DeleteClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DeleteClusterNodesInput struct {
	Cluster *string   `json:"cluster" name:"cluster" location:"params"` // Required
	Force   *int      `json:"force" name:"force" location:"params"`
	Nodes   []*string `json:"nodes" name:"nodes" location:"params"` // Required
}

func (v *DeleteClusterNodesInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "DeleteClusterNodesInput",
		}
	}

	if len(v.Nodes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Nodes",
			ParentName:    "DeleteClusterNodesInput",
		}
	}

	return nil
}

type DeleteClusterNodesOutput struct {
	Message        *string   `json:"message" name:"message"`
	Action         *string   `json:"action" name:"action" location:"elements"`
	ClusterID      *string   `json:"cluster_id" name:"cluster_id" location:"elements"`
	DeletedNodeIDs []*string `json:"deleted_node_ids" name:"deleted_node_ids" location:"elements"`
	JobID          *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode        *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/delete_clusters.html
func (s *ClusterService) DeleteClusters(i *DeleteClustersInput) (*DeleteClustersOutput, error) {
	if i == nil {
		i = &DeleteClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteClusters",
		RequestMethod: "GET",
	}

	x := &DeleteClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DeleteClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"` // Required
	Force    *int      `json:"force" name:"force" location:"params"`
}

func (v *DeleteClustersInput) Validate() error {

	if len(v.Clusters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Clusters",
			ParentName:    "DeleteClustersInput",
		}
	}

	return nil
}

type DeleteClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_cluster_display_tabs.html
func (s *ClusterService) DescribeClusterDisplayTabs(i *DescribeClusterDisplayTabsInput) (*DescribeClusterDisplayTabsOutput, error) {
	if i == nil {
		i = &DescribeClusterDisplayTabsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusterDisplayTabs",
		RequestMethod: "GET",
	}

	x := &DescribeClusterDisplayTabsOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClusterDisplayTabsInput struct {
	Cluster     *string `json:"cluster" name:"cluster" location:"params"`           // Required
	DisplayTabs *string `json:"display_tabs" name:"display_tabs" location:"params"` // Required
	Role        *string `json:"role" name:"role" location:"params"`
}

func (v *DescribeClusterDisplayTabsInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "DescribeClusterDisplayTabsInput",
		}
	}

	if v.DisplayTabs == nil {
		return errors.ParameterRequiredError{
			ParameterName: "DisplayTabs",
			ParentName:    "DescribeClusterDisplayTabsInput",
		}
	}

	return nil
}

type DescribeClusterDisplayTabsOutput struct {
	Message     *string            `json:"message" name:"message"`
	Action      *string            `json:"action" name:"action" location:"elements"`
	DisplayTabs map[string]*string `json:"display_tabs" name:"display_tabs" location:"elements"`
	RetCode     *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_cluster_nodes.html
func (s *ClusterService) DescribeClusterNodes(i *DescribeClusterNodesInput) (*DescribeClusterNodesOutput, error) {
	if i == nil {
		i = &DescribeClusterNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusterNodes",
		RequestMethod: "GET",
	}

	x := &DescribeClusterNodesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClusterNodesInput struct {
	Cluster    *string   `json:"cluster" name:"cluster" location:"params"`
	Console    *string   `json:"console" name:"console" location:"params"`
	Limit      *int      `json:"limit" name:"limit" location:"params"`
	Nodes      []*string `json:"nodes" name:"nodes" location:"params"`
	Offset     *int      `json:"offset" name:"offset" location:"params"`
	Owner      *string   `json:"owner" name:"owner" location:"params"`
	Reverse    *int      `json:"reverse" name:"reverse" location:"params"`
	Role       *string   `json:"role" name:"role" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	SortKey    *string   `json:"sort_key" name:"sort_key" location:"params"`
	Status     *string   `json:"status" name:"status" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeClusterNodesInput) Validate() error {

	return nil
}

type DescribeClusterNodesOutput struct {
	Message    *string        `json:"message" name:"message"`
	Action     *string        `json:"action" name:"action" location:"elements"`
	NodeSet    []*ClusterNode `json:"node_set" name:"node_set" location:"elements"`
	RetCode    *int           `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int           `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_cluster_users.html
func (s *ClusterService) DescribeClusterUsers(i *DescribeClusterUsersInput) (*DescribeClusterUsersOutput, error) {
	if i == nil {
		i = &DescribeClusterUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusterUsers",
		RequestMethod: "GET",
	}

	x := &DescribeClusterUsersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClusterUsersInput struct {
	AppVersions   []*string `json:"app_versions" name:"app_versions" location:"params"`
	Apps          []*string `json:"apps" name:"apps" location:"params"` // Required
	ClusterStatus []*string `json:"cluster_status" name:"cluster_status" location:"params"`
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Users         []*string `json:"users" name:"users" location:"params"`
	Zones         []*string `json:"zones" name:"zones" location:"params"` // Required
}

func (v *DescribeClusterUsersInput) Validate() error {

	if len(v.Apps) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Apps",
			ParentName:    "DescribeClusterUsersInput",
		}
	}

	if len(v.Zones) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Zones",
			ParentName:    "DescribeClusterUsersInput",
		}
	}

	return nil
}

type DescribeClusterUsersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	Apps    []*string          `json:"apps" name:"apps" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
	Users   map[string]*string `json:"users" name:"users" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/describe_clusters.html
func (s *ClusterService) DescribeClusters(i *DescribeClustersInput) (*DescribeClustersOutput, error) {
	if i == nil {
		i = &DescribeClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeClusters",
		RequestMethod: "GET",
	}

	x := &DescribeClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeClustersInput struct {
	AppVersions       []*string `json:"app_versions" name:"app_versions" location:"params"`
	Apps              []*string `json:"apps" name:"apps" location:"params"`
	CfgmgmtID         *string   `json:"cfgmgmt_id" name:"cfgmgmt_id" location:"params"`
	Clusters          []*string `json:"clusters" name:"clusters" location:"params"`
	Console           *string   `json:"console" name:"console" location:"params"`
	ExternalClusterID *string   `json:"external_cluster_id" name:"external_cluster_id" location:"params"`
	Limit             *int      `json:"limit" name:"limit" location:"params"`
	Link              *string   `json:"link" name:"link" location:"params"`
	Name              *string   `json:"name" name:"name" location:"params"`
	Offset            *int      `json:"offset" name:"offset" location:"params"`
	Owner             *string   `json:"owner" name:"owner" location:"params"`
	Reverse           *int      `json:"reverse" name:"reverse" location:"params"`
	Role              *string   `json:"role" name:"role" location:"params"`
	// Scope's available values: all, cfgmgmt
	Scope            *string   `json:"scope" name:"scope" location:"params"`
	SearchWord       *string   `json:"search_word" name:"search_word" location:"params"`
	SortKey          *string   `json:"sort_key" name:"sort_key" location:"params"`
	Status           *string   `json:"status" name:"status" location:"params"`
	TransitionStatus *string   `json:"transition_status" name:"transition_status" location:"params"`
	Users            []*string `json:"users" name:"users" location:"params"`
	Verbose          *int      `json:"verbose" name:"verbose" location:"params"`
	VxNet            *string   `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *DescribeClustersInput) Validate() error {

	if v.Scope != nil {
		scopeValidValues := []string{"all", "cfgmgmt"}
		scopeParameterValue := fmt.Sprint(*v.Scope)

		scopeIsValid := false
		for _, value := range scopeValidValues {
			if value == scopeParameterValue {
				scopeIsValid = true
			}
		}

		if !scopeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Scope",
				ParameterValue: scopeParameterValue,
				AllowedValues:  scopeValidValues,
			}
		}
	}

	return nil
}

type DescribeClustersOutput struct {
	Message    *string    `json:"message" name:"message"`
	Action     *string    `json:"action" name:"action" location:"elements"`
	ClusterSet []*Cluster `json:"cluster_set" name:"cluster_set" location:"elements"`
	RetCode    *int       `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int       `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/dissociate_eip_from_cluster_node.html
func (s *ClusterService) DissociateEIPFromClusterNode(i *DissociateEIPFromClusterNodeInput) (*DissociateEIPFromClusterNodeOutput, error) {
	if i == nil {
		i = &DissociateEIPFromClusterNodeInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DissociateEipFromClusterNode",
		RequestMethod: "GET",
	}

	x := &DissociateEIPFromClusterNodeOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DissociateEIPFromClusterNodeInput struct {
	EIPs []*string `json:"eips" name:"eips" location:"params"` // Required
}

func (v *DissociateEIPFromClusterNodeInput) Validate() error {

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "DissociateEIPFromClusterNodeInput",
		}
	}

	return nil
}

type DissociateEIPFromClusterNodeOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/modify_cluster_attributes.html
func (s *ClusterService) ModifyClusterAttributes(i *ModifyClusterAttributesInput) (*ModifyClusterAttributesOutput, error) {
	if i == nil {
		i = &ModifyClusterAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyClusterAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyClusterAttributesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ModifyClusterAttributesInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	Cluster        *string `json:"cluster" name:"cluster" location:"params"` // Required
	Description    *string `json:"description" name:"description" location:"params"`
	Name           *string `json:"name" name:"name" location:"params"`
}

func (v *ModifyClusterAttributesInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "ModifyClusterAttributesInput",
		}
	}

	return nil
}

type ModifyClusterAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/modify_cluster_node_attributes.html
func (s *ClusterService) ModifyClusterNodeAttributes(i *ModifyClusterNodeAttributesInput) (*ModifyClusterNodeAttributesOutput, error) {
	if i == nil {
		i = &ModifyClusterNodeAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyClusterNodeAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyClusterNodeAttributesOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ModifyClusterNodeAttributesInput struct {
	ClusterNode *string `json:"cluster_node" name:"cluster_node" location:"params"` // Required
	Name        *string `json:"name" name:"name" location:"params"`
}

func (v *ModifyClusterNodeAttributesInput) Validate() error {

	if v.ClusterNode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ClusterNode",
			ParentName:    "ModifyClusterNodeAttributesInput",
		}
	}

	return nil
}

type ModifyClusterNodeAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/recover_clusters.html
func (s *ClusterService) RecoverClusters(i *RecoverClustersInput) (*RecoverClustersOutput, error) {
	if i == nil {
		i = &RecoverClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "Lease",
		RequestMethod: "GET",
	}

	x := &RecoverClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type RecoverClustersInput struct {
	Resources []*string `json:"resources" name:"resources" location:"params"`
}

func (v *RecoverClustersInput) Validate() error {

	return nil
}

type RecoverClustersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/resize_cluster.html
func (s *ClusterService) ResizeCluster(i *ResizeClusterInput) (*ResizeClusterOutput, error) {
	if i == nil {
		i = &ResizeClusterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeCluster",
		RequestMethod: "GET",
	}

	x := &ResizeClusterOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type ResizeClusterInput struct {
	Cluster     *string `json:"cluster" name:"cluster" location:"params"` // Required
	CPU         *int    `json:"cpu" name:"cpu" location:"params"`
	Gpu         *int    `json:"gpu" name:"gpu" location:"params"`
	Memory      *int    `json:"memory" name:"memory" location:"params"`
	NodeRole    *string `json:"node_role" name:"node_role" location:"params"`
	StorageSize *int    `json:"storage_size" name:"storage_size" location:"params"`
}

func (v *ResizeClusterInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "ResizeClusterInput",
		}
	}

	return nil
}

type ResizeClusterOutput struct {
	Message     *string `json:"message" name:"message"`
	Action      *string `json:"action" name:"action" location:"elements"`
	ClusterID   *string `json:"cluster_id" name:"cluster_id" location:"elements"`
	CPU         *int    `json:"cpu" name:"cpu" location:"elements"`
	Gpu         *int    `json:"gpu" name:"gpu" location:"elements"`
	JobID       *string `json:"job_id" name:"job_id" location:"elements"`
	Memory      *int    `json:"memory" name:"memory" location:"elements"`
	RetCode     *int    `json:"ret_code" name:"ret_code" location:"elements"`
	Role        *string `json:"role" name:"role" location:"elements"`
	StorageSize *int    `json:"storage_size" name:"storage_size" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/restart_cluster_service.html
func (s *ClusterService) RestartClusterService(i *RestartClusterServiceInput) (*RestartClusterServiceOutput, error) {
	if i == nil {
		i = &RestartClusterServiceInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestartClusterService",
		RequestMethod: "GET",
	}

	x := &RestartClusterServiceOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type RestartClusterServiceInput struct {
	Cluster *string `json:"cluster" name:"cluster" location:"params"`
	Role    *string `json:"role" name:"role" location:"params"`
}

func (v *RestartClusterServiceInput) Validate() error {

	return nil
}

type RestartClusterServiceOutput struct {
	Message   *string `json:"message" name:"message"`
	Action    *string `json:"action" name:"action" location:"elements"`
	ClusterID *string `json:"cluster_id" name:"cluster_id" location:"elements"`
	JobID     *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int    `json:"ret_code" name:"ret_code" location:"elements"`
	Role      *string `json:"role" name:"role" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/restore_cluster_from_snapshot.html
func (s *ClusterService) RestoreClusterFromSnapshot(i *RestoreClusterFromSnapshotInput) (*RestoreClusterFromSnapshotOutput, error) {
	if i == nil {
		i = &RestoreClusterFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestoreClusterFromSnapshot",
		RequestMethod: "GET",
	}

	x := &RestoreClusterFromSnapshotOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type RestoreClusterFromSnapshotInput struct {
	Cluster       *string `json:"cluster" name:"cluster" location:"params"` // Required
	ServiceParams *string `json:"service_params" name:"service_params" location:"params"`
	Snapshot      *string `json:"snapshot" name:"snapshot" location:"params"` // Required
}

func (v *RestoreClusterFromSnapshotInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "RestoreClusterFromSnapshotInput",
		}
	}

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "RestoreClusterFromSnapshotInput",
		}
	}

	return nil
}

type RestoreClusterFromSnapshotOutput struct {
	Message       *string `json:"message" name:"message"`
	Action        *string `json:"action" name:"action" location:"elements"`
	ClusterID     *string `json:"cluster_id" name:"cluster_id" location:"elements"`
	JobID         *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode       *int    `json:"ret_code" name:"ret_code" location:"elements"`
	ServiceParams *string `json:"service_params" name:"service_params" location:"elements"`
	SnapshotID    *string `json:"snapshot_id" name:"snapshot_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/run_cluster_custom_service.html
func (s *ClusterService) RunClusterCustomService(i *RunClusterCustomServiceInput) (*RunClusterCustomServiceOutput, error) {
	if i == nil {
		i = &RunClusterCustomServiceInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RunClusterCustomService",
		RequestMethod: "GET",
	}

	x := &RunClusterCustomServiceOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type RunClusterCustomServiceInput struct {
	Cluster       *string `json:"cluster" name:"cluster" location:"params"` // Required
	Role          *string `json:"role" name:"role" location:"params"`
	Service       *string `json:"service" name:"service" location:"params"` // Required
	ServiceParams *string `json:"service_params" name:"service_params" location:"params"`
}

func (v *RunClusterCustomServiceInput) Validate() error {

	if v.Cluster == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cluster",
			ParentName:    "RunClusterCustomServiceInput",
		}
	}

	if v.Service == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Service",
			ParentName:    "RunClusterCustomServiceInput",
		}
	}

	return nil
}

type RunClusterCustomServiceOutput struct {
	Message   *string `json:"message" name:"message"`
	Action    *string `json:"action" name:"action" location:"elements"`
	ClusterID *string `json:"cluster_id" name:"cluster_id" location:"elements"`
	JobID     *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int    `json:"ret_code" name:"ret_code" location:"elements"`
	Role      *string `json:"role" name:"role" location:"elements"`
	Service   *string `json:"service" name:"service" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/start_clusters.html
func (s *ClusterService) StartClusters(i *StartClustersInput) (*StartClustersOutput, error) {
	if i == nil {
		i = &StartClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartClusters",
		RequestMethod: "GET",
	}

	x := &StartClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type StartClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"` // Required
}

func (v *StartClustersInput) Validate() error {

	if len(v.Clusters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Clusters",
			ParentName:    "StartClustersInput",
		}
	}

	return nil
}

type StartClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/stop_clusters.html
func (s *ClusterService) StopClusters(i *StopClustersInput) (*StopClustersOutput, error) {
	if i == nil {
		i = &StopClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopClusters",
		RequestMethod: "GET",
	}

	x := &StopClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type StopClustersInput struct {
	Clusters []*string `json:"clusters" name:"clusters" location:"params"` // Required
	Force    *int      `json:"force" name:"force" location:"params"`
}

func (v *StopClustersInput) Validate() error {

	if len(v.Clusters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Clusters",
			ParentName:    "StopClustersInput",
		}
	}

	return nil
}

type StopClustersOutput struct {
	Message *string            `json:"message" name:"message"`
	Action  *string            `json:"action" name:"action" location:"elements"`
	JobIDs  map[string]*string `json:"job_ids" name:"job_ids" location:"elements"`
	RetCode *int               `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/update_cluster_environment.html
func (s *ClusterService) UpdateClusterEnvironment(i *UpdateClusterEnvironmentInput) (*UpdateClusterEnvironmentOutput, error) {
	if i == nil {
		i = &UpdateClusterEnvironmentInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateClusterEnvironment",
		RequestMethod: "GET",
	}

	x := &UpdateClusterEnvironmentOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type UpdateClusterEnvironmentInput struct {
	Cluster *string     `json:"cluster" name:"cluster" location:"params"`
	Env     interface{} `json:"env" name:"env" location:"params"`
	Roles   []*string   `json:"roles" name:"roles" location:"params"`
}

func (v *UpdateClusterEnvironmentInput) Validate() error {

	return nil
}

type UpdateClusterEnvironmentOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cluster/upgrade_clusters.html
func (s *ClusterService) UpgradeClusters(i *UpgradeClustersInput) (*UpgradeClustersOutput, error) {
	if i == nil {
		i = &UpgradeClustersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpgradeClusters",
		RequestMethod: "GET",
	}

	x := &UpgradeClustersOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type UpgradeClustersInput struct {
	AppVersion    *string   `json:"app_version" name:"app_version" location:"params"`
	Clusters      []*string `json:"clusters" name:"clusters" location:"params"`
	ServiceParams *string   `json:"service_params" name:"service_params" location:"params"`
}

func (v *UpgradeClustersInput) Validate() error {

	return nil
}

type UpgradeClustersOutput struct {
	Message   *string   `json:"message" name:"message"`
	Action    *string   `json:"action" name:"action" location:"elements"`
	ClusterID []*string `json:"cluster_id" name:"cluster_id" location:"elements"`
	RetCode   *int      `json:"ret_code" name:"ret_code" location:"elements"`
}
