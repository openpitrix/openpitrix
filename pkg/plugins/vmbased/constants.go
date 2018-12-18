// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

const (
	ActionRunInstances       = "RunInstances"
	ActionStartInstances     = "StartInstances"
	ActionStopInstances      = "StopInstances"
	ActionTerminateInstances = "TerminateInstances"
	ActionResizeInstances    = "ResizeInstances"

	ActionCreateVolumes = "CreateVolumes"
	ActionAttachVolumes = "AttachVolumes"
	ActionDetachVolumes = "DetachVolumes"
	ActionDeleteVolumes = "DeleteVolumes"
	ActionResizeVolumes = "ResizeVolumes"

	ActionFormatAndMountVolume         = "FormatAndMountVolume"
	ActionWaitFrontgateAvailable       = "WaitFrontgateAvailable"
	ActionRegisterMetadata             = "RegisterMetadata"
	ActionRegisterMetadataMapping      = "RegisterMetadataMapping"
	ActionRegisterNodesMetadata        = "RegisterNodesMetadata"
	ActionRegisterEnvMetadata          = "RegisterEnvMetadata"
	ActionRegisterNodesMetadataMapping = "RegisterNodesMetadataMapping"
	ActionDeregisterMetadata           = "DeregisterMetadata"
	ActionDeregisterMetadataMapping    = "DeregisterMetadataMapping"
	ActionRegisterCmd                  = "RegisterCmd"
	ActionDeregisterCmd                = "DeregisterCmd"
	ActionStartConfd                   = "StartConfd"
	ActionStopConfd                    = "StopConfd"
	ActionSetFrontgateConfig           = "SetFrontgateConfig"
	ActionSetDroneConfig               = "SetDroneConfig"
	ActionPingDrone                    = "PingDrone"
	ActionPingFrontgate                = "PingFrontgate"
	PingMetadataBackend                = "PingMetadataBackend"
	ActionRunCommandOnDrone            = "RunCommandOnDrone"
	ActionRemoveContainerOnDrone       = "RemoveContainerOnDrone"
	ActionRemoveContainerOnFrontgate   = "RemoveContainerOnFrontgate"
	ActionRunCommandOnFrontgateNode    = "RunCommandOnFrontgateNode"
)

const (
	RegisterClustersRootPath = "clusters"
	RegisterNodeHosts        = "hosts"
	RegisterNodeHost         = "host"
	RegisterNodeCluster      = "cluster"
	RegisterNodeEnv          = "env"
	RegisterNodeLoadbalancer = "loadbalancer"
	RegisterNodeSelf         = "self"
	RegisterNodeCmd          = "cmd"
	RegisterNodeCmdId        = "id"
	RegisterNodeCmdTimeout   = "timeout"
	RegisterNodeEndpoint     = "endpoints"
	RegisterNodeLinks        = "links"
	RegisterNodeAdding       = "adding-hosts"
	RegisterNodeDeleting     = "deleting-hosts"
	RegisterNodeScaling      = "scaling-hosts"
	RegisterNodeStopping     = "stopping-hosts"
	RegisterNodeStarting     = "starting-hosts"
)

const (
	// second
	TimeoutStartConfd           = 60
	TimeoutStopConfd            = 60
	TimeoutDeregister           = 60
	TimeoutRegister             = 60
	TimeoutFormatAndMountVolume = 600
	TimeoutUmountVolume         = 120
	TimeoutSshKeygen            = 120
	TimeoutRemoveContainer      = 120
	TimeoutKeyPair              = 60
)

const (
	OpenPitrixBasePath = "/opt/openpitrix/"
	OpenPitrixExecFile = "/etc/rc.local"
	OpenPitrixConfPath = OpenPitrixBasePath + "conf/"
	OpenPitrixSbinPath = OpenPitrixBasePath + "sbin/"
	OpenPitrixConfFile = "openpitrix.conf"
	DroneConfFile      = "drone.conf"
	FrontgateConfFile  = "frontgate.conf"
	UpdateFstabFile    = "update_fstab.sh"
	ConfdPath          = "/etc/confd/"
	MetadataLogLevel   = "debug"
	ConfdBackendType   = "libconfd-backend-metad"
	ConfdCmdLogPath    = "/opt/openpitrix/log/cmd.log"
	HostCmdPrefix      = "nsenter -t 1 -m -u -n -i sh -c"
	MetadataNodeName   = "metadata"
	DefaultNodeName    = "node"
	EtcdPort           = 2379
	MetadPort          = 80
)

const (
	InstanceSize           = 20
	DefaultMountPoint      = "/data"
	Ext4FileSystem         = "ext4"
	XfsFileSystem          = "xfs"
	DefaultExt4MountOption = "defaults,noatime"
	DefaultXfsMountOption  = "rw,noatime,inode64,allocsize=16m"
)
