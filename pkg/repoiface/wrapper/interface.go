// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package wrapper

import (
	"time"

	"github.com/Masterminds/semver"
)

type VersionInterface interface {
	GetName() string
	GetVersion() string
	GetAppVersion() string
	GetDescription() string
	GetUrls() string
	GetVersionName() string
	GetIcon() string
	GetHome() string
	GetSources() string
	GetKeywords() string
	GetMaintainers() string
	GetScreenshots() string
	GetPackageName() string
	GetCreateTime() time.Time
}

type VersionInterfaces []VersionInterface

// Len returns the length.
func (c VersionInterfaces) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c VersionInterfaces) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c VersionInterfaces) Less(a, b int) bool {
	// Failed parse pushes to the back.
	i, err := semver.NewVersion(c[a].GetVersion())
	if err != nil {
		return true
	}
	j, err := semver.NewVersion(c[b].GetVersion())
	if err != nil {
		return false
	}
	return i.LessThan(j)
}

type IndexInterface interface {
	GetEntries() map[string]VersionInterfaces
}
