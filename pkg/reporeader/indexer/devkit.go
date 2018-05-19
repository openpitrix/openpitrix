// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package indexer

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/reporeader"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type devkitIndexer struct {
	indexer
}

func NewDevkitIndexer(repo *pb.Repo, reader reporeader.Reader) *devkitIndexer {
	return &devkitIndexer{
		indexer: indexer{repo: repo, reader: reader},
	}
}

func (i *devkitIndexer) getIndexFile() (*app.IndexFile, error) {
	var indexFile = new(app.IndexFile)
	content, err := i.reader.GetIndexYaml()
	if err != nil {
		return nil, errors.Wrap(err, "get index yaml failed")
	}
	err = yamlutil.Decode(content, indexFile)
	if err != nil {
		logger.Debug("%s", string(content))
		return nil, errors.Wrap(err, "decode yaml failed")
	}
	indexFile.SortEntries()
	return indexFile, nil
}

type appVersionWrapper struct {
	*app.Version
}

func (h appVersionWrapper) GetVersion() string     { return h.Version.GetVersion() }
func (h appVersionWrapper) GetAppVersion() string  { return h.Version.GetAppVersion() }
func (h appVersionWrapper) GetDescription() string { return h.Version.GetDescription() }
func (h appVersionWrapper) GetUrls() []string      { return h.Version.GetUrls() }
func (h appVersionWrapper) GetKeywords() []string  { return h.Version.GetKeywords() }
func (h appVersionWrapper) GetMaintainers() string {
	return jsonutil.ToString(h.Version.GetMaintainers())
}

func (i *devkitIndexer) IndexRepo() error {
	indexFile, err := i.getIndexFile()
	if err != nil {
		return err
	}
	for appName, appVersions := range indexFile.Entries {
		var appId string
		logger.Debug("Start index app [%s]", appName)
		logger.Debug("App [%s] has [%d] versions", appName, appVersions.Len())
		if len(appVersions) == 0 {
			return fmt.Errorf("failed to sync app [%s], no versions", appName)
		}
		appId, err = i.syncAppInfo(appVersionWrapper{appVersions[0]})
		if err != nil {
			logger.Error("Failed to sync app [%s] to app info", appName)
			return err
		}
		logger.Info("Sync chart [%s] to app [%s] success", appName, appId)
		sort.Sort(appVersions)
		for index, appVersion := range appVersions {
			var versionId string
			versionId, err = i.syncAppVersionInfo(appId, appVersion, index)
			if err != nil {
				logger.Error("Failed to sync app version [%s] to app version", appVersion.GetAppVersion())
				return err
			}
			logger.Debug("App version [%s] sync to app version [%s]", appVersion.GetVersion(), versionId)
		}
	}
	return err
}
