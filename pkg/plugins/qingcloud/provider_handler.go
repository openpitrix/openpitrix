// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"

	"openpitrix.io/openpitrix/pkg/utils/jsontool"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
)

var MyProvider = constants.ProviderQingCloud

type ProviderHandler struct {
	vmbased.FrameHandler
}

func (p *ProviderHandler) initQingCloudService(runtimeUrl, runtimeCredential, zone string) (*qcservice.QingCloudService, error) {
	credential := new(Credential)
	err := json.Unmarshal([]byte(runtimeCredential), credential)
	if err != nil {
		logger.Errorf("Parse [%s] credential failed: %v", MyProvider, err)
		return nil, err
	}
	conf, err := qcconfig.New(credential.AccessKeyId, credential.SecretAccessKey)
	if err != nil {
		return nil, err
	}
	conf.Zone = zone
	if strings.HasPrefix(runtimeUrl, "https://") {
		runtimeUrl = strings.Split(runtimeUrl, "https://")[1]
	}
	urlAndPort := strings.Split(runtimeUrl, ":")
	if len(urlAndPort) == 2 {
		conf.Port, err = strconv.Atoi(urlAndPort[1])
	}
	conf.Host = urlAndPort[0]
	if err != nil {
		logger.Errorf("Parse [%s] runtimeUrl [%s] failed: %+v", MyProvider, runtimeUrl, err)
		return nil, err
	}
	return qcservice.Init(conf)
}

func (p *ProviderHandler) initService(runtimeId string) (*qcservice.QingCloudService, error) {
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	return p.initQingCloudService(runtime.RuntimeUrl, runtime.Credential, runtime.Zone)
}

func (p *ProviderHandler) RunInstances(task *models.Task) error {

	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", MyProvider, err)
		return err
	}

	input := &qcservice.RunInstancesInput{
		ImageID:       qcservice.String(instance.ImageId),
		CPU:           qcservice.Int(instance.Cpu),
		Memory:        qcservice.Int(instance.Memory),
		InstanceName:  qcservice.String(instance.Name),
		InstanceClass: qcservice.Int(DefaultInstanceClass),
		VxNets:        qcservice.StringSlice([]string{instance.Subnet}),
		LoginMode:     qcservice.String(DefaultLoginMode),
		LoginPasswd:   qcservice.String(DefaultLoginPassword),
		// GPU:     qcservice.Int(instance.Gpu),
	}
	if instance.VolumeId != "" {
		input.Volumes = qcservice.StringSlice([]string{instance.VolumeId})
	}
	if instance.UserdataPath != "" {
		input.UserdataPath = qcservice.String(instance.UserdataPath)
	}
	if instance.UserDataValue != "" {
		input.UserdataValue = qcservice.String(instance.UserDataValue)
	}
	logger.Debugf("RunInstances with input: %s", jsontool.ToString(input))
	output, err := instanceService.RunInstances(input)
	if err != nil {
		logger.Errorf("Send RunInstances to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send RunInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send RunInstances to %s failed: %s", MyProvider, message)
	}

	if len(output.Instances) == 0 {
		logger.Errorf("Send RunInstances to %s failed with 0 output instances", MyProvider)
		return fmt.Errorf("send RunInstances to %s failed with 0 output instances", MyProvider)
	}

	instance.InstanceId = qcservice.StringValue(output.Instances[0])
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", MyProvider, err)
		return err
	}

	volumeOutput, err := volumeService.DescribeVolumes(
		&qcservice.DescribeVolumesInput{
			Volumes: qcservice.StringSlice([]string{instance.VolumeId}),
		},
	)
	if err != nil {
		logger.Errorf("Send DescribeVolumes to %s failed: %v", MyProvider, err)
		return err
	}

	volumeRetCode := qcservice.IntValue(volumeOutput.RetCode)
	if volumeRetCode != 0 {
		volumeMessage := qcservice.StringValue(volumeOutput.Message)
		logger.Errorf("Send DescribeVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, volumeRetCode, volumeMessage)
		return fmt.Errorf("send DescribeVolumes to %s failed: %s", MyProvider, volumeMessage)
	}

	if len(volumeOutput.VolumeSet) == 0 {
		logger.Errorf("Send DescribeVolumes to %s failed with 0 output volumes", MyProvider)
		return fmt.Errorf("send DescribeVolumes to %s failed with 0 output volumes", MyProvider)
	}

	// Such as /dev/sdc
	instance.Device = qcservice.StringValue(volumeOutput.VolumeSet[0].Device)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = directive

	return nil
}

