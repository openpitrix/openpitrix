// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

import (
	"context"
	"strings"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	accessclient "openpitrix.io/openpitrix/pkg/client/access"
	accountclient "openpitrix.io/openpitrix/pkg/client/account"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	nfclient "openpitrix.io/openpitrix/pkg/client/notification"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (s *Server) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendor verify info: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	pbVendorSet := models.VendorVerifyInfoSetToPbSet(vendors)

	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: pbVendorSet,
		TotalCount:          vendorCount,
	}
	return res, nil
}

func (s *Server) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	sender := ctxutil.GetSender(ctx)
	err := s.validateSubmitParams(ctx, req)
	if err != nil {
		return nil, err
	}

	appVendorUserId := req.UserId

	accountClient, err := accountclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	accessClient, err := accessclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	vendorUser, err := accountClient.GetUser(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to get vendor user [%s]: %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	isExist, err := s.checkIsExist(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to get vendor [%s] verify info: %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	vendor := models.ReqToVendorVerifyInfo(ctx, req)
	if isExist {
		_, err := CheckAppVendorPermission(ctx, appVendorUserId)
		if err != nil {
			return nil, err
		}
		err = UpdateVendorVerifyInfo(ctx, req)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendor [%s] verify info: %+v", appVendorUserId, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
		}
	} else {
		count, err := GetVendorVerifyInfoCountByCompanyName(ctx, vendor.CompanyName)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
		if count > 0 {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCompanyNameExists, vendor.CompanyName)
		}
		err = CreateVendorVerifyInfo(ctx, vendor)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
		logger.Debug(ctx, "Create new vendor verify info: [%+v]", vendor)
	}

	if !stringutil.StringIn(sender.UserId, constants.InternalUsers) {
		var emailNotifications []*models.EmailNotification

		// notify isv_auth users
		systemCtx := clientutil.SetSystemUserToContext(ctx)
		actionBundleUsers, _ := accessClient.GetActionBundleUsers(systemCtx, []string{constants.ActionBundleIsvAuth})
		platformName := pi.Global().GlobalConfig().BasicCfg.PlatformName
		platformUrl := pi.Global().GlobalConfig().BasicCfg.PlatformUrl
		for _, user := range actionBundleUsers {
			emailNotifications = append(emailNotifications, &models.EmailNotification{
				Title:       constants.SubmitVendorNotifyAdminTitle.GetDefaultMessage(platformName, vendor.CompanyName),
				Content:     constants.SubmitVendorNotifyAdminContent.GetDefaultMessage(platformName, user.GetUsername().GetValue(), vendor.CompanyName, platformUrl, platformUrl, platformUrl),
				Owner:       sender.UserId,
				ContentType: constants.NfContentTypeVerify,
				Addresses:   []string{user.GetEmail().GetValue()},
			})
		}

		// notify isv
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.SubmitVendorNotifyIsvTitle.GetDefaultMessage(platformName),
			Content:     constants.SubmitVendorNotifyIsvContent.GetDefaultMessage(platformName, vendorUser.GetUsername().GetValue(), platformUrl, platformUrl, platformUrl),
			Owner:       sender.UserId,
			ContentType: constants.NfContentTypeVerify,
			Addresses:   []string{vendorUser.GetEmail().GetValue()},
		})

		// send notifications
		nfclient.SendEmailNotification(ctx, emailNotifications)
	}
	res := &pb.SubmitVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(appVendorUserId),
	}
	return res, nil
}

func (s *Server) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	sender := ctxutil.GetSender(ctx)
	appVendorUserId := req.UserId

	appVendor, err := CheckAppVendorPermission(ctx, appVendorUserId)
	if err != nil {
		return nil, err
	}

	accessClient, err := accessclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	accountClient, err := accountclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// Use system vendorUser to change the role of vendor to isv
	systemCtx := clientutil.SetSystemUserToContext(context.Background())
	_, err = accessClient.BindUserRole(systemCtx, &pb.BindUserRoleRequest{
		RoleId: []string{constants.RoleIsv},
		UserId: []string{appVendorUserId},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	err = PassVendorVerifyInfo(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to pass vendor [%s] verify info: %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	if !stringutil.StringIn(sender.UserId, constants.InternalUsers) {
		// prepared notifications
		var emailNotifications []*models.EmailNotification

		vendorUser, err := accountClient.GetUser(ctx, appVendorUserId)
		if err != nil {
			logger.Error(ctx, "Failed to get vendor user [%s]: %+v", appVendorUserId, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		// notify isv
		platformName := pi.Global().GlobalConfig().BasicCfg.PlatformName
		platformUrl := pi.Global().GlobalConfig().BasicCfg.PlatformUrl
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.PassVendorNotifyTitle.GetDefaultMessage(platformName, appVendor.CompanyName),
			Content:     constants.PassVendorNotifyContent.GetDefaultMessage(platformName, vendorUser.GetUsername().GetValue(), appVendor.CompanyName, platformUrl, platformUrl, platformUrl),
			Owner:       sender.UserId,
			ContentType: constants.NfContentTypeVerify,
			Addresses:   []string{vendorUser.GetEmail().GetValue()},
		})
		// send notifications
		nfclient.SendEmailNotification(ctx, emailNotifications)
	}

	res := &pb.PassVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(appVendorUserId),
	}
	return res, nil
}

func (s *Server) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	appVendorUserId := req.GetUserId()
	appVendor, err := CheckAppVendorPermission(ctx, appVendorUserId)
	if err != nil {
		return nil, err
	}

	sender := ctxutil.GetSender(ctx)
	approver := sender.UserId
	rejectMsg := req.GetRejectMessage().GetValue()
	err = RejectVendorVerifyInfo(ctx, appVendorUserId, rejectMsg, approver)
	if err != nil {
		logger.Error(ctx, "Failed to reject vendor [%s] verify info: %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	accountClient, err := accountclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	if !stringutil.StringIn(sender.UserId, constants.InternalUsers) {
		var emailNotifications []*models.EmailNotification

		vendorUser, err := accountClient.GetUser(ctx, appVendorUserId)
		if err != nil {
			logger.Error(ctx, "Failed to get vendor user [%s]: %+v", appVendorUserId, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		platformName := pi.Global().GlobalConfig().BasicCfg.PlatformName
		platformUrl := pi.Global().GlobalConfig().BasicCfg.PlatformUrl

		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.RejectVendorNotifyTitle.GetDefaultMessage(platformName, appVendor.CompanyName),
			Content:     constants.RejectVendorNotifyContent.GetDefaultMessage(platformName, vendorUser.GetUsername().GetValue(), appVendor.CompanyName, platformUrl, platformUrl, platformUrl),
			Owner:       sender.UserId,
			ContentType: constants.NfContentTypeVerify,
			Addresses:   []string{vendorUser.GetEmail().GetValue()},
		})

		// send notifications
		nfclient.SendEmailNotification(ctx, emailNotifications)
	}

	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(appVendorUserId),
	}
	return res, nil
}

