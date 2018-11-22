// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package aliyun

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var MyProvider = constants.ProviderAliyun
var DevicePattern = regexp.MustCompile("/dev/xvd(.)")

type ProviderHandler struct {
	vmbased.FrameHandler
}

func GetProviderHandler(ctx context.Context) *ProviderHandler {
	providerHandler := new(ProviderHandler)
	providerHandler.Ctx = ctx
	return providerHandler
}

func (p *ProviderHandler) initInstanceService(runtimeId string) (*ecs.Client, error) {
	runtime, err := runtimeclient.NewRuntime(p.Ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return p.initInstanceServiceWithCredential(runtime.RuntimeUrl, runtime.Credential, runtime.Zone)
}

func (p *ProviderHandler) initInstanceServiceWithCredential(runtimeUrl, runtimeCredential, zone string) (*ecs.Client, error) {
	credential := new(vmbased.Credential)
	err := jsonutil.Decode([]byte(runtimeCredential), credential)
	if err != nil {
		logger.Error(p.Ctx, "Parse [%s] credential failed: %+v", MyProvider, err)
		return nil, err
	}

	ecsClient, err := ecs.NewClientWithAccessKey(zone, credential.AccessKeyId, credential.SecretAccessKey)
	if err != nil {
		return nil, err
	}

	return ecsClient, nil
}

func (p *ProviderHandler) RunInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceType, err := ConvertToInstanceType(instance.Cpu, instance.Memory)
	if err != nil {
		logger.Error(p.Ctx, "Could not find an aliyun instance type: %+v", err)
		return err
	}

	logger.Info(p.Ctx, "RunInstances with name [%s] instance type [%s]", instance.Name, instanceType)

	input := ecs.CreateCreateInstanceRequest()
	input.InstanceName = instance.Name
	input.ImageId = instance.ImageId
	input.InstanceType = instanceType
	input.VSwitchId = instance.Subnet
	input.ZoneId = instance.Zone
	input.Password = DefaultLoginPassword

	if instance.NeedUserData == 1 {
		input.UserData = instance.UserDataValue
	}

	logger.Debug(p.Ctx, "RunInstances with input: %s", jsonutil.ToString(input))
	output, err := instanceService.CreateInstance(input)
	if err != nil {
		logger.Error(p.Ctx, "Send RunInstances to %s failed: %+v", MyProvider, err)
		return err
	}
	logger.Debug(p.Ctx, "RunInstances get output: %s", jsonutil.ToString(output))

	instance.InstanceId = output.InstanceId

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StopInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeInstancesRequest()
	describeInput.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
	describeOutput, err := instanceService.DescribeInstances(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Instances.Instance) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := describeOutput.Instances.Instance[0].Status

	if status == strings.Title(constants.StatusStopped) {
		logger.Warn(p.Ctx, "Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	logger.Info(p.Ctx, "StopInstances [%s]", instance.Name)

	input := ecs.CreateStopInstanceRequest()
	input.InstanceId = instance.InstanceId

	_, err = instanceService.StopInstance(input)
	if err != nil {
		logger.Error(p.Ctx, "Send StopInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StartInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeInstancesRequest()
	describeInput.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
	describeOutput, err := instanceService.DescribeInstances(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Instances.Instance) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := describeOutput.Instances.Instance[0].Status

	if status == strings.Title(constants.StatusRunning) {
		logger.Warn(p.Ctx, "Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	logger.Info(p.Ctx, "StartInstances [%s]", instance.Name)

	input := ecs.CreateStartInstanceRequest()
	input.InstanceId = instance.InstanceId

	_, err = instanceService.StartInstance(input)
	if err != nil {
		logger.Error(p.Ctx, "Send StartInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) DeleteInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeInstancesRequest()
	describeInput.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
	describeOutput, err := instanceService.DescribeInstances(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Instances.Instance) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := describeOutput.Instances.Instance[0].Status
	if status == strings.Title(constants.StatusRunning) {
		logger.Info(p.Ctx, "StopInstance [%s] before delete it", instance.Name)
		err := p.StopInstances(task)
		if err != nil {
			logger.Error(p.Ctx, "Send StopInstances to %s failed: %+v", MyProvider, err)
			return err
		}

		err = p.WaitStopInstances(task)
		if err != nil {
			return err
		}
	}

	logger.Info(p.Ctx, "DeleteInstance [%s]", instance.Name)

	input := ecs.CreateDeleteInstanceRequest()
	input.InstanceId = instance.InstanceId

	_, err = instanceService.DeleteInstance(input)
	if err != nil {
		logger.Error(p.Ctx, "Send DeleteInstance to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) ResizeInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeInstancesRequest()
	describeInput.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
	describeOutput, err := instanceService.DescribeInstances(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Instances.Instance) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := describeOutput.Instances.Instance[0].Status

	if status != strings.Title(constants.StatusStopped) {
		logger.Warn(p.Ctx, "Instance [%s] is in status [%s], can not resize", instance.InstanceId, status)
		return fmt.Errorf("instance [%s] is in status [%s], can not resize", instance.InstanceId, status)
	}

	instanceType, err := ConvertToInstanceType(instance.Cpu, instance.Memory)
	if err != nil {
		logger.Error(p.Ctx, "Could not find an aliyun instance type: %+v", err)
		return err
	}

	logger.Info(p.Ctx, "ResizeInstances [%s] with instance type [%s]", instance.Name, instanceType)

	input := ecs.CreateModifyInstanceSpecRequest()
	input.InstanceId = instance.InstanceId
	input.InstanceType = instanceType

	_, err = instanceService.ModifyInstanceSpec(input)
	if err != nil {
		logger.Error(p.Ctx, "Send ResizeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)
	return nil
}

func (p *ProviderHandler) CreateVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	volumeType, err := ConvertToVolumeType(DefaultVolumeClass)
	if err != nil {
		return err
	}

	logger.Info(p.Ctx, "CreateVolumes with name [%s] volume type [%s] size [%d]", volume.Name, volumeType, volume.Size)

	input := ecs.CreateCreateDiskRequest()
	input.DiskName = volume.Name
	input.ZoneId = volume.Zone
	input.Size = requests.NewInteger(volume.Size)
	input.DiskCategory = volumeType

	output, err := instanceService.CreateDisk(input)
	if err != nil {
		logger.Error(p.Ctx, "Send CreateVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	volume.VolumeId = output.DiskId

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DetachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume id")
		return nil
	}
	if volume.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeDisksRequest()
	describeInput.DiskIds = fmt.Sprintf("[\"%s\"]", volume.VolumeId)
	describeOutput, err := instanceService.DescribeDisks(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Disks.Disk) == 0 {
		return fmt.Errorf("volume with id [%s] not exist", volume.VolumeId)
	}

	status := describeOutput.Disks.Disk[0].Status

	if status == strings.Title(constants.StatusAvailable) {
		logger.Warn(p.Ctx, "Volume [%s] is in status [%s], no need to detach.", volume.VolumeId, status)
		return nil
	}

	logger.Info(p.Ctx, "DetachVolume [%s] from instance with id [%s]", volume.Name, volume.InstanceId)

	input := ecs.CreateDetachDiskRequest()
	input.InstanceId = volume.InstanceId
	input.DiskId = volume.VolumeId

	_, err = instanceService.DetachDisk(input)
	if err != nil {
		logger.Error(p.Ctx, "Send DetachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) AttachVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume id")
		return nil
	}
	if volume.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	logger.Info(p.Ctx, "AttachVolume [%s] to instance with id [%s]", volume.Name, volume.InstanceId)

	input := ecs.CreateAttachDiskRequest()
	input.InstanceId = volume.InstanceId
	input.DiskId = volume.VolumeId

	_, err = instanceService.AttachDisk(input)
	if err != nil {
		logger.Error(p.Ctx, "Send AttachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume id")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeDisksRequest()
	describeInput.DiskIds = fmt.Sprintf("[\"%s\"]", volume.VolumeId)
	describeOutput, err := instanceService.DescribeDisks(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Disks.Disk) == 0 {
		return fmt.Errorf("volume with id [%s] not exist", volume.VolumeId)
	}

	disk := describeOutput.Disks.Disk[0]

	logger.Info(p.Ctx, "DeleteVolume [%s] with status [%s]", volume.Name, disk.Status)
	if disk.Status == strings.Title(constants.StatusInUse2) {
		err := p.WaitVolumeState(task, strings.Title(constants.StatusAvailable))
		if err != nil {
			return err
		}
	}

	input := ecs.CreateDeleteDiskRequest()
	input.DiskId = volume.VolumeId

	_, err = instanceService.DeleteDisk(input)
	if err != nil {
		logger.Error(p.Ctx, "Send DeleteVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) ResizeVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeInput := ecs.CreateDescribeDisksRequest()
	describeInput.DiskIds = fmt.Sprintf("[\"%s\"]", volume.VolumeId)
	describeOutput, err := instanceService.DescribeDisks(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Disks.Disk) == 0 {
		return fmt.Errorf("volume with id [%s] not exist", volume.VolumeId)
	}

	status := describeOutput.Disks.Disk[0].Status
	if status != strings.Title(constants.StatusAvailable) {
		logger.Warn(p.Ctx, "Volume [%s] is in status [%s], can not resize.", volume.VolumeId, status)
		return fmt.Errorf("volume [%s] is in status [%s], can not resize", volume.VolumeId, status)
	}

	logger.Info(p.Ctx, "ResizeVolumes [%s] with size [%d]", volume.Name, volume.Size)

	input := ecs.CreateResizeDiskRequest()
	input.DiskId = volume.VolumeId
	input.NewSize = requests.NewInteger(volume.Size)

	_, err = instanceService.ResizeDisk(input)
	if err != nil {
		logger.Error(p.Ctx, "Send ResizeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)
	return nil
}

func (p *ProviderHandler) waitInstanceVolume(instanceService *ecs.Client, task *models.Task, instance *models.Instance) error {
	logger.Debug(p.Ctx, "Waiting for volume [%s] attached to Instance [%s]", instance.VolumeId, instance.InstanceId)

	err := p.AttachVolumes(task)
	if err != nil {
		logger.Error(p.Ctx, "Attach volume [%s] to Instance [%s] failed: %+v", instance.VolumeId, instance.InstanceId, err)
		return err
	}

	err = p.WaitAttachVolumes(task)
	if err != nil {
		logger.Error(p.Ctx, "Waiting for volume [%s] attached to Instance [%s] failed: %+v", instance.VolumeId, instance.InstanceId, err)
		return err
	}

	describeInput := ecs.CreateDescribeDisksRequest()
	describeInput.DiskIds = fmt.Sprintf("[\"%s\"]", instance.VolumeId)
	describeOutput, err := instanceService.DescribeDisks(describeInput)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Disks.Disk) == 0 {
		return fmt.Errorf("volume with id [%s] not exist", instance.VolumeId)
	}

	vol := describeOutput.Disks.Disk[0]
	instance.Device = vol.Device

	describeInput2 := ecs.CreateDescribeInstancesRequest()
	describeInput2.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
	describeOutput2, err := instanceService.DescribeInstances(describeInput2)
	if err != nil {
		logger.Error(p.Ctx, "Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput2.Instances.Instance) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	ins := describeOutput2.Instances.Instance[0]

	if ins.IoOptimized {
		instance.Device = DevicePattern.ReplaceAllString(instance.Device, "/dev/vd$1")
	}

	logger.Info(p.Ctx, "Instance [%s] with io optimized [%t] attached volume [%s] as device [%s]", instance.InstanceId, ins.IoOptimized, instance.VolumeId, instance.Device)
	return nil
}

func (p *ProviderHandler) waitInstanceNetwork(instanceService *ecs.Client, instance *models.Instance, timeout time.Duration, waitInterval time.Duration) error {
	err := funcutil.WaitForSpecificOrError(func() (bool, error) {
		describeInput := ecs.CreateDescribeInstancesRequest()
		describeInput.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
		describeOutput, err := instanceService.DescribeInstances(describeInput)
		if err != nil {
			return false, err
		}

		if len(describeOutput.Instances.Instance) == 0 {
			return false, fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
		}

		ins := describeOutput.Instances.Instance[0]

		if len(ins.VpcAttributes.PrivateIpAddress.IpAddress) == 0 {
			return false, nil
		}

		instance.PrivateIp = ins.VpcAttributes.PrivateIpAddress.IpAddress[0]
		instance.Eip = ins.EipAddress.IpAddress
		return true, nil
	}, timeout, waitInterval)

	logger.Info(p.Ctx, "Instance [%s] get private IP address [%s]", instance.InstanceId, instance.PrivateIp)

	if instance.Eip != "" {
		logger.Info(p.Ctx, "Instance [%s] get EIP address [%s]", instance.InstanceId, instance.Eip)
	}

	return err
}

func (p *ProviderHandler) WaitRunInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = p.WaitInstanceState(task, strings.Title(constants.StatusStopped))
	if err != nil {
		logger.Error(p.Ctx, "Wait %s job [%s] failed: %+v", MyProvider, instance.TargetJobId, err)
		return err
	}

	if instance.VolumeId != "" {
		err := p.waitInstanceVolume(instanceService, task, instance)
		if err != nil {
			logger.Error(p.Ctx, "Attach volume [%s] to Instance [%s] failed: %+v", instance.VolumeId, instance.InstanceId, err)
			return err
		}
	}

	err = p.waitInstanceNetwork(instanceService, instance, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		logger.Error(p.Ctx, "Wait %s instance [%s] network failed: %+v", MyProvider, instance.InstanceId, err)
		return err
	}

	err = p.StartInstances(task)
	if err != nil {
		logger.Error(p.Ctx, "Start %s instance [%s] failed: %+v", MyProvider, instance.InstanceId, err)
		return err
	}

	err = p.WaitInstanceState(task, strings.Title(constants.StatusRunning))
	if err != nil {
		logger.Error(p.Ctx, "Wait %s job [%s] failed: %+v", MyProvider, instance.TargetJobId, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	logger.Debug(p.Ctx, "WaitRunInstances task [%s] directive: %s", task.TaskId, task.Directive)

	return nil
}

func (p *ProviderHandler) WaitInstanceState(task *models.Task, state string) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ecs.CreateDescribeInstancesRequest()
		input.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
		output, err := instanceService.DescribeInstances(input)
		if err != nil {
			return true, err
		}

		if len(output.Instances.Instance) == 0 {
			return true, fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
		}

		if output.Instances.Instance[0].Status == state {
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		logger.Error(p.Ctx, "Wait %s instance [%s] status become to [%s] failed: %+v", MyProvider, instance.InstanceId, state, err)
		return err
	}

	logger.Info(p.Ctx, "Wait %s instance [%s] status become to [%s] success", MyProvider, instance.InstanceId, state)

	return nil
}

func (p *ProviderHandler) WaitVolumeState(task *models.Task, state string) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume id")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ecs.CreateDescribeDisksRequest()
		input.DiskIds = fmt.Sprintf("[\"%s\"]", volume.VolumeId)

		output, err := instanceService.DescribeDisks(input)
		if err != nil {
			return true, err
		}

		if len(output.Disks.Disk) == 0 {
			return true, fmt.Errorf("volume [%s] not found", volume.VolumeId)
		}

		if output.Disks.Disk[0].Status == state {
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		logger.Error(p.Ctx, "Wait %s volume [%s] status become to [%s] failed: %+v", MyProvider, volume.VolumeId, state, err)
		return err
	}

	logger.Info(p.Ctx, "Wait %s volume [%s] status become to [%s] success", MyProvider, volume.VolumeId, state)

	return nil
}

func (p *ProviderHandler) WaitStopInstances(task *models.Task) error {
	return p.WaitInstanceState(task, strings.Title(constants.StatusStopped))
}

func (p *ProviderHandler) WaitStartInstances(task *models.Task) error {
	return p.WaitInstanceState(task, strings.Title(constants.StatusRunning))
}

func (p *ProviderHandler) WaitDeleteInstances(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		logger.Warn(p.Ctx, "Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ecs.CreateDescribeInstancesRequest()
		input.InstanceIds = fmt.Sprintf("[\"%s\"]", instance.InstanceId)
		output, err := instanceService.DescribeInstances(input)
		if err != nil {
			return true, err
		}

		if len(output.Instances.Instance) == 0 {
			logger.Info(p.Ctx, "Wait %s instance [%s] to be deleted successfully", MyProvider, instance.InstanceId)
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	return err
}

func (p *ProviderHandler) WaitResizeInstances(task *models.Task) error {
	return p.WaitInstanceState(task, strings.Title(constants.StatusStopped))
}

func (p *ProviderHandler) WaitCreateVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, strings.Title(constants.StatusAvailable))
}

func (p *ProviderHandler) WaitAttachVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, strings.Title(constants.StatusInUse2))
}

func (p *ProviderHandler) WaitDetachVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, strings.Title(constants.StatusAvailable))
}

func (p *ProviderHandler) WaitDeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		logger.Warn(p.Ctx, "Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		logger.Warn(p.Ctx, "Skip task without volume")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ecs.CreateDescribeDisksRequest()
		input.DiskIds = fmt.Sprintf("[\"%s\"]", volume.VolumeId)
		output, err := instanceService.DescribeDisks(input)
		if err != nil {
			return true, err
		}

		if len(output.Disks.Disk) == 0 {
			logger.Info(p.Ctx, "Wait %s volume [%s] to be deleted successfully", MyProvider, volume.VolumeId)
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	return err
}

func (p *ProviderHandler) WaitResizeVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, strings.Title(constants.StatusAvailable))
}

func (p *ProviderHandler) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	instanceService, err := p.initInstanceService(req.GetRuntimeId().GetValue())
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	input := ecs.CreateDescribeVSwitchesRequest()

	if len(req.GetZone()) == 1 {
		input.ZoneId = req.GetZone()[0]
	}

	if len(req.GetSubnetId()) > 0 {
		input.VSwitchId = strings.Join(req.GetSubnetId(), ",")
	}

	output, err := instanceService.DescribeVSwitches(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeVSwitches to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	if len(output.VSwitches.VSwitch) == 0 {
		logger.Error(p.Ctx, "Send DescribeVSwitches to %s failed with 0 output subnets", MyProvider)
		return nil, fmt.Errorf("send DescribeVSwitches to %s failed with 0 output subnets", MyProvider)
	}

	response := new(pb.DescribeSubnetsResponse)

	for _, vs := range output.VSwitches.VSwitch {
		subnet := &pb.Subnet{
			SubnetId: pbutil.ToProtoString(vs.VSwitchId),
			Name:     pbutil.ToProtoString(vs.VSwitchName),
			VpcId:    pbutil.ToProtoString(vs.VpcId),
			Zone:     pbutil.ToProtoString(vs.ZoneId),
		}
		response.SubnetSet = append(response.SubnetSet, subnet)
	}

	response.TotalCount = uint32(len(response.SubnetSet))

	return response, nil
}

func (p *ProviderHandler) CheckResourceQuotas(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	roleCount := make(map[string]int)
	for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
		role := clusterNode.Role
		_, isExist := roleCount[role]
		if isExist {
			roleCount[role] = roleCount[role] + 1
		} else {
			roleCount[role] = 1
		}
	}

	return nil
}

func (p *ProviderHandler) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	instanceService, err := p.initInstanceService(runtimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	input := ecs.CreateDescribeVpcsRequest()
	input.VpcId = vpcId

	output, err := instanceService.DescribeVpcs(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeVpcs to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	if len(output.Vpcs.Vpc) == 0 {
		logger.Error(p.Ctx, "Send DescribeVpcs to %s failed with 0 output instances", MyProvider)
		return nil, fmt.Errorf("send DescribeVpcs to %s failed with 0 output instances", MyProvider)
	}

	vpc := output.Vpcs.Vpc[0]

	return &models.Vpc{
		VpcId:   vpc.VpcId,
		Name:    vpc.VpcName,
		Status:  vpc.Status,
		Subnets: vpc.VSwitchIds.VSwitchId,
	}, nil
}

func (p *ProviderHandler) DescribeZones(url, credential string) ([]string, error) {
	zone := DefaultZone
	instanceService, err := p.initInstanceServiceWithCredential(url, credential, zone)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	input := ecs.CreateDescribeRegionsRequest()
	output, err := instanceService.DescribeRegions(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeRegions to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var zones []string
	for _, zone := range output.Regions.Region {
		zones = append(zones, zone.RegionId)
	}
	return zones, nil
}

func (p *ProviderHandler) DescribeKeyPairs(url, credential, zone string) ([]string, error) {
	instanceService, err := p.initInstanceServiceWithCredential(url, credential, zone)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	input := ecs.CreateDescribeKeyPairsRequest()
	output, err := instanceService.DescribeKeyPairs(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeKeyPairs to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var keys []string
	for _, key := range output.KeyPairs.KeyPair {
		keys = append(keys, key.KeyPairName)
	}
	return keys, nil
}

func (p *ProviderHandler) DescribeImage(runtimeId, imageName string) (string, error) {
	instanceService, err := p.initInstanceService(runtimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return "", err
	}

	input := ecs.CreateDescribeImagesRequest()
	input.ImageName = imageName

	output, err := instanceService.DescribeImages(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeImages to %s failed: %+v", MyProvider, err)
		return "", err
	}

	if len(output.Images.Image) == 0 {
		return "", fmt.Errorf("image with name [%s] not exist", imageName)
	}

	imageId := output.Images.Image[0].ImageId

	return imageId, nil
}

func (p *ProviderHandler) DescribeAvailabilityZoneBySubnetId(runtimeId, subnetId string) (string, error) {
	instanceService, err := p.initInstanceService(runtimeId)
	if err != nil {
		logger.Error(p.Ctx, "Init %s api service failed: %+v", MyProvider, err)
		return "", err
	}

	input := ecs.CreateDescribeVSwitchesRequest()
	input.VSwitchId = subnetId
	output, err := instanceService.DescribeVSwitches(input)
	if err != nil {
		logger.Error(p.Ctx, "DescribeSubnets to %s failed: %+v", MyProvider, err)
		return "", err
	}

	if len(output.VSwitches.VSwitch) == 0 {
		return "", fmt.Errorf("subnet with id [%s] not exist", subnetId)
	}

	zone := output.VSwitches.VSwitch[0].ZoneId

	return zone, nil
}
