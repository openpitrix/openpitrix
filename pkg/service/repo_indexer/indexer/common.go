package indexer

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

type Indexer interface {
	IndexRepo() error
}

func GetIndexer(repo *pb.Repo) Indexer {
	var i Indexer
	providers := repo.GetProviders()
	if utils.StringIn(constants.ProviderKubernetes, providers) {
		i = NewHelmIndexer(repo)
	} else {
		i = NewDevkitIndexer(repo)
	}
	return i
}

type indexer struct {
	repo *pb.Repo
}
type appInterface interface {
	GetName() string
	GetDescription() string
	GetIcon() string
	GetHome() string
	GetSources() []string
}
type versionInterface interface {
	GetVersion() string
	GetAppVersion() string
	GetDescription() string
	GetUrls() []string
}

func (i *indexer) syncAppInfo(app appInterface) (string, error) {
	chartName := app.GetName()
	repoId := i.repo.GetRepoId().GetValue()
	owner := i.repo.GetOwner().GetValue()

	var appId string
	ctx := client.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
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
	description = utils.ToProtoString(app.GetDescription())
	icon = utils.ToProtoString(app.GetIcon())
	home = utils.ToProtoString(app.GetHome())
	sources = utils.ToProtoString(strings.Join(app.GetSources(), ","))
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

func (i *indexer) syncAppVersionInfo(appId string, version versionInterface) (string, error) {
	owner := i.repo.GetOwner().GetValue()

	var versionId string
	ctx := client.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
	if err != nil {
		return versionId, err
	}
	appVersionName := version.GetVersion()
	if version.GetAppVersion() != "" {
		appVersionName += fmt.Sprintf(" [%s]", version.GetAppVersion())
	}
	packageName := version.GetUrls()[0]
	description := version.GetDescription()
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
