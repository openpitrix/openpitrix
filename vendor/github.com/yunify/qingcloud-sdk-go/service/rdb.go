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

type RDBService struct {
	Config     *config.Config
	Properties *RDBServiceProperties
}

type RDBServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) RDB(zone string) (*RDBService, error) {
	properties := &RDBServiceProperties{
		Zone: &zone,
	}

	return &RDBService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/apply_rdb_parameter_group.html
func (s *RDBService) ApplyRDBParameterGroup(i *ApplyRDBParameterGroupInput) (*ApplyRDBParameterGroupOutput, error) {
	if i == nil {
		i = &ApplyRDBParameterGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ApplyRDBParameterGroup",
		RequestMethod: "GET",
	}

	x := &ApplyRDBParameterGroupOutput{}
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

type ApplyRDBParameterGroupInput struct {
	RDB *string `json:"rdb" name:"rdb" location:"params"` // Required
}

func (v *ApplyRDBParameterGroupInput) Validate() error {

	if v.RDB == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDB",
			ParentName:    "ApplyRDBParameterGroupInput",
		}
	}

	return nil
}

type ApplyRDBParameterGroupOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RDB     *string `json:"rdb" name:"rdb" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/cease_rdb_instance.html
func (s *RDBService) CeaseRDBInstance(i *CeaseRDBInstanceInput) (*CeaseRDBInstanceOutput, error) {
	if i == nil {
		i = &CeaseRDBInstanceInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CeaseRDBInstance",
		RequestMethod: "GET",
	}

	x := &CeaseRDBInstanceOutput{}
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

type CeaseRDBInstanceInput struct {
	RDB         *string `json:"rdb" name:"rdb" location:"params"`                   // Required
	RDBInstance *string `json:"rdb_instance" name:"rdb_instance" location:"params"` // Required
}

func (v *CeaseRDBInstanceInput) Validate() error {

	if v.RDB == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDB",
			ParentName:    "CeaseRDBInstanceInput",
		}
	}

	if v.RDBInstance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBInstance",
			ParentName:    "CeaseRDBInstanceInput",
		}
	}

	return nil
}

type CeaseRDBInstanceOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/copy_rdb_instance_files_to_ftp.html
func (s *RDBService) CopyRDBInstanceFilesToFTP(i *CopyRDBInstanceFilesToFTPInput) (*CopyRDBInstanceFilesToFTPOutput, error) {
	if i == nil {
		i = &CopyRDBInstanceFilesToFTPInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CopyRDBInstanceFilesToFTP",
		RequestMethod: "GET",
	}

	x := &CopyRDBInstanceFilesToFTPOutput{}
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

type CopyRDBInstanceFilesToFTPInput struct {
	Files       []*string `json:"files" name:"files" location:"params"`               // Required
	RDBInstance *string   `json:"rdb_instance" name:"rdb_instance" location:"params"` // Required
}

func (v *CopyRDBInstanceFilesToFTPInput) Validate() error {

	if len(v.Files) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Files",
			ParentName:    "CopyRDBInstanceFilesToFTPInput",
		}
	}

	if v.RDBInstance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBInstance",
			ParentName:    "CopyRDBInstanceFilesToFTPInput",
		}
	}

	return nil
}

type CopyRDBInstanceFilesToFTPOutput struct {
	Message     *string `json:"message" name:"message"`
	Action      *string `json:"action" name:"action" location:"elements"`
	JobID       *string `json:"job_id" name:"job_id" location:"elements"`
	RDBInstance *string `json:"rdb_instance" name:"rdb_instance" location:"elements"`
	RetCode     *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/create_rdb.html
func (s *RDBService) CreateRDB(i *CreateRDBInput) (*CreateRDBOutput, error) {
	if i == nil {
		i = &CreateRDBInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateRDB",
		RequestMethod: "GET",
	}

	x := &CreateRDBOutput{}
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

type CreateRDBInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	Description    *string `json:"description" name:"description" location:"params"`
	// EngineVersion's available values: mysql,5.5, mysql,5.6, mysql,5.7, psql,9.3, psql,9.4
	EngineVersion *string         `json:"engine_version" name:"engine_version" default:"mysql,5.7" location:"params"`
	NodeCount     *int            `json:"node_count" name:"node_count" location:"params"`
	PrivateIPs    []*RDBPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	ProxyCount    *int            `json:"proxy_count" name:"proxy_count" location:"params"`
	RDBClass      *int            `json:"rdb_class" name:"rdb_class" location:"params"`
	// RDBEngine's available values: mysql, psql
	RDBEngine   *string `json:"rdb_engine" name:"rdb_engine" default:"mysql" location:"params"`
	RDBName     *string `json:"rdb_name" name:"rdb_name" location:"params"`
	RDBPassword *string `json:"rdb_password" name:"rdb_password" location:"params"` // Required
	// RDBType's available values: 1, 2, 4, 8, 16, 32
	RDBType     *int    `json:"rdb_type" name:"rdb_type" location:"params"`         // Required
	RDBUsername *string `json:"rdb_username" name:"rdb_username" location:"params"` // Required
	StorageSize *int    `json:"storage_size" name:"storage_size" location:"params"` // Required
	VxNet       *string `json:"vxnet" name:"vxnet" location:"params"`               // Required
}

func (v *CreateRDBInput) Validate() error {

	if v.EngineVersion != nil {
		engineVersionValidValues := []string{"mysql,5.5", "mysql,5.6", "mysql,5.7", "psql,9.3", "psql,9.4"}
		engineVersionParameterValue := fmt.Sprint(*v.EngineVersion)

		engineVersionIsValid := false
		for _, value := range engineVersionValidValues {
			if value == engineVersionParameterValue {
				engineVersionIsValid = true
			}
		}

		if !engineVersionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "EngineVersion",
				ParameterValue: engineVersionParameterValue,
				AllowedValues:  engineVersionValidValues,
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

	if v.RDBEngine != nil {
		rdbEngineValidValues := []string{"mysql", "psql"}
		rdbEngineParameterValue := fmt.Sprint(*v.RDBEngine)

		rdbEngineIsValid := false
		for _, value := range rdbEngineValidValues {
			if value == rdbEngineParameterValue {
				rdbEngineIsValid = true
			}
		}

		if !rdbEngineIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RDBEngine",
				ParameterValue: rdbEngineParameterValue,
				AllowedValues:  rdbEngineValidValues,
			}
		}
	}

	if v.RDBPassword == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBPassword",
			ParentName:    "CreateRDBInput",
		}
	}

	if v.RDBType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBType",
			ParentName:    "CreateRDBInput",
		}
	}

	if v.RDBType != nil {
		rdbTypeValidValues := []string{"1", "2", "4", "8", "16", "32"}
		rdbTypeParameterValue := fmt.Sprint(*v.RDBType)

		rdbTypeIsValid := false
		for _, value := range rdbTypeValidValues {
			if value == rdbTypeParameterValue {
				rdbTypeIsValid = true
			}
		}

		if !rdbTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RDBType",
				ParameterValue: rdbTypeParameterValue,
				AllowedValues:  rdbTypeValidValues,
			}
		}
	}

	if v.RDBUsername == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBUsername",
			ParentName:    "CreateRDBInput",
		}
	}

	if v.StorageSize == nil {
		return errors.ParameterRequiredError{
			ParameterName: "StorageSize",
			ParentName:    "CreateRDBInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateRDBInput",
		}
	}

	return nil
}

type CreateRDBOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RDB     *string `json:"rdb" name:"rdb" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/create_rdb_from_snapshot.html
func (s *RDBService) CreateRDBFromSnapshot(i *CreateRDBFromSnapshotInput) (*CreateRDBFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateRDBFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateRDBFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateRDBFromSnapshotOutput{}
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

type CreateRDBFromSnapshotInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	Description    *string `json:"description" name:"description" location:"params"`
	// EngineVersion's available values: mysql,5.5, mysql,5.6, mysql,5.7, psql,9.3, psql,9.4
	EngineVersion *string         `json:"engine_version" name:"engine_version" default:"mysql,5.7" location:"params"`
	NodeCount     *int            `json:"node_count" name:"node_count" location:"params"`
	PrivateIPs    []*RDBPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	ProxyCount    *int            `json:"proxy_count" name:"proxy_count" location:"params"`
	// RDBEngine's available values: mysql, psql
	RDBEngine *string `json:"rdb_engine" name:"rdb_engine" default:"mysql" location:"params"`
	RDBName   *string `json:"rdb_name" name:"rdb_name" location:"params"`
	// RDBType's available values: 1, 2, 4, 8, 16, 32
	RDBType     *int    `json:"rdb_type" name:"rdb_type" location:"params"` // Required
	Snapshot    *string `json:"snapshot" name:"snapshot" location:"params"` // Required
	StorageSize *int    `json:"storage_size" name:"storage_size" location:"params"`
	VxNet       *string `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *CreateRDBFromSnapshotInput) Validate() error {

	if v.EngineVersion != nil {
		engineVersionValidValues := []string{"mysql,5.5", "mysql,5.6", "mysql,5.7", "psql,9.3", "psql,9.4"}
		engineVersionParameterValue := fmt.Sprint(*v.EngineVersion)

		engineVersionIsValid := false
		for _, value := range engineVersionValidValues {
			if value == engineVersionParameterValue {
				engineVersionIsValid = true
			}
		}

		if !engineVersionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "EngineVersion",
				ParameterValue: engineVersionParameterValue,
				AllowedValues:  engineVersionValidValues,
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

	if v.RDBEngine != nil {
		rdbEngineValidValues := []string{"mysql", "psql"}
		rdbEngineParameterValue := fmt.Sprint(*v.RDBEngine)

		rdbEngineIsValid := false
		for _, value := range rdbEngineValidValues {
			if value == rdbEngineParameterValue {
				rdbEngineIsValid = true
			}
		}

		if !rdbEngineIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RDBEngine",
				ParameterValue: rdbEngineParameterValue,
				AllowedValues:  rdbEngineValidValues,
			}
		}
	}

	if v.RDBType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBType",
			ParentName:    "CreateRDBFromSnapshotInput",
		}
	}

	if v.RDBType != nil {
		rdbTypeValidValues := []string{"1", "2", "4", "8", "16", "32"}
		rdbTypeParameterValue := fmt.Sprint(*v.RDBType)

		rdbTypeIsValid := false
		for _, value := range rdbTypeValidValues {
			if value == rdbTypeParameterValue {
				rdbTypeIsValid = true
			}
		}

		if !rdbTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RDBType",
				ParameterValue: rdbTypeParameterValue,
				AllowedValues:  rdbTypeValidValues,
			}
		}
	}

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "CreateRDBFromSnapshotInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateRDBFromSnapshotInput",
		}
	}

	return nil
}

type CreateRDBFromSnapshotOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RDB     *string `json:"rdb" name:"rdb" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/create_temp_rdb_instance_from_snapshot.html
func (s *RDBService) CreateTempRDBInstanceFromSnapshot(i *CreateTempRDBInstanceFromSnapshotInput) (*CreateTempRDBInstanceFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateTempRDBInstanceFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateTempRDBInstanceFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateTempRDBInstanceFromSnapshotOutput{}
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

type CreateTempRDBInstanceFromSnapshotInput struct {
	RDB      *string `json:"rdb" name:"rdb" location:"params"`           // Required
	Snapshot *string `json:"snapshot" name:"snapshot" location:"params"` // Required
}

func (v *CreateTempRDBInstanceFromSnapshotInput) Validate() error {

	if v.RDB == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDB",
			ParentName:    "CreateTempRDBInstanceFromSnapshotInput",
		}
	}

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "CreateTempRDBInstanceFromSnapshotInput",
		}
	}

	return nil
}

type CreateTempRDBInstanceFromSnapshotOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RDB     *string `json:"rdb" name:"rdb" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/delete_rdbs.html
func (s *RDBService) DeleteRDBs(i *DeleteRDBsInput) (*DeleteRDBsOutput, error) {
	if i == nil {
		i = &DeleteRDBsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteRDBs",
		RequestMethod: "GET",
	}

	x := &DeleteRDBsOutput{}
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

type DeleteRDBsInput struct {
	RDBs []*string `json:"rdbs" name:"rdbs" location:"params"` // Required
}

func (v *DeleteRDBsInput) Validate() error {

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "DeleteRDBsInput",
		}
	}

	return nil
}

type DeleteRDBsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RDBs    []*string `json:"rdbs" name:"rdbs" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/describe_rdb_parameters.html
func (s *RDBService) DescribeRDBParameters(i *DescribeRDBParametersInput) (*DescribeRDBParametersOutput, error) {
	if i == nil {
		i = &DescribeRDBParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRDBParameters",
		RequestMethod: "GET",
	}

	x := &DescribeRDBParametersOutput{}
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

type DescribeRDBParametersInput struct {
	Limit          *int    `json:"limit" name:"limit" location:"params"`
	Offset         *int    `json:"offset" name:"offset" location:"params"`
	ParameterGroup *string `json:"parameter_group" name:"parameter_group" location:"params"`
	RDB            *string `json:"rdb" name:"rdb" location:"params"` // Required
}

func (v *DescribeRDBParametersInput) Validate() error {

	if v.RDB == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDB",
			ParentName:    "DescribeRDBParametersInput",
		}
	}

	return nil
}

