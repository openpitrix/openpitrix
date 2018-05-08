// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package indexer

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/httputil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type helmIndexer struct {
	indexer
}

func NewHelmIndexer(repo *pb.Repo) *helmIndexer {
	return &helmIndexer{
		indexer: indexer{repo: repo},
	}
}

type helmVersionWrapper struct {
	*repo.ChartVersion
}

func (h helmVersionWrapper) GetVersion() string     { return h.ChartVersion.GetVersion() }
func (h helmVersionWrapper) GetAppVersion() string  { return h.ChartVersion.GetAppVersion() }
func (h helmVersionWrapper) GetDescription() string { return h.ChartVersion.GetDescription() }
func (h helmVersionWrapper) GetUrls() []string      { return h.ChartVersion.URLs }

func (i *helmIndexer) IndexRepo() error {
	indexFile, err := i.getIndexFile()
	if err != nil {
		return err
	}
	for chartName, chartVersions := range indexFile.Entries {
		var appId string
		logger.Debug("Start index chart [%s]", chartName)
		logger.Debug("Chart [%s] has [%d] versions", chartName, chartVersions.Len())
		if len(chartVersions) == 0 {
			return fmt.Errorf("failed to sync chart [%s], no versions", chartName)
		}
		appId, err = i.syncAppInfo(chartVersions[0])
		if err != nil {
			logger.Error("Failed to sync chart [%s] to app info", chartName)
			return err
		}
		logger.Info("Sync chart [%s] to app [%s] success", chartName, appId)
		sort.Sort(chartVersions)
		for index, chartVersion := range chartVersions {
			var versionId string
			v := helmVersionWrapper{ChartVersion: chartVersion}
			versionId, err = i.syncAppVersionInfo(appId, v, index)
			if err != nil {
				logger.Error("Failed to sync chart version [%s] to app version", chartVersion.GetAppVersion())
				return err
			}
			logger.Debug("Chart version [%s] sync to app version [%s]", chartVersion.GetVersion(), versionId)
		}
	}
	return err
}

// Reference: https://sourcegraph.com/github.com/kubernetes/helm@fe9d365/-/blob/pkg/repo/chartrepo.go#L111:27
func (i *helmIndexer) getIndexFile() (indexFile *repo.IndexFile, err error) {
	repoUrl := i.repo.GetUrl().GetValue()
	var indexURL string
	indexFile = &repo.IndexFile{}
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		return
	}
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") + "/index.yaml"
	indexURL = parsedURL.String()
	resp, err := httputil.HttpGet(indexURL)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = yamlutil.Decode(content, indexFile)
	if err != nil {
		return
	}
	indexFile.SortEntries()
	return
}

// Reference: https://sourcegraph.com/github.com/kubernetes/helm@fe9d365/-/blob/pkg/downloader/chart_downloader.go#L225:35
func GetPackageFile(chartVersion *repo.ChartVersion, repoUrl string) (*chart.Chart, error) {
	if len(chartVersion.URLs) == 0 {
		return nil, fmt.Errorf("chart [%s] has no downloadable URLs", chartVersion.Name)
	}
	u, err := url.Parse(chartVersion.URLs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid chart URL format: %v", chartVersion.URLs)
	}

	// If the URL is relative (no scheme), prepend the chart repo's base URL
	if !u.IsAbs() {
		repoURL, err := url.Parse(repoUrl)
		if err != nil {
			return nil, err
		}
		q := repoURL.Query()
		// We need a trailing slash for ResolveReference to work, but make sure there isn't already one
		repoURL.Path = strings.TrimSuffix(repoURL.Path, "/") + "/"
		u = repoURL.ResolveReference(u)
		u.RawQuery = q.Encode()
	}
	resp, err := httputil.HttpGet(u.String())
	if err != nil {
		return nil, err
	}
	return chartutil.LoadArchive(resp.Body)
}
