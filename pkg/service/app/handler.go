// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/client/attachment"
	repoClient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/archiveutil"
	"openpitrix.io/openpitrix/pkg/util/imageutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

var (
	_ pb.AppManagerServer = &Server{}
)

func (p *Server) SyncRepo(ctx context.Context, req *pb.SyncRepoRequest) (*pb.SyncRepoResponse, error) {
	var res = &pb.SyncRepoResponse{}
	var repoId = req.GetRepoId()
	if len(repoId) == 0 {
		return res, nil
	}
	var failed = func(reason string) (*pb.SyncRepoResponse, error) {
		res.Failed = true
		res.Result = reason
		return res, nil
	}

	repoManagerClient, err := repoClient.NewRepoManagerClient()
	if err != nil {
		return failed("internal error")
	}
	describeRepoReq := pb.DescribeReposRequest{
		RepoId: []string{repoId},
	}
	describeRepoRes, err := repoManagerClient.DescribeRepos(ctx, &describeRepoReq)
	if err != nil {
		logger.Error(ctx, "Failed to describe repo [%s], %+v", repoId, err)
		return failed("internal error")
	}
	if describeRepoRes.TotalCount == 0 {
		logger.Error(ctx, "Failed to describe repo [%s], repo not exists", repoId)
		return failed("internal error")
	}
	repo := describeRepoRes.RepoSet[0]
	err = newRepoProxy(repo).SyncRepo(ctx)
	if err != nil {
		return failed(err.Error())
	}
	return res, nil
}

func (p *Server) DescribeActiveApps(ctx context.Context, req *pb.DescribeAppsRequest) (*pb.DescribeAppsResponse, error) {
	return p.describeApps(ctx, req, true)
}

func (p *Server) DescribeApps(ctx context.Context, req *pb.DescribeAppsRequest) (*pb.DescribeAppsResponse, error) {
	return p.describeApps(ctx, req, false)
}

func (p *Server) describeApps(ctx context.Context, req *pb.DescribeAppsRequest, active bool) (*pb.DescribeAppsResponse, error) {
	var apps []*models.App
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	categoryIds := req.GetCategoryId()

	query := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(constants.TableApp).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableApp)).
		Where(db.Eq(constants.ColumnActive, active))

	if len(categoryIds) > 0 {
		subqueryStmt := pi.Global().DB(ctx).
			Select(constants.ColumnResouceId).
			From(constants.TableCategoryResource).
			Where(db.Eq(constants.ColumnStatus, constants.StatusEnabled)).
			Where(db.Eq(constants.ColumnCategoryId, categoryIds))
		query = query.Where(db.Eq(constants.ColumnAppId, []*db.SelectQuery{subqueryStmt}))
	}
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&apps)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	appSet, err := formatAppSet(ctx, apps, active)
	if err != nil {
		return nil, err
	}

	res := &pb.DescribeAppsResponse{
		AppSet:     appSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) ValidatePackage(ctx context.Context, req *pb.ValidatePackageRequest) (*pb.ValidatePackageResponse, error) {
	res := &pb.ValidatePackageResponse{}
	v, err := repoiface.LoadPackage(ctx, req.GetVersionType(), req.GetVersionPackage())
	if err != nil {

		matchPackageFailedError(err, res)

		if res.Error == nil && len(res.ErrorDetails) == 0 {
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
		}
	} else {
		res.Name = pbutil.ToProtoString(v.GetName())
		res.VersionName = pbutil.ToProtoString(v.GetVersionName())
	}
	return res, nil
}