type DescribeRDBParametersOutput struct {
	Message      *string         `json:"message" name:"message"`
	Action       *string         `json:"action" name:"action" location:"elements"`
	ParameterSet []*RDBParameter `json:"parameter_set" name:"parameter_set" location:"elements"`
	RetCode      *int            `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int            `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/describe_rdbs.html
func (s *RDBService) DescribeRDBs(i *DescribeRDBsInput) (*DescribeRDBsOutput, error) {
	if i == nil {
		i = &DescribeRDBsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRDBs",
		RequestMethod: "GET",
	}

	x := &DescribeRDBsOutput{}
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

type DescribeRDBsInput struct {
	Limit      *int      `json:"limit" name:"limit" location:"params"`
	Offset     *int      `json:"offset" name:"offset" location:"params"`
	RDBEngine  *string   `json:"rdb_engine" name:"rdb_engine" location:"params"`
	RDBName    *string   `json:"rdb_name" name:"rdb_name" location:"params"`
	RDBs       []*string `json:"rdbs" name:"rdbs" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeRDBsInput) Validate() error {

	return nil
}

type DescribeRDBsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	RDBSet     []*RDB  `json:"rdb_set" name:"rdb_set" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/get_rdb_instance_files.html
func (s *RDBService) GetRDBInstanceFiles(i *GetRDBInstanceFilesInput) (*GetRDBInstanceFilesOutput, error) {
	if i == nil {
		i = &GetRDBInstanceFilesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetRDBInstanceFiles",
		RequestMethod: "GET",
	}

	x := &GetRDBInstanceFilesOutput{}
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

type GetRDBInstanceFilesInput struct {
	RDBInstance *string `json:"rdb_instance" name:"rdb_instance" location:"params"` // Required
}

func (v *GetRDBInstanceFilesInput) Validate() error {

	if v.RDBInstance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBInstance",
			ParentName:    "GetRDBInstanceFilesInput",
		}
	}

	return nil
}

type GetRDBInstanceFilesOutput struct {
	Message     *string  `json:"message" name:"message"`
	Action      *string  `json:"action" name:"action" location:"elements"`
	Files       *RDBFile `json:"files" name:"files" location:"elements"`
	RDBInstance *string  `json:"rdb_instance" name:"rdb_instance" location:"elements"`
	RetCode     *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_rdb_monitor.html
func (s *RDBService) GetRDBMonitor(i *GetRDBMonitorInput) (*GetRDBMonitorOutput, error) {
	if i == nil {
		i = &GetRDBMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetRDBMonitor",
		RequestMethod: "GET",
	}

	x := &GetRDBMonitorOutput{}
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

type GetRDBMonitorInput struct {
	EndTime     *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"` // Required
	Meters      []*string  `json:"meters" name:"meters" location:"params"`                       // Required
	RDBEngine   *string    `json:"rdb_engine" name:"rdb_engine" location:"params"`               // Required
	RDBInstance *string    `json:"rdb_instance" name:"rdb_instance" location:"params"`
	Resource    *string    `json:"resource" name:"resource" location:"params"`                       // Required
	Role        *string    `json:"role" name:"role" location:"params"`                               // Required
	StartTime   *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetRDBMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetRDBMonitorInput",
		}
	}

	if v.RDBEngine == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDBEngine",
			ParentName:    "GetRDBMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetRDBMonitorInput",
		}
	}

	if v.Role == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Role",
			ParentName:    "GetRDBMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetRDBMonitorInput",
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

type GetRDBMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/modify_rdb_parameters.html
func (s *RDBService) ModifyRDBParameters(i *ModifyRDBParametersInput) (*ModifyRDBParametersOutput, error) {
	if i == nil {
		i = &ModifyRDBParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyRDBParameters",
		RequestMethod: "GET",
	}

	x := &ModifyRDBParametersOutput{}
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

type ModifyRDBParametersInput struct {
	Parameters []*RDBParameters `json:"parameters" name:"parameters" location:"params"`
	RDB        *string          `json:"rdb" name:"rdb" location:"params"` // Required
}

func (v *ModifyRDBParametersInput) Validate() error {

	if len(v.Parameters) > 0 {
		for _, property := range v.Parameters {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.RDB == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RDB",
			ParentName:    "ModifyRDBParametersInput",
		}
	}

	return nil
}

type ModifyRDBParametersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RDB     *string `json:"rdb" name:"rdb" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/rdbs_join_vxnet.html
func (s *RDBService) RDBsJoinVxNet(i *RDBsJoinVxNetInput) (*RDBsJoinVxNetOutput, error) {
	if i == nil {
		i = &RDBsJoinVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RDBsJoinVxnet",
		RequestMethod: "GET",
	}

	x := &RDBsJoinVxNetOutput{}
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

type RDBsJoinVxNetInput struct {
	RDBs  []*string `json:"rdbs" name:"rdbs" location:"params"`   // Required
	VxNet *string   `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *RDBsJoinVxNetInput) Validate() error {

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "RDBsJoinVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "RDBsJoinVxNetInput",
		}
	}

	return nil
}

type RDBsJoinVxNetOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RDBs    []*string `json:"rdbs" name:"rdbs" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNet   *string   `json:"vxnet" name:"vxnet" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/rdbs_leave_vxnet.html
func (s *RDBService) RDBsLeaveVxNet(i *RDBsLeaveVxNetInput) (*RDBsLeaveVxNetOutput, error) {
	if i == nil {
		i = &RDBsLeaveVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RDBsLeaveVxnet",
		RequestMethod: "GET",
	}

	x := &RDBsLeaveVxNetOutput{}
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

type RDBsLeaveVxNetInput struct {
	RDBs  []*string `json:"rdbs" name:"rdbs" location:"params"`   // Required
	VxNet *string   `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *RDBsLeaveVxNetInput) Validate() error {

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "RDBsLeaveVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "RDBsLeaveVxNetInput",
		}
	}

	return nil
}

type RDBsLeaveVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/resize_rdbs.html
func (s *RDBService) ResizeRDBs(i *ResizeRDBsInput) (*ResizeRDBsOutput, error) {
	if i == nil {
		i = &ResizeRDBsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeRDBs",
		RequestMethod: "GET",
	}

	x := &ResizeRDBsOutput{}
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

type ResizeRDBsInput struct {

	// RDBType's available values: 1, 2, 4, 8, 16, 32
	RDBType     *int      `json:"rdb_type" name:"rdb_type" location:"params"`
	RDBs        []*string `json:"rdbs" name:"rdbs" location:"params"` // Required
	StorageSize *int      `json:"storage_size" name:"storage_size" location:"params"`
}

func (v *ResizeRDBsInput) Validate() error {

	if v.RDBType != nil {
		rdbTypeValidValues := []string{"1", "2", "4", "8", "16", "32"}
		rdbTypeParameterValue := fmt.Sprint(*v.RDBType)

		rdbTypeIsValid := false
		for _, value := range rdbTypeValidValues {
			if value == rdbTypeParameterValue {
				rdbTypeIsValid = true
			}
		}

		if !rdbTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RDBType",
				ParameterValue: rdbTypeParameterValue,
				AllowedValues:  rdbTypeValidValues,
			}
		}
	}

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "ResizeRDBsInput",
		}
	}

	return nil
}

type ResizeRDBsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RDBs    []*string `json:"rdbs" name:"rdbs" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/start_rdbs.html
func (s *RDBService) StartRDBs(i *StartRDBsInput) (*StartRDBsOutput, error) {
	if i == nil {
		i = &StartRDBsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartRDBs",
		RequestMethod: "GET",
	}

	x := &StartRDBsOutput{}
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

type StartRDBsInput struct {
	RDBs []*string `json:"rdbs" name:"rdbs" location:"params"` // Required
}

func (v *StartRDBsInput) Validate() error {

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "StartRDBsInput",
		}
	}

	return nil
}

type StartRDBsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RDBs    []*string `json:"rdbs" name:"rdbs" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/rdb/stop_rdbs.html
func (s *RDBService) StopRDBs(i *StopRDBsInput) (*StopRDBsOutput, error) {
	if i == nil {
		i = &StopRDBsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopRDBs",
		RequestMethod: "GET",
	}

	x := &StopRDBsOutput{}
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

type StopRDBsInput struct {
	RDBs []*string `json:"rdbs" name:"rdbs" location:"params"` // Required
}

func (v *StopRDBsInput) Validate() error {

	if len(v.RDBs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RDBs",
			ParentName:    "StopRDBsInput",
		}
	}

	return nil
}

type StopRDBsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RDBs    []*string `json:"rdbs" name:"rdbs" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}
