// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package wrapper

import (
	"fmt"
	"strings"

	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type HelmVersionWrapper struct {
	*repo.ChartVersion
}

func (h HelmVersionWrapper) GetVersion() string     { return h.ChartVersion.GetVersion() }
func (h HelmVersionWrapper) GetAppVersion() string  { return h.ChartVersion.GetAppVersion() }
func (h HelmVersionWrapper) GetDescription() string { return h.ChartVersion.GetDescription() }
func (h HelmVersionWrapper) GetUrls() string {
	return h.ChartVersion.URLs[0]
}

func (h HelmVersionWrapper) GetSources() string {
	if len(h.ChartVersion.GetSources()) == 0 {
		return ""
	}
	return jsonutil.ToString(h.ChartVersion.GetSources())
}

func (h HelmVersionWrapper) GetKeywords() string {
	return strings.Join(h.ChartVersion.GetKeywords(), ",")
}

func (h HelmVersionWrapper) GetMaintainers() string {
	if len(h.ChartVersion.GetMaintainers()) == 0 {
		return ""
	}
	return jsonutil.ToString(h.ChartVersion.GetMaintainers())
}

func (h HelmVersionWrapper) GetScreenshots() string {
	return ""
}

func (h HelmVersionWrapper) GetVersionName() string {
	versionName := h.GetVersion()
	if h.GetAppVersion() != "" {
		versionName += fmt.Sprintf(" [%s]", h.GetAppVersion())
	}
	return versionName
}

func (h HelmVersionWrapper) GetPackageName() string {
	return fmt.Sprintf("%s-%s.tgz", h.Name, h.Version)
}