func (p *Server) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	name := req.GetName().GetValue()

	pkg := req.GetVersionPackage().GetValue()

	v, err := repoiface.LoadPackage(ctx, req.GetVersionType().GetValue(), pkg)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
	}
	if v.GetName() != name {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorAppNameConflictWithPackage)
	}

	attachmentContent, err := archiveutil.Load(bytes.NewReader(pkg))
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
	}

	attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	uploadAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
		AttachmentContent: attachmentContent,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var iconAttachmentId string
	if req.GetIcon() != nil {
		icon := req.GetIcon().GetValue()
		content, err := imageutil.Thumbnail(ctx, icon)
		if err != nil {
			logger.Error(ctx, "Make thumbnail failed: %+v", err)
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorImageDecodeFailed)
		}
		createAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
			AttachmentContent: content,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		iconAttachmentId = createAttachmentRes.AttachmentId
	}

	attachmentId := uploadAttachmentRes.AttachmentId

	s := senderutil.GetSenderFromContext(ctx)
	newApp := models.NewApp(
		name,
		s.UserId)

	if len(iconAttachmentId) > 0 {
		newApp.Icon = iconAttachmentId
	}

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableApp).
		Record(newApp).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateAppResponse{
		AppId: pbutil.ToProtoString(newApp.AppId),
	}

	versionName := req.GetVersionName().GetValue()
	if versionName == "" {
		versionName = v.GetVersionName()
	}
	version := models.NewAppVersion(newApp.AppId, versionName, newApp.Description, newApp.Owner)
	version.PackageName = attachmentId
	version.Type = req.GetVersionType().GetValue()

	err = insertVersion(ctx, version)
	if err != nil {
		return nil, err
	}

	res.VersionId = pbutil.ToProtoString(version.VersionId)

	return res, nil
}

func (p *Server) ModifyApp(ctx context.Context, req *pb.ModifyAppRequest) (*pb.ModifyAppResponse, error) {
	appId := req.GetAppId().GetValue()
	app, err := CheckAppPermission(ctx, appId)
	if err != nil {
		return nil, err
	}

	if app.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, appId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "description", "home", "maintainers", "sources", "readme", "keywords")

	err = updateApp(ctx, appId, attributes)
	if err != nil {
		return nil, err
	}

	if req.GetCategoryId() != nil {
		err = categoryutil.SyncResourceCategories(
			ctx,
			pi.Global().DB(ctx),
			appId,
			categoryutil.DecodeCategoryIds(req.GetCategoryId().GetValue()),
		)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}

	res := &pb.ModifyAppResponse{
		AppId: req.GetAppId(),
	}
	return res, nil
}

func (p *Server) UploadAppAttachment(ctx context.Context, req *pb.UploadAppAttachmentRequest) (*pb.UploadAppAttachmentResponse, error) {
	appId := req.GetAppId().GetValue()
	app, err := CheckAppPermission(ctx, appId)
	if err != nil {
		return nil, err
	}

	if app.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, appId)
	}

	attributes := make(map[string]interface{})

	attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	var res = &pb.UploadAppAttachmentResponse{
		AppId: req.GetAppId(),
	}

	switch req.Type {
	case pb.UploadAppAttachmentRequest_icon:
		content, err := imageutil.Thumbnail(ctx, req.GetAttachmentContent().GetValue())
		if err != nil {
			logger.Error(ctx, "Make thumbnail failed: %+v", err)
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorImageDecodeFailed)
		}
		if app.Icon == "" {
			createAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
				AttachmentContent: content,
			})
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
			attributes[constants.ColumnIcon] = createAttachmentRes.AttachmentId
		} else {
			_, err := attachmentManagerClient.ReplaceAttachment(ctx, &pb.ReplaceAttachmentRequest{
				AttachmentId:      app.Icon,
				AttachmentContent: content,
			})
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
		return res, nil
	case pb.UploadAppAttachmentRequest_screenshot:
		var screenshots = strings.Split(app.Screenshots, ",")
		var isDelete = len(req.GetAttachmentContent().GetValue()) == 0
		var seq = int(req.GetSequence().GetValue())
		if seq > 5 || seq < 0 {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorUnsupportedParameterValue, fmt.Sprint(seq))
		}
		if isDelete {
			if len(screenshots) == 0 || len(screenshots) < seq {
				return res, nil
			}
			screenshots = append(screenshots[:seq], screenshots[seq+1:]...)
		} else {
			content, err := imageutil.Thumbnail(ctx, req.GetAttachmentContent().GetValue())
			if err != nil {
				logger.Error(ctx, "Make thumbnail failed: %+v", err)
				return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorImageDecodeFailed)
			}
			if len(screenshots) > seq {
				// if len(screenshots) == 5
				//    seq == 4 , replace screenshots[4]
				_, err := attachmentManagerClient.ReplaceAttachment(ctx, &pb.ReplaceAttachmentRequest{
					AttachmentId:      screenshots[seq],
					AttachmentContent: content,
				})
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
				}
			} else {
				createAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
					AttachmentContent: content,
				})
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
				}
				screenshots = append(screenshots, createAttachmentRes.AttachmentId)
			}
		}
		attributes[constants.ColumnScreenshots] = strings.Join(screenshots, ",")
		return res, nil
	}

	if len(attributes) > 0 {
		err = updateApp(ctx, appId, attributes)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (p *Server) DeleteApps(ctx context.Context, req *pb.DeleteAppsRequest) (*pb.DeleteAppsResponse, error) {
	appIds := req.GetAppId()
	_, err := CheckAppsPermission(ctx, appIds)
	if err != nil {
		return nil, err
	}
	// check permission
	for _, appId := range appIds {
		count, err := pi.Global().DB(ctx).
			Select().
			From(constants.TableAppVersion).
			Where(db.Eq(constants.ColumnAppId, appId)).
			Where(db.Eq(constants.ColumnActive, false)).
			Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
			Count()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
		if count > 0 {
			return nil, gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorExistsNoDeleteVersions, appId)
		}
	}
	_, err = pi.Global().DB(ctx).
		Update(constants.TableApp).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(constants.ColumnAppId, appIds)).
		Where(db.Eq(constants.ColumnActive, false)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pb.DeleteAppsResponse{
		AppId: appIds,
	}, nil
}

