package indexer

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var IndexYaml = "index.yaml"

type Indexer interface {
	IndexRepo() error
	DeleteRepo() error
}

func GetIndexer(ctx context.Context, repo *pb.Repo) Indexer {
	var i Indexer
	providers := repo.GetProviders()
	repoInterface, err := repoiface.New(ctx, repo.Type.GetValue(), repo.Url.GetValue(), repo.Credential.GetValue())
	if err != nil {
		panic(fmt.Sprintf("failed to get repo interface from repo [%s]", repo.RepoId.GetValue()))
	}

	if stringutil.StringIn(constants.ProviderKubernetes, providers) {
		i = NewHelmIndexer(newIndexer(ctx, repo, repoInterface))
	} else {
		i = NewDevkitIndexer(newIndexer(ctx, repo, repoInterface))
	}
	return i
}

type indexer struct {
	ctx           context.Context
	repo          *pb.Repo
	repoInterface repoiface.RepoInterface
}

func newIndexer(ctx context.Context, repo *pb.Repo, repoInterface repoiface.RepoInterface) indexer {
	return indexer{
		ctx:           ctx,
		repo:          repo,
		repoInterface: repoInterface,
	}
}

type appInterface interface {
	GetName() string
	GetDescription() string
	GetIcon() string
	GetHome() string
	GetSources() string
	GetKeywords() string
	GetMaintainers() string
	GetScreenshots() string
	GetStatus() string
}
type versionInterface interface {
	GetVersion() string
	GetAppVersion() string
	GetDescription() string
	GetUrls() string
	GetStatus() string
}

func (i *indexer) syncAppInfo(app appInterface) (string, error) {
	chartName := app.GetName()
	repoId := i.repo.GetRepoId().GetValue()
	owner := i.repo.GetOwner().GetValue()

	var appId string
	ctx := client.SetSystemUserToContext(i.ctx)
	appManagerClient, err := appclient.NewAppManagerClient()
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
	description := pbutil.ToProtoString(app.GetDescription())
	icon := pbutil.ToProtoString(app.GetIcon())
	home := pbutil.ToProtoString(app.GetHome())
	sources := pbutil.ToProtoString(app.GetSources())
	keywords := pbutil.ToProtoString(app.GetKeywords())
	maintainers := pbutil.ToProtoString(app.GetMaintainers())
	screenshots := pbutil.ToProtoString(app.GetScreenshots())
	status := pbutil.ToProtoString(app.GetStatus())

	var enabledCategoryIds []string
	var disabledCategoryIds []string

	for _, c := range i.repo.GetCategorySet() {
		switch c.Status.GetValue() {
		case constants.StatusEnabled:
			enabledCategoryIds = append(enabledCategoryIds, c.CategoryId.GetValue())
		case constants.StatusDisabled:
			disabledCategoryIds = append(disabledCategoryIds, c.CategoryId.GetValue())
		}
	}
	if len(enabledCategoryIds) == 0 {
		enabledCategoryIds = append(enabledCategoryIds, models.UncategorizedId)
	}

	if res.TotalCount == 0 {
		createReq := pb.CreateAppRequest{}
		createReq.RepoId = pbutil.ToProtoString(repoId)
		createReq.ChartName = pbutil.ToProtoString(chartName)
		createReq.Name = pbutil.ToProtoString(chartName)
		createReq.Description = description
		createReq.Icon = icon
		createReq.Home = home
		createReq.Sources = sources
		createReq.Keywords = keywords
		createReq.Maintainers = maintainers
		createReq.Screenshots = screenshots
		createReq.CategoryId = pbutil.ToProtoString(strings.Join(enabledCategoryIds, ","))
		createReq.Status = status

		createRes, err := appManagerClient.CreateApp(ctx, &createReq)
		if err != nil {
			return appId, err
		}
		appId = createRes.GetAppId().GetValue()
		return appId, err
	} else {
		app := res.AppSet[0]
		var categoryMap = make(map[string]bool)
		for _, c := range app.GetCategorySet() {
			categoryId := c.GetCategoryId().GetValue()
			// app follow repo's categories:
			// if repo *disable* some categories, app MUST *disable* it
			// if repo *enable*  some categories, app MUST *enable*  it
			if c.GetStatus().GetValue() == constants.StatusEnabled {
				if !stringutil.StringIn(categoryId, disabledCategoryIds) {
					categoryMap[categoryId] = true
				}
			}
		}
		for _, c := range enabledCategoryIds {
			categoryMap[c] = true
		}
		var categoryIds []string
		for c := range categoryMap {
			if c == models.UncategorizedId && len(categoryMap) > 1 {
				continue
			}
			categoryIds = append(categoryIds, c)
		}

		modifyReq := pb.ModifyAppRequest{}
		modifyReq.AppId = app.AppId
		modifyReq.Name = pbutil.ToProtoString(chartName)
		modifyReq.ChartName = pbutil.ToProtoString(chartName)
		modifyReq.Description = description
		modifyReq.Icon = icon
		modifyReq.Home = home
		modifyReq.Sources = sources
		modifyReq.Keywords = keywords
		modifyReq.Maintainers = maintainers
		modifyReq.Screenshots = screenshots
		modifyReq.CategoryId = pbutil.ToProtoString(strings.Join(categoryIds, ","))

		modifyRes, err := appManagerClient.ModifyApp(ctx, &modifyReq)
		if err != nil {
			return appId, err
		}
		appId = modifyRes.GetAppId().GetValue()
		return appId, err
	}
}

