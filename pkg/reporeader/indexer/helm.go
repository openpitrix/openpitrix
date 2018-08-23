// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package indexer

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type helmIndexer struct {
	indexer
}

func NewHelmIndexer(i indexer) *helmIndexer {
	return &helmIndexer{
		indexer: i,
	}
}

type helmVersionWrapper struct {
	*repo.ChartVersion
}

func (h helmVersionWrapper) GetVersion() string     { return h.ChartVersion.GetVersion() }
func (h helmVersionWrapper) GetAppVersion() string  { return h.ChartVersion.GetAppVersion() }
func (h helmVersionWrapper) GetDescription() string { return h.ChartVersion.GetDescription() }
func (h helmVersionWrapper) GetUrls() []string      { return h.ChartVersion.URLs }
func (h helmVersionWrapper) GetKeywords() []string  { return h.ChartVersion.GetKeywords() }
func (h helmVersionWrapper) GetMaintainers() string {
	return jsonutil.ToString(h.ChartVersion.GetMaintainers())
}
func (h helmVersionWrapper) GetScreenshots() string {
	return ""
}

func (i *helmIndexer) IndexRepo() error {
	indexFile, err := i.getIndexFile()
	if err != nil {
		return err
	}
	for chartName, chartVersions := range indexFile.Entries {
		var appId string
		logger.Debug(i.ctx, "Start index chart [%s]", chartName)
		logger.Debug(i.ctx, "Chart [%s] has [%d] versions", chartName, chartVersions.Len())
		if len(chartVersions) == 0 {
			return fmt.Errorf("failed to sync chart [%s], no versions", chartName)
		}
		appId, err = i.syncAppInfo(helmVersionWrapper{chartVersions[0]})
		if err != nil {
			logger.Error(i.ctx, "Failed to sync chart [%s] to app info", chartName)
			return err
		}
		logger.Info(i.ctx, "Sync chart [%s] to app [%s] success", chartName, appId)
		sort.Sort(chartVersions)
		for index, chartVersion := range chartVersions {
			var versionId string
			v := helmVersionWrapper{ChartVersion: chartVersion}
			versionId, err = i.syncAppVersionInfo(appId, v, index)
			if err != nil {
				logger.Error(i.ctx, "Failed to sync chart version [%s] to app version", chartVersion.GetAppVersion())
				return err
			}
			logger.Debug(i.ctx, "Chart version [%s] sync to app version [%s]", chartVersion.GetVersion(), versionId)
		}
	}
	return err
}

// Reference: https://sourcegraph.com/github.com/kubernetes/helm@fe9d365/-/blob/pkg/repo/chartrepo.go#L111:27
func (i *helmIndexer) getIndexFile() (*repo.IndexFile, error) {
	var indexFile = new(repo.IndexFile)
	content, err := i.reader.GetIndexYaml()
	if err != nil {
		return nil, errors.Wrap(err, "get index yaml failed")
	}
	err = yamlutil.Decode(content, indexFile)
	if err != nil {
		return nil, errors.Wrap(err, "decode yaml failed")
	}
	indexFile.SortEntries()
	return indexFile, nil
}