func (p *Server) CreateAppVersion(ctx context.Context, req *pb.CreateAppVersionRequest) (*pb.CreateAppVersionResponse, error) {
	pkg := req.GetPackage().GetValue()

	_, err := repoiface.LoadPackage(ctx, req.GetType().GetValue(), pkg)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
	}

	attachmentContent, err := archiveutil.Load(bytes.NewReader(pkg))
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
	}

	attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	uploadAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
		AttachmentContent: attachmentContent,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	attachmentId := uploadAttachmentRes.AttachmentId

	s := senderutil.GetSenderFromContext(ctx)
	newAppVersion := models.NewAppVersion(
		req.GetAppId().GetValue(),
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		s.UserId)

	newAppVersion.PackageName = attachmentId
	newAppVersion.Type = req.GetType().GetValue()
	err = insertVersion(ctx, newAppVersion)
	if err != nil {
		return nil, err
	}

	res := &pb.CreateAppVersionResponse{
		VersionId: pbutil.ToProtoString(newAppVersion.VersionId),
	}
	return res, nil

}

func (p *Server) DescribeAppVersionAudits(ctx context.Context, req *pb.DescribeAppVersionAuditsRequest) (*pb.DescribeAppVersionAuditsResponse, error) {
	_, err := CheckAppPermission(ctx, req.GetAppId())
	if err != nil {
		return nil, err
	}

	_, err = CheckAppVersionPermission(ctx, req.GetVersionId())
	if err != nil {
		return nil, err
	}

	var versionAudits []*models.AppVersionAudit
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().DB(ctx).
		Select(models.AppVersionAuditColumns...).
		From(constants.TableAppVersionAudit).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableAppVersionAudit))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnStatusTime)
	_, err = query.Load(&versionAudits)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeAppVersionAuditsResponse{
		AppVersionAuditSet: models.AppVersionAuditsToPbs(versionAudits),
		TotalCount:         count,
	}
	return res, nil
}

func (p *Server) DescribeActiveAppVersions(ctx context.Context, req *pb.DescribeAppVersionsRequest) (*pb.DescribeAppVersionsResponse, error) {
	return p.describeAppVersions(ctx, req, true)
}

