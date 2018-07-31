// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package sshutil

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/util/passphraseless"
)

func MakeSSHKeyPair(keyType string) (string, string, error) {
	switch keyType {
	case "ssh-rsa":
		return passphraseless.MakeSSHRsaKeyPair()
	case "ssh-dsa":
		return passphraseless.MakeSSHDsaKeyPair()
	default:
		return "", "", fmt.Errorf("wrong key type [%s]", keyType)
	}
}

func prepareDirAndFileCmd() string {
	path := "/root/.ssh"
	file := fmt.Sprintf("%s/authorized_keys", path)
	cmd := fmt.Sprintf("if [ ! -d %s ];then mkdir -p %s && chmod 0755 %s && touch %s && chmod 0644 %s;fi",
		path, path, path, file, file)
	return cmd
}

func DoAttachCmd(keyPair string) string {
	path := "/root/.ssh"
	file := fmt.Sprintf("%s/authorized_keys", path)
	tmpFile := fmt.Sprintf("%s/authorized_keys.tmp", path)
	tmpFileCmd := fmt.Sprintf("touch %s && chmod 0644 %s", tmpFile, tmpFile)
	cmd := fmt.Sprintf("while read line; do if [ \\\"$line\\\" != \\\"%s\\\" ] && [ -n \\\"$line\\\" ];then echo $line >> %s;fi;done < %s && echo %s >> %s",
		keyPair, tmpFile, file, keyPair, tmpFile)
	mvCmd := fmt.Sprintf("mv %s %s", tmpFile, file)
	return fmt.Sprintf("%s && %s && %s && %s", prepareDirAndFileCmd(), tmpFileCmd, cmd, mvCmd)
}

func DoDetachCmd(keyPair string) string {
	path := "/root/.ssh"
	file := fmt.Sprintf("%s/authorized_keys", path)
	tmpFile := fmt.Sprintf("%s/authorized_keys.tmp", path)
	tmpFileCmd := fmt.Sprintf("touch %s && chmod 0644 %s", tmpFile, tmpFile)
	cmd := fmt.Sprintf("while read line; do if [ \\\"$line\\\" != \\\"%s\\\" ] && [ -n \\\"$line\\\" ];then echo $line >> %s;fi;done < %s",
		keyPair, tmpFile, file)
	mvCmd := fmt.Sprintf("mv %s %s", tmpFile, file)
	return fmt.Sprintf("%s && %s && %s && %s", prepareDirAndFileCmd(), tmpFileCmd, cmd, mvCmd)
}
