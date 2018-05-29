// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import "fmt"

func formatVolumeCmd(device, fileSystem string) string {
	mkfs := fmt.Sprintf("(echo y) | mkfs -t %s %s > /dev/null", fileSystem, device)
	return mkfs
}

func updateTabCmd(device, mountPoint, fileSystem, mountOptions string) string {
	cmd := fmt.Sprintf("%s%s %s %s %s %s", OpenPitrixSbinPath, UpdateFstabFile, fileSystem, device, mountPoint, mountOptions)
	return cmd
}

func mountVolumeCmd(device, mountPoint, fileSystem, mountOptions string) string {
	fileSystemCmd := fileSystem
	if fileSystem != "" {
		fileSystemCmd = fmt.Sprintf("-t %s", fileSystem)
	}

	mountOptionsCmd := mountOptions
	if mountOptions != "" {
		mountOptionsCmd = fmt.Sprintf("-o %s", mountOptions)
	}
	mount := fmt.Sprintf("mkdir -p %s && mount %s %s %s %s",
		mountPoint, fileSystemCmd, mountOptionsCmd, device, mountPoint)

	updateTab := updateTabCmd(device, mountPoint, fileSystem, mountOptions)

	return fmt.Sprintf("%s && %s", mount, updateTab)
}

func FormatAndMountVolumeCmd(device, mountPoint, fileSystem, mountOptions string) string {
	formatCmd := formatVolumeCmd(device, fileSystem)
	mountCmd := mountVolumeCmd(device, mountPoint, fileSystem, mountOptions)
	return fmt.Sprintf("%s \"%s && %s\"", HostCmdPrefix, formatCmd, mountCmd)
}

func UmountVolumeCmd(mountPoint string) string {
	umount := fmt.Sprintf("%s \"fuser -ck %s; umount %s\"", HostCmdPrefix, mountPoint, mountPoint)
	return umount
}
