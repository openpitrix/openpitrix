package models

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type RoleResource struct {
	Role         string
	Cpu          uint32
	Gpu          uint32
	Memory       uint32
	InstanceSize uint32
	StorageSize  uint32
}

func PbToRoleResource(pbRoleResource *pb.RoleResource) *RoleResource {
	roleResource := &RoleResource{}
	roleResource.Role = pbRoleResource.GetRole().GetValue()
	roleResource.Cpu = pbRoleResource.GetCpu().GetValue()
	roleResource.Gpu = pbRoleResource.GetGpu().GetValue()
	roleResource.Memory = pbRoleResource.GetMemory().GetValue()
	roleResource.InstanceSize = pbRoleResource.GetInstanceSize().GetValue()
	roleResource.StorageSize = pbRoleResource.GetStorageSize().GetValue()

	return roleResource
}

func (r *RoleResource) IsSame(clusterRole *ClusterRole) (bool, *RoleResizeResource) {
	if r.Role != clusterRole.Role {
		logger.Error(nil, "OperatorType resource [%s] not match cluster role [%s]", r.Role, clusterRole.Role)
		return false, nil
	}

	roleResizeResource := &RoleResizeResource{
		Role: clusterRole.Role,
	}
	if r.Cpu != 0 && r.Cpu != clusterRole.Cpu {
		roleResizeResource.Cpu = true
		clusterRole.Cpu = r.Cpu
	}
	if r.Gpu != 0 && r.Gpu != clusterRole.Gpu {
		// roleResizeResource.Gpu = true
		// clusterRole.Gpu = r.Gpu
	}
	if r.Memory != 0 && r.Memory != clusterRole.Memory {
		roleResizeResource.Memory = true
		clusterRole.Memory = r.Memory
	}
	if r.StorageSize > clusterRole.StorageSize {
		roleResizeResource.StorageSize = true
		clusterRole.StorageSize = r.StorageSize
	}
	if r.InstanceSize > clusterRole.InstanceSize {
		// roleResizeResource.InstanceSize = true
		// clusterRole.InstanceSize = r.InstanceSize
	}

	if !roleResizeResource.Cpu &&
		!roleResizeResource.Gpu &&
		!roleResizeResource.Memory &&
		!roleResizeResource.StorageSize &&
		!roleResizeResource.InstanceSize {
		return true, nil
	} else {
		return false, roleResizeResource
	}
}

type RoleResizeResource struct {
	Role         string
	Cpu          bool
	Gpu          bool
	Memory       bool
	InstanceSize bool
	StorageSize  bool
}

type RoleResizeResources []*RoleResizeResource

func NewRoleResizeResources(data string) (RoleResizeResources, error) {
	var roleResizeResources RoleResizeResources
	err := jsonutil.Decode([]byte(data), &roleResizeResources)
	if err != nil {
		logger.Error(nil, "Decode [%s] into role resize resources failed: %+v", data, err)
	}
	return roleResizeResources, err
}
