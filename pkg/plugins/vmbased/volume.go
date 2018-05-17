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
	// `blkid /dev/vdc` outputs like
	// /dev/vdc: UUID="b73ab0cd-f976-4559-b9f6-a6fbf013570d" TYPE="ext4"
	uuid := fmt.Sprintf("uuid=`blkid %s | awk -F '\"' '{print $2}'`", device)
	content := fmt.Sprintf("UUID=$uuid %s %s %s 0 2", mountPoint, fileSystem, mountOptions)
	fstab := fmt.Sprintf("sed -i \"s/^UUID=$uuid .*//g\" /etc/fstab && echo \"%s\" >> /etc/fstab", content)
	return fmt.Sprintf("%s && %s", uuid, fstab)
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
	return fmt.Sprintf("%s && %s", formatCmd, mountCmd)
}

func UmountVolumeCmd(mountPoint string) string {
	umount := fmt.Sprintf("fuser -ck %s; umount %s", mountPoint, mountPoint)
	return umount
}
