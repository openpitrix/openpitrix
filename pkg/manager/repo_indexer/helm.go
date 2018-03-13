// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/utils"
)

// Reference: https://sourcegraph.com/github.com/kubernetes/helm@fe9d365/-/blob/pkg/repo/chartrepo.go#L117:2
func GetIndexFile(repoUrl string) (indexFile repo.IndexFile, err error) {
	var indexURL string
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
	err = yaml.Unmarshal(content, &indexFile)
	if err != nil {
		return
	}
	//indexFile.SortEntries()
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
	resp, err := utils.HttpGet(u.String())
	if err != nil {
		return nil, err
	}
	return chartutil.LoadArchive(resp.Body)
}