func (p *ProviderHandler) StopInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", MyProvider, err)
		return err
	}

	output, err := instanceService.StopInstances(
		&qcservice.StopInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("Send StopInstances to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send StopInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send StopInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) StartInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", MyProvider, err)
		return err
	}

	output, err := instanceService.StartInstances(
		&qcservice.StartInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("Send StartInstances to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send StartInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send StartInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) DeleteInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", MyProvider, err)
		return err
	}

	output, err := instanceService.TerminateInstances(
		&qcservice.TerminateInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("Send TerminateInstances to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send TerminateInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send TerminateInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) CreateVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", MyProvider, err)
		return err
	}

	output, err := volumeService.CreateVolumes(
		&qcservice.CreateVolumesInput{
			Size:       qcservice.Int(volume.Size),
			VolumeName: qcservice.String(volume.Name),
			VolumeType: qcservice.Int(DefaultVolumeClass),
		},
	)
	if err != nil {
		logger.Errorf("Send CreateVolumes to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send CreateVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send CreateVolumes to %s failed: %s", MyProvider, message)
	}
	volume.VolumeId = qcservice.StringValue(output.Volumes[0])
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := volume.ToString()
	if err != nil {
		return err
	}
	task.Directive = directive

	return nil
}

func (p *ProviderHandler) DetachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", MyProvider, err)
		return err
	}

	output, err := volumeService.DetachVolumes(
		&qcservice.DetachVolumesInput{
			Instance: qcservice.String(volume.InstanceId),
			Volumes:  qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Errorf("Send DetachVolumes to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DetachVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send DetachVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := volume.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) AttachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", MyProvider, err)
		return err
	}

	output, err := volumeService.AttachVolumes(
		&qcservice.AttachVolumesInput{
			Instance: qcservice.String(volume.InstanceId),
			Volumes:  qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Errorf("Send AttachVolumes to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send AttachVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send AttachVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := volume.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) DeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", MyProvider, err)
		return err
	}

	output, err := volumeService.DeleteVolumes(
		&qcservice.DeleteVolumesInput{
			Volumes: qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Errorf("Send DeleteVolumes to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DeleteVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send DeleteVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := volume.ToString()
	if err != nil {
		return err
	}
	task.Directive = string(directive)

	return nil
}

func (p *ProviderHandler) WaitRunInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", MyProvider, instance.TargetJobId, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", MyProvider, err)
		return err
	}

	output, err := instanceService.DescribeInstances(
		&qcservice.DescribeInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("DescribeInstances to %s failed: %v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DescribeInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send DescribeInstances to %s failed: %s", MyProvider, message)
	}

	if len(output.InstanceSet) == 0 {
		logger.Errorf("Send DescribeInstances to %s failed with 0 output instances", MyProvider)
		return fmt.Errorf("send DescribeInstances to %s failed with 0 output instances", MyProvider)
	}

	outputInstance := output.InstanceSet[0]
	instance.PrivateIp = qcservice.StringValue(outputInstance.PrivateIP)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = directive

	return nil
}

func (p *ProviderHandler) WaitInstanceTask(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", MyProvider, instance.TargetJobId, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitVolumeTask(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, volume.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", MyProvider, volume.TargetJobId, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitStopInstances(task *models.Task) error {
	return p.WaitInstanceTask(task)
}

func (p *ProviderHandler) WaitStartInstances(task *models.Task) error {
	return p.WaitInstanceTask(task)
}

func (p *ProviderHandler) WaitDeleteInstances(task *models.Task) error {
	return p.WaitInstanceTask(task)
}

func (p *ProviderHandler) WaitCreateVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) WaitAttachVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) WaitDetachVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) WaitDeleteVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) DescribeSubnet(runtimeId, subnetId string) (*models.Subnet, error) {
	qingcloudService, err := p.initService(runtimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return nil, err
	}

	vxnetService, err := qingcloudService.VxNet(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", MyProvider, err)
		return nil, err
	}

	output, err := vxnetService.DescribeVxNets(
		&qcservice.DescribeVxNetsInput{
			VxNets: qcservice.StringSlice([]string{subnetId}),
		},
	)
	if err != nil {
		logger.Errorf("DescribeVxNets to %s failed: %v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DescribeVxNets to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeVxNets to %s failed: %s", MyProvider, message)
	}

	if len(output.VxNetSet) == 0 {
		logger.Errorf("Send DescribeVxNets to %s failed with 0 output instances", MyProvider)
		return nil, fmt.Errorf("send DescribeVxNets to %s failed with 0 output instances", MyProvider)
	}

	vxnet := output.VxNetSet[0]
	return &models.Subnet{
		SubnetId:    qcservice.StringValue(vxnet.VxNetID),
		Name:        qcservice.StringValue(vxnet.VxNetName),
		CreateTime:  qcservice.TimeValue(vxnet.CreateTime),
		Description: qcservice.StringValue(vxnet.Description),
		InstanceIds: qcservice.StringValueSlice(vxnet.InstanceIDs),
		VpcId:       qcservice.StringValue(vxnet.VpcRouterID),
	}, nil
}

