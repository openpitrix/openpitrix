// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"encoding/json"
	"fmt"
	"time"

	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
)

var provider = constants.ProviderQingCloud

type ProviderHandler struct {
	vmbased.FrameHandler
}

func init() {
	vmbased.RegisterProviderHandler(constants.ProviderQingCloud, new(ProviderHandler))
}

func (p *ProviderHandler) initService(runtimeId string) (*qcservice.QingCloudService, error) {
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	credential := new(Credential)
	err = json.Unmarshal([]byte(runtime.Credential), credential)
	if err != nil {
		logger.Errorf("Parse [%s] credential failed: %v", provider, err)
		return nil, err
	}
	conf, err := qcconfig.New(credential.AccessKeyId, credential.SecretAccessKey)
	if err != nil {
		return nil, err
	}
	conf.Zone = runtime.Zone
	conf.URI = runtime.RuntimeUrl
	return qcservice.Init(conf)
}

func (p *ProviderHandler) RunInstances(task *models.Task) error {

	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", provider, err)
		return err
	}

	output, err := instanceService.RunInstances(
		&qcservice.RunInstancesInput{
			ImageID:       qcservice.String(instance.ImageId),
			CPU:           qcservice.Int(instance.Cpu),
			Memory:        qcservice.Int(instance.Memory),
			InstanceName:  qcservice.String(instance.Name),
			InstanceClass: qcservice.Int(DefaultInstanceClass),
			Volumes:       qcservice.StringSlice([]string{instance.VolumeId}),
			VxNets:        qcservice.StringSlice([]string{instance.Subnet}),
			LoginMode:     qcservice.String(DefaultLoginMode),
			LoginPasswd:   qcservice.String(DefaultLoginPassword),
			UserdataValue: qcservice.String(instance.UserDataValue),
			UserdataPath:  qcservice.String(instance.UserdataPath),
			// GPU:     qcservice.Int(instance.Gpu),
		},
	)
	if err != nil {
		logger.Errorf("Send RunInstances to %s failed: %v", provider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send RunInstances to %s failed with return code [%d], message [%p]",
			provider, retCode, message)
		return fmt.Errorf("send RunInstances to %s failed: %p", provider, message)
	}

	if len(output.Instances) == 0 {
		logger.Errorf("Send RunInstances to %s failed with 0 output instances", provider)
		return fmt.Errorf("send RunInstances to %s failed with 0 output instances", provider)
	}

	instance.InstanceId = qcservice.StringValue(output.Instances[0])
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	directive, err := instance.ToString()
	if err != nil {
		return err
	}
	task.Directive = directive

	return nil
}

func (p *ProviderHandler) WaitRunInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", provider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout)*time.Second,
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", provider, instance.TargetJobId, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", provider, err)
		return err
	}

	output, err := instanceService.DescribeInstances(
		&qcservice.DescribeInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("DescribeInstances to %s failed: %v", provider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DescribeInstances to %s failed with return code [%d], message [%p]",
			provider, retCode, message)
		return fmt.Errorf("send DescribeInstances to %s failed: %p", provider, message)
	}

	if len(output.InstanceSet) == 0 {
		logger.Errorf("Send DescribeInstances to %s failed with 0 output instances", provider)
		return fmt.Errorf("send DescribeInstances to %s failed with 0 output instances", provider)
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
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", provider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout)*time.Second,
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", provider, instance.TargetJobId, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) CreateVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", provider, err)
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
		logger.Errorf("Send CreateVolumes to %s failed: %v", provider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send CreateVolumes to %s failed with return code [%d], message [%p]",
			provider, retCode, message)
		return fmt.Errorf("send CreateVolumes to %s failed: %p", provider, message)
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

func (p *ProviderHandler) WaitVolumeTask(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s job api service failed: %v", provider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, volume.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout)*time.Second,
		constants.WaitTaskInterval)
	if err != nil {
		logger.Errorf("Wait %s job [%s] failed: %v", provider, volume.TargetJobId, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitCreateVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) DetachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s volume api service failed: %v", provider, err)
		return err
	}

	output, err := volumeService.DetachVolumes(
		&qcservice.DetachVolumesInput{
			Instance: qcservice.String(volume.InstanceId),
			Volumes:  qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Errorf("Send DetachVolumes to %s failed: %v", provider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send DetachVolumes to %s failed with return code [%d], message [%p]",
			provider, retCode, message)
		return fmt.Errorf("send DetachVolumes to %s failed: %p", provider, message)
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

func (p *ProviderHandler) WaitDetachVolumes(task *models.Task) error {
	return p.WaitVolumeTask(task)
}

func (p *ProviderHandler) StopInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warnf("Skip empty task [%p] directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Errorf("Init %s api service failed: %v", provider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Errorf("Init %s instance api service failed: %v", provider, err)
		return err
	}

	output, err := instanceService.StopInstances(
		&qcservice.StopInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Errorf("Send StopInstances to %s failed: %v", provider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Errorf("Send RunInstances to %s failed with return code [%d], message [%p]",
			provider, retCode, message)
		return fmt.Errorf("send RunInstances to %s failed: %p", provider, message)
	}

	return nil
}

func (p *ProviderHandler) WaitStopInstances(task *models.Task) error {
	return p.WaitInstanceTask(task)
}
