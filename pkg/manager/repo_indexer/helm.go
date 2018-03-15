// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager/app"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
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
	// TODO: SortEntries will panic, fix this bug
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

func SyncAppInfo(repoId, owner, chartName string, chartVersions *repo.ChartVersions) (string, error) {
	var appId string
	logger.Debugf("chart [%s] has [%d] versions", chartName, chartVersions.Len())
	ctx := sender.NewContext(context.Background(), sender.GetSystemUser())
	appManagerClient, err := app.NewAppManagerClient(ctx)
	if err != nil {
		return appId, err
	}
	req := pb.DescribeAppsRequest{}
	req.RepoId = []string{repoId}
	req.Owner = []string{owner}
	req.ChartName = []string{chartName}
	res, err := appManagerClient.DescribeApps(ctx, &req)
	if err != nil {
		return appId, err
	}
	if res.TotalCount == 0 {
		createReq := pb.CreateAppRequest{}
		createReq.RepoId = utils.ToProtoString(repoId)
		createReq.ChartName = utils.ToProtoString(chartName)
		createReq.Name = utils.ToProtoString(chartName)
		createRes, err := appManagerClient.CreateApp(ctx, &createReq)
		if err != nil {
			return appId, err
		}
		appId = createRes.GetApp().GetAppId().GetValue()
		return appId, err

	} else {
		modifyReq := pb.ModifyAppRequest{}
		modifyReq.AppId = res.AppSet[0].AppId
		modifyReq.Name = utils.ToProtoString(chartName)
		modifyReq.ChartName = utils.ToProtoString(chartName)
		modifyRes, err := appManagerClient.ModifyApp(ctx, &modifyReq)
		if err != nil {
			return appId, err
		}
		appId = modifyRes.GetApp().GetAppId().GetValue()
		return appId, err
	}
}
