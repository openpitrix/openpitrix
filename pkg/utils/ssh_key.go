// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package utils

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/utils/passphraseless"
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
