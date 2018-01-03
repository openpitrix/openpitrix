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

	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

type Cache struct {
	AutoBackupTime *int `json:"auto_backup_time" name:"auto_backup_time"`
	// CacheClass's available values: 0, 1
	CacheClass            *int    `json:"cache_class" name:"cache_class"`
	CacheID               *string `json:"cache_id" name:"cache_id"`
	CacheName             *string `json:"cache_name" name:"cache_name"`
	CacheParameterGroupID *string `json:"cache_parameter_group_id" name:"cache_parameter_group_id"`
	CachePort             *int    `json:"cache_port" name:"cache_port"`
	CacheSize             *int    `json:"cache_size" name:"cache_size"`
	// CacheType's available values: Redis2.8.17, Memcached1.4.13
	CacheType    *string    `json:"cache_type" name:"cache_type"`
	CacheVersion *string    `json:"cache_version" name:"cache_version"`
	CreateTime   *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description  *string    `json:"description" name:"description"`
	// IsApplied's available values: 0, 1
	IsApplied       *int         `json:"is_applied" name:"is_applied"`
	MasterCount     *int         `json:"master_count" name:"master_count"`
	MaxMemory       *int         `json:"max_memory" name:"max_memory"`
	NodeCount       *int         `json:"node_count" name:"node_count"`
	Nodes           []*CacheNode `json:"nodes" name:"nodes"`
	ReplicateCount  *int         `json:"replicate_count" name:"replicate_count"`
	SecurityGroupID *string      `json:"security_group_id" name:"security_group_id"`
	// Status's available values: pending, active, stopped, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, starting, stopping, updating, suspending, resuming, deleting
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VxNet            *VxNet  `json:"vxnet" name:"vxnet"`
}

func (v *Cache) Validate() error {

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

	if v.CacheType != nil {
		cacheTypeValidValues := []string{"Redis2.8.17", "Memcached1.4.13"}
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

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	if len(v.Nodes) > 0 {
		for _, property := range v.Nodes {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "stopped", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "starting", "stopping", "updating", "suspending", "resuming", "deleting"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.VxNet != nil {
		if err := v.VxNet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type CacheNode struct {
	AlarmStatus   *string `json:"alarm_status" name:"alarm_status"`
	CacheID       *string `json:"cache_id" name:"cache_id"`
	CacheNodeID   *string `json:"cache_node_id" name:"cache_node_id"`
	CacheNodeName *string `json:"cache_node_name" name:"cache_node_name"`
	// CacheRole's available values: master, slave
	CacheRole  *string    `json:"cache_role" name:"cache_role"`
	CacheType  *string    `json:"cache_type" name:"cache_type"`
	CreateTime *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	PrivateIP  *string    `json:"private_ip" name:"private_ip"`
	Slaveof    *string    `json:"slaveof" name:"slaveof"`
	// Status's available values: pending, active, down, suspended
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	// TransitionStatus's available values: creating, starting, stopping, updating, suspending, resuming, deleting
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
}

func (v *CacheNode) Validate() error {

	if v.CacheRole != nil {
		cacheRoleValidValues := []string{"master", "slave"}
		cacheRoleParameterValue := fmt.Sprint(*v.CacheRole)

		cacheRoleIsValid := false
		for _, value := range cacheRoleValidValues {
			if value == cacheRoleParameterValue {
				cacheRoleIsValid = true
			}
		}

		if !cacheRoleIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CacheRole",
				ParameterValue: cacheRoleParameterValue,
				AllowedValues:  cacheRoleValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "down", "suspended"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "starting", "stopping", "updating", "suspending", "resuming", "deleting"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	return nil
}

type CacheParameter struct {
	CacheParameterName  *string `json:"cache_parameter_name" name:"cache_parameter_name"` // Required
	CacheParameterType  *string `json:"cache_parameter_type" name:"cache_parameter_type"`
	CacheParameterValue *string `json:"cache_parameter_value" name:"cache_parameter_value"` // Required
	CacheType           *string `json:"cache_type" name:"cache_type"`
	// IsReadonly's available values: 0, 1
	IsReadonly      *int    `json:"is_readonly" name:"is_readonly"`
	IsStatic        *int    `json:"is_static" name:"is_static"`
	OPTName         *string `json:"opt_name" name:"opt_name"`
	ParameterType   *string `json:"parameter_type" name:"parameter_type"`
	ResourceVersion *string `json:"resource_version" name:"resource_version"`
	ValueRange      *string `json:"value_range" name:"value_range"`
}

func (v *CacheParameter) Validate() error {

	if v.CacheParameterName == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterName",
			ParentName:    "CacheParameter",
		}
	}

	if v.CacheParameterValue == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CacheParameterValue",
			ParentName:    "CacheParameter",
		}
	}

	if v.IsReadonly != nil {
		isReadonlyValidValues := []string{"0", "1"}
		isReadonlyParameterValue := fmt.Sprint(*v.IsReadonly)

		isReadonlyIsValid := false
		for _, value := range isReadonlyValidValues {
			if value == isReadonlyParameterValue {
				isReadonlyIsValid = true
			}
		}

		if !isReadonlyIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsReadonly",
				ParameterValue: isReadonlyParameterValue,
				AllowedValues:  isReadonlyValidValues,
			}
		}
	}

	return nil
}

type CacheParameterGroup struct {
	CacheParameterGroupID   *string    `json:"cache_parameter_group_id" name:"cache_parameter_group_id"`
	CacheParameterGroupName *string    `json:"cache_parameter_group_name" name:"cache_parameter_group_name"`
	CacheType               *string    `json:"cache_type" name:"cache_type"`
	CreateTime              *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description             *string    `json:"description" name:"description"`
	// IsApplied's available values: 0, 1
	IsApplied *int        `json:"is_applied" name:"is_applied"`
	IsDefault *int        `json:"is_default" name:"is_default"`
	Resources []*Resource `json:"resources" name:"resources"`
}

