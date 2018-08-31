package indexer

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/repoiface/wrapper"
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
	repoReader, err := repoiface.NewReader(ctx, repo)
	if err != nil {
		panic(fmt.Sprintf("failed to get repo interface from repo [%s]", repo.RepoId.GetValue()))
	}

	if stringutil.StringIn(constants.ProviderKubernetes, providers) {
		i = NewHelmIndexer(newIndexer(ctx, repo, repoReader))
	} else {
		i = NewDevkitIndexer(newIndexer(ctx, repo, repoReader))
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

func (i *indexer) syncAppInfo(app wrapper.VersionInterface) (string, error) {
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
		createReq.CategoryId = pbutil.ToProtoString(strings.Join(enabledCategoryIds, ","))

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
		modifyReq.CategoryId = pbutil.ToProtoString(strings.Join(categoryIds, ","))

		modifyRes, err := appManagerClient.ModifyApp(ctx, &modifyReq)
		if err != nil {
			return appId, err
		}
		appId = modifyRes.GetAppId().GetValue()
		return appId, err
	}
}

func (i *indexer) syncAppVersionInfo(appId string, version wrapper.VersionInterface, index int) (string, error) {

	var versionId string
	ctx := client.SetSystemUserToContext(i.ctx)
	appManagerClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return versionId, err
	}
	appVersionName := version.GetVersionName()

	owner := pbutil.ToProtoString(i.repo.GetOwner().GetValue())
	name := pbutil.ToProtoString(appVersionName)
	packageName := pbutil.ToProtoString(repoiface.GetFileName(version.GetUrls()))
	description := pbutil.ToProtoString(version.GetDescription())
	sequence := pbutil.ToProtoUInt32(uint32(index))
	icon := pbutil.ToProtoString(version.GetIcon())
	home := pbutil.ToProtoString(version.GetHome())
	sources := pbutil.ToProtoString(version.GetSources())
	keywords := pbutil.ToProtoString(version.GetKeywords())
	maintainers := pbutil.ToProtoString(version.GetMaintainers())
	screenshots := pbutil.ToProtoString(version.GetScreenshots())
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

		createReq.Home = home
		createReq.Icon = icon
		createReq.Screenshots = screenshots
		createReq.Maintainers = maintainers
		createReq.Keywords = keywords
		createReq.Sources = sources

		createRes, err := appManagerClient.CreateAppVersion(ctx, &createReq)
		if err != nil {
			return versionId, err
		}
		versionId = createRes.GetVersionId().GetValue()
		return versionId, err
	} else {
		existVersion := res.AppVersionSet[0]
		modifyReq := pb.ModifyAppVersionRequest{}
		if existVersion.PackageName.GetValue() != packageName.GetValue() {
			modifyReq.PackageName = packageName
		}
		if existVersion.Description.GetValue() != description.GetValue() {
			modifyReq.Description = description
		}
		if existVersion.Sequence.GetValue() != sequence.GetValue() {
			modifyReq.Sequence = sequence
		}
		if existVersion.Home.GetValue() != home.GetValue() {
			modifyReq.Home = home
		}
		if existVersion.Icon.GetValue() != icon.GetValue() {
			modifyReq.Icon = icon
		}
		if existVersion.Screenshots.GetValue() != screenshots.GetValue() {
			modifyReq.Screenshots = screenshots
		}
		if existVersion.Maintainers.GetValue() != maintainers.GetValue() {
			modifyReq.Maintainers = maintainers
		}
		if existVersion.Keywords.GetValue() != keywords.GetValue() {
			modifyReq.Keywords = keywords
		}
		if existVersion.Sources.GetValue() != sources.GetValue() {
			modifyReq.Sources = sources
		}
		if proto.Size(&modifyReq) == 0 {
			return versionId, err
		}

		modifyReq.VersionId = existVersion.VersionId

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
