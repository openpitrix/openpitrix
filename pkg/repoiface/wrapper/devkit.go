// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package wrapper

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type OpVersionWrapper struct {
	*opapp.OpVersion
}

func (h OpVersionWrapper) GetVersion() string { return h.OpVersion.GetVersion() }

func (h OpVersionWrapper) GetAppVersion() string { return h.OpVersion.GetAppVersion() }

func (h OpVersionWrapper) GetDescription() string { return h.OpVersion.GetDescription() }

func (h OpVersionWrapper) GetCreateTime() time.Time { return h.OpVersion.Created }

func (h OpVersionWrapper) GetUrls() string {
	return h.OpVersion.GetUrls()[0]
}

func (h OpVersionWrapper) GetSources() string {
	if len(h.OpVersion.GetSources()) == 0 {
		return ""
	}
	return jsonutil.ToString(h.OpVersion.GetSources())
}

func (h OpVersionWrapper) GetKeywords() string {
	return strings.Join(h.OpVersion.GetKeywords(), ",")
}

func (h OpVersionWrapper) GetMaintainers() string {
	if len(h.OpVersion.GetMaintainers()) == 0 {
		return ""
	}
	return jsonutil.ToString(h.OpVersion.GetMaintainers())
}

func (h OpVersionWrapper) GetScreenshots() string {
	if len(h.OpVersion.GetScreenshots()) == 0 {
		return ""
	}
	return jsonutil.ToString(h.OpVersion.GetScreenshots())
}

func (h OpVersionWrapper) GetVersionName() string {
	versionName := h.GetVersion()
	if h.GetAppVersion() != "" {
		versionName += fmt.Sprintf(" [%s]", h.GetAppVersion())
	}
	return versionName
}

func (h OpVersionWrapper) GetPackageName() string {
	return h.OpVersion.GetPackageName()
}

type OpIndexWrapper struct {
	*opapp.IndexFile
}

func (h OpIndexWrapper) GetEntries() map[string]VersionInterfaces {
	var entries = make(map[string]VersionInterfaces)
	for chartName, chartVersions := range h.Entries {
		var versions VersionInterfaces
		sort.Sort(chartVersions)
		for _, v := range chartVersions {
			versions = append(versions, OpVersionWrapper{OpVersion: v})
		}
		entries[chartName] = versions
	}
	return entries
}
