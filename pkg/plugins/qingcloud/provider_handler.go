// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

var MyProvider = constants.ProviderQingCloud

type ProviderHandler struct {
	vmbased.FrameHandler
}

func (p *ProviderHandler) initQingCloudService(runtimeUrl, runtimeCredential, zone string) (*qcservice.QingCloudService, error) {
	credential := new(Credential)
	err := jsonutil.Decode([]byte(runtimeCredential), credential)
	if err != nil {
		logger.Error("Parse [%s] credential failed: %+v", MyProvider, err)
		return nil, err
	}
	conf, err := qcconfig.New(credential.AccessKeyId, credential.SecretAccessKey)
	if err != nil {
		return nil, err
	}
	conf.Zone = zone
	if strings.HasPrefix(runtimeUrl, "https://") {
		runtimeUrl = strings.Split(runtimeUrl, "https://")[1]
	} else if strings.HasPrefix(runtimeUrl, "http://") {
		runtimeUrl = strings.Split(runtimeUrl, "http://")[1]
	}
	urlAndPort := strings.Split(runtimeUrl, ":")
	if len(urlAndPort) == 2 {
		conf.Port, err = strconv.Atoi(urlAndPort[1])
	}
	conf.Host = urlAndPort[0]
	if err != nil {
		logger.Error("Parse [%s] runtimeUrl [%s] failed: %+v", MyProvider, runtimeUrl, err)
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

func (p *ProviderHandler) waitInstanceNetworkAndVolume(instanceService *qcservice.InstanceService, instanceId string, needVolume bool, timeout time.Duration, waitInterval time.Duration) (ins *qcservice.Instance, err error) {
	logger.Debug("Waiting for IP address to be assigned and volume attached to Instance [%s]", instanceId)
	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		describeOutput, err := instanceService.DescribeInstances(
			&qcservice.DescribeInstancesInput{
				Instances: qcservice.StringSlice([]string{instanceId}),
			},
		)
		if err != nil {
			return false, err
		}

		describeRetCode := qcservice.IntValue(describeOutput.RetCode)
		if describeRetCode != 0 {
			return false, err
		}
		if len(describeOutput.InstanceSet) == 0 {
			return false, fmt.Errorf("Instance with id [%s] not exist", instanceId)
		}
		instance := describeOutput.InstanceSet[0]
		if len(instance.VxNets) == 0 || instance.VxNets[0].PrivateIP == nil || *instance.VxNets[0].PrivateIP == "" {
			return false, nil
		}
		if needVolume {
			if len(instance.Volumes) == 0 || instance.Volumes[0].Device == nil || *instance.Volumes[0].Device == "" {
				return false, nil
			}
		}
		ins = instance
		logger.Debug("Instance [%s] get IP address [%s]", instanceId, *ins.VxNets[0].PrivateIP)
		return true, nil
	}, timeout, waitInterval)
	return
}