func (v *CacheParameterGroup) Validate() error {

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	if len(v.Resources) > 0 {
		for _, property := range v.Resources {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type CachePrivateIP struct {
	CacheNodeID *string `json:"cache_node_id" name:"cache_node_id"`
	// CacheRole's available values: master, slave
	CacheRole  *string `json:"cache_role" name:"cache_role"`
	PrivateIPs *string `json:"private_ips" name:"private_ips"`
}

func (v *CachePrivateIP) Validate() error {

	if v.CacheRole != nil {
		cacheRoleValidValues := []string{"master", "slave"}
		cacheRoleParameterValue := fmt.Sprint(*v.CacheRole)

		cacheRoleIsValid := false
		for _, value := range cacheRoleValidValues {
			if value == cacheRoleParameterValue {
				cacheRoleIsValid = true
			}
		}

		if !cacheRoleIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CacheRole",
				ParameterValue: cacheRoleParameterValue,
				AllowedValues:  cacheRoleValidValues,
			}
		}
	}

	return nil
}

type Cluster struct {
	AdvancedActions            map[string]*string `json:"advanced_actions" name:"advanced_actions"`
	AppID                      *string            `json:"app_id" name:"app_id"`
	AppInfo                    interface{}        `json:"app_info" name:"app_info"`
	AppVersion                 *string            `json:"app_version" name:"app_version"`
	AppVersionInfo             interface{}        `json:"app_version_info" name:"app_version_info"`
	AutoBackupTime             *int               `json:"auto_backup_time" name:"auto_backup_time"`
	Backup                     map[string]*bool   `json:"backup" name:"backup"`
	BackupPolicy               *string            `json:"backup_policy" name:"backup_policy"`
	BackupService              interface{}        `json:"backup_service" name:"backup_service"`
	CfgmgmtID                  *string            `json:"cfgmgmt_id" name:"cfgmgmt_id"`
	ClusterID                  *string            `json:"cluster_id" name:"cluster_id"`
	ClusterType                *int               `json:"cluster_type" name:"cluster_type"`
	ConsoleID                  *string            `json:"console_id" name:"console_id"`
	Controller                 *string            `json:"controller" name:"controller"`
	CreateTime                 *time.Time         `json:"create_time" name:"create_time" format:"ISO 8601"`
	CustomService              interface{}        `json:"custom_service" name:"custom_service"`
	Debug                      *bool              `json:"debug" name:"debug"`
	Description                *string            `json:"description" name:"description"`
	DisplayTabs                interface{}        `json:"display_tabs" name:"display_tabs"`
	Endpoints                  interface{}        `json:"endpoints" name:"endpoints"`
	GlobalUUID                 *string            `json:"global_uuid" name:"global_uuid"`
	HealthCheckEnablement      map[string]*bool   `json:"health_check_enablement" name:"health_check_enablement"`
	IncrementalBackupSupported *bool              `json:"incremental_backup_supported" name:"incremental_backup_supported"`
	LatestSnapshotTime         *string            `json:"latest_snapshot_time" name:"latest_snapshot_time"`
	Links                      map[string]*string `json:"links" name:"links"`
	MetadataRootAccess         *int               `json:"metadata_root_access" name:"metadata_root_access"`
	Name                       *string            `json:"name" name:"name"`
	NodeCount                  *int               `json:"node_count" name:"node_count"`
	Nodes                      []*ClusterNode     `json:"nodes" name:"nodes"`
	Owner                      *string            `json:"owner" name:"owner"`
	PartnerAccess              *bool              `json:"partner_access" name:"partner_access"`
	RestoreService             interface{}        `json:"restore_service" name:"restore_service"`
	ReuseHyper                 *int               `json:"reuse_hyper" name:"reuse_hyper"`
	RoleCount                  map[string]*int    `json:"role_count" name:"role_count"`
	Roles                      []*string          `json:"roles" name:"roles"`
	RootUserID                 *string            `json:"root_user_id" name:"root_user_id"`
	SecurityGroupID            *string            `json:"security_group_id" name:"security_group_id"`
	Status                     *string            `json:"status" name:"status"`
	StatusTime                 *time.Time         `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode                    *int               `json:"sub_code" name:"sub_code"`
	TransitionStatus           *string            `json:"transition_status" name:"transition_status"`
	UpgradePolicy              []*string          `json:"upgrade_policy" name:"upgrade_policy"`
	UpgradeStatus              *string            `json:"upgrade_status" name:"upgrade_status"`
	UpgradeTime                *time.Time         `json:"upgrade_time" name:"upgrade_time" format:"ISO 8601"`
	VxNet                      *VxNet             `json:"vxnet" name:"vxnet"`
}

func (v *Cluster) Validate() error {

	if len(v.Nodes) > 0 {
		for _, property := range v.Nodes {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.VxNet != nil {
		if err := v.VxNet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ClusterNode struct {
	AdvancedActions            *string     `json:"advanced_actions" name:"advanced_actions"`
	AgentInstalled             *bool       `json:"agent_installed" name:"agent_installed"`
	AlarmStatus                *string     `json:"alarm_status" name:"alarm_status"`
	AppID                      *string     `json:"app_id" name:"app_id"`
	AppVersion                 *string     `json:"app_version" name:"app_version"`
	AutoBackup                 *int        `json:"auto_backup" name:"auto_backup"`
	BackupPolicy               *string     `json:"backup_policy" name:"backup_policy"`
	BackupService              interface{} `json:"backup_service" name:"backup_service"`
	ClusterID                  *string     `json:"cluster_id" name:"cluster_id"`
	ConsoleID                  *string     `json:"console_id" name:"console_id"`
	Controller                 *string     `json:"controller" name:"controller"`
	CPU                        *int        `json:"cpu" name:"cpu"`
	CreateTime                 *time.Time  `json:"create_time" name:"create_time" format:"ISO 8601"`
	CustomMetadataScript       interface{} `json:"custom_metadata_script" name:"custom_metadata_script"`
	CustomService              interface{} `json:"custom_service" name:"custom_service"`
	Debug                      *bool       `json:"debug" name:"debug"`
	DestroyService             interface{} `json:"destroy_service" name:"destroy_service"`
	DisplayTabs                interface{} `json:"display_tabs" name:"display_tabs"`
	EIP                        *string     `json:"eip" name:"eip"`
	Env                        *string     `json:"env" name:"env"`
	GlobalServerID             *int        `json:"global_server_id" name:"global_server_id"`
	Gpu                        *int        `json:"gpu" name:"gpu"`
	GpuClass                   *int        `json:"gpu_class" name:"gpu_class"`
	GroupID                    *int        `json:"group_id" name:"group_id"`
	HealthCheck                interface{} `json:"health_check" name:"health_check"`
	HealthStatus               *string     `json:"health_status" name:"health_status"`
	Hypervisor                 *string     `json:"hypervisor" name:"hypervisor"`
	ImageID                    *string     `json:"image_id" name:"image_id"`
	IncrementalBackupSupported *bool       `json:"incremental_backup_supported" name:"incremental_backup_supported"`
	InitService                interface{} `json:"init_service" name:"init_service"`
	InstanceID                 *string     `json:"instance_id" name:"instance_id"`
	IsBackup                   *int        `json:"is_backup" name:"is_backup"`
	Memory                     *int        `json:"memory" name:"memory"`
	Monitor                    interface{} `json:"monitor" name:"monitor"`
	Name                       *string     `json:"name" name:"name"`
	NodeID                     *string     `json:"node_id" name:"node_id"`
	Owner                      *string     `json:"owner" name:"owner"`
	Passphraseless             *string     `json:"passphraseless" name:"passphraseless"`
	PrivateIP                  *string     `json:"private_ip" name:"private_ip"`
	Repl                       *string     `json:"repl" name:"repl"`
	ResourceClass              *int        `json:"resource_class" name:"resource_class"`
	RestartService             interface{} `json:"restart_service" name:"restart_service"`
	RestoreService             interface{} `json:"restore_service" name:"restore_service"`
	Role                       *string     `json:"role" name:"role"`
	RootUserID                 *string     `json:"root_user_id" name:"root_user_id"`
	ScaleInService             interface{} `json:"scale_in_service" name:"scale_in_service"`
	ScaleOutService            interface{} `json:"scale_out_service" name:"scale_out_service"`
	SecurityGroup              *string     `json:"security_group" name:"security_group"`
	ServerID                   *int        `json:"server_id" name:"server_id"`
	ServerIDUpperBound         *int        `json:"server_id_upper_bound" name:"server_id_upper_bound"`
	SingleNodeRepl             *string     `json:"single_node_repl" name:"single_node_repl"`
	StartService               interface{} `json:"start_service" name:"start_service"`
	Status                     *string     `json:"status" name:"status"`
	StatusTime                 *time.Time  `json:"status_time" name:"status_time" format:"ISO 8601"`
	StopService                interface{} `json:"stop_service" name:"stop_service"`
	StorageSize                *int        `json:"storage_size" name:"storage_size"`
	TransitionStatus           *string     `json:"transition_status" name:"transition_status"`
	UserAccess                 *int        `json:"user_access" name:"user_access"`
	VerticalScalingPolicy      *string     `json:"vertical_scaling_policy" name:"vertical_scaling_policy"`
	VolumeIDs                  *string     `json:"volume_ids" name:"volume_ids"`
	VolumeType                 *int        `json:"volume_type" name:"volume_type"`
	VxNetID                    *string     `json:"vxnet_id" name:"vxnet_id"`
}

func (v *ClusterNode) Validate() error {

	return nil
}

type Data struct {
	Data  *string `json:"data" name:"data"`
	EIPID *string `json:"eip_id" name:"eip_id"`
}

func (v *Data) Validate() error {

	return nil
}

type DHCPOption struct {
	RouterStaticID *string `json:"router_static_id" name:"router_static_id"`
	Val2           *string `json:"val2" name:"val2"`
}

func (v *DHCPOption) Validate() error {

	return nil
}

type DNSAlias struct {
	CreateTime   *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description  *string    `json:"description" name:"description"`
	DNSAliasID   *string    `json:"dns_alias_id" name:"dns_alias_id"`
	DNSAliasName *string    `json:"dns_alias_name" name:"dns_alias_name"`
	DomainName   *string    `json:"domain_name" name:"domain_name"`
	ResourceID   *string    `json:"resource_id" name:"resource_id"`
	Status       *string    `json:"status" name:"status"`
}

func (v *DNSAlias) Validate() error {

	return nil
}

type EIP struct {
	AlarmStatus   *string `json:"alarm_status" name:"alarm_status"`
	AssociateMode *int    `json:"associate_mode" name:"associate_mode"`
	Bandwidth     *int    `json:"bandwidth" name:"bandwidth"`
	// BillingMode's available values: bandwidth, traffic
	BillingMode *string      `json:"billing_mode" name:"billing_mode"`
	CreateTime  *time.Time   `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string      `json:"description" name:"description"`
	EIPAddr     *string      `json:"eip_addr" name:"eip_addr"`
	EIPGroup    *EIPGroup    `json:"eip_group" name:"eip_group"`
	EIPID       *string      `json:"eip_id" name:"eip_id"`
	EIPName     *string      `json:"eip_name" name:"eip_name"`
	ICPCodes    *string      `json:"icp_codes" name:"icp_codes"`
	NeedICP     *int         `json:"need_icp" name:"need_icp"`
	Resource    *EIPResource `json:"resource" name:"resource"`
	// Status's available values: pending, available, associated, suspended, released, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: associating, dissociating, suspending, resuming, releasing
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
}

func (v *EIP) Validate() error {

	if v.BillingMode != nil {
		billingModeValidValues := []string{"bandwidth", "traffic"}
		billingModeParameterValue := fmt.Sprint(*v.BillingMode)

		billingModeIsValid := false
		for _, value := range billingModeValidValues {
			if value == billingModeParameterValue {
				billingModeIsValid = true
			}
		}

		if !billingModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "BillingMode",
				ParameterValue: billingModeParameterValue,
				AllowedValues:  billingModeValidValues,
			}
		}
	}

	if v.EIPGroup != nil {
		if err := v.EIPGroup.Validate(); err != nil {
			return err
		}
	}

	if v.Resource != nil {
		if err := v.Resource.Validate(); err != nil {
			return err
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "available", "associated", "suspended", "released", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"associating", "dissociating", "suspending", "resuming", "releasing"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	return nil
}

type EIPGroup struct {
	EIPGroupID   *string `json:"eip_group_id" name:"eip_group_id"`
	EIPGroupName *string `json:"eip_group_name" name:"eip_group_name"`
}

func (v *EIPGroup) Validate() error {

	return nil
}

type EIPResource struct {
	ResourceID   *string `json:"resource_id" name:"resource_id"`
	ResourceName *string `json:"resource_name" name:"resource_name"`
	ResourceType *string `json:"resource_type" name:"resource_type"`
}

func (v *EIPResource) Validate() error {

	return nil
}

type Extra struct {
	BlockBus   *string `json:"block_bus" name:"block_bus"`
	BootDev    *string `json:"boot_dev" name:"boot_dev"`
	CPUMax     *int    `json:"cpu_max" name:"cpu_max"`
	CPUModel   *string `json:"cpu_model" name:"cpu_model"`
	Features   *int    `json:"features" name:"features"`
	Hypervisor *string `json:"hypervisor" name:"hypervisor"`
	MemMax     *int    `json:"mem_max" name:"mem_max"`
	NICMqueue  *int    `json:"nic_mqueue" name:"nic_mqueue"`
	NoLimit    *int    `json:"no_limit" name:"no_limit"`
	NoRestrict *int    `json:"no_restrict" name:"no_restrict"`
	OSDiskSize *int    `json:"os_disk_size" name:"os_disk_size"`
	USB        *int    `json:"usb" name:"usb"`
}

func (v *Extra) Validate() error {

	return nil
}

type File struct {
	File       *string `json:"file" name:"file"`
	LastModify *string `json:"last_modify" name:"last_modify"`
	Size       *int    `json:"size" name:"size"`
}

func (v *File) Validate() error {

	return nil
}

type Image struct {
	AppBillingID  *string    `json:"app_billing_id" name:"app_billing_id"`
	Architecture  *string    `json:"architecture" name:"architecture"`
	BillingID     *string    `json:"billing_id" name:"billing_id"`
	CreateTime    *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	DefaultPasswd *string    `json:"default_passwd" name:"default_passwd"`
	DefaultUser   *string    `json:"default_user" name:"default_user"`
	Description   *string    `json:"description" name:"description"`
	FResetpwd     *int       `json:"f_resetpwd" name:"f_resetpwd"`
	Feature       *int       `json:"feature" name:"feature"`
	Features      *int       `json:"features" name:"features"`
	Hypervisor    *string    `json:"hypervisor" name:"hypervisor"`
	ImageID       *string    `json:"image_id" name:"image_id"`
	ImageName     *string    `json:"image_name" name:"image_name"`
	InstanceIDs   []*string  `json:"instance_ids" name:"instance_ids"`
	OSFamily      *string    `json:"os_family" name:"os_family"`
	Owner         *string    `json:"owner" name:"owner"`
	// Platform's available values: linux, windows
	Platform *string `json:"platform" name:"platform"`
	// ProcessorType's available values: 64bit, 32bit
	ProcessorType *string `json:"processor_type" name:"processor_type"`
	// Provider's available values: system, self
	Provider        *string `json:"provider" name:"provider"`
	RecommendedType *string `json:"recommended_type" name:"recommended_type"`
	RootID          *string `json:"root_id" name:"root_id"`
	Size            *int    `json:"size" name:"size"`
	// Status's available values: pending, available, deprecated, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	// TransitionStatus's available values: creating, suspending, resuming, deleting, recovering
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	UIType           *string `json:"ui_type" name:"ui_type"`
	// Visibility's available values: public, private
	Visibility *string `json:"visibility" name:"visibility"`
}

func (v *Image) Validate() error {

	if v.Platform != nil {
		platformValidValues := []string{"linux", "windows"}
		platformParameterValue := fmt.Sprint(*v.Platform)

		platformIsValid := false
		for _, value := range platformValidValues {
			if value == platformParameterValue {
				platformIsValid = true
			}
		}

		if !platformIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Platform",
				ParameterValue: platformParameterValue,
				AllowedValues:  platformValidValues,
			}
		}
	}

	if v.ProcessorType != nil {
		processorTypeValidValues := []string{"64bit", "32bit"}
		processorTypeParameterValue := fmt.Sprint(*v.ProcessorType)

		processorTypeIsValid := false
		for _, value := range processorTypeValidValues {
			if value == processorTypeParameterValue {
				processorTypeIsValid = true
			}
		}

		if !processorTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "ProcessorType",
				ParameterValue: processorTypeParameterValue,
				AllowedValues:  processorTypeValidValues,
			}
		}
	}

	if v.Provider != nil {
		providerValidValues := []string{"system", "self"}
		providerParameterValue := fmt.Sprint(*v.Provider)

		providerIsValid := false
		for _, value := range providerValidValues {
			if value == providerParameterValue {
				providerIsValid = true
			}
		}

		if !providerIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Provider",
				ParameterValue: providerParameterValue,
				AllowedValues:  providerValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "available", "deprecated", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "suspending", "resuming", "deleting", "recovering"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.Visibility != nil {
		visibilityValidValues := []string{"public", "private"}
		visibilityParameterValue := fmt.Sprint(*v.Visibility)

		visibilityIsValid := false
		for _, value := range visibilityValidValues {
			if value == visibilityParameterValue {
				visibilityIsValid = true
			}
		}

		if !visibilityIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Visibility",
				ParameterValue: visibilityParameterValue,
				AllowedValues:  visibilityValidValues,
			}
		}
	}

	return nil
}