func (s *Server) DescribeAppVendorStatistics(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorStatisticsResponse, error) {
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendorStatisticsSet []*models.VendorStatistics
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendor verify info: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	for _, vendor := range vendors {
		vendorStatistics := &models.VendorStatistics{
			UserId:            vendor.UserId,
			CompanyName:       vendor.CompanyName,
			ActiveAppCount:    0,
			ClusterCountMonth: 0,
			ClusterCountTotal: 0,
		}

		vendorCtx, err := clientutil.SetUserToContext(ctx, vendor.UserId, "DescribeApps")
		if err != nil {
			logger.Error(ctx, "Failed to set [%s] as sender: %+v", vendor.UserId, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		describeAppsReq := &pb.DescribeAppsRequest{
			Status: []string{constants.StatusActive},
		}
		// Use vendor's sender to describe apps
		describeAllResponses, err := pbutil.DescribeAllResponses(vendorCtx, new(appclient.DescribeAppsApi), describeAppsReq)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		var appIds []string
		for _, response := range describeAllResponses {
			switch r := response.(type) {
			case *pb.DescribeAppsResponse:
				for _, app := range r.AppSet {
					appIds = append(appIds, app.GetAppId().GetValue())
				}
			default:
				return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorDescribeResourcesFailed)
			}
		}
		vendorStatistics.ActiveAppCount = uint32(len(appIds))

		if vendorStatistics.ActiveAppCount > 0 {
			monthResponse, err := clusterClient.DescribeAppClusters(ctx, &pb.DescribeAppClustersRequest{
				AppId:          appIds,
				CreatedDate:    pbutil.ToProtoUInt32(30),
				DisplayColumns: []string{},
			})
			if err != nil {
				logger.Error(ctx, "Describe app clusters with app id [%s] failed: %+v", strings.Join(appIds, ","), err)
				return nil, err
			}
			vendorStatistics.ClusterCountMonth = monthResponse.TotalCount

			totalResponse, err := clusterClient.DescribeAppClusters(ctx, &pb.DescribeAppClustersRequest{
				AppId:          appIds,
				DisplayColumns: []string{},
			})
			if err != nil {
				logger.Error(ctx, "Describe app clusters failed: %+v", err)
				return nil, err
			}
			vendorStatistics.ClusterCountTotal = totalResponse.TotalCount
		}

		vendorStatisticsSet = append(vendorStatisticsSet, vendorStatistics)
	}

	pbVendorStatisticsSet := models.VendorStatisticsSetToPbSet(vendorStatisticsSet)

	res := &pb.DescribeVendorStatisticsResponse{
		VendorVerifyStatisticsSet: pbVendorStatisticsSet,
		TotalCount:                vendorCount,
	}
	return res, nil
}

func (s *Server) GetVendorVerifyInfo(ctx context.Context, req *pb.GetVendorVerifyInfoRequest) (*pb.GetVendorVerifyInfoResponse, error) {
	appVendorUserID := req.GetUserId().GetValue()
	appVendor, err := GetVendorVerifyInfo(ctx, appVendorUserID)

	var sender = ctxutil.GetSender(ctx)
	if sender != nil {
		if !appVendor.OwnerPath.CheckPermission(sender) {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorResourceAccessDenied, appVendor.UserId)
		}
	}
	if err != nil {
		logger.Error(ctx, "Failed to get vendor [%s] verify info: %+v", appVendorUserID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, appVendorUserID)
	}

	pbVendor := models.VendorVerifyInfoToPb(appVendor)
	res := &pb.GetVendorVerifyInfoResponse{
		VendorVerifyInfo: pbVendor,
	}
	return res, err
}

func (s *Server) checkIsExist(ctx context.Context, userID string) (bool, error) {
	info, err := GetVendorVerifyInfo(ctx, userID)
	if info == nil && err != nil {
		return false, nil
	} else if info != nil && err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (s *Server) validateSubmitParams(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) error {
	url := req.CompanyWebsite.GetValue()
	if err := VerifyUrl(ctx, url); err != nil {
		return err
	}

	email := req.AuthorizerEmail.GetValue()
	if err := VerifyEmail(ctx, email); err != nil {
		return err
	}

	phoneNumber := req.AuthorizerPhone.GetValue()
	if err := VerifyPhoneNumber(ctx, phoneNumber); err != nil {
		return err
	}

	bankAccountNumber := req.BankAccountNumber.GetValue()
	if err := VerifyBankAccountNumber(ctx, bankAccountNumber); err != nil {
		return err
	}
	return nil
}
