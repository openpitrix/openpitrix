// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package wrapper

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
}