func (p *Server) DescribeAppVersions(ctx context.Context, req *pb.DescribeAppVersionsRequest) (*pb.DescribeAppVersionsResponse, error) {
	return p.describeAppVersions(ctx, req, false)
}

func (p *Server) describeAppVersions(ctx context.Context, req *pb.DescribeAppVersionsRequest, active bool) (*pb.DescribeAppVersionsResponse, error) {
	var versions []*models.AppVersion
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableAppVersion)).
		Where(db.Eq(constants.ColumnActive, active))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnSequence)
	_, err := query.Load(&versions)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeAppVersionsResponse{
		AppVersionSet: models.AppVersionsToPbs(versions),
		TotalCount:    count,
	}
	return res, nil

}

func (p *Server) ModifyAppVersion(ctx context.Context, req *pb.ModifyAppVersionRequest) (*pb.ModifyAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}

	err = checkAppVersionHandlePermission(ctx, Modify, version)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req, "name", "description")

	if version.Status == constants.StatusRejected {
		attributes[constants.ColumnStatus] = constants.StatusDraft
		defer addAppVersionAudit(ctx, version, constants.StatusDraft, constants.RoleDeveloper, "")
	}

	err = updateVersion(ctx, versionId, attributes)
	if err != nil {
		return nil, err
	}

	pkg := req.GetPackage().GetValue()
	if len(pkg) > 0 {
		_, err = repoiface.LoadPackage(ctx, version.Type, pkg)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
		}

		attachmentContent, err := archiveutil.Load(bytes.NewReader(pkg))
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorPackageParseFailed)
		}

		attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		_, err = attachmentManagerClient.ReplaceAttachment(ctx, &pb.ReplaceAttachmentRequest{
			AttachmentId:      version.PackageName,
			AttachmentContent: attachmentContent,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	res := &pb.ModifyAppVersionResponse{
		VersionId: req.GetVersionId(),
	}
	return res, nil

}

func (p *Server) GetAppVersionPackage(ctx context.Context, req *pb.GetAppVersionPackageRequest) (*pb.GetAppVersionPackageResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := getAppVersion(ctx, versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}
	logger.Debug(ctx, "Got app version: [%+v]", version)
	content, err := newVersionProxy(version).GetPackageFile(ctx)
	if err != nil {
		return nil, err
	}
	archiveFiles, err := archiveutil.Save(content, versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	return &pb.GetAppVersionPackageResponse{
		Package:   archiveFiles,
		AppId:     pbutil.ToProtoString(version.AppId),
		VersionId: req.GetVersionId(),
	}, nil
}

func (p *Server) GetAppVersionPackageFiles(ctx context.Context, req *pb.GetAppVersionPackageFilesRequest) (*pb.GetAppVersionPackageFilesResponse, error) {
	versionId := req.GetVersionId().GetValue()
	includeFiles := req.Files
	version, err := getAppVersion(ctx, versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}

	archiveFiles, err := newVersionProxy(version).GetPackageFile(ctx, includeFiles...)
	return &pb.GetAppVersionPackageFilesResponse{
		Files:     archiveFiles,
		VersionId: req.GetVersionId(),
	}, nil
}

type appStatistic struct {
	Date  string `db:"DATE_FORMAT(create_time, '%Y-%m-%d')"`
	Count uint32 `db:"COUNT(app_id)"`
}
type repoStatistic struct {
	RepoId string `db:"repo_id"`
	Count  uint32 `db:"COUNT(app_id)"`
}

func (p *Server) GetAppStatistics(ctx context.Context, req *pb.GetAppStatisticsRequest) (*pb.GetAppStatisticsResponse, error) {
	res := &pb.GetAppStatisticsResponse{
		LastTwoWeekCreated: make(map[string]uint32),
		TopTenRepos:        make(map[string]uint32),
	}
	appCount, err := pi.Global().DB(ctx).
		Select(constants.ColumnAppId).
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Where(db.Eq(constants.ColumnActive, false)).
		Count()
	if err != nil {
		logger.Error(ctx, "Failed to get app count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.AppCount = appCount

	err = pi.Global().DB(ctx).
		Select("COUNT(DISTINCT repo_id)").
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Where(db.Eq(constants.ColumnActive, false)).
		LoadOne(&res.RepoCount)
	if err != nil {
		logger.Error(ctx, "Failed to get repo count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var as []*appStatistic
	_, err = pi.Global().DB(ctx).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(app_id)").
		From(constants.TableApp).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(constants.ColumnCreateTime, time2week)).
		Where(db.Eq(constants.ColumnActive, false)).
		Limit(14).Load(&as)

	if err != nil {
		logger.Error(ctx, "Failed to get app statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range as {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var rs []*repoStatistic
	_, err = pi.Global().DB(ctx).
		Select("repo_id", "COUNT(app_id)").
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Where(db.Eq(constants.ColumnActive, false)).
		GroupBy(constants.ColumnRepoId).
		OrderDir("COUNT(app_id)", false).
		Limit(10).Load(&rs)

	if err != nil {
		logger.Error(ctx, "Failed to get repo statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.TopTenRepos[a.RepoId] = a.Count
	}

	return res, nil
}

func (p *Server) SubmitAppVersion(ctx context.Context, req *pb.SubmitAppVersionRequest) (*pb.SubmitAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}

	err = checkAppVersionHandlePermission(ctx, Submit, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusSubmitted)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusSubmitted, constants.RoleDeveloper, "")
	if err != nil {
		return nil, err
	}
	res := pb.SubmitAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) CancelAppVersion(ctx context.Context, req *pb.CancelAppVersionRequest) (*pb.CancelAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Cancel, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusDraft)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusDraft, constants.RoleDeveloper, "")
	if err != nil {
		return nil, err
	}
	res := pb.CancelAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) ReleaseAppVersion(ctx context.Context, req *pb.ReleaseAppVersionRequest) (*pb.ReleaseAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Release, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusActive)
	if err != nil {
		return nil, err
	}
	err = syncAppVersion(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusActive, constants.RoleDeveloper, "")
	if err != nil {
		return nil, err
	}
	res := pb.ReleaseAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) DeleteAppVersion(ctx context.Context, req *pb.DeleteAppVersionRequest) (*pb.DeleteAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	// check permission
	err = checkAppVersionHandlePermission(ctx, Delete, version)
	if err != nil {
		return nil, err
	}

	err = updateVersionStatus(ctx, version, constants.StatusDeleted)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusDeleted, constants.RoleDeveloper, "")
	if err != nil {
		return nil, err
	}
	res := pb.DeleteAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) PassAppVersion(ctx context.Context, req *pb.PassAppVersionRequest) (*pb.PassAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Pass, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusPassed)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusPassed, constants.RoleGlobalAdmin, "")
	if err != nil {
		return nil, err
	}
	res := pb.PassAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) RejectAppVersion(ctx context.Context, req *pb.RejectAppVersionRequest) (*pb.RejectAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Reject, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusRejected, map[string]interface{}{
		constants.ColumnMessage: req.GetMessage().GetValue(),
	})
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusRejected, constants.RoleGlobalAdmin, req.GetMessage().GetValue())
	if err != nil {
		return nil, err
	}
	res := pb.RejectAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) SuspendAppVersion(ctx context.Context, req *pb.SuspendAppVersionRequest) (*pb.SuspendAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Suspend, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusSuspended)
	if err != nil {
		return nil, err
	}
	err = syncAppVersion(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusSuspended, constants.RoleGlobalAdmin, "")
	if err != nil {
		return nil, err
	}
	res := pb.SuspendAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) RecoverAppVersion(ctx context.Context, req *pb.RecoverAppVersionRequest) (*pb.RecoverAppVersionResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := CheckAppVersionPermission(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = checkAppVersionHandlePermission(ctx, Recover, version)
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusActive)
	if err != nil {
		return nil, err
	}
	err = syncAppVersion(ctx, versionId)
	if err != nil {
		return nil, err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusActive, constants.RoleGlobalAdmin, "")
	if err != nil {
		return nil, err
	}
	res := pb.RecoverAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}
