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

type CacheService struct {
	Config     *config.Config
	Properties *CacheServiceProperties
}

type CacheServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Cache(zone string) (*CacheService, error) {
	properties := &CacheServiceProperties{
		Zone: &zone,
	}

	return &CacheService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/cache/add_cache_nodes.html
func (s *CacheService) AddCacheNodes(i *AddCacheNodesInput) (*AddCacheNodesOutput, error) {
	if i == nil {
		i = &AddCacheNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddCacheNodes",
		RequestMethod: "GET",
	}

	x := &AddCacheNodesOutput{}
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

type AddCacheNodesInput struct {
	Cache      *string           `json:"cache" name:"cache" location:"params"`           // Required
	NodeCount  *int              `json:"node_count" name:"node_count" location:"params"` // Required
	PrivateIPs []*CachePrivateIP `json:"private_ips" name:"private_ips" location:"params"`
}

func (v *AddCacheNodesInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "AddCacheNodesInput",
		}
	}

	if v.NodeCount == nil {
		return errors.ParameterRequiredError{
			ParameterName: "NodeCount",
			ParentName:    "AddCacheNodesInput",
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AddCacheNodesOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/apply_cache_parameter_group.html
func (s *CacheService) ApplyCacheParameterGroup(i *ApplyCacheParameterGroupInput) (*ApplyCacheParameterGroupOutput, error) {
	if i == nil {
		i = &ApplyCacheParameterGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ApplyCacheParameterGroup",
		RequestMethod: "GET",
	}

	x := &ApplyCacheParameterGroupOutput{}
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

type ApplyCacheParameterGroupInput struct {
	CacheParameterGroup *string   `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"` // Required
	Caches              []*string `json:"caches" name:"caches" location:"params"`
}

func (v *ApplyCacheParameterGroupInput) Validate() error {

	if v.CacheParameterGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroup",
			ParentName:    "ApplyCacheParameterGroupInput",
		}
	}

	return nil
}

type ApplyCacheParameterGroupOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/change_cache_vxnet.html
func (s *CacheService) ChangeCacheVxNet(i *ChangeCacheVxNetInput) (*ChangeCacheVxNetOutput, error) {
	if i == nil {
		i = &ChangeCacheVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeCacheVxnet",
		RequestMethod: "GET",
	}

	x := &ChangeCacheVxNetOutput{}
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

type ChangeCacheVxNetInput struct {
	Cache      *string           `json:"cache" name:"cache" location:"params"` // Required
	PrivateIPs []*CachePrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	VxNet      *string           `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *ChangeCacheVxNetInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "ChangeCacheVxNetInput",
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "ChangeCacheVxNetInput",
		}
	}

	return nil
}

type ChangeCacheVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	CacheID *string `json:"cache_id" name:"cache_id" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
	VxNetID *string `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/create_cache.html
func (s *CacheService) CreateCache(i *CreateCacheInput) (*CreateCacheOutput, error) {
	if i == nil {
		i = &CreateCacheInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateCache",
		RequestMethod: "GET",
	}

	x := &CreateCacheOutput{}
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

type CreateCacheInput struct {
	AutoBackupTime *int `json:"auto_backup_time" name:"auto_backup_time" default:"-1" location:"params"`
	// CacheClass's available values: 0, 1
	CacheClass          *int              `json:"cache_class" name:"cache_class" location:"params"`
	CacheName           *string           `json:"cache_name" name:"cache_name" location:"params"`
	CacheParameterGroup *string           `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"`
	CacheSize           *int              `json:"cache_size" name:"cache_size" location:"params"` // Required
	CacheType           *string           `json:"cache_type" name:"cache_type" location:"params"` // Required
	MasterCount         *int              `json:"master_count" name:"master_count" location:"params"`
	NetworkType         *int              `json:"network_type" name:"network_type" location:"params"`
	NodeCount           *int              `json:"node_count" name:"node_count" default:"1" location:"params"`
	PrivateIPs          []*CachePrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	ReplicateCount      *int              `json:"replicate_count" name:"replicate_count" location:"params"`
	VxNet               *string           `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *CreateCacheInput) Validate() error {

	if v.CacheClass != nil {
		cacheClassValidValues := []string{"0", "1"}
		cacheClassParameterValue := fmt.Sprint(*v.CacheClass)

		cacheClassIsValid := false
		for _, value := range cacheClassValidValues {
			if value == cacheClassParameterValue {
				cacheClassIsValid = true
			}
		}

		if !cacheClassIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CacheClass",
				ParameterValue: cacheClassParameterValue,
				AllowedValues:  cacheClassValidValues,
			}
		}
	}

	if v.CacheSize == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheSize",
			ParentName:    "CreateCacheInput",
		}
	}

	if v.CacheType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheType",
			ParentName:    "CreateCacheInput",
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateCacheInput",
		}
	}

	return nil
}

type CreateCacheOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	CacheID    *string   `json:"cache_id" name:"cache_id" location:"elements"`
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/create_cache_from_snapshot.html
func (s *CacheService) CreateCacheFromSnapshot(i *CreateCacheFromSnapshotInput) (*CreateCacheFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateCacheFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateCacheFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateCacheFromSnapshotOutput{}
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

type CreateCacheFromSnapshotInput struct {
	AutoBackupTime *int `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	// CacheClass's available values: 0, 1
	CacheClass          *int              `json:"cache_class" name:"cache_class" location:"params"`
	CacheName           *string           `json:"cache_name" name:"cache_name" location:"params"`
	CacheParameterGroup *string           `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"`
	CacheSize           *int              `json:"cache_size" name:"cache_size" location:"params"`
	CacheType           *string           `json:"cache_type" name:"cache_type" location:"params"`
	NetworkType         *int              `json:"network_type" name:"network_type" location:"params"`
	NodeCount           *int              `json:"node_count" name:"node_count" location:"params"`
	PrivateIPs          []*CachePrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	Snapshot            *string           `json:"snapshot" name:"snapshot" location:"params"` // Required
	VxNet               *string           `json:"vxnet" name:"vxnet" location:"params"`       // Required
}

func (v *CreateCacheFromSnapshotInput) Validate() error {

	if v.CacheClass != nil {
		cacheClassValidValues := []string{"0", "1"}
		cacheClassParameterValue := fmt.Sprint(*v.CacheClass)

		cacheClassIsValid := false
		for _, value := range cacheClassValidValues {
			if value == cacheClassParameterValue {
				cacheClassIsValid = true
			}
		}

		if !cacheClassIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CacheClass",
				ParameterValue: cacheClassParameterValue,
				AllowedValues:  cacheClassValidValues,
			}
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "CreateCacheFromSnapshotInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateCacheFromSnapshotInput",
		}
	}

	return nil
}

type CreateCacheFromSnapshotOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	CacheID    *string   `json:"cache_id" name:"cache_id" location:"elements"`
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/create_cache_parameter_group.html
func (s *CacheService) CreateCacheParameterGroup(i *CreateCacheParameterGroupInput) (*CreateCacheParameterGroupOutput, error) {
	if i == nil {
		i = &CreateCacheParameterGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateCacheParameterGroup",
		RequestMethod: "GET",
	}

	x := &CreateCacheParameterGroupOutput{}
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

type CreateCacheParameterGroupInput struct {
	CacheParameterGroupName *string `json:"cache_parameter_group_name" name:"cache_parameter_group_name" location:"params"`
	// CacheType's available values: redis2.8.17, memcached1.4.13
	CacheType *string `json:"cache_type" name:"cache_type" location:"params"` // Required
}

func (v *CreateCacheParameterGroupInput) Validate() error {

	if v.CacheType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheType",
			ParentName:    "CreateCacheParameterGroupInput",
		}
	}

	if v.CacheType != nil {
		cacheTypeValidValues := []string{"redis2.8.17", "memcached1.4.13"}
		cacheTypeParameterValue := fmt.Sprint(*v.CacheType)

		cacheTypeIsValid := false
		for _, value := range cacheTypeValidValues {
			if value == cacheTypeParameterValue {
				cacheTypeIsValid = true
			}
		}

		if !cacheTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CacheType",
				ParameterValue: cacheTypeParameterValue,
				AllowedValues:  cacheTypeValidValues,
			}
		}
	}

	return nil
}