func (i *indexer) syncAppVersionInfo(appId string, version versionInterface, index int) (string, error) {

	var versionId string
	ctx := client.SetSystemUserToContext(i.ctx)
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return versionId, err
	}
	appVersionName := version.GetVersion()
	if version.GetAppVersion() != "" {
		appVersionName += fmt.Sprintf(" [%s]", version.GetAppVersion())
	}

	owner := pbutil.ToProtoString(i.repo.GetOwner().GetValue())
	name := pbutil.ToProtoString(appVersionName)
	packageName := pbutil.ToProtoString(repoiface.GetFileName(version.GetUrls()))
	description := pbutil.ToProtoString(version.GetDescription())
	sequence := pbutil.ToProtoUInt32(uint32(index))
	status := pbutil.ToProtoString(version.GetStatus())
	req := pb.DescribeAppVersionsRequest{}
	req.AppId = []string{appId}
	req.Owner = []string{owner.Value}
	req.Name = []string{appVersionName}
	res, err := appManagerClient.DescribeAppVersions(ctx, &req)
	if err != nil {
		return versionId, err
	}
	if res.TotalCount == 0 {
		createReq := pb.CreateAppVersionRequest{}
		createReq.AppId = pbutil.ToProtoString(appId)
		createReq.Owner = owner
		createReq.Name = name
		createReq.PackageName = packageName
		createReq.Description = description
		createReq.Sequence = sequence
		createReq.Status = status

		createRes, err := appManagerClient.CreateAppVersion(ctx, &createReq)
		if err != nil {
			return versionId, err
		}
		versionId = createRes.GetVersionId().GetValue()
		return versionId, err
	} else {
		modifyReq := pb.ModifyAppVersionRequest{}
		modifyReq.VersionId = res.AppVersionSet[0].VersionId
		modifyReq.PackageName = packageName
		modifyReq.Description = description
		modifyReq.Sequence = sequence

		modifyRes, err := appManagerClient.ModifyAppVersion(ctx, &modifyReq)
		if err != nil {
			return versionId, err
		}
		versionId = modifyRes.GetVersionId().GetValue()
		return versionId, err
	}
}

func (i *indexer) DeleteRepo() error {
	ctx := client.SetSystemUserToContext(i.ctx)
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return err
	}
	limit := 50
	for {
		req := pb.DescribeAppsRequest{}
		req.RepoId = []string{i.repo.GetRepoId().GetValue()}
		req.Status = []string{constants.StatusActive}
		req.Limit = uint32(limit)
		res, err := appManagerClient.DescribeApps(ctx, &req)
		if err != nil {
			return err
		}
		if len(res.GetAppSet()) == 0 {
			break
		}
		var appIds []string
		for _, app := range res.GetAppSet() {
			appIds = append(appIds, app.GetAppId().GetValue())
		}
		deleteReq := pb.DeleteAppsRequest{}
		deleteReq.AppId = appIds
		_, err = appManagerClient.DeleteApps(ctx, &deleteReq)
		if err != nil {
			return err
		}
	}
	return nil
}
