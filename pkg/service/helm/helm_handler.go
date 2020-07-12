package helm

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"

	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func getActionConfig(ctx context.Context, runtimeId, driver string, getter CredentialGetter) (*action.Configuration, error) {
	file, err := ioutil.TempFile("", "config")
	defer os.Remove(file.Name())
	if err != nil {
		logger.Debug(ctx, "get helm config error: [%s]", err.Error())
		return nil, err
	}

	ns, credentialContent, err := getter(ctx, runtimeId)
	if len(credentialContent) == 0 {
		return nil, err
	}

	_, err = file.Write([]byte(credentialContent))
	if err != nil {
		logger.Debug(ctx, "write crendential content error: [%s]", err.Error())
		return nil, err
	}
	kubeConfigPath := file.Name()
	actionConfig := new(action.Configuration)

	// todo
	var FMT = func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}
	//todo context
	if err := actionConfig.Init(kube.GetConfig(kubeConfigPath, "", ns), ns, driver, FMT); err != nil {
		logger.Debug(ctx, "Init ActionConfig Error: [%s]", err.Error())
		return nil, err
	}
	return actionConfig, nil

}

func (server *HelmServer) CreateRelease(ctx context.Context, req *pb.CreateReleaseRequest) (*pb.CreateReleaseResponse, error) {
	releaseName := req.GetReleaseName().GetValue()
	namespace := req.GetNamespace().GetValue()
	//appId := req.GetAppId()
	versionId := req.GetVersionId()
	runtimeId := req.GetRuntimeId().GetValue()
	cfg, err := getActionConfig(ctx, runtimeId, Driver, DefaultCredentialGetter)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	appClient, err := app.NewAppManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	pkgReq := &pb.GetAppVersionPackageRequest{
		VersionId: versionId,
	}
	pkgResp, err := appClient.GetAppVersionPackage(ctx, pkgReq)
	pkg := pkgResp.GetPackage()
	r := bytes.NewReader(pkg)
	chrt, err := loader.LoadArchive(r)
	releaseInfo, err := CreateRelease(cfg, releaseName, namespace, chrt, nil)
	if err != nil {
		logger.Error(ctx, "Create release failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	createReleaseResp := &pb.CreateReleaseResponse{
		ReleaseName: pbutil.ToProtoString(releaseInfo.Name),
	}
	return createReleaseResp, nil
}

func (server *HelmServer) ListReleases(ctx context.Context, req *pb.ListReleasesRequest) (*pb.ListReleaseResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	cfg, err := getActionConfig(ctx, runtimeId, Driver, DefaultCredentialGetter)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	releaseName := req.GetReleaseName().GetValue()
	releaseInfo, err := GetRelease(cfg, releaseName)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	pbRelease := &pb.Release{
		ReleaseName:       pbutil.ToProtoString(releaseInfo.Name),
		Version:           pbutil.ToProtoInt32(int32(releaseInfo.Version)),
		Namespace:         pbutil.ToProtoString(releaseInfo.Namespace),
		Status:            pbutil.ToProtoString(releaseInfo.Info.Status.String()),
		Description:       pbutil.ToProtoString(releaseInfo.Info.Description),
		FirstDeployedTime: pbutil.ToProtoTimestamp(releaseInfo.Info.FirstDeployed.Time),
		LastDeployedTime:  pbutil.ToProtoTimestamp(releaseInfo.Info.LastDeployed.Time),
		DeletedTime:       pbutil.ToProtoTimestamp(releaseInfo.Info.Deleted.Time),
	}
	listReleaseResp := &pb.ListReleaseResponse{
		ReleaseSet: []*pb.Release{pbRelease},
	}
	return listReleaseResp, nil
}

func (server *HelmServer) UpgradeRelease(ctx context.Context, req *pb.UpgradeReleaseRequest) (*pb.UpgradeReleaseResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	versionId := req.GetVersionId().GetValue()
	releaseName := req.GetReleaseName().GetValue()
	cfg, err := getActionConfig(ctx, runtimeId, Driver, DefaultCredentialGetter)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	appClient, err := app.NewAppManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	pkgReq := &pb.GetAppVersionPackageRequest{
		VersionId: pbutil.ToProtoString(versionId),
	}
	pkgResp, err := appClient.GetAppVersionPackage(ctx, pkgReq)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	pkg := pkgResp.GetPackage()
	r := bytes.NewReader(pkg)
	chrt, err := loader.LoadArchive(r)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	releaseInfo, err := UpgradeRelease(cfg, releaseName, chrt, nil)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	upgradeReleaseResp := &pb.UpgradeReleaseResponse{
		ReleaseName: pbutil.ToProtoString(releaseInfo.Name),
	}
	return upgradeReleaseResp, nil
}

func (server *HelmServer) RollbackRelease(ctx context.Context, req *pb.RollbackReleaseRequest) (*pb.RollbackReleaseResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	releaseName := req.GetReleaseName().GetValue()
	version := req.GetVersion().GetValue()
	cfg, err := getActionConfig(ctx, runtimeId, Driver, DefaultCredentialGetter)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	releaseInfo, err := RollbackRelease(cfg, releaseName, version)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	rollbackReleaseResp := &pb.RollbackReleaseResponse{
		ReleaseName: pbutil.ToProtoString(releaseInfo.Name),
		Version:     pbutil.ToProtoInt32(int32(releaseInfo.Version)),
	}

	return rollbackReleaseResp, nil

}

func (server *HelmServer) DeleteRelease(ctx context.Context, req *pb.DeleteReleaseRequest) (*pb.DeleteReleaseResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	releaseName := req.GetReleaseName().GetValue()
	keepHistory := !req.GetPurge().Value
	cfg, err := getActionConfig(ctx, runtimeId, Driver, DefaultCredentialGetter)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	err = DeleteRelease(cfg, releaseName, keepHistory)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	deleteResp := &pb.DeleteReleaseResponse{
		ReleaseName: pbutil.ToProtoString(releaseName),
	}

	return deleteResp, nil
}
