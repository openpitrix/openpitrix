// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package indexer

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/yaml"
)

type devkitIndexer struct {
	indexer
}

func NewDevkitIndexer(repo *pb.Repo) *devkitIndexer {
	return &devkitIndexer{
		indexer: indexer{repo: repo},
	}
}

func (i *devkitIndexer) getIndexFile() (indexFile *app.IndexFile, err error) {
	repoUrl := i.repo.GetUrl().GetValue()
	var indexURL string
	indexFile = &app.IndexFile{}
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return
	}
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") + "/index.yaml"
	indexURL = parsedURL.String()
	resp, err := utils.HttpGet(indexURL)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = yaml.Decode(content, indexFile)
	if err != nil {
		return
	}
	indexFile.SortEntries()
	return
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
		appId, err = i.syncAppInfo(appVersions[0])
		if err != nil {
			logger.Error("Failed to sync app [%s] to app info", appName)
			return err
		}
		logger.Info("Sync chart [%s] to app [%s] success", appName, appId)
		for _, appVersion := range appVersions {
			var versionId string
			versionId, err = i.syncAppVersionInfo(appId, appVersion)
			if err != nil {
				logger.Error("Failed to sync app version [%s] to app version", appVersion.GetAppVersion())
				return err
			}
			logger.Debug("App version [%s] sync to app version [%s]", appVersion.GetVersion(), versionId)
		}
	}
	return err
}