func (p *ProviderHandler) RunInstances(task *models.Task) error {

	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s instance api service failed: %+v", MyProvider, err)
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
		NeedUserdata:  qcservice.Int(instance.NeedUserData),
		Hostname:      qcservice.String(instance.Hostname),
		// GPU:     qcservice.Int(instance.Gpu),
	}
	if instance.VolumeId != "" {
		input.Volumes = qcservice.StringSlice([]string{instance.VolumeId})
	}
	if instance.UserdataFile != "" {
		input.UserdataFile = qcservice.String(instance.UserdataFile)
	}
	if instance.UserDataValue != "" {
		input.UserdataValue = qcservice.String(instance.UserDataValue)
		input.UserdataType = qcservice.String(DefaultUserDataType)
	}
	logger.Debug("RunInstances with input: %s", jsonutil.ToString(input))
	output, err := instanceService.RunInstances(input)
	if err != nil {
		logger.Error("Send RunInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send RunInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send RunInstances to %s failed: %s", MyProvider, message)
	}

	if len(output.Instances) == 0 {
		logger.Error("Send RunInstances to %s failed with 0 output instances", MyProvider)
		return fmt.Errorf("send RunInstances to %s failed with 0 output instances", MyProvider)
	}

	instance.InstanceId = qcservice.StringValue(output.Instances[0])
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StopInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn("Skip task [%s] without instance", task.TaskId)
		return nil
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s instance api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&qcservice.DescribeInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	describeRetCode := qcservice.IntValue(describeOutput.RetCode)
	if describeRetCode != 0 {
		message := qcservice.StringValue(describeOutput.Message)
		logger.Error("Send DescribeInstances to %s failed with return code [%d], message [%s]",
			MyProvider, describeRetCode, message)
		return fmt.Errorf("send DescribeInstances to %s failed: %s", MyProvider, message)
	}
	if len(describeOutput.InstanceSet) == 0 {
		return fmt.Errorf("Instance with id [%s] not exist", instance.InstanceId)
	}

	status := qcservice.StringValue(describeOutput.InstanceSet[0].Status)

	if status == constants.StatusStopped {
		logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	output, err := instanceService.StopInstances(
		&qcservice.StopInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send StopInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send StopInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send StopInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StartInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn("Skip task [%s] without instance", task.TaskId)
		return nil
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s instance api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&qcservice.DescribeInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	describeRetCode := qcservice.IntValue(describeOutput.RetCode)
	if describeRetCode != 0 {
		message := qcservice.StringValue(describeOutput.Message)
		logger.Error("Send DescribeInstances to %s failed with return code [%d], message [%s]",
			MyProvider, describeRetCode, message)
		return fmt.Errorf("send DescribeInstances to %s failed: %s", MyProvider, message)
	}
	if len(describeOutput.InstanceSet) == 0 {
		return fmt.Errorf("Instance with id [%s] not exist", instance.InstanceId)
	}

	status := qcservice.StringValue(describeOutput.InstanceSet[0].Status)

	if status == constants.StatusActive {
		logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	output, err := instanceService.StartInstances(
		&qcservice.StartInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send StartInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send StartInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send StartInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) DeleteInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn("Skip task [%s] without instance", task.TaskId)
		return nil
	}

	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s instance api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&qcservice.DescribeInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	describeRetCode := qcservice.IntValue(describeOutput.RetCode)
	if describeRetCode != 0 {
		message := qcservice.StringValue(describeOutput.Message)
		logger.Error("Send DescribeInstances to %s failed with return code [%d], message [%s]",
			MyProvider, describeRetCode, message)
		return fmt.Errorf("send DescribeInstances to %s failed: %s", MyProvider, message)
	}
	if len(describeOutput.InstanceSet) == 0 {
		return fmt.Errorf("Instance with id [%s] not exist", instance.InstanceId)
	}

	status := qcservice.StringValue(describeOutput.InstanceSet[0].Status)

	if status == constants.StatusDeleted || status == constants.StatusCeased {
		logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	output, err := instanceService.TerminateInstances(
		&qcservice.TerminateInstancesInput{
			Instances: qcservice.StringSlice([]string{instance.InstanceId}),
		},
	)
	if err != nil {
		logger.Error("Send TerminateInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send TerminateInstances to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send TerminateInstances to %s failed: %s", MyProvider, message)
	}
	instance.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(instance)
	return nil
}

func (p *ProviderHandler) CreateVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s volume api service failed: %+v", MyProvider, err)
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
		logger.Error("Send CreateVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send CreateVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send CreateVolumes to %s failed: %s", MyProvider, message)
	}
	volume.VolumeId = qcservice.StringValue(output.Volumes[0])
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DetachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s volume api service failed: %+v", MyProvider, err)
		return err
	}

	output, err := volumeService.DetachVolumes(
		&qcservice.DetachVolumesInput{
			Instance: qcservice.String(volume.InstanceId),
			Volumes:  qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Error("Send DetachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send DetachVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send DetachVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) AttachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}

	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s volume api service failed: %+v", MyProvider, err)
		return err
	}

	output, err := volumeService.AttachVolumes(
		&qcservice.AttachVolumesInput{
			Instance: qcservice.String(volume.InstanceId),
			Volumes:  qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Error("Send AttachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send AttachVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send AttachVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}

	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn("Skip task [%s] without volume", task.TaskId)
		return nil
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	volumeService, err := qingcloudService.Volume(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s volume api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := volumeService.DescribeVolumes(
		&qcservice.DescribeVolumesInput{
			Volumes: qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Error("Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	describeRetCode := qcservice.IntValue(describeOutput.RetCode)
	if describeRetCode != 0 {
		message := qcservice.StringValue(describeOutput.Message)
		logger.Error("Send DescribeVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, describeRetCode, message)
		return fmt.Errorf("send DescribeVolumes to %s failed: %s", MyProvider, message)
	}
	if len(describeOutput.VolumeSet) == 0 {
		return fmt.Errorf("Volume with id [%s] not exist", volume.VolumeId)
	}

	status := qcservice.StringValue(describeOutput.VolumeSet[0].Status)

	if status == constants.StatusDeleted || status == constants.StatusCeased {
		logger.Warn("Volume [%s] has already been [%s], do nothing", volume.VolumeId, status)
		return nil
	}

	output, err := volumeService.DeleteVolumes(
		&qcservice.DeleteVolumesInput{
			Volumes: qcservice.StringSlice([]string{volume.VolumeId}),
		},
	)
	if err != nil {
		logger.Error("Send DeleteVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send DeleteVolumes to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return fmt.Errorf("send DeleteVolumes to %s failed: %s", MyProvider, message)
	}
	volume.TargetJobId = qcservice.StringValue(output.JobID)

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) WaitRunInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.TargetJobId == "" {
		logger.Warn("Skip task [%s] without target job id", task.TaskId)
		return nil
	}

	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s job api service failed: %+v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Error("Wait %s job [%s] failed: %+v", MyProvider, instance.TargetJobId, err)
		return err
	}

	instanceService, err := qingcloudService.Instance(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s instance api service failed: %+v", MyProvider, err)
		return err
	}

	needVolume := false
	if instance.VolumeId != "" {
		needVolume = true
	}

	output, err := p.waitInstanceNetworkAndVolume(instanceService, instance.InstanceId, needVolume,
		task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		logger.Error("Wait %s instance [%s] network failed: %+v", MyProvider, instance.InstanceId, err)
		return err
	}

	instance.PrivateIp = qcservice.StringValue(output.VxNets[0].PrivateIP)
	if len(output.Volumes) > 0 {
		instance.Device = qcservice.StringValue(output.Volumes[0].Device)
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	logger.Debug("WaitRunInstances task [%s] directive: %s", task.TaskId, task.Directive)

	return nil
}

func (p *ProviderHandler) WaitInstanceTask(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.TargetJobId == "" {
		logger.Warn("Skip task [%s] without target job id", task.TaskId)
		return nil
	}
	qingcloudService, err := p.initService(instance.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s job api service failed: %+v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, instance.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Error("Wait %s job [%s] failed: %+v", MyProvider, instance.TargetJobId, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitVolumeTask(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn("Skip task [%s] without directive", task.TaskId)
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.TargetJobId == "" {
		logger.Warn("Skip task [%s] without target job id", task.TaskId)
		return nil
	}
	qingcloudService, err := p.initService(volume.RuntimeId)
	if err != nil {
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	jobService, err := qingcloudService.Job(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s job api service failed: %+v", MyProvider, err)
		return err
	}

	err = qcclient.WaitJob(jobService, volume.TargetJobId, task.GetTimeout(constants.WaitTaskTimeout),
		constants.WaitTaskInterval)
	if err != nil {
		logger.Error("Wait %s job [%s] failed: %+v", MyProvider, volume.TargetJobId, err)
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
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	vxnetService, err := qingcloudService.VxNet(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s job api service failed: %+v", MyProvider, err)
		return nil, err
	}

	output, err := vxnetService.DescribeVxNets(
		&qcservice.DescribeVxNetsInput{
			VxNets: qcservice.StringSlice([]string{subnetId}),
		},
	)
	if err != nil {
		logger.Error("DescribeVxNets to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send DescribeVxNets to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeVxNets to %s failed: %s", MyProvider, message)
	}

	if len(output.VxNetSet) == 0 {
		logger.Error("Send DescribeVxNets to %s failed with 0 output instances", MyProvider)
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
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	routerService, err := qingcloudService.Router(qingcloudService.Config.Zone)
	if err != nil {
		logger.Error("Init %s job api service failed: %+v", MyProvider, err)
		return nil, err
	}

	output, err := routerService.DescribeRouters(
		&qcservice.DescribeRoutersInput{
			Routers: qcservice.StringSlice([]string{vpcId}),
		},
	)
	if err != nil {
		logger.Error("DescribeRouters to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send DescribeRouters to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeRouters to %s failed: %s", MyProvider, message)
	}

	if len(output.RouterSet) == 0 {
		logger.Error("Send DescribeRouters to %s failed with 0 output instances", MyProvider)
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
		logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	output, err := qingcloudService.DescribeZones(
		&qcservice.DescribeZonesInput{
			Status: qcservice.StringSlice([]string{"active"}),
		},
	)
	if err != nil {
		logger.Error("DescribeZones to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	retCode := qcservice.IntValue(output.RetCode)
	if retCode != 0 {
		message := qcservice.StringValue(output.Message)
		logger.Error("Send DescribeZones to %s failed with return code [%d], message [%s]",
			MyProvider, retCode, message)
		return nil, fmt.Errorf("send DescribeZones to %s failed: %s", MyProvider, message)
	}

	var zones []string
	for _, zone := range output.ZoneSet {
		zones = append(zones, qcservice.StringValue(zone.ZoneID))
	}
	return zones, nil
}