type CreateCacheParameterGroupOutput struct {
	Message               *string `json:"message" name:"message"`
	Action                *string `json:"action" name:"action" location:"elements"`
	CacheParameterGroupID *string `json:"cache_parameter_group_id" name:"cache_parameter_group_id" location:"elements"`
	RetCode               *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/delete_cache_nodes.html
func (s *CacheService) DeleteCacheNodes(i *DeleteCacheNodesInput) (*DeleteCacheNodesOutput, error) {
	if i == nil {
		i = &DeleteCacheNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteCacheNodes",
		RequestMethod: "GET",
	}

	x := &DeleteCacheNodesOutput{}
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

type DeleteCacheNodesInput struct {
	Cache      *string   `json:"cache" name:"cache" location:"params"`             // Required
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"params"` // Required
}

func (v *DeleteCacheNodesInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "DeleteCacheNodesInput",
		}
	}

	if len(v.CacheNodes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "CacheNodes",
			ParentName:    "DeleteCacheNodesInput",
		}
	}

	return nil
}

type DeleteCacheNodesOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"elements"`
	JobID      *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/delete_cache_parameter_groups.html
func (s *CacheService) DeleteCacheParameterGroups(i *DeleteCacheParameterGroupsInput) (*DeleteCacheParameterGroupsOutput, error) {
	if i == nil {
		i = &DeleteCacheParameterGroupsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteCacheParameterGroups",
		RequestMethod: "GET",
	}

	x := &DeleteCacheParameterGroupsOutput{}
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

type DeleteCacheParameterGroupsInput struct {
	CacheParameterGroups []*string `json:"cache_parameter_groups" name:"cache_parameter_groups" location:"params"` // Required
}

func (v *DeleteCacheParameterGroupsInput) Validate() error {

	if len(v.CacheParameterGroups) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroups",
			ParentName:    "DeleteCacheParameterGroupsInput",
		}
	}

	return nil
}

type DeleteCacheParameterGroupsOutput struct {
	Message         *string   `json:"message" name:"message"`
	Action          *string   `json:"action" name:"action" location:"elements"`
	ParameterGroups []*string `json:"parameter_groups" name:"parameter_groups" location:"elements"`
	RetCode         *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/delete_caches.html
func (s *CacheService) DeleteCaches(i *DeleteCachesInput) (*DeleteCachesOutput, error) {
	if i == nil {
		i = &DeleteCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteCaches",
		RequestMethod: "GET",
	}

	x := &DeleteCachesOutput{}
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

type DeleteCachesInput struct {
	Caches []*string `json:"caches" name:"caches" location:"params"` // Required
}

func (v *DeleteCachesInput) Validate() error {

	if len(v.Caches) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Caches",
			ParentName:    "DeleteCachesInput",
		}
	}

	return nil
}

type DeleteCachesOutput struct {
	Message  *string   `json:"message" name:"message"`
	Action   *string   `json:"action" name:"action" location:"elements"`
	CacheIDs []*string `json:"cache_ids" name:"cache_ids" location:"elements"`
	JobID    *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/describe_cache_nodes.html
func (s *CacheService) DescribeCacheNodes(i *DescribeCacheNodesInput) (*DescribeCacheNodesOutput, error) {
	if i == nil {
		i = &DescribeCacheNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeCacheNodes",
		RequestMethod: "GET",
	}

	x := &DescribeCacheNodesOutput{}
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

type DescribeCacheNodesInput struct {
	Cache      *string   `json:"cache" name:"cache" location:"params"`
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"params"`
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeCacheNodesInput) Validate() error {

	return nil
}

type DescribeCacheNodesOutput struct {
	Message      *string      `json:"message" name:"message"`
	Action       *string      `json:"action" name:"action" location:"elements"`
	CacheNodeSet []*CacheNode `json:"cache_node_set" name:"cache_node_set" location:"elements"`
	RetCode      *int         `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int         `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/describe_cache_parameter_groups.html
func (s *CacheService) DescribeCacheParameterGroups(i *DescribeCacheParameterGroupsInput) (*DescribeCacheParameterGroupsOutput, error) {
	if i == nil {
		i = &DescribeCacheParameterGroupsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeCacheParameterGroups",
		RequestMethod: "GET",
	}

	x := &DescribeCacheParameterGroupsOutput{}
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

type DescribeCacheParameterGroupsInput struct {
	CacheParameterGroups []*string `json:"cache_parameter_groups" name:"cache_parameter_groups" location:"params"`
	CacheType            *string   `json:"cache_type" name:"cache_type" location:"params"`
	Limit                *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset               *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord           *string   `json:"search_word" name:"search_word" location:"params"`
	Verbose              *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeCacheParameterGroupsInput) Validate() error {

	return nil
}

type DescribeCacheParameterGroupsOutput struct {
	Message                *string                `json:"message" name:"message"`
	Action                 *string                `json:"action" name:"action" location:"elements"`
	CacheParameterGroupSet []*CacheParameterGroup `json:"cache_parameter_group_set" name:"cache_parameter_group_set" location:"elements"`
	RetCode                *int                   `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount             *int                   `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/describe_cache_parameters.html
func (s *CacheService) DescribeCacheParameters(i *DescribeCacheParametersInput) (*DescribeCacheParametersOutput, error) {
	if i == nil {
		i = &DescribeCacheParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeCacheParameters",
		RequestMethod: "GET",
	}

	x := &DescribeCacheParametersOutput{}
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

type DescribeCacheParametersInput struct {
	CacheParameterGroup *string `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"` // Required
	Verbose             *int    `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeCacheParametersInput) Validate() error {

	if v.CacheParameterGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroup",
			ParentName:    "DescribeCacheParametersInput",
		}
	}

	return nil
}

type DescribeCacheParametersOutput struct {
	Message           *string           `json:"message" name:"message"`
	Action            *string           `json:"action" name:"action" location:"elements"`
	CacheParameterSet []*CacheParameter `json:"cache_parameter_set" name:"cache_parameter_set" location:"elements"`
	RetCode           *int              `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount        *int              `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/describe_caches.html
func (s *CacheService) DescribeCaches(i *DescribeCachesInput) (*DescribeCachesOutput, error) {
	if i == nil {
		i = &DescribeCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeCaches",
		RequestMethod: "GET",
	}

	x := &DescribeCachesOutput{}
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

type DescribeCachesInput struct {
	CacheType  []*string `json:"cache_type" name:"cache_type" location:"params"`
	Caches     []*string `json:"caches" name:"caches" location:"params"`
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeCachesInput) Validate() error {

	return nil
}

type DescribeCachesOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	CacheSet   []*Cache `json:"cache_set" name:"cache_set" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int     `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_cache_monitor.html
func (s *CacheService) GetCacheMonitor(i *GetCacheMonitorInput) (*GetCacheMonitorOutput, error) {
	if i == nil {
		i = &GetCacheMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetCacheMonitor",
		RequestMethod: "GET",
	}

	x := &GetCacheMonitorOutput{}
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

type GetCacheMonitorInput struct {
	EndTime   *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"`     // Required
	Meters    []*string  `json:"meters" name:"meters" location:"params"`                           // Required
	Resource  *string    `json:"resource" name:"resource" location:"params"`                       // Required
	StartTime *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetCacheMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetCacheMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetCacheMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetCacheMonitorInput",
		}
	}

	if v.Step != nil {
		stepValidValues := []string{"5m", "15m", "2h", "1d"}
		stepParameterValue := fmt.Sprint(*v.Step)

		stepIsValid := false
		for _, value := range stepValidValues {
			if value == stepParameterValue {
				stepIsValid = true
			}
		}

		if !stepIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Step",
				ParameterValue: stepParameterValue,
				AllowedValues:  stepValidValues,
			}
		}
	}

	return nil
}

type GetCacheMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/modify_cache_attributes.html
func (s *CacheService) ModifyCacheAttributes(i *ModifyCacheAttributesInput) (*ModifyCacheAttributesOutput, error) {
	if i == nil {
		i = &ModifyCacheAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyCacheAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyCacheAttributesOutput{}
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

type ModifyCacheAttributesInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" default:"99" location:"params"`
	Cache          *string `json:"cache" name:"cache" location:"params"` // Required
	CacheName      *string `json:"cache_name" name:"cache_name" location:"params"`
	Description    *string `json:"description" name:"description" location:"params"`
}

func (v *ModifyCacheAttributesInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "ModifyCacheAttributesInput",
		}
	}

	return nil
}

type ModifyCacheAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/modify_cache_node_attributes.html
func (s *CacheService) ModifyCacheNodeAttributes(i *ModifyCacheNodeAttributesInput) (*ModifyCacheNodeAttributesOutput, error) {
	if i == nil {
		i = &ModifyCacheNodeAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyCacheNodeAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyCacheNodeAttributesOutput{}
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

type ModifyCacheNodeAttributesInput struct {
	CacheNode     *string `json:"cache_node" name:"cache_node" location:"params"` // Required
	CacheNodeName *string `json:"cache_node_name" name:"cache_node_name" location:"params"`
}

func (v *ModifyCacheNodeAttributesInput) Validate() error {

	if v.CacheNode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheNode",
			ParentName:    "ModifyCacheNodeAttributesInput",
		}
	}

	return nil
}

type ModifyCacheNodeAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/modify_cache_parameter_group_attributes.html
func (s *CacheService) ModifyCacheParameterGroupAttributes(i *ModifyCacheParameterGroupAttributesInput) (*ModifyCacheParameterGroupAttributesOutput, error) {
	if i == nil {
		i = &ModifyCacheParameterGroupAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyCacheParameterGroupAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyCacheParameterGroupAttributesOutput{}
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

type ModifyCacheParameterGroupAttributesInput struct {
	CacheParameterGroup     *string `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"` // Required
	CacheParameterGroupName *string `json:"cache_parameter_group_name" name:"cache_parameter_group_name" location:"params"`
	Description             *string `json:"description" name:"description" location:"params"`
}

func (v *ModifyCacheParameterGroupAttributesInput) Validate() error {

	if v.CacheParameterGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroup",
			ParentName:    "ModifyCacheParameterGroupAttributesInput",
		}
	}

	return nil
}

type ModifyCacheParameterGroupAttributesOutput struct {
	Message               *string `json:"message" name:"message"`
	Action                *string `json:"action" name:"action" location:"elements"`
	CacheParameterGroupID *string `json:"cache_parameter_group_id" name:"cache_parameter_group_id" location:"elements"`
	RetCode               *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/reset_cache_parameters.html
func (s *CacheService) ResetCacheParameters(i *ResetCacheParametersInput) (*ResetCacheParametersOutput, error) {
	if i == nil {
		i = &ResetCacheParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResetCacheParameters",
		RequestMethod: "GET",
	}

	x := &ResetCacheParametersOutput{}
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

type ResetCacheParametersInput struct {
	CacheParameterGroup *string   `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"` // Required
	CacheParameterNames []*string `json:"cache_parameter_names" name:"cache_parameter_names" location:"params"`
}

func (v *ResetCacheParametersInput) Validate() error {

	if v.CacheParameterGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroup",
			ParentName:    "ResetCacheParametersInput",
		}
	}

	return nil
}

type ResetCacheParametersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/resize_cache.html
func (s *CacheService) ResizeCaches(i *ResizeCachesInput) (*ResizeCachesOutput, error) {
	if i == nil {
		i = &ResizeCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeCaches",
		RequestMethod: "GET",
	}

	x := &ResizeCachesOutput{}
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

type ResizeCachesInput struct {
	CacheSize *int      `json:"cache_size" name:"cache_size" location:"params"` // Required
	Caches    []*string `json:"caches" name:"caches" location:"params"`         // Required
}

func (v *ResizeCachesInput) Validate() error {

	if v.CacheSize == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheSize",
			ParentName:    "ResizeCachesInput",
		}
	}

	if len(v.Caches) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Caches",
			ParentName:    "ResizeCachesInput",
		}
	}

	return nil
}

type ResizeCachesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/restart_cache_nodes.html
func (s *CacheService) RestartCacheNodes(i *RestartCacheNodesInput) (*RestartCacheNodesOutput, error) {
	if i == nil {
		i = &RestartCacheNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestartCacheNodes",
		RequestMethod: "GET",
	}

	x := &RestartCacheNodesOutput{}
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

type RestartCacheNodesInput struct {
	Cache      *string   `json:"cache" name:"cache" location:"params"`             // Required
	CacheNodes []*string `json:"cache_nodes" name:"cache_nodes" location:"params"` // Required
}

func (v *RestartCacheNodesInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "RestartCacheNodesInput",
		}
	}

	if len(v.CacheNodes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "CacheNodes",
			ParentName:    "RestartCacheNodesInput",
		}
	}

	return nil
}

type RestartCacheNodesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// RestartCaches: Only available for memcached.
// Documentation URL: https://docs.qingcloud.com/api/cache/restart_caches.html
func (s *CacheService) RestartCaches(i *RestartCachesInput) (*RestartCachesOutput, error) {
	if i == nil {
		i = &RestartCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestartCaches",
		RequestMethod: "GET",
	}

	x := &RestartCachesOutput{}
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

type RestartCachesInput struct {
	Caches []*string `json:"caches" name:"caches" location:"params"` // Required
}

func (v *RestartCachesInput) Validate() error {

	if len(v.Caches) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Caches",
			ParentName:    "RestartCachesInput",
		}
	}

	return nil
}

type RestartCachesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/start_caches.html
func (s *CacheService) StartCaches(i *StartCachesInput) (*StartCachesOutput, error) {
	if i == nil {
		i = &StartCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartCaches",
		RequestMethod: "GET",
	}

	x := &StartCachesOutput{}
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

type StartCachesInput struct {
	Caches []*string `json:"caches" name:"caches" location:"params"` // Required
}

func (v *StartCachesInput) Validate() error {

	if len(v.Caches) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Caches",
			ParentName:    "StartCachesInput",
		}
	}

	return nil
}

type StartCachesOutput struct {
	Message  *string   `json:"message" name:"message"`
	Action   *string   `json:"action" name:"action" location:"elements"`
	CacheIDs []*string `json:"cache_ids" name:"cache_ids" location:"elements"`
	JobID    *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/stop_caches.html
func (s *CacheService) StopCaches(i *StopCachesInput) (*StopCachesOutput, error) {
	if i == nil {
		i = &StopCachesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopCaches",
		RequestMethod: "GET",
	}

	x := &StopCachesOutput{}
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

type StopCachesInput struct {
	Caches []*string `json:"caches" name:"caches" location:"params"` // Required
}

func (v *StopCachesInput) Validate() error {

	if len(v.Caches) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Caches",
			ParentName:    "StopCachesInput",
		}
	}

	return nil
}

type StopCachesOutput struct {
	Message  *string   `json:"message" name:"message"`
	Action   *string   `json:"action" name:"action" location:"elements"`
	CacheIDs []*string `json:"cache_ids" name:"cache_ids" location:"elements"`
	JobID    *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/update_cache.html
func (s *CacheService) UpdateCache(i *UpdateCacheInput) (*UpdateCacheOutput, error) {
	if i == nil {
		i = &UpdateCacheInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateCache",
		RequestMethod: "GET",
	}

	x := &UpdateCacheOutput{}
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

type UpdateCacheInput struct {
	Cache      *string           `json:"cache" name:"cache" location:"params"` // Required
	PrivateIPs []*CachePrivateIP `json:"private_ips" name:"private_ips" location:"params"`
}

func (v *UpdateCacheInput) Validate() error {

	if v.Cache == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Cache",
			ParentName:    "UpdateCacheInput",
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type UpdateCacheOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/cache/update_cache_parameters.html
func (s *CacheService) UpdateCacheParameters(i *UpdateCacheParametersInput) (*UpdateCacheParametersOutput, error) {
	if i == nil {
		i = &UpdateCacheParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateCacheParameters",
		RequestMethod: "GET",
	}

	x := &UpdateCacheParametersOutput{}
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

type UpdateCacheParametersInput struct {
	CacheParameterGroup *string         `json:"cache_parameter_group" name:"cache_parameter_group" location:"params"` // Required
	Parameters          *CacheParameter `json:"parameters" name:"parameters" location:"params"`                       // Required
}

func (v *UpdateCacheParametersInput) Validate() error {

	if v.CacheParameterGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterGroup",
			ParentName:    "UpdateCacheParametersInput",
		}
	}

	if v.Parameters != nil {
		if err := v.Parameters.Validate(); err != nil {
			return err
		}
	}

	if v.Parameters == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Parameters",
			ParentName:    "UpdateCacheParametersInput",
		}
	}

	return nil
}

type UpdateCacheParametersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
