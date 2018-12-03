// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"bytes"
	"context"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/repoiface/wrapper"
)

const (
	Helm    = "helm"
	Vmbased = "vmbased"
)

var SupportedPackageType = []string{Helm, Vmbased}

func LoadPackage(ctx context.Context, t string, pkg []byte) (wrapper.VersionInterface, error) {
	if t == Helm {
		p, err := chartutil.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			logger.Error(ctx, "Failed to load package, error: %+v", err)
			return nil, err
		}

		pkgVersion := wrapper.HelmVersionWrapper{ChartVersion: &repo.ChartVersion{Metadata: p.Metadata}}
		return pkgVersion, nil
	} else {
		p, err := devkit.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			logger.Error(ctx, "Failed to load package, error: %+v", err)
			return nil, err
		}

		pkgVersion := wrapper.OpVersionWrapper{OpVersion: &opapp.OpVersion{Metadata: p.Metadata}}
		return pkgVersion, nil
	}
}
