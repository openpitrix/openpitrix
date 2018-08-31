// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"fmt"
	"path/filepath"

	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/urlutil"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

func IndexDirectory(dir, baseURL string) (*opapp.IndexFile, error) {
	archives, err := filepath.Glob(filepath.Join(dir, "*.tgz"))
	if err != nil {
		return nil, err
	}

	index := opapp.NewIndexFile()
	for _, arch := range archives {
		fname, err := filepath.Rel(dir, arch)
		if err != nil {
			return index, err
		}

		var parentDir string
		parentDir, fname = filepath.Split(fname)
		parentURL, err := urlutil.URLJoin(baseURL, parentDir)
		if err != nil {
			parentURL = filepath.Join(baseURL, parentDir)
		}

		c, err := Load(arch)
		if err != nil {
			fmt.Printf("Load file [%s] error: %s\n", fname, err)
			continue
		}
		hash, err := provenance.DigestFile(arch)
		if err != nil {
			return index, err
		}
		index.Add(c.Metadata, fname, parentURL, hash)
	}
	return index, nil
}