func (p *ProviderHandler) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	qingcloudService, err := p.initService(runtimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return nil, err
	}

	routerService, err := qingcloudService.Router(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", MyProvider, err)
		return nil, err
	}

	output, err := routerService.DescribeRouters(
		&qcservice.DescribeRoutersInput{
			Routers: qcservice.StringSlice([]string{vpcId}),
		},
	)
	if err != nil {
		logger.Errorf("DescribeRouters to %s failed: %v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DescribeRouters to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeRouters to %s failed: %s", MyProvider, message)
	}

	if len(output.RouterSet) == 0 {
		logger.Errorf("Send DescribeRouters to %s failed with 0 output instances", MyProvider)
		return nil, fmt.Errorf("send DescribeRouters to %s failed with 0 output instances", MyProvider)
	}

	vpc := output.RouterSet[0]

	var subnets []string
	for _, subnet := range vpc.VxNets {
		subnets = append(subnets, qcservice.StringValue(subnet.VxNetID))
	}

	return &models.Vpc{
		VpcId:            qcservice.StringValue(vpc.RouterID),
		Name:             qcservice.StringValue(vpc.RouterName),
		CreateTime:       qcservice.TimeValue(vpc.CreateTime),
		Description:      qcservice.StringValue(vpc.Description),
		Status:           qcservice.StringValue(vpc.Status),
		TransitionStatus: qcservice.StringValue(vpc.TransitionStatus),
		Subnets:          subnets,
		Eip: &models.Eip{
			EipId: qcservice.StringValue(vpc.EIP.EIPID),
			Name:  qcservice.StringValue(vpc.EIP.EIPName),
			Addr:  qcservice.StringValue(vpc.EIP.EIPAddr),
		},
	}, nil
}

func (p *ProviderHandler) DescribeZones(url, credential string) ([]string, error) {
	qingcloudService, err := p.initQingCloudService(url, credential, "")
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", MyProvider, err)
		return nil, err
	}

	output, err := qingcloudService.DescribeZones(
		&qcservice.DescribeZonesInput{
			Status: qcservice.StringSlice([]string{"active"}),
		},
	)
	if err != nil {
		logger.Errorf("DescribeZones to %s failed: %v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DescribeZones to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeZones to %s failed: %s", MyProvider, message)
	}

	var zones []string
	for _, zone := range output.ZoneSet {
		zones = append(zones, qcservice.StringValue(zone.ZoneID))
	}
	return zones, nil
}
