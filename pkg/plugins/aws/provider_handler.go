// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

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

var MyProvider = constants.ProviderAWS

type ProviderHandler struct {
	vmbased.FrameHandler
}

func GetProviderHandler(Logger *logger.Logger) *ProviderHandler {
	providerHandler := new(ProviderHandler)
	if Logger == nil {
		providerHandler.Logger = logger.NewLogger()
	} else {
		providerHandler.Logger = Logger
	}
	return providerHandler
}

func (p *ProviderHandler) initAWSSession(runtimeUrl, runtimeCredential, zone string) (*session.Session, error) {
	credential := new(Credential)
	err := jsonutil.Decode([]byte(runtimeCredential), credential)
	if err != nil {
		p.Logger.Error("Parse [%s] credential failed: %+v", MyProvider, err)
		return nil, err
	}

	creds := credentials.NewStaticCredentials(credential.AccessKeyId, credential.SecretAccessKey, "")
	config := &aws.Config{
		Region:      aws.String(zone),
		Endpoint:    aws.String(""),
		Credentials: creds,
	}

	return session.NewSession(config)
}

func (p *ProviderHandler) initSession(runtimeId string) (*session.Session, error) {
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	return p.initAWSSession(runtime.RuntimeUrl, runtime.Credential, runtime.Zone)
}

func (p *ProviderHandler) initInstanceService(runtimeId string) (*ec2.EC2, error) {
	awsSession, err := p.initSession(runtimeId)
	if err != nil {
		p.Logger.Error("Init %s api session failed: %+v", MyProvider, err)
		return nil, err
	}

	return ec2.New(awsSession), nil
}

