// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"net/url"
	"path/filepath"
	"strings"
)

func GetFileName(packageName string) string {
	// Ref: https://sourcegraph.com/github.com/kubernetes/helm@8478fb4fc723885b155c924d1c8c410b7a9444e6/-/blob/pkg/repo/index.go#L110:14
	_, file := filepath.Split(packageName)
	return file
}

func URLJoin(repoUrl, fileName string) string {
	if !strings.HasSuffix(repoUrl, "/") {
		repoUrl += "/"
	}
	base, _ := url.Parse(repoUrl)
	n, _ := url.Parse(fileName)
	return base.ResolveReference(n).String()
}

func getBucketPrefix(path string) (bucket, prefix string) {
	p := strings.Split(path, "/")
	bucket = p[1]
	prefix = strings.Join(p[2:], "/")
	return
}
