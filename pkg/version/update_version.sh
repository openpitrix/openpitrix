#!/bin/sh
# Copyright 2017 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

VERSION=`git describe --tags --always --dirty="-dev"`
GIT_SHA1=`git show --quiet --pretty=format:%H`
BUILD_DATE=`date +%Y-%m-%d`

cat <<EOF | gofmt > z_update_version.go
// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package version

func init() {
    ShortVersion = "$VERSION"
    GitSha1Version = "$GIT_SHA1"
    BuildDate = "$BUILD_DATE"
}
EOF