type ImageUser struct {
	CreateTime *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	ImageID    *string    `json:"image_id" name:"image_id"`
	User       *User      `json:"user" name:"user"`
}

func (v *ImageUser) Validate() error {

	if v.User != nil {
		if err := v.User.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Instance struct {
	AlarmStatus      *string        `json:"alarm_status" name:"alarm_status"`
	CPUTopology      *string        `json:"cpu_topology" name:"cpu_topology"`
	CreateTime       *time.Time     `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description      *string        `json:"description" name:"description"`
	Device           *string        `json:"device" name:"device"`
	DHCPOptions      *DHCPOption    `json:"dhcp_options" name:"dhcp_options"`
	DNSAliases       []*DNSAlias    `json:"dns_aliases" name:"dns_aliases"`
	EIP              *EIP           `json:"eip" name:"eip"`
	Extra            *Extra         `json:"extra" name:"extra"`
	GraphicsPasswd   *string        `json:"graphics_passwd" name:"graphics_passwd"`
	GraphicsProtocol *string        `json:"graphics_protocol" name:"graphics_protocol"`
	Image            *Image         `json:"image" name:"image"`
	ImageID          *string        `json:"image_id" name:"image_id"`
	InstanceClass    *int           `json:"instance_class" name:"instance_class"`
	InstanceID       *string        `json:"instance_id" name:"instance_id"`
	InstanceName     *string        `json:"instance_name" name:"instance_name"`
	InstanceType     *string        `json:"instance_type" name:"instance_type"`
	KeyPairIDs       []*string      `json:"keypair_ids" name:"keypair_ids"`
	MemoryCurrent    *int           `json:"memory_current" name:"memory_current"`
	PrivateIP        *string        `json:"private_ip" name:"private_ip"`
	SecurityGroup    *SecurityGroup `json:"security_group" name:"security_group"`
	// Status's available values: pending, running, stopped, suspended, terminated, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, starting, stopping, restarting, suspending, resuming, terminating, recovering, resetting
	TransitionStatus *string          `json:"transition_status" name:"transition_status"`
	VCPUsCurrent     *int             `json:"vcpus_current" name:"vcpus_current"`
	VolumeIDs        []*string        `json:"volume_ids" name:"volume_ids"`
	Volumes          []*Volume        `json:"volumes" name:"volumes"`
	VxNets           []*InstanceVxNet `json:"vxnets" name:"vxnets"`
}

func (v *Instance) Validate() error {

	if v.DHCPOptions != nil {
		if err := v.DHCPOptions.Validate(); err != nil {
			return err
		}
	}

	if len(v.DNSAliases) > 0 {
		for _, property := range v.DNSAliases {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.EIP != nil {
		if err := v.EIP.Validate(); err != nil {
			return err
		}
	}

	if v.Extra != nil {
		if err := v.Extra.Validate(); err != nil {
			return err
		}
	}

	if v.Image != nil {
		if err := v.Image.Validate(); err != nil {
			return err
		}
	}

	if v.SecurityGroup != nil {
		if err := v.SecurityGroup.Validate(); err != nil {
			return err
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "running", "stopped", "suspended", "terminated", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "starting", "stopping", "restarting", "suspending", "resuming", "terminating", "recovering", "resetting"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if len(v.Volumes) > 0 {
		for _, property := range v.Volumes {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.VxNets) > 0 {
		for _, property := range v.VxNets {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type InstanceType struct {
	Description      *string `json:"description" name:"description"`
	InstanceTypeID   *string `json:"instance_type_id" name:"instance_type_id"`
	InstanceTypeName *string `json:"instance_type_name" name:"instance_type_name"`
	MemoryCurrent    *int    `json:"memory_current" name:"memory_current"`
	// Status's available values: available, deprecated
	Status       *string `json:"status" name:"status"`
	VCPUsCurrent *int    `json:"vcpus_current" name:"vcpus_current"`
	ZoneID       *string `json:"zone_id" name:"zone_id"`
}

func (v *InstanceType) Validate() error {

	if v.Status != nil {
		statusValidValues := []string{"available", "deprecated"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	return nil
}

type InstanceVxNet struct {
	NICID     *string `json:"nic_id" name:"nic_id"`
	PrivateIP *string `json:"private_ip" name:"private_ip"`
	Role      *int    `json:"role" name:"role"`
	VxNetID   *string `json:"vxnet_id" name:"vxnet_id"`
	VxNetName *string `json:"vxnet_name" name:"vxnet_name"`
	// VxNetType's available values: 0, 1
	VxNetType *int `json:"vxnet_type" name:"vxnet_type"`
}

func (v *InstanceVxNet) Validate() error {

	if v.VxNetType != nil {
		vxnetTypeValidValues := []string{"0", "1"}
		vxnetTypeParameterValue := fmt.Sprint(*v.VxNetType)

		vxnetTypeIsValid := false
		for _, value := range vxnetTypeValidValues {
			if value == vxnetTypeParameterValue {
				vxnetTypeIsValid = true
			}
		}

		if !vxnetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VxNetType",
				ParameterValue: vxnetTypeParameterValue,
				AllowedValues:  vxnetTypeValidValues,
			}
		}
	}

	return nil
}

type Job struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	JobAction   *string    `json:"job_action" name:"job_action"`
	JobID       *string    `json:"job_id" name:"job_id"`
	Owner       *string    `json:"owner" name:"owner"`
	ResourceIDs *string    `json:"resource_ids" name:"resource_ids"`
	// Status's available values: pending, working, failed, successful, done with failure
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
}

func (v *Job) Validate() error {

	if v.Status != nil {
		statusValidValues := []string{"pending", "working", "failed", "successful", "done with failure"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	return nil
}

type KeyPair struct {
	Description *string `json:"description" name:"description"`
	// EncryptMethod's available values: ssh-rsa, ssh-dss
	EncryptMethod *string   `json:"encrypt_method" name:"encrypt_method"`
	InstanceIDs   []*string `json:"instance_ids" name:"instance_ids"`
	KeyPairID     *string   `json:"keypair_id" name:"keypair_id"`
	KeyPairName   *string   `json:"keypair_name" name:"keypair_name"`
	PubKey        *string   `json:"pub_key" name:"pub_key"`
	Tags          []*Tag    `json:"tags" name:"tags"`
}

func (v *KeyPair) Validate() error {

	if v.EncryptMethod != nil {
		encryptMethodValidValues := []string{"ssh-rsa", "ssh-dss"}
		encryptMethodParameterValue := fmt.Sprint(*v.EncryptMethod)

		encryptMethodIsValid := false
		for _, value := range encryptMethodValidValues {
			if value == encryptMethodParameterValue {
				encryptMethodIsValid = true
			}
		}

		if !encryptMethodIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "EncryptMethod",
				ParameterValue: encryptMethodParameterValue,
				AllowedValues:  encryptMethodValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type LoadBalancer struct {
	Cluster     []*EIP     `json:"cluster" name:"cluster"`
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	EIPs        []*EIP     `json:"eips" name:"eips"`
	// IsApplied's available values: 0, 1
	IsApplied        *int                    `json:"is_applied" name:"is_applied"`
	Listeners        []*LoadBalancerListener `json:"listeners" name:"listeners"`
	LoadBalancerID   *string                 `json:"loadbalancer_id" name:"loadbalancer_id"`
	LoadBalancerName *string                 `json:"loadbalancer_name" name:"loadbalancer_name"`
	// LoadBalancerType's available values: 0, 1, 2, 3, 4, 5
	LoadBalancerType *int      `json:"loadbalancer_type" name:"loadbalancer_type"`
	NodeCount        *int      `json:"node_count" name:"node_count"`
	PrivateIPs       []*string `json:"private_ips" name:"private_ips"`
	SecurityGroupID  *string   `json:"security_group_id" name:"security_group_id"`
	// Status's available values: pending, active, stopped, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, starting, stopping, updating, suspending, resuming, deleting
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VxNetID          *string `json:"vxnet_id" name:"vxnet_id"`
}

func (v *LoadBalancer) Validate() error {

	if len(v.Cluster) > 0 {
		for _, property := range v.Cluster {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.EIPs) > 0 {
		for _, property := range v.EIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	if len(v.Listeners) > 0 {
		for _, property := range v.Listeners {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.LoadBalancerType != nil {
		loadBalancerTypeValidValues := []string{"0", "1", "2", "3", "4", "5"}
		loadBalancerTypeParameterValue := fmt.Sprint(*v.LoadBalancerType)

		loadBalancerTypeIsValid := false
		for _, value := range loadBalancerTypeValidValues {
			if value == loadBalancerTypeParameterValue {
				loadBalancerTypeIsValid = true
			}
		}

		if !loadBalancerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "LoadBalancerType",
				ParameterValue: loadBalancerTypeParameterValue,
				AllowedValues:  loadBalancerTypeValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "stopped", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "starting", "stopping", "updating", "suspending", "resuming", "deleting"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	return nil
}

type LoadBalancerBackend struct {
	CreateTime              *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	LoadBalancerBackendID   *string    `json:"loadbalancer_backend_id" name:"loadbalancer_backend_id"`
	LoadBalancerBackendName *string    `json:"loadbalancer_backend_name" name:"loadbalancer_backend_name"`
	LoadBalancerID          *string    `json:"loadbalancer_id" name:"loadbalancer_id"`
	LoadBalancerListenerID  *string    `json:"loadbalancer_listener_id" name:"loadbalancer_listener_id"`
	LoadBalancerPolicyID    *string    `json:"loadbalancer_policy_id" name:"loadbalancer_policy_id"`
	Port                    *int       `json:"port" name:"port"`
	ResourceID              *string    `json:"resource_id" name:"resource_id"`
	Status                  *string    `json:"status" name:"status"`
	Weight                  *int       `json:"weight" name:"weight"`
}

func (v *LoadBalancerBackend) Validate() error {

	return nil
}

type LoadBalancerListener struct {
	BackendProtocol *string                `json:"backend_protocol" name:"backend_protocol"`
	Backends        []*LoadBalancerBackend `json:"backends" name:"backends"`
	// BalanceMode's available values: roundrobin, leastconn, source
	BalanceMode              *string    `json:"balance_mode" name:"balance_mode"`
	CreateTime               *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Forwardfor               *int       `json:"forwardfor" name:"forwardfor"`
	HealthyCheckMethod       *string    `json:"healthy_check_method" name:"healthy_check_method"`
	HealthyCheckOption       *string    `json:"healthy_check_option" name:"healthy_check_option" default:"10|5|2|5"`
	ListenerOption           *int       `json:"listener_option" name:"listener_option"`
	ListenerPort             *int       `json:"listener_port" name:"listener_port"`
	ListenerProtocol         *string    `json:"listener_protocol" name:"listener_protocol"`
	LoadBalancerID           *string    `json:"loadbalancer_id" name:"loadbalancer_id"`
	LoadBalancerListenerID   *string    `json:"loadbalancer_listener_id" name:"loadbalancer_listener_id"`
	LoadBalancerListenerName *string    `json:"loadbalancer_listener_name" name:"loadbalancer_listener_name"`
	ServerCertificateID      *string    `json:"server_certificate_id" name:"server_certificate_id"`
	SessionSticky            *string    `json:"session_sticky" name:"session_sticky"`
}

func (v *LoadBalancerListener) Validate() error {

	if len(v.Backends) > 0 {
		for _, property := range v.Backends {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.BalanceMode != nil {
		balanceModeValidValues := []string{"roundrobin", "leastconn", "source"}
		balanceModeParameterValue := fmt.Sprint(*v.BalanceMode)

		balanceModeIsValid := false
		for _, value := range balanceModeValidValues {
			if value == balanceModeParameterValue {
				balanceModeIsValid = true
			}
		}

		if !balanceModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "BalanceMode",
				ParameterValue: balanceModeParameterValue,
				AllowedValues:  balanceModeValidValues,
			}
		}
	}

	return nil
}

type LoadBalancerPolicy struct {
	CreateTime *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	// IsApplied's available values: 0, 1
	IsApplied              *int      `json:"is_applied" name:"is_applied"`
	LoadBalancerIDs        []*string `json:"loadbalancer_ids" name:"loadbalancer_ids"`
	LoadBalancerPolicyID   *string   `json:"loadbalancer_policy_id" name:"loadbalancer_policy_id"`
	LoadBalancerPolicyName *string   `json:"loadbalancer_policy_name" name:"loadbalancer_policy_name"`
}

func (v *LoadBalancerPolicy) Validate() error {

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	return nil
}

type LoadBalancerPolicyRule struct {
	LoadBalancerPolicyRuleID   *string `json:"loadbalancer_policy_rule_id" name:"loadbalancer_policy_rule_id"`
	LoadBalancerPolicyRuleName *string `json:"loadbalancer_policy_rule_name" name:"loadbalancer_policy_rule_name"`
	RuleType                   *string `json:"rule_type" name:"rule_type"`
	Val                        *string `json:"val" name:"val"`
}

func (v *LoadBalancerPolicyRule) Validate() error {

	return nil
}

type Meter struct {
	Data     interface{}   `json:"data" name:"data"`
	DataSet  []interface{} `json:"data_set" name:"data_set"`
	MeterID  *string       `json:"meter_id" name:"meter_id"`
	Sequence *int          `json:"sequence" name:"sequence"`
	VxNetID  *string       `json:"vxnet_id" name:"vxnet_id"`
}

func (v *Meter) Validate() error {

	return nil
}

type Mongo struct {
	// AlarmStatus's available values: ok, alarm, insufficient
	AlarmStatus         *string    `json:"alarm_status" name:"alarm_status"`
	AutoBackupTime      *int       `json:"auto_backup_time" name:"auto_backup_time"`
	AutoMinorVerUpgrade *int       `json:"auto_minor_ver_upgrade" name:"auto_minor_ver_upgrade"`
	CreateTime          *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description         *string    `json:"description" name:"description"`
	LatestSnapshotTime  *time.Time `json:"latest_snapshot_time" name:"latest_snapshot_time" format:"ISO 8601"`
	MongoID             *string    `json:"mongo_id" name:"mongo_id"`
	MongoName           *string    `json:"mongo_name" name:"mongo_name"`
	MongoType           *int       `json:"mongo_type" name:"mongo_type"`
	MongoVersion        *string    `json:"mongo_version" name:"mongo_version"`
	// Status's available values: pending, active, stopped, deleted, suspended, ceased
	Status      *string    `json:"status" name:"status"`
	StatusTime  *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	StorageSize *int       `json:"storage_size" name:"storage_size"`
	Tags        []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, stopping, starting, deleting, resizing, suspending, vxnet-changing, snapshot-creating, instances-adding, instances-removing, pg-applying
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VxNet            *VxNet  `json:"vxnet" name:"vxnet"`
}

func (v *Mongo) Validate() error {

	if v.AlarmStatus != nil {
		alarmStatusValidValues := []string{"ok", "alarm", "insufficient"}
		alarmStatusParameterValue := fmt.Sprint(*v.AlarmStatus)

		alarmStatusIsValid := false
		for _, value := range alarmStatusValidValues {
			if value == alarmStatusParameterValue {
				alarmStatusIsValid = true
			}
		}

		if !alarmStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "AlarmStatus",
				ParameterValue: alarmStatusParameterValue,
				AllowedValues:  alarmStatusValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "stopped", "deleted", "suspended", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "stopping", "starting", "deleting", "resizing", "suspending", "vxnet-changing", "snapshot-creating", "instances-adding", "instances-removing", "pg-applying"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.VxNet != nil {
		if err := v.VxNet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type MongoNode struct {
	IP          *string `json:"ip" name:"ip"`
	MongoID     *string `json:"mongo_id" name:"mongo_id"`
	MongoNodeID *string `json:"mongo_node_id" name:"mongo_node_id"`
	Primary     *int    `json:"primary" name:"primary"`
	Status      *string `json:"status" name:"status"`
	VxNetID     *string `json:"vxnet_id" name:"vxnet_id"`
}

func (v *MongoNode) Validate() error {

	return nil
}

type MongoParameter struct {
	// IsReadonly's available values: 0, 1
	IsReadonly *int `json:"is_readonly" name:"is_readonly"`
	// IsStatic's available values: 0, 1
	IsStatic      *int    `json:"is_static" name:"is_static"`
	OPTName       *string `json:"opt_name" name:"opt_name"`
	ParameterName *string `json:"parameter_name" name:"parameter_name"`
	// ParameterType's available values: string, int, bool
	ParameterType  *string `json:"parameter_type" name:"parameter_type"`
	ParameterValue *string `json:"parameter_value" name:"parameter_value"`
	ResourceType   *string `json:"resource_type" name:"resource_type"`
}

func (v *MongoParameter) Validate() error {

	if v.IsReadonly != nil {
		isReadonlyValidValues := []string{"0", "1"}
		isReadonlyParameterValue := fmt.Sprint(*v.IsReadonly)

		isReadonlyIsValid := false
		for _, value := range isReadonlyValidValues {
			if value == isReadonlyParameterValue {
				isReadonlyIsValid = true
			}
		}

		if !isReadonlyIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsReadonly",
				ParameterValue: isReadonlyParameterValue,
				AllowedValues:  isReadonlyValidValues,
			}
		}
	}

	if v.IsStatic != nil {
		isStaticValidValues := []string{"0", "1"}
		isStaticParameterValue := fmt.Sprint(*v.IsStatic)

		isStaticIsValid := false
		for _, value := range isStaticValidValues {
			if value == isStaticParameterValue {
				isStaticIsValid = true
			}
		}

		if !isStaticIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsStatic",
				ParameterValue: isStaticParameterValue,
				AllowedValues:  isStaticValidValues,
			}
		}
	}

	if v.ParameterType != nil {
		parameterTypeValidValues := []string{"string", "int", "bool"}
		parameterTypeParameterValue := fmt.Sprint(*v.ParameterType)

		parameterTypeIsValid := false
		for _, value := range parameterTypeValidValues {
			if value == parameterTypeParameterValue {
				parameterTypeIsValid = true
			}
		}

		if !parameterTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "ParameterType",
				ParameterValue: parameterTypeParameterValue,
				AllowedValues:  parameterTypeValidValues,
			}
		}
	}

	return nil
}

type MongoPrivateIP struct {
	Priority0 *string `json:"priority0" name:"priority0"`
	Replica   *string `json:"replica" name:"replica"`
}

func (v *MongoPrivateIP) Validate() error {

	return nil
}

type NIC struct {
	CreateTime    *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	InstanceID    *string    `json:"instance_id" name:"instance_id"`
	NICID         *string    `json:"nic_id" name:"nic_id"`
	NICName       *string    `json:"nic_name" name:"nic_name"`
	Owner         *string    `json:"owner" name:"owner"`
	PrivateIP     *string    `json:"private_ip" name:"private_ip"`
	Role          *int       `json:"role" name:"role"`
	RootUserID    *string    `json:"root_user_id" name:"root_user_id"`
	SecurityGroup *string    `json:"security_group" name:"security_group"`
	Sequence      *int       `json:"sequence" name:"sequence"`
	// Status's available values: available, in-use
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	VxNetID    *string    `json:"vxnet_id" name:"vxnet_id"`
}

func (v *NIC) Validate() error {

	if v.Status != nil {
		statusValidValues := []string{"available", "in-use"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type NICIP struct {
	NICID     *string `json:"nic_id" name:"nic_id"`
	PrivateIP *string `json:"private_ip" name:"private_ip"`
}

func (v *NICIP) Validate() error {

	return nil
}

type RDB struct {
	// AlarmStatus's available values: ok, alarm, insufficient
	AlarmStatus         *string    `json:"alarm_status" name:"alarm_status"`
	AutoBackupTime      *int       `json:"auto_backup_time" name:"auto_backup_time"`
	AutoMinorVerUpgrade *int       `json:"auto_minor_ver_upgrade" name:"auto_minor_ver_upgrade"`
	CreateTime          *string    `json:"create_time" name:"create_time"`
	Description         *string    `json:"description" name:"description"`
	EngineVersion       *string    `json:"engine_version" name:"engine_version"`
	LatestSnapshotTime  *time.Time `json:"latest_snapshot_time" name:"latest_snapshot_time" format:"ISO 8601"`
	MasterIP            *string    `json:"master_ip" name:"master_ip"`
	RDBEngine           *string    `json:"rdb_engine" name:"rdb_engine"`
	RDBID               *string    `json:"rdb_id" name:"rdb_id"`
	RDBName             *string    `json:"rdb_name" name:"rdb_name"`
	RDBType             *int       `json:"rdb_type" name:"rdb_type"`
	// Status's available values: pending, active, stopped, deleted, suspended, ceased
	Status      *string `json:"status" name:"status"`
	StatusTime  *string `json:"status_time" name:"status_time"`
	StorageSize *int    `json:"storage_size" name:"storage_size"`
	Tags        []*Tag  `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, stopping, starting, deleting, backup-creating, temp-creating, configuring, switching, invalid-tackling, resizing, suspending, ceasing, instance-ceasing, vxnet-leaving, vxnet-joining
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VxNet            *VxNet  `json:"vxnet" name:"vxnet"`
}

func (v *RDB) Validate() error {

	if v.AlarmStatus != nil {
		alarmStatusValidValues := []string{"ok", "alarm", "insufficient"}
		alarmStatusParameterValue := fmt.Sprint(*v.AlarmStatus)

		alarmStatusIsValid := false
		for _, value := range alarmStatusValidValues {
			if value == alarmStatusParameterValue {
				alarmStatusIsValid = true
			}
		}

		if !alarmStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "AlarmStatus",
				ParameterValue: alarmStatusParameterValue,
				AllowedValues:  alarmStatusValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "stopped", "deleted", "suspended", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "stopping", "starting", "deleting", "backup-creating", "temp-creating", "configuring", "switching", "invalid-tackling", "resizing", "suspending", "ceasing", "instance-ceasing", "vxnet-leaving", "vxnet-joining"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.VxNet != nil {
		if err := v.VxNet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type RDBFile struct {
	BinaryLog []*File `json:"binary_log" name:"binary_log"`
	ErrorLog  []*File `json:"error_log" name:"error_log"`
	SlowLog   []*File `json:"slow_log" name:"slow_log"`
}

func (v *RDBFile) Validate() error {

	if len(v.BinaryLog) > 0 {
		for _, property := range v.BinaryLog {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.ErrorLog) > 0 {
		for _, property := range v.ErrorLog {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.SlowLog) > 0 {
		for _, property := range v.SlowLog {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type RDBParameter struct {
	Family *string `json:"family" name:"family"`
	// IsReadonly's available values: 0, 1
	IsReadonly *int `json:"is_readonly" name:"is_readonly"`
	// IsStatic's available values: 0, 1
	IsStatic    *int    `json:"is_static" name:"is_static"`
	MaxValue    *int    `json:"max_value" name:"max_value"`
	MinValue    *int    `json:"min_value" name:"min_value"`
	OPTName     *string `json:"opt_name" name:"opt_name"`
	SectionName *string `json:"section_name" name:"section_name"`
	VarName     *string `json:"var_name" name:"var_name"`
	VarType     *string `json:"var_type" name:"var_type"`
	VarValue    *string `json:"var_value" name:"var_value"`
}

func (v *RDBParameter) Validate() error {

	if v.IsReadonly != nil {
		isReadonlyValidValues := []string{"0", "1"}
		isReadonlyParameterValue := fmt.Sprint(*v.IsReadonly)

		isReadonlyIsValid := false
		for _, value := range isReadonlyValidValues {
			if value == isReadonlyParameterValue {
				isReadonlyIsValid = true
			}
		}

		if !isReadonlyIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsReadonly",
				ParameterValue: isReadonlyParameterValue,
				AllowedValues:  isReadonlyValidValues,
			}
		}
	}

	if v.IsStatic != nil {
		isStaticValidValues := []string{"0", "1"}
		isStaticParameterValue := fmt.Sprint(*v.IsStatic)

		isStaticIsValid := false
		for _, value := range isStaticValidValues {
			if value == isStaticParameterValue {
				isStaticIsValid = true
			}
		}

		if !isStaticIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsStatic",
				ParameterValue: isStaticParameterValue,
				AllowedValues:  isStaticValidValues,
			}
		}
	}

	return nil
}

type RDBParameters struct {
	BindAddress               *string `json:"bind_address" name:"bind_address"`
	BinlogFormat              *string `json:"binlog_format" name:"binlog_format"`
	CharacterSetServer        *string `json:"character_set_server" name:"character_set_server"`
	DataDir                   *string `json:"datadir" name:"datadir"`
	DefaultStorageEngine      *string `json:"default_storage_engine" name:"default_storage_engine"`
	ExpireLogsDays            *int    `json:"expire_logs_days" name:"expire_logs_days"`
	InnoDB                    *string `json:"innodb" name:"innodb"`
	InnoDBBufferPoolInstances *int    `json:"innodb_buffer_pool_instances" name:"innodb_buffer_pool_instances"`
	InnoDBBufferPoolSize      *string `json:"innodb_buffer_pool_size" name:"innodb_buffer_pool_size"`
	InnoDBFilePerTable        *int    `json:"innodb_file_per_table" name:"innodb_file_per_table"`
	InnoDBFlushLogAtTRXCommit *int    `json:"innodb_flush_log_at_trx_commit" name:"innodb_flush_log_at_trx_commit"`
	InnoDBFlushMethod         *string `json:"innodb_flush_method" name:"innodb_flush_method"`
	InnoDBIOCapacity          *int    `json:"innodb_io_capacity" name:"innodb_io_capacity"`
	InnoDBLogBufferSize       *string `json:"innodb_log_buffer_size" name:"innodb_log_buffer_size"`
	InnoDBLogFileSize         *string `json:"innodb_log_file_size" name:"innodb_log_file_size"`
	InnoDBLogFilesInGroup     *int    `json:"innodb_log_files_in_group" name:"innodb_log_files_in_group"`
	InnoDBMaxDirtyPagesPct    *int    `json:"innodb_max_dirty_pages_pct" name:"innodb_max_dirty_pages_pct"`
	InnoDBReadIOThreads       *int    `json:"innodb_read_io_threads" name:"innodb_read_io_threads"`
	InnoDBWriteIOThreads      *int    `json:"innodb_write_io_threads" name:"innodb_write_io_threads"`
	InteractiveTimeout        *int    `json:"interactive_timeout" name:"interactive_timeout"`
	KeyBufferSize             *string `json:"key_buffer_size" name:"key_buffer_size"`
	LogBinIndex               *string `json:"log-bin-index" name:"log-bin-index"`
	LogBin                    *string `json:"log_bin" name:"log_bin"`
	LogError                  *string `json:"log_error" name:"log_error"`
	LogQueriesNotUsingIndexes *string `json:"log_queries_not_using_indexes" name:"log_queries_not_using_indexes"`
	LogSlaveUpdates           *int    `json:"log_slave_updates" name:"log_slave_updates"`
	LongQueryTime             *int    `json:"long_query_time" name:"long_query_time"`
	LowerCaseTableNames       *int    `json:"lower_case_table_names" name:"lower_case_table_names"`
	MaxAllowedPacket          *string `json:"max_allowed_packet" name:"max_allowed_packet"`
	MaxConnectErrors          *int    `json:"max_connect_errors" name:"max_connect_errors"`
	MaxConnections            *int    `json:"max_connections" name:"max_connections"`
	MaxHeapTableSize          *string `json:"max_heap_table_size" name:"max_heap_table_size"`
	OpenFilesLimit            *int    `json:"open_files_limit" name:"open_files_limit"`
	Port                      *int    `json:"port" name:"port"`
	QueryCacheSize            *int    `json:"query_cache_size" name:"query_cache_size"`
	QueryCacheType            *int    `json:"query_cache_type" name:"query_cache_type"`
	RelayLog                  *string `json:"relay_log" name:"relay_log"`
	RelayLogIndex             *string `json:"relay_log_index" name:"relay_log_index"`
	SkipSlaveStart            *int    `json:"skip-slave-start" name:"skip-slave-start"`
	SkipNameResolve           *int    `json:"skip_name_resolve" name:"skip_name_resolve"`
	SlaveExecMode             *string `json:"slave_exec_mode" name:"slave_exec_mode"`
	SlaveNetTimeout           *int    `json:"slave_net_timeout" name:"slave_net_timeout"`
	SlowQueryLog              *int    `json:"slow_query_log" name:"slow_query_log"`
	SlowQueryLogFile          *string `json:"slow_query_log_file" name:"slow_query_log_file"`
	SQLMode                   *string `json:"sql_mode" name:"sql_mode"`
	SyncBinlog                *int    `json:"sync_binlog" name:"sync_binlog"`
	SyncMasterInfo            *int    `json:"sync_master_info" name:"sync_master_info"`
	SyncRelayLog              *int    `json:"sync_relay_log" name:"sync_relay_log"`
	SyncRelayLogInfo          *int    `json:"sync_relay_log_info" name:"sync_relay_log_info"`
	TableOpenCache            *int    `json:"table_open_cache" name:"table_open_cache"`
	ThreadCacheSize           *int    `json:"thread_cache_size" name:"thread_cache_size"`
	TMPTableSize              *string `json:"tmp_table_size" name:"tmp_table_size"`
	TMPDir                    *string `json:"tmpdir" name:"tmpdir"`
	User                      *string `json:"user" name:"user"`
	WaitTimeout               *int    `json:"wait_timeout" name:"wait_timeout"`
}

func (v *RDBParameters) Validate() error {

	return nil
}

type RDBPrivateIP struct {
	Master   *string `json:"master" name:"master"`
	TopSlave *string `json:"topslave" name:"topslave"`
}

func (v *RDBPrivateIP) Validate() error {

	return nil
}

type Resource struct {
	ResourceID   *string `json:"resource_id" name:"resource_id"`
	ResourceName *string `json:"resource_name" name:"resource_name"`
	ResourceType *string `json:"resource_type" name:"resource_type"`
}

func (v *Resource) Validate() error {

	return nil
}

type ResourceTagPair struct {
	ResourceID   *string    `json:"resource_id" name:"resource_id"`
	ResourceType *string    `json:"resource_type" name:"resource_type"`
	Status       *string    `json:"status" name:"status"`
	StatusTime   *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	TagID        *string    `json:"tag_id" name:"tag_id"`
}

func (v *ResourceTagPair) Validate() error {

	return nil
}

type ResourceTypeCount struct {
	Count        *int    `json:"count" name:"count"`
	ResourceType *string `json:"resource_type" name:"resource_type"`
}

func (v *ResourceTypeCount) Validate() error {

	return nil
}

type Router struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	DYNIPEnd    *string    `json:"dyn_ip_end" name:"dyn_ip_end"`
	DYNIPStart  *string    `json:"dyn_ip_start" name:"dyn_ip_start"`
	EIP         *EIP       `json:"eip" name:"eip"`
	IPNetwork   *string    `json:"ip_network" name:"ip_network"`
	// IsApplied's available values: 0, 1
	IsApplied  *int    `json:"is_applied" name:"is_applied"`
	ManagerIP  *string `json:"manager_ip" name:"manager_ip"`
	Mode       *int    `json:"mode" name:"mode"`
	PrivateIP  *string `json:"private_ip" name:"private_ip"`
	RouterID   *string `json:"router_id" name:"router_id"`
	RouterName *string `json:"router_name" name:"router_name"`
	// RouterType's available values: 1
	RouterType      *int    `json:"router_type" name:"router_type"`
	SecurityGroupID *string `json:"security_group_id" name:"security_group_id"`
	// Status's available values: pending, active, poweroffed, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, updating, suspending, resuming, poweroffing, poweroning, deleting
	TransitionStatus *string  `json:"transition_status" name:"transition_status"`
	VxNets           []*VxNet `json:"vxnets" name:"vxnets"`
}

func (v *Router) Validate() error {

	if v.EIP != nil {
		if err := v.EIP.Validate(); err != nil {
			return err
		}
	}

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	if v.RouterType != nil {
		routerTypeValidValues := []string{"1"}
		routerTypeParameterValue := fmt.Sprint(*v.RouterType)

		routerTypeIsValid := false
		for _, value := range routerTypeValidValues {
			if value == routerTypeParameterValue {
				routerTypeIsValid = true
			}
		}

		if !routerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RouterType",
				ParameterValue: routerTypeParameterValue,
				AllowedValues:  routerTypeValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "poweroffed", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "updating", "suspending", "resuming", "poweroffing", "poweroning", "deleting"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if len(v.VxNets) > 0 {
		for _, property := range v.VxNets {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type RouterStatic struct {
	CreateTime       *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	RouterID         *string    `json:"router_id" name:"router_id"`
	RouterStaticID   *string    `json:"router_static_id" name:"router_static_id"`
	RouterStaticName *string    `json:"router_static_name" name:"router_static_name"`
	// StaticType's available values: 1, 2, 3, 4, 5, 6, 7, 8
	StaticType *int    `json:"static_type" name:"static_type"`
	Val1       *string `json:"val1" name:"val1"`
	Val2       *string `json:"val2" name:"val2"`
	Val3       *string `json:"val3" name:"val3"`
	Val4       *string `json:"val4" name:"val4"`
	Val5       *string `json:"val5" name:"val5"`
	Val6       *string `json:"val6" name:"val6"`
	VxNetID    *string `json:"vxnet_id" name:"vxnet_id"`
}

func (v *RouterStatic) Validate() error {

	if v.StaticType != nil {
		staticTypeValidValues := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
		staticTypeParameterValue := fmt.Sprint(*v.StaticType)

		staticTypeIsValid := false
		for _, value := range staticTypeValidValues {
			if value == staticTypeParameterValue {
				staticTypeIsValid = true
			}
		}

		if !staticTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "StaticType",
				ParameterValue: staticTypeParameterValue,
				AllowedValues:  staticTypeValidValues,
			}
		}
	}

	return nil
}

type RouterStaticEntry struct {
	RouterID              *string `json:"router_id" name:"router_id"`
	RouterStaticEntryID   *string `json:"router_static_entry_id" name:"router_static_entry_id"`
	RouterStaticEntryName *string `json:"router_static_entry_name" name:"router_static_entry_name"`
	Val1                  *string `json:"val1" name:"val1"`
	Val2                  *string `json:"val2" name:"val2"`
}

func (v *RouterStaticEntry) Validate() error {

	return nil
}

type RouterVxNet struct {
	CreateTime *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	DYNIPEnd   *string    `json:"dyn_ip_end" name:"dyn_ip_end"`
	DYNIPStart *string    `json:"dyn_ip_start" name:"dyn_ip_start"`
	Features   *int       `json:"features" name:"features"`
	IPNetwork  *string    `json:"ip_network" name:"ip_network"`
	ManagerIP  *string    `json:"manager_ip" name:"manager_ip"`
	RouterID   *string    `json:"router_id" name:"router_id"`
	VxNetID    *string    `json:"vxnet_id" name:"vxnet_id"`
}

func (v *RouterVxNet) Validate() error {

	return nil
}

type S2DefaultParameters struct {
	DefaultValue *string `json:"default_value" name:"default_value"`
	Description  *string `json:"description" name:"description"`
	ParamName    *string `json:"param_name" name:"param_name"`
	ServiceType  *string `json:"service_type" name:"service_type"`
	TargetType   *string `json:"target_type" name:"target_type"`
}

func (v *S2DefaultParameters) Validate() error {

	return nil
}

type S2Server struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	// IsApplied's available values: 0, 1
	IsApplied  *int    `json:"is_applied" name:"is_applied"`
	Name       *string `json:"name" name:"name"`
	PrivateIP  *string `json:"private_ip" name:"private_ip"`
	S2ServerID *string `json:"s2_server_id" name:"s2_server_id"`
	// S2ServerType's available values: 0, 1, 2, 3
	S2ServerType *int `json:"s2_server_type" name:"s2_server_type"`
	// ServiceType's available values: vsan
	ServiceType *string `json:"service_type" name:"service_type"`
	// Status's available values: pending, active, poweroffed, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, updating, suspending, resuming, poweroffing
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VxNet            *VxNet  `json:"vxnet" name:"vxnet"`
}

func (v *S2Server) Validate() error {

	if v.IsApplied != nil {
		isAppliedValidValues := []string{"0", "1"}
		isAppliedParameterValue := fmt.Sprint(*v.IsApplied)

		isAppliedIsValid := false
		for _, value := range isAppliedValidValues {
			if value == isAppliedParameterValue {
				isAppliedIsValid = true
			}
		}

		if !isAppliedIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsApplied",
				ParameterValue: isAppliedParameterValue,
				AllowedValues:  isAppliedValidValues,
			}
		}
	}

	if v.S2ServerType != nil {
		s2ServerTypeValidValues := []string{"0", "1", "2", "3"}
		s2ServerTypeParameterValue := fmt.Sprint(*v.S2ServerType)

		s2ServerTypeIsValid := false
		for _, value := range s2ServerTypeValidValues {
			if value == s2ServerTypeParameterValue {
				s2ServerTypeIsValid = true
			}
		}

		if !s2ServerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "S2ServerType",
				ParameterValue: s2ServerTypeParameterValue,
				AllowedValues:  s2ServerTypeValidValues,
			}
		}
	}

	if v.ServiceType != nil {
		serviceTypeValidValues := []string{"vsan"}
		serviceTypeParameterValue := fmt.Sprint(*v.ServiceType)

		serviceTypeIsValid := false
		for _, value := range serviceTypeValidValues {
			if value == serviceTypeParameterValue {
				serviceTypeIsValid = true
			}
		}

		if !serviceTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "ServiceType",
				ParameterValue: serviceTypeParameterValue,
				AllowedValues:  serviceTypeValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "active", "poweroffed", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "updating", "suspending", "resuming", "poweroffing"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.VxNet != nil {
		if err := v.VxNet.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type S2SharedTarget struct {
	CreateTime       *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description      *string    `json:"description" name:"description"`
	ExportName       *string    `json:"export_name" name:"export_name"`
	S2ServerID       *string    `json:"s2_server_id" name:"s2_server_id"`
	S2SharedTargetID *string    `json:"s2_shared_target_id" name:"s2_shared_target_id"`
	StatusTime       *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	// TargetType's available values: ISCSI, NFS
	TargetType *string `json:"target_type" name:"target_type"`
}

func (v *S2SharedTarget) Validate() error {

	if v.TargetType != nil {
		targetTypeValidValues := []string{"ISCSI", "NFS"}
		targetTypeParameterValue := fmt.Sprint(*v.TargetType)

		targetTypeIsValid := false
		for _, value := range targetTypeValidValues {
			if value == targetTypeParameterValue {
				targetTypeIsValid = true
			}
		}

		if !targetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TargetType",
				ParameterValue: targetTypeParameterValue,
				AllowedValues:  targetTypeValidValues,
			}
		}
	}

	return nil
}

type SecurityGroup struct {
	CreateTime        *time.Time  `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description       *string     `json:"description" name:"description"`
	IsApplied         *int        `json:"is_applied" name:"is_applied"`
	IsDefault         *int        `json:"is_default" name:"is_default"`
	Resources         []*Resource `json:"resources" name:"resources"`
	SecurityGroupID   *string     `json:"security_group_id" name:"security_group_id"`
	SecurityGroupName *string     `json:"security_group_name" name:"security_group_name"`
	Tags              []*Tag      `json:"tags" name:"tags"`
}

func (v *SecurityGroup) Validate() error {

	if len(v.Resources) > 0 {
		for _, property := range v.Resources {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type SecurityGroupIPSet struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	// IPSetType's available values: 0, 1
	IPSetType              *int    `json:"ipset_type" name:"ipset_type"`
	SecurityGroupIPSetID   *string `json:"security_group_ipset_id" name:"security_group_ipset_id"`
	SecurityGroupIPSetName *string `json:"security_group_ipset_name" name:"security_group_ipset_name"`
	Val                    *string `json:"val" name:"val"`
}

func (v *SecurityGroupIPSet) Validate() error {

	if v.IPSetType != nil {
		ipSetTypeValidValues := []string{"0", "1"}
		ipSetTypeParameterValue := fmt.Sprint(*v.IPSetType)

		ipSetTypeIsValid := false
		for _, value := range ipSetTypeValidValues {
			if value == ipSetTypeParameterValue {
				ipSetTypeIsValid = true
			}
		}

		if !ipSetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IPSetType",
				ParameterValue: ipSetTypeParameterValue,
				AllowedValues:  ipSetTypeValidValues,
			}
		}
	}

	return nil
}

type SecurityGroupRule struct {
	// Action's available values: accept, drop
	Action *string `json:"action" name:"action"`
	// Direction's available values: 0, 1
	Direction             *int    `json:"direction" name:"direction"`
	Priority              *int    `json:"priority" name:"priority"`
	Protocol              *string `json:"protocol" name:"protocol"`
	SecurityGroupID       *string `json:"security_group_id" name:"security_group_id"`
	SecurityGroupRuleID   *string `json:"security_group_rule_id" name:"security_group_rule_id"`
	SecurityGroupRuleName *string `json:"security_group_rule_name" name:"security_group_rule_name"`
	Val1                  *string `json:"val1" name:"val1"`
	Val2                  *string `json:"val2" name:"val2"`
	Val3                  *string `json:"val3" name:"val3"`
}

func (v *SecurityGroupRule) Validate() error {

	if v.Action != nil {
		actionValidValues := []string{"accept", "drop"}
		actionParameterValue := fmt.Sprint(*v.Action)

		actionIsValid := false
		for _, value := range actionValidValues {
			if value == actionParameterValue {
				actionIsValid = true
			}
		}

		if !actionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Action",
				ParameterValue: actionParameterValue,
				AllowedValues:  actionValidValues,
			}
		}
	}

	if v.Direction != nil {
		directionValidValues := []string{"0", "1"}
		directionParameterValue := fmt.Sprint(*v.Direction)

		directionIsValid := false
		for _, value := range directionValidValues {
			if value == directionParameterValue {
				directionIsValid = true
			}
		}

		if !directionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Direction",
				ParameterValue: directionParameterValue,
				AllowedValues:  directionValidValues,
			}
		}
	}

	return nil
}

type SecurityGroupSnapshot struct {
	GroupID                 *string              `json:"group_id" name:"group_id"`
	Rules                   []*SecurityGroupRule `json:"rules" name:"rules"`
	SecurityGroupSnapshotID *string              `json:"security_group_snapshot_id" name:"security_group_snapshot_id"`
}

func (v *SecurityGroupSnapshot) Validate() error {

	if len(v.Rules) > 0 {
		for _, property := range v.Rules {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type ServerCertificate struct {
	CertificateContent    *string    `json:"certificate_content" name:"certificate_content"`
	CreateTime            *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description           *string    `json:"description" name:"description"`
	PrivateKey            *string    `json:"private_key" name:"private_key"`
	ServerCertificateID   *string    `json:"server_certificate_id" name:"server_certificate_id"`
	ServerCertificateName *string    `json:"server_certificate_name" name:"server_certificate_name"`
}

func (v *ServerCertificate) Validate() error {

	return nil
}

type Snapshot struct {
	CreateTime  *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description *string    `json:"description" name:"description"`
	HeadChain   *string    `json:"head_chain" name:"head_chain"`
	// IsHead's available values: 0, 1
	IsHead *int `json:"is_head" name:"is_head"`
	// IsTaken's available values: 0, 1
	IsTaken            *int              `json:"is_taken" name:"is_taken"`
	LatestSnapshotTime *time.Time        `json:"latest_snapshot_time" name:"latest_snapshot_time" format:"ISO 8601"`
	ParentID           *string           `json:"parent_id" name:"parent_id"`
	Provider           *string           `json:"provider" name:"provider"`
	Resource           *Resource         `json:"resource" name:"resource"`
	RootID             *string           `json:"root_id" name:"root_id"`
	Size               *int              `json:"size" name:"size"`
	SnapshotID         *string           `json:"snapshot_id" name:"snapshot_id"`
	SnapshotName       *string           `json:"snapshot_name" name:"snapshot_name"`
	SnapshotResource   *SnapshotResource `json:"snapshot_resource" name:"snapshot_resource"`
	SnapshotTime       *time.Time        `json:"snapshot_time" name:"snapshot_time" format:"ISO 8601"`
	// SnapshotType's available values: 0, 1
	SnapshotType *string `json:"snapshot_type" name:"snapshot_type"`
	// Status's available values: pending, available, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	TotalCount *int       `json:"total_count" name:"total_count"`
	TotalSize  *int       `json:"total_size" name:"total_size"`
	// TransitionStatus's available values: creating, suspending, resuming, deleting, recovering
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VirtualSize      *int    `json:"virtual_size" name:"virtual_size"`
	Visibility       *string `json:"visibility" name:"visibility"`
}

func (v *Snapshot) Validate() error {

	if v.IsHead != nil {
		isHeadValidValues := []string{"0", "1"}
		isHeadParameterValue := fmt.Sprint(*v.IsHead)

		isHeadIsValid := false
		for _, value := range isHeadValidValues {
			if value == isHeadParameterValue {
				isHeadIsValid = true
			}
		}

		if !isHeadIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsHead",
				ParameterValue: isHeadParameterValue,
				AllowedValues:  isHeadValidValues,
			}
		}
	}

	if v.IsTaken != nil {
		isTakenValidValues := []string{"0", "1"}
		isTakenParameterValue := fmt.Sprint(*v.IsTaken)

		isTakenIsValid := false
		for _, value := range isTakenValidValues {
			if value == isTakenParameterValue {
				isTakenIsValid = true
			}
		}

		if !isTakenIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsTaken",
				ParameterValue: isTakenParameterValue,
				AllowedValues:  isTakenValidValues,
			}
		}
	}

	if v.Resource != nil {
		if err := v.Resource.Validate(); err != nil {
			return err
		}
	}

	if v.SnapshotResource != nil {
		if err := v.SnapshotResource.Validate(); err != nil {
			return err
		}
	}

	if v.SnapshotType != nil {
		snapshotTypeValidValues := []string{"0", "1"}
		snapshotTypeParameterValue := fmt.Sprint(*v.SnapshotType)

		snapshotTypeIsValid := false
		for _, value := range snapshotTypeValidValues {
			if value == snapshotTypeParameterValue {
				snapshotTypeIsValid = true
			}
		}

		if !snapshotTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "SnapshotType",
				ParameterValue: snapshotTypeParameterValue,
				AllowedValues:  snapshotTypeValidValues,
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "available", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "suspending", "resuming", "deleting", "recovering"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	return nil
}

type SnapshotResource struct {
	OSFamily *string `json:"os_family" name:"os_family"`
	Platform *string `json:"platform" name:"platform"`
}

func (v *SnapshotResource) Validate() error {

	return nil
}

type Tag struct {
	Color             *string              `json:"color" name:"color"`
	CreateTime        *time.Time           `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description       *string              `json:"description" name:"description"`
	Owner             *string              `json:"owner" name:"owner"`
	ResourceCount     *int                 `json:"resource_count" name:"resource_count"`
	ResourceTagPairs  []*ResourceTagPair   `json:"resource_tag_pairs" name:"resource_tag_pairs"`
	ResourceTypeCount []*ResourceTypeCount `json:"resource_type_count" name:"resource_type_count"`
	TagID             *string              `json:"tag_id" name:"tag_id"`
	TagKey            *string              `json:"tag_key" name:"tag_key"`
	TagName           *string              `json:"tag_name" name:"tag_name"`
}

func (v *Tag) Validate() error {

	if len(v.ResourceTagPairs) > 0 {
		for _, property := range v.ResourceTagPairs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if len(v.ResourceTypeCount) > 0 {
		for _, property := range v.ResourceTypeCount {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type User struct {
	Email  *string `json:"email" name:"email"`
	UserID *string `json:"user_id" name:"user_id"`
}

func (v *User) Validate() error {

	return nil
}

type Volume struct {
	CreateTime         *time.Time  `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description        *string     `json:"description" name:"description"`
	Device             *string     `json:"device" name:"device"`
	Instance           *Instance   `json:"instance" name:"instance"`
	Instances          []*Instance `json:"instances" name:"instances"`
	LatestSnapshotTime *time.Time  `json:"latest_snapshot_time" name:"latest_snapshot_time" format:"ISO 8601"`
	Owner              *string     `json:"owner" name:"owner"`
	PlaceGroupID       *string     `json:"place_group_id" name:"place_group_id"`
	Size               *int        `json:"size" name:"size"`
	// Status's available values: pending, available, in-use, suspended, deleted, ceased
	Status     *string    `json:"status" name:"status"`
	StatusTime *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
	SubCode    *int       `json:"sub_code" name:"sub_code"`
	Tags       []*Tag     `json:"tags" name:"tags"`
	// TransitionStatus's available values: creating, attaching, detaching, suspending, resuming, deleting, recovering
	TransitionStatus *string `json:"transition_status" name:"transition_status"`
	VolumeID         *string `json:"volume_id" name:"volume_id"`
	VolumeName       *string `json:"volume_name" name:"volume_name"`
	// VolumeType's available values: 0, 1, 2, 3
	VolumeType *int `json:"volume_type" name:"volume_type"`
}

func (v *Volume) Validate() error {

	if v.Instance != nil {
		if err := v.Instance.Validate(); err != nil {
			return err
		}
	}

	if len(v.Instances) > 0 {
		for _, property := range v.Instances {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.Status != nil {
		statusValidValues := []string{"pending", "available", "in-use", "suspended", "deleted", "ceased"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.TransitionStatus != nil {
		transitionStatusValidValues := []string{"creating", "attaching", "detaching", "suspending", "resuming", "deleting", "recovering"}
		transitionStatusParameterValue := fmt.Sprint(*v.TransitionStatus)

		transitionStatusIsValid := false
		for _, value := range transitionStatusValidValues {
			if value == transitionStatusParameterValue {
				transitionStatusIsValid = true
			}
		}

		if !transitionStatusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TransitionStatus",
				ParameterValue: transitionStatusParameterValue,
				AllowedValues:  transitionStatusValidValues,
			}
		}
	}

	if v.VolumeType != nil {
		volumeTypeValidValues := []string{"0", "1", "2", "3"}
		volumeTypeParameterValue := fmt.Sprint(*v.VolumeType)

		volumeTypeIsValid := false
		for _, value := range volumeTypeValidValues {
			if value == volumeTypeParameterValue {
				volumeTypeIsValid = true
			}
		}

		if !volumeTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VolumeType",
				ParameterValue: volumeTypeParameterValue,
				AllowedValues:  volumeTypeValidValues,
			}
		}
	}

	return nil
}

type VxNet struct {
	AvailableIPCount *int       `json:"available_ip_count" name:"available_ip_count"`
	CreateTime       *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
	Description      *string    `json:"description" name:"description"`
	InstanceIDs      []*string  `json:"instance_ids" name:"instance_ids"`
	Owner            *string    `json:"owner" name:"owner"`
	Router           *Router    `json:"router" name:"router"`
	Tags             []*Tag     `json:"tags" name:"tags"`
	VpcRouterID      *string    `json:"vpc_router_id" name:"vpc_router_id"`
	VxNetID          *string    `json:"vxnet_id" name:"vxnet_id"`
	VxNetName        *string    `json:"vxnet_name" name:"vxnet_name"`
	// VxNetType's available values: 0, 1
	VxNetType *int `json:"vxnet_type" name:"vxnet_type"`
}

func (v *VxNet) Validate() error {

	if v.Router != nil {
		if err := v.Router.Validate(); err != nil {
			return err
		}
	}

	if len(v.Tags) > 0 {
		for _, property := range v.Tags {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.VxNetType != nil {
		vxnetTypeValidValues := []string{"0", "1"}
		vxnetTypeParameterValue := fmt.Sprint(*v.VxNetType)

		vxnetTypeIsValid := false
		for _, value := range vxnetTypeValidValues {
			if value == vxnetTypeParameterValue {
				vxnetTypeIsValid = true
			}
		}

		if !vxnetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VxNetType",
				ParameterValue: vxnetTypeParameterValue,
				AllowedValues:  vxnetTypeValidValues,
			}
		}
	}

	return nil
}

type Zone struct {
	// Status's available values: active, faulty, defunct
	Status *string `json:"status" name:"status"`
	ZoneID *string `json:"zone_id" name:"zone_id"`
}

func (v *Zone) Validate() error {

	if v.Status != nil {
		statusValidValues := []string{"active", "faulty", "defunct"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
			}
		}
	}

	return nil
}
