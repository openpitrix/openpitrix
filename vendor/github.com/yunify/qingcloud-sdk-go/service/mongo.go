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

type MongoService struct {
	Config     *config.Config
	Properties *MongoServiceProperties
}

type MongoServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Mongo(zone string) (*MongoService, error) {
	properties := &MongoServiceProperties{
		Zone: &zone,
	}

	return &MongoService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/add_mongo_instances.html
func (s *MongoService) AddMongoInstances(i *AddMongoInstancesInput) (*AddMongoInstancesOutput, error) {
	if i == nil {
		i = &AddMongoInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddMongoInstances",
		RequestMethod: "GET",
	}

	x := &AddMongoInstancesOutput{}
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

type AddMongoInstancesInput struct {
	Mongo      *string           `json:"mongo" name:"mongo" location:"params"`
	NodeCount  *int              `json:"node_count" name:"node_count" location:"params"`
	PrivateIPs []*MongoPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
}

func (v *AddMongoInstancesInput) Validate() error {

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AddMongoInstancesOutput struct {
	Message   *string   `json:"message" name:"message"`
	Action    *string   `json:"action" name:"action" location:"elements"`
	JobID     *string   `json:"job_id" name:"job_id" location:"elements"`
	Mongo     *string   `json:"mongo" name:"mongo" location:"elements"`
	MongoNode []*string `json:"mongo_node" name:"mongo_node" location:"elements"`
	RetCode   *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/change_mongo_vxnet.html
func (s *MongoService) ChangeMongoVxNet(i *ChangeMongoVxNetInput) (*ChangeMongoVxNetOutput, error) {
	if i == nil {
		i = &ChangeMongoVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeMongoVxnet",
		RequestMethod: "GET",
	}

	x := &ChangeMongoVxNetOutput{}
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

type ChangeMongoVxNetInput struct {
	Mongo      *string           `json:"mongo" name:"mongo" location:"params"` // Required
	PrivateIPs []*MongoPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	VxNet      *string           `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *ChangeMongoVxNetInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "ChangeMongoVxNetInput",
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
			ParentName:    "ChangeMongoVxNetInput",
		}
	}

	return nil
}

type ChangeMongoVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/create_mongo.html
func (s *MongoService) CreateMongo(i *CreateMongoInput) (*CreateMongoOutput, error) {
	if i == nil {
		i = &CreateMongoInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateMongo",
		RequestMethod: "GET",
	}

	x := &CreateMongoOutput{}
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

type CreateMongoInput struct {
	AutoBackupTime *int              `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	Description    *string           `json:"description" name:"description" location:"params"`
	MongoName      *string           `json:"mongo_name" name:"mongo_name" location:"params"`
	MongoPassword  *string           `json:"mongo_password" name:"mongo_password" location:"params"`
	MongoType      *int              `json:"mongo_type" name:"mongo_type" location:"params"` // Required
	MongoUsername  *string           `json:"mongo_username" name:"mongo_username" location:"params"`
	MongoVersion   *string           `json:"mongo_version" name:"mongo_version" location:"params"`
	PrivateIPs     []*MongoPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
	ResourceClass  *int              `json:"resource_class" name:"resource_class" location:"params"`
	StorageSize    *int              `json:"storage_size" name:"storage_size" location:"params"` // Required
	VxNet          *string           `json:"vxnet" name:"vxnet" location:"params"`               // Required
}

func (v *CreateMongoInput) Validate() error {

	if v.MongoType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "MongoType",
			ParentName:    "CreateMongoInput",
		}
	}

	if len(v.PrivateIPs) > 0 {
		for _, property := range v.PrivateIPs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.StorageSize == nil {
		return errors.ParameterRequiredError{
			ParameterName: "StorageSize",
			ParentName:    "CreateMongoInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateMongoInput",
		}
	}

	return nil
}

type CreateMongoOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/create_mongo_from_snapshot.html
func (s *MongoService) CreateMongoFromSnapshot(i *CreateMongoFromSnapshotInput) (*CreateMongoFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateMongoFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateMongoFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateMongoFromSnapshotOutput{}
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

type CreateMongoFromSnapshotInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	MongoName      *string `json:"mongo_name" name:"mongo_name" location:"params"`
	MongoType      *int    `json:"mongo_type" name:"mongo_type" location:"params"`
	MongoVersion   *int    `json:"mongo_version" name:"mongo_version" location:"params"`
	ResourceClass  *int    `json:"resource_class" name:"resource_class" location:"params"`
	Snapshot       *string `json:"snapshot" name:"snapshot" location:"params"`
	StorageSize    *int    `json:"storage_size" name:"storage_size" location:"params"`
	VxNet          *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *CreateMongoFromSnapshotInput) Validate() error {

	return nil
}

type CreateMongoFromSnapshotOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/delete_mongos.html
func (s *MongoService) DeleteMongos(i *DeleteMongosInput) (*DeleteMongosOutput, error) {
	if i == nil {
		i = &DeleteMongosInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteMongos",
		RequestMethod: "GET",
	}

	x := &DeleteMongosOutput{}
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

type DeleteMongosInput struct {
	Mongos []*string `json:"mongos" name:"mongos" location:"params"` // Required
}

func (v *DeleteMongosInput) Validate() error {

	if len(v.Mongos) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Mongos",
			ParentName:    "DeleteMongosInput",
		}
	}

	return nil
}

type DeleteMongosOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	Mongos  []*string `json:"mongos" name:"mongos" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/describe_mongo_nodes.html
func (s *MongoService) DescribeMongoNodes(i *DescribeMongoNodesInput) (*DescribeMongoNodesOutput, error) {
	if i == nil {
		i = &DescribeMongoNodesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeMongoNodes",
		RequestMethod: "GET",
	}

	x := &DescribeMongoNodesOutput{}
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

type DescribeMongoNodesInput struct {
	Limit  *int      `json:"limit" name:"limit" location:"params"`
	Mongo  *string   `json:"mongo" name:"mongo" location:"params"` // Required
	Offset *int      `json:"offset" name:"offset" location:"params"`
	Status []*string `json:"status" name:"status" location:"params"`
}

func (v *DescribeMongoNodesInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "DescribeMongoNodesInput",
		}
	}

	return nil
}

type DescribeMongoNodesOutput struct {
	Message      *string      `json:"message" name:"message"`
	Action       *string      `json:"action" name:"action" location:"elements"`
	MongoNodeSet []*MongoNode `json:"mongo_node_set" name:"mongo_node_set" location:"elements"`
	RetCode      *int         `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int         `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/describe_mongo_parameters.html
func (s *MongoService) DescribeMongoParameters(i *DescribeMongoParametersInput) (*DescribeMongoParametersOutput, error) {
	if i == nil {
		i = &DescribeMongoParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeMongoParameters",
		RequestMethod: "GET",
	}

	x := &DescribeMongoParametersOutput{}
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

type DescribeMongoParametersInput struct {
	Limit  *int    `json:"limit" name:"limit" default:"20" location:"params"`
	Mongo  *string `json:"mongo" name:"mongo" location:"params"` // Required
	Offset *int    `json:"offset" name:"offset" default:"0" location:"params"`
}

func (v *DescribeMongoParametersInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "DescribeMongoParametersInput",
		}
	}

	return nil
}

type DescribeMongoParametersOutput struct {
	Message      *string           `json:"message" name:"message"`
	Action       *string           `json:"action" name:"action" location:"elements"`
	ParameterSet []*MongoParameter `json:"parameter_set" name:"parameter_set" location:"elements"`
	RetCode      *int              `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int              `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/describe_mongos.html
func (s *MongoService) DescribeMongos(i *DescribeMongosInput) (*DescribeMongosOutput, error) {
	if i == nil {
		i = &DescribeMongosInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeMongos",
		RequestMethod: "GET",
	}

	x := &DescribeMongosOutput{}
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

type DescribeMongosInput struct {
	Limit     *int      `json:"limit" name:"limit" default:"20" location:"params"`
	MongoName *string   `json:"mongo_name" name:"mongo_name" location:"params"`
	Mongos    []*string `json:"mongos" name:"mongos" location:"params"`
	Offset    *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Status    []*string `json:"status" name:"status" location:"params"`
	Tags      []*string `json:"tags" name:"tags" location:"params"`
	Verbose   *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeMongosInput) Validate() error {

	return nil
}

type DescribeMongosOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MongoSet   []*Mongo `json:"mongo_set" name:"mongo_set" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int     `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_mongo_monitor.html
func (s *MongoService) GetMongoMonitor(i *GetMongoMonitorInput) (*GetMongoMonitorOutput, error) {
	if i == nil {
		i = &GetMongoMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetMongoMonitor",
		RequestMethod: "GET",
	}

	x := &GetMongoMonitorOutput{}
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

type GetMongoMonitorInput struct {
	EndTime   *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"`     // Required
	Meters    []*string  `json:"meters" name:"meters" location:"params"`                           // Required
	Resource  *string    `json:"resource" name:"resource" location:"params"`                       // Required
	StartTime *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetMongoMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetMongoMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetMongoMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetMongoMonitorInput",
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

type GetMongoMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/modify_mongo_attributes.html
func (s *MongoService) ModifyMongoAttributes(i *ModifyMongoAttributesInput) (*ModifyMongoAttributesOutput, error) {
	if i == nil {
		i = &ModifyMongoAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyMongoAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyMongoAttributesOutput{}
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

type ModifyMongoAttributesInput struct {
	AutoBackupTime *int    `json:"auto_backup_time" name:"auto_backup_time" location:"params"`
	Description    *string `json:"description" name:"description" location:"params"`
	Mongo          *string `json:"mongo" name:"mongo" location:"params"` // Required
	MongoName      *string `json:"mongo_name" name:"mongo_name" location:"params"`
}

func (v *ModifyMongoAttributesInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "ModifyMongoAttributesInput",
		}
	}

	return nil
}

type ModifyMongoAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/modify_mongo_instances.html
func (s *MongoService) ModifyMongoInstances(i *ModifyMongoInstancesInput) (*ModifyMongoInstancesOutput, error) {
	if i == nil {
		i = &ModifyMongoInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyMongoInstances",
		RequestMethod: "GET",
	}

	x := &ModifyMongoInstancesOutput{}
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

type ModifyMongoInstancesInput struct {
	Mongo      *string           `json:"mongo" name:"mongo" location:"params"` // Required
	PrivateIPs []*MongoPrivateIP `json:"private_ips" name:"private_ips" location:"params"`
}

func (v *ModifyMongoInstancesInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "ModifyMongoInstancesInput",
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

type ModifyMongoInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/remove_mongo_instances.html
func (s *MongoService) RemoveMongoInstances(i *RemoveMongoInstancesInput) (*RemoveMongoInstancesOutput, error) {
	if i == nil {
		i = &RemoveMongoInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RemoveMongoInstances",
		RequestMethod: "GET",
	}

	x := &RemoveMongoInstancesOutput{}
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

type RemoveMongoInstancesInput struct {
	Mongo          *string   `json:"mongo" name:"mongo" location:"params"`                     // Required
	MongoInstances []*string `json:"mongo_instances" name:"mongo_instances" location:"params"` // Required
}

func (v *RemoveMongoInstancesInput) Validate() error {

	if v.Mongo == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongo",
			ParentName:    "RemoveMongoInstancesInput",
		}
	}

	if len(v.MongoInstances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "MongoInstances",
			ParentName:    "RemoveMongoInstancesInput",
		}
	}

	return nil
}

type RemoveMongoInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	Mongo   *string `json:"mongo" name:"mongo" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/resize_mongos.html
func (s *MongoService) ResizeMongos(i *ResizeMongosInput) (*ResizeMongosOutput, error) {
	if i == nil {
		i = &ResizeMongosInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeMongos",
		RequestMethod: "GET",
	}

	x := &ResizeMongosOutput{}
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

type ResizeMongosInput struct {
	MongoType   *int      `json:"mongo_type" name:"mongo_type" location:"params"`
	Mongos      []*string `json:"mongos" name:"mongos" location:"params"` // Required
	StorageSize *int      `json:"storage_size" name:"storage_size" location:"params"`
}

func (v *ResizeMongosInput) Validate() error {

	if len(v.Mongos) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Mongos",
			ParentName:    "ResizeMongosInput",
		}
	}

	return nil
}

type ResizeMongosOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	Mongos  []*string `json:"mongos" name:"mongos" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/start_mongos.html
func (s *MongoService) StartMongos(i *StartMongosInput) (*StartMongosOutput, error) {
	if i == nil {
		i = &StartMongosInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartMongos",
		RequestMethod: "GET",
	}

	x := &StartMongosOutput{}
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

type StartMongosInput struct {
	Mongos *string `json:"mongos" name:"mongos" location:"params"` // Required
}

func (v *StartMongosInput) Validate() error {

	if v.Mongos == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Mongos",
			ParentName:    "StartMongosInput",
		}
	}

	return nil
}

type StartMongosOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/mongo/stop_mongos.html
func (s *MongoService) StopMongos(i *StopMongosInput) (*StopMongosOutput, error) {
	if i == nil {
		i = &StopMongosInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopMongos",
		RequestMethod: "GET",
	}

	x := &StopMongosOutput{}
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

type StopMongosInput struct {
	Mongos []*string `json:"mongos" name:"mongos" location:"params"` // Required
}

func (v *StopMongosInput) Validate() error {

	if len(v.Mongos) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Mongos",
			ParentName:    "StopMongosInput",
		}
	}

	return nil
}

type StopMongosOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
