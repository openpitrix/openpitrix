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

	"github.com/golang/protobuf/ptypes/wrappers"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
	"openpitrix.io/openpitrix/pkg/utils/yaml"
)

// Reference: https://sourcegraph.com/github.com/kubernetes/helm@fe9d365/-/blob/pkg/repo/chartrepo.go#L111:27
func GetIndexFile(repoUrl string) (indexFile *repo.IndexFile, err error) {
	var indexURL string
	indexFile = &repo.IndexFile{}
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
	err = yaml.Unmarshal(content, indexFile)
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
	var description, icon, home, sources *wrappers.StringValue
	if chartVersions.Len() > 0 {
		chartVersion := (*chartVersions)[0]
		description = utils.ToProtoString(chartVersion.GetDescription())
		icon = utils.ToProtoString(chartVersion.GetIcon())
		home = utils.ToProtoString(chartVersion.GetHome())
		sources = utils.ToProtoString(strings.Join(chartVersion.Sources, ","))
	}
	if res.TotalCount == 0 {
		createReq := pb.CreateAppRequest{}
		createReq.RepoId = utils.ToProtoString(repoId)
		createReq.ChartName = utils.ToProtoString(chartName)
		createReq.Name = utils.ToProtoString(chartName)
		createReq.Description = description
		createReq.Icon = icon
		createReq.Home = home
		createReq.Sources = sources

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
		modifyReq.Description = description
		modifyReq.Icon = icon
		modifyReq.Home = home
		modifyReq.Sources = sources

		modifyRes, err := appManagerClient.ModifyApp(ctx, &modifyReq)
		if err != nil {
			return appId, err
		}
		appId = modifyRes.GetApp().GetAppId().GetValue()
		return appId, err
	}
}

func SyncAppVersionInfo(appId, owner string, chartVersion *repo.ChartVersion) (string, error) {
	var versionId string
	ctx := sender.NewContext(context.Background(), sender.GetSystemUser())
	appManagerClient, err := app.NewAppManagerClient(ctx)
	if err != nil {
		return versionId, err
	}
	appVersionName := chartVersion.GetVersion()
	if chartVersion.GetAppVersion() != "" {
		appVersionName += fmt.Sprintf(" [%s]", chartVersion.GetAppVersion())
	}
	packageName := chartVersion.URLs[0]
	description := chartVersion.GetDescription()
	req := pb.DescribeAppVersionsRequest{}
	req.AppId = []string{appId}
	req.Owner = []string{owner}
	req.Name = []string{appVersionName}
	res, err := appManagerClient.DescribeAppVersions(ctx, &req)
	if err != nil {
		return versionId, err
	}
	if res.TotalCount == 0 {
		createReq := pb.CreateAppVersionRequest{}
		createReq.AppId = utils.ToProtoString(appId)
		createReq.Owner = utils.ToProtoString(owner)
		createReq.Name = utils.ToProtoString(appVersionName)
		createReq.PackageName = utils.ToProtoString(packageName)
		createReq.Description = utils.ToProtoString(description)

		createRes, err := appManagerClient.CreateAppVersion(ctx, &createReq)
		if err != nil {
			return versionId, err
		}
		versionId = createRes.GetAppVersion().GetVersionId().GetValue()
		return versionId, err
	} else {
		modifyReq := pb.ModifyAppVersionRequest{}
		modifyReq.VersionId = res.AppVersionSet[0].VersionId
		modifyReq.PackageName = utils.ToProtoString(packageName)
		modifyReq.Description = utils.ToProtoString(description)

		modifyRes, err := appManagerClient.ModifyAppVersion(ctx, &modifyReq)
		if err != nil {
			return versionId, err
		}
		versionId = modifyRes.GetAppVersion().GetVersionId().GetValue()
		return versionId, err
	}
}
