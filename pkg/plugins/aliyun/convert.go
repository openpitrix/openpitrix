// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package aliyun

import "fmt"

const (
	G = 1024
)

type CpuAndMemory struct {
	Cpu    int
	Memory int
}

var (
	InstanceTypeMap = map[CpuAndMemory]string{
		{1, 1 * G}: "ecs.t5-lc1m1.small",
		{1, 2 * G}: "ecs.t5-lc1m2.small",

		{2, 4 * G}: "ecs.t5-lc1m2.large",
		{2, 8 * G}: "ecs.t5-lc1m4.large",
		{2, 2 * G}: "ecs.t5-c1m1.large",

		{4, 4 * G}:  "ecs.t5-c1m1.xlarge",
		{4, 8 * G}:  "ecs.t5-c1m2.xlarge",
		{4, 16 * G}: "ecs.t5-c1m4.xlarge",

		{8, 8 * G}:  "ecs.t5-c1m1.2xlarge",
		{8, 16 * G}: "ecs.t5-c1m2.2xlarge",
		{8, 32 * G}: "ecs.t5-c1m4.2xlarge",

		{16, 16 * G}: "ecs.t5-c1m1.4xlarge",
		{16, 32 * G}: "ecs.t5-c1m2.4xlarge",
	}

	// The volume type. This can be gp2 for General Purpose SSD, io1 for Provisioned
	// IOPS SSD, st1 for Throughput Optimized HDD, sc1 for Cold HDD, or standard
	// for Magnetic volumes.
	//
	// Defaults: If no volume type is specified, the default is standard in us-east-1,
	// eu-west-1, eu-central-1, us-west-2, us-west-1, sa-east-1, ap-northeast-1,
	// ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-south-1, us-gov-west-1,
	// and cn-north-1. In all other regions, EBS defaults to gp2.
	volumeTypeMap = map[int]string{
		1: "cloud",
		2: "cloud_efficiency",
		3: "cloud_ssd",
	}
)

func ConvertToInstanceType(cpu, memory int) (string, error) {
	instanceType, ok := InstanceTypeMap[CpuAndMemory{cpu, memory}]
	if !ok {
		return "", fmt.Errorf("no aws instance type matched with cpu [%d] memory [%d]", cpu, memory)
	}

	return instanceType, nil
}

func ConvertToVolumeType(volumeClass int) (string, error) {
	volumeType, ok := volumeTypeMap[volumeClass]
	if !ok {
		return "", fmt.Errorf("no aws volume type matched with volume class [%d]", volumeClass)
	}

	return volumeType, nil
}