func (p *ProviderHandler) RunInstances(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	instanceType, err := ConvertToInstanceType(instance.Cpu, instance.Memory)
	if err != nil {
		p.Logger.Error("Could not find an aws instance type: %+v", err)
		return err
	}

	tag := ec2.Tag{
		Key:   aws.String("Name"),
		Value: aws.String(instance.Name),
	}

	tagSpec := ec2.TagSpecification{
		ResourceType: aws.String("instance"),
		Tags:         []*ec2.Tag{&tag},
	}

	input := ec2.RunInstancesInput{
		KeyName:           aws.String(DefaultKeyName),
		ImageId:           aws.String(instance.ImageId),
		InstanceType:      aws.String(instanceType),
		TagSpecifications: []*ec2.TagSpecification{&tagSpec},
		SubnetId:          aws.String(instance.Subnet),
		Placement:         &ec2.Placement{AvailabilityZone: aws.String(instance.Zone)},
		MaxCount:          aws.Int64(1),
		MinCount:          aws.Int64(1),
	}
	if instance.NeedUserData == 1 {
		input.UserData = aws.String(instance.UserDataValue)
	}

	p.Logger.Debug("RunInstances with input: %s", jsonutil.ToString(input))
	output, err := instanceService.RunInstances(&input)
	if err != nil {
		p.Logger.Error("Send RunInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(output.Instances) == 0 {
		p.Logger.Error("Send RunInstances to %s failed with 0 output instances", MyProvider)
		return fmt.Errorf("send RunInstances to %s failed with 0 output instances", MyProvider)
	}

	instance.InstanceId = aws.StringValue(output.Instances[0].InstanceId)

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StopInstances(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		p.Logger.Warn("Skip task without instance")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Reservations) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	if len(describeOutput.Reservations[0].Instances) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := aws.StringValue(describeOutput.Reservations[0].Instances[0].State.Name)

	if status == constants.StatusStopped {
		p.Logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	_, err = instanceService.StopInstances(
		&ec2.StopInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send StopInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) StartInstances(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		p.Logger.Warn("Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Reservations) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	if len(describeOutput.Reservations[0].Instances) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := aws.StringValue(describeOutput.Reservations[0].Instances[0].State.Name)

	if status == constants.StatusRunning {
		p.Logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	_, err = instanceService.StartInstances(
		&ec2.StartInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send StartInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) DeleteInstances(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	if instance.InstanceId == "" {
		p.Logger.Warn("Skip task without instance id")
		return nil
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeInstances(
		&ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send DescribeInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Reservations) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	if len(describeOutput.Reservations[0].Instances) == 0 {
		return fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
	}

	status := aws.StringValue(describeOutput.Reservations[0].Instances[0].State.Name)

	if status == constants.StatusTerminated {
		p.Logger.Warn("Instance [%s] has already been [%s], do nothing", instance.InstanceId, status)
		return nil
	}

	_, err = instanceService.TerminateInstances(
		&ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice([]string{instance.InstanceId}),
		})
	if err != nil {
		p.Logger.Error("Send TerminateInstances to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	return nil
}

func (p *ProviderHandler) CreateVolumes(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	tag := ec2.Tag{
		Key:   aws.String("Name"),
		Value: aws.String(volume.Name),
	}

	tagSpec := ec2.TagSpecification{
		ResourceType: aws.String("volume"),
		Tags:         []*ec2.Tag{&tag},
	}

	volumeType, err := ConvertToVolumeType(DefaultVolumeClass)
	if err != nil {
		return err
	}

	input := ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String(volume.Zone),
		Size:              aws.Int64(int64(volume.Size)),
		VolumeType:        aws.String(volumeType),
		TagSpecifications: []*ec2.TagSpecification{&tagSpec},
	}

	output, err := instanceService.CreateVolume(&input)
	if err != nil {
		p.Logger.Error("Send CreateVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	volume.VolumeId = aws.StringValue(output.VolumeId)

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DetachVolumes(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	_, err = instanceService.DetachVolume(
		&ec2.DetachVolumeInput{
			InstanceId: aws.String(volume.InstanceId),
			VolumeId:   aws.String(volume.VolumeId),
		})
	if err != nil {
		p.Logger.Error("Send DetachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) AttachVolumes(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	_, err = instanceService.AttachVolume(
		&ec2.AttachVolumeInput{
			InstanceId: aws.String(volume.InstanceId),
			VolumeId:   aws.String(volume.VolumeId),
			Device:     aws.String(DefaultDevice),
		})
	if err != nil {
		p.Logger.Error("Send AttachVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) DeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	if volume.VolumeId == "" {
		p.Logger.Warn("Skip task without volume")
		return nil
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	describeOutput, err := instanceService.DescribeVolumes(
		&ec2.DescribeVolumesInput{
			VolumeIds: aws.StringSlice([]string{volume.VolumeId}),
		})
	if err != nil {
		p.Logger.Error("Send DescribeVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	if len(describeOutput.Volumes) == 0 {
		return fmt.Errorf("volume with id [%s] not exist", volume.VolumeId)
	}

	_, err = instanceService.DeleteVolume(
		&ec2.DeleteVolumeInput{
			VolumeId: aws.String(volume.VolumeId),
		})
	if err != nil {
		p.Logger.Error("Send DeleteVolumes to %s failed: %+v", MyProvider, err)
		return err
	}

	// write back
	task.Directive = jsonutil.ToString(volume)

	return nil
}

func (p *ProviderHandler) waitInstanceVolumeAndNetwork(instanceService *ec2.EC2, task *models.Task, instanceId, volumeId string, timeout time.Duration, waitInterval time.Duration) (ins *ec2.Instance, err error) {
	p.Logger.Debug("Waiting for volume [%s] attached to Instance [%s]", volumeId, instanceId)
	err = p.AttachVolumes(task)
	if err != nil {
		p.Logger.Debug("Attach volume [%s] to Instance [%s] failed: %+v", volumeId, instanceId, err)
		return nil, err
	}

	err = p.WaitAttachVolumes(task)
	if err != nil {
		p.Logger.Debug("Waiting for volume [%s] attached to Instance [%s] failed: %+v", volumeId, instanceId, err)
		return nil, err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		describeOutput, err := instanceService.DescribeInstances(
			&ec2.DescribeInstancesInput{
				InstanceIds: aws.StringSlice([]string{instanceId}),
			},
		)
		if err != nil {
			return false, err
		}

		if len(describeOutput.Reservations) == 0 {
			return false, fmt.Errorf("instance with id [%s] not exist", instanceId)
		}
		if len(describeOutput.Reservations[0].Instances) == 0 {
			return false, fmt.Errorf("instance with id [%s] not exist", instanceId)
		}

		instance := describeOutput.Reservations[0].Instances[0]

		if instance.PrivateIpAddress == nil || aws.StringValue(instance.PrivateIpAddress) == "" {
			return false, nil
		}
		if instance.PublicIpAddress == nil || aws.StringValue(instance.PublicIpAddress) == "" {
			return false, nil
		}
		if volumeId != "" {
			if len(instance.BlockDeviceMappings) == 0 {
				return false, nil
			}

			found := false
			for _, dev := range instance.BlockDeviceMappings {
				if aws.StringValue(dev.Ebs.VolumeId) == volumeId {
					found = true
				}
			}

			if !found {
				return false, nil
			}
		}

		ins = instance
		p.Logger.Debug("Instance [%s] get IP address [%s]", instanceId, *ins.PrivateIpAddress)
		return true, nil
	}, timeout, waitInterval)
	return
}

func (p *ProviderHandler) WaitRunInstances(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = p.WaitInstanceState(task, constants.StatusRunning)
	if err != nil {
		p.Logger.Error("Wait %s job [%s] failed: %+v", MyProvider, instance.TargetJobId, err)
		return err
	}

	output, err := p.waitInstanceVolumeAndNetwork(instanceService, task, instance.InstanceId, instance.VolumeId, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		p.Logger.Error("Wait %s instance [%s] network failed: %+v", MyProvider, instance.InstanceId, err)
		return err
	}

	instance.PrivateIp = aws.StringValue(output.PrivateIpAddress)
	instance.EIP = aws.StringValue(output.PublicIpAddress)
	if len(output.BlockDeviceMappings) > 0 {
		for _, dev := range output.BlockDeviceMappings {
			if aws.StringValue(dev.Ebs.VolumeId) == instance.VolumeId {
				instance.Device = aws.StringValue(dev.DeviceName)
			}
		}
	}

	// write back
	task.Directive = jsonutil.ToString(instance)

	p.Logger.Debug("WaitRunInstances task [%s] directive: %s", task.TaskId, task.Directive)

	return nil
}

func (p *ProviderHandler) WaitInstanceState(task *models.Task, state string) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	instance, err := models.NewInstance(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(instance.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instance.InstanceId)},
		}

		output, err := instanceService.DescribeInstances(&input)
		if err != nil {
			return true, err
		}

		if len(output.Reservations) == 0 {
			return true, fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
		}

		if len(output.Reservations[0].Instances) == 0 {
			return true, fmt.Errorf("instance with id [%s] not exist", instance.InstanceId)
		}

		if aws.StringValue(output.Reservations[0].Instances[0].State.Name) == state {
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		p.Logger.Error("Wait %s instance [%s] status become to [%s] failed: %+v", MyProvider, instance.InstanceId, state, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitVolumeState(task *models.Task, state string) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	err = funcutil.WaitForSpecificOrError(func() (bool, error) {
		input := ec2.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(volume.VolumeId)},
		}

		output, err := instanceService.DescribeVolumes(&input)
		if err != nil {
			return true, err
		}

		if len(output.Volumes) == 0 {
			return true, fmt.Errorf("volume [%s] not found", volume.VolumeId)
		}

		if aws.StringValue(output.Volumes[0].State) == state {
			return true, nil
		}

		return false, nil
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
	if err != nil {
		p.Logger.Error("Wait %s volume [%s] status become to [%s] failed: %+v", MyProvider, volume.VolumeId, state, err)
		return err
	}

	return nil
}

func (p *ProviderHandler) WaitStopInstances(task *models.Task) error {
	return p.WaitInstanceState(task, constants.StatusStopped)
}

func (p *ProviderHandler) WaitStartInstances(task *models.Task) error {
	return p.WaitInstanceState(task, constants.StatusRunning)
}

func (p *ProviderHandler) WaitDeleteInstances(task *models.Task) error {
	return p.WaitInstanceState(task, constants.StatusTerminated)
}

func (p *ProviderHandler) WaitCreateVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, constants.StatusAvailable)
}

func (p *ProviderHandler) WaitAttachVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, constants.StatusInUse)
}

func (p *ProviderHandler) WaitDetachVolumes(task *models.Task) error {
	return p.WaitVolumeState(task, constants.StatusAvailable)
}

func (p *ProviderHandler) WaitDeleteVolumes(task *models.Task) error {
	if task.Directive == "" {
		p.Logger.Warn("Skip task without directive")
		return nil
	}
	volume, err := models.NewVolume(task.Directive)
	if err != nil {
		return err
	}
	instanceService, err := p.initInstanceService(volume.RuntimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return err
	}

	input2 := ec2.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(volume.VolumeId)},
	}
	return instanceService.WaitUntilVolumeDeleted(&input2)
}

func (p *ProviderHandler) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	instanceService, err := p.initInstanceService(req.GetRuntimeId().GetValue())
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	filter := ec2.Filter{
		Name:   aws.String("availabilityZone"),
		Values: aws.StringSlice(req.GetZone()),
	}

	input := new(ec2.DescribeSubnetsInput)
	if len(req.GetSubnetId()) > 0 {
		input.SubnetIds = aws.StringSlice(req.GetSubnetId())
		input.Filters = []*ec2.Filter{&filter}
	}

	output, err := instanceService.DescribeSubnets(input)
	if err != nil {
		p.Logger.Error("DescribeSubnets to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	if len(output.Subnets) == 0 {
		p.Logger.Error("Send DescribeVxNets to %s failed with 0 output subnets", MyProvider)
		return nil, fmt.Errorf("send DescribeVxNets to %s failed with 0 output subnets", MyProvider)
	}

	response := new(pb.DescribeSubnetsResponse)

	for _, sn := range output.Subnets {

		name := ""
		for _, tag := range sn.Tags {
			if aws.StringValue(tag.Key) == "Name" {
				name = aws.StringValue(tag.Value)
			}
		}

		subnet := &pb.Subnet{
			SubnetId: pbutil.ToProtoString(aws.StringValue(sn.SubnetId)),
			Name:     pbutil.ToProtoString(name),
			VpcId:    pbutil.ToProtoString(aws.StringValue(sn.VpcId)),
			Zone:     pbutil.ToProtoString(aws.StringValue(sn.AvailabilityZone)),
		}
		response.SubnetSet = append(response.SubnetSet, subnet)
	}

	response.TotalCount = uint32(len(response.SubnetSet))

	return response, nil
}

func (p *ProviderHandler) CheckResourceQuotas(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	roleCount := make(map[string]int)
	for _, clusterNode := range clusterWrapper.ClusterNodes {
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
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	output, err := instanceService.DescribeVpcs(
		&ec2.DescribeVpcsInput{
			VpcIds: aws.StringSlice([]string{vpcId}),
		})
	if err != nil {
		p.Logger.Error("DescribeVpcs to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	if len(output.Vpcs) == 0 {
		p.Logger.Error("Send DescribeVpcs to %s failed with 0 output instances", MyProvider)
		return nil, fmt.Errorf("send DescribeVpcs to %s failed with 0 output instances", MyProvider)
	}

	vpc := output.Vpcs[0]

	filter := &ec2.Filter{
		Name:   aws.String("vpc-id"),
		Values: []*string{vpc.VpcId},
	}

	subnetOutput, err := instanceService.DescribeSubnets(
		&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{filter},
		})
	if err != nil {
		p.Logger.Error("DescribeSubnets to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var subnets []string
	for _, subnet := range subnetOutput.Subnets {
		subnets = append(subnets, aws.StringValue(subnet.SubnetId))
	}

	name := ""
	for _, tag := range vpc.Tags {
		if aws.StringValue(tag.Key) == "Name" {
			name = aws.StringValue(tag.Value)
		}
	}

	return &models.Vpc{
		VpcId:   aws.StringValue(vpc.VpcId),
		Name:    name,
		Status:  aws.StringValue(vpc.State),
		Subnets: subnets,
	}, nil
}

func (p *ProviderHandler) DescribeAvailabilityZones(url, credential, zone string) ([]string, error) {
	awsSession, err := p.initAWSSession(url, credential, zone)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	instanceService := ec2.New(awsSession)

	input := ec2.DescribeAvailabilityZonesInput{}

	output, err := instanceService.DescribeAvailabilityZones(&input)
	if err != nil {
		p.Logger.Error("DescribeAvailabilityZones to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var zones []string
	for _, zone := range output.AvailabilityZones {
		zones = append(zones, aws.StringValue(zone.ZoneName))
	}
	return zones, nil
}

func (p *ProviderHandler) DescribeZones(url, credential string) ([]string, error) {
	zone := DefaultZone
	awsSession, err := p.initAWSSession(url, credential, zone)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	instanceService := ec2.New(awsSession)

	input := ec2.DescribeRegionsInput{}

	output, err := instanceService.DescribeRegions(&input)
	if err != nil {
		p.Logger.Error("DescribeRegions to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var zones []string
	for _, zone := range output.Regions {
		zones = append(zones, aws.StringValue(zone.RegionName))
	}
	return zones, nil
}

func (p *ProviderHandler) DescribeKeyPairs(url, credential, zone string) ([]string, error) {
	awsSession, err := p.initAWSSession(url, credential, zone)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return nil, err
	}

	instanceService := ec2.New(awsSession)

	input := ec2.DescribeKeyPairsInput{}

	output, err := instanceService.DescribeKeyPairs(&input)
	if err != nil {
		p.Logger.Error("DescribeKeyPairs to %s failed: %+v", MyProvider, err)
		return nil, err
	}

	var keys []string
	for _, key := range output.KeyPairs {
		keys = append(keys, aws.StringValue(key.KeyName))
	}
	return keys, nil
}

func (p *ProviderHandler) DescribeImage(runtimeId, name string) (string, error) {
	instanceService, err := p.initInstanceService(runtimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return "", err
	}

	filter := &ec2.Filter{Name: aws.String("name"), Values: aws.StringSlice([]string{name})}

	input := ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{filter},
	}

	output, err := instanceService.DescribeImages(&input)
	if err != nil {
		p.Logger.Error("DescribeImages to %s failed: %+v", MyProvider, err)
		return "", err
	}

	if len(output.Images) == 0 {
		return "", fmt.Errorf("image with name [%s] not exist", name)
	}

	imageId := aws.StringValue(output.Images[0].ImageId)

	return imageId, nil
}

func (p *ProviderHandler) DescribeAvailabilityZoneBySubnetId(runtimeId, subnetId string) (string, error) {
	instanceService, err := p.initInstanceService(runtimeId)
	if err != nil {
		p.Logger.Error("Init %s api service failed: %+v", MyProvider, err)
		return "", err
	}

	input := ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice([]string{subnetId}),
	}

	output, err := instanceService.DescribeSubnets(&input)
	if err != nil {
		p.Logger.Error("DescribeSubnets to %s failed: %+v", MyProvider, err)
		return "", err
	}

	if len(output.Subnets) == 0 {
		return "", fmt.Errorf("subnet with id [%s] not exist", subnetId)
	}

	zone := aws.StringValue(output.Subnets[0].AvailabilityZone)

	return zone, nil
}
