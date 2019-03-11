// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

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
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendor models.VendorVerifyInfo //need use a appvendor object to call function
	vendorPbSet := vendor.ParseVendorSet2PbSet(ctx, vendors)

	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: vendorPbSet,
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

	ifExist, err := s.checkIsExist(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], %+v", req.UserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if ifExist {
		_, err := CheckAppVendorPermission(ctx, appVendorUserId)
		if err != nil {
			return nil, err
		}
		appVendorUserId, err = UpdateVendorVerifyInfo(ctx, req)

		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%s], %+v", appVendorUserId, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
		}
	} else {
		vendor := &models.VendorVerifyInfo{}
		vendor = vendor.ParseReq2Vendor(ctx, req)
		appVendorUserId, err = CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%+v], %+v", vendor, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
		logger.Debug(ctx, "vendorVerifyInfo does not exit, create new vendorVerifyInfo verify info, [%+v]", vendor)

		var emailNotifications []*models.EmailNotification
		// notify admin
		adminUsers, err := accountclient.GetRoleUsers(ctx, []string{constants.RoleGlobalAdmin})
		if err != nil {
			logger.Error(ctx, "Failed to describe role [%s] users: %+v", constants.RoleGlobalAdmin, err)
		} else {
			for _, adminUser := range adminUsers {
				emailNotifications = append(emailNotifications, &models.EmailNotification{
					Title:       constants.SubmitVendorNotifyAdminTitle.GetDefaultMessage(vendor.CompanyName),
					Content:     constants.SubmitVendorNotifyAdminContent.GetDefaultMessage(adminUser.GetUsername().GetValue(), vendor.CompanyName),
					Owner:       sender.UserId,
					ContentType: constants.NfContentTypeVerify,
					Addresses:   []string{adminUser.GetEmail().GetValue()},
				})
			}
		}

		// notify isv
		users, err := accountclient.GetUsers(ctx, []string{appVendorUserId})
		if err != nil || len(users) != 1 {
			logger.Error(ctx, "Failed to describe users [%s]: %+v", appVendorUserId, err)
		} else {
			emailNotifications = append(emailNotifications, &models.EmailNotification{
				Title:       constants.SubmitVendorNotifyIsvTitle.GetDefaultMessage(),
				Content:     constants.SubmitVendorNotifyIsvContent.GetDefaultMessage(users[0].GetUsername().GetValue()),
				Owner:       sender.UserId,
				ContentType: constants.NfContentTypeVerify,
				Addresses:   []string{users[0].GetEmail().GetValue()},
			})
		}

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

	systemCtx := clientutil.SetSystemUserToContext(context.Background())
	// Use system user to change the role of appvendor to isv
	accessClient, err := accessclient.NewClient()

	_, err = accessClient.BindUserRole(systemCtx, &pb.BindUserRoleRequest{
		RoleId: []string{constants.RoleIsv},
		UserId: []string{appVendorUserId},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = PassVendorVerifyInfo(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to pass vendorVerifyInfo [%s], %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	// prepared notifications
	var emailNotifications []*models.EmailNotification
	users, err := accountclient.GetUsers(ctx, []string{appVendorUserId})
	if err != nil || len(users) != 1 {
		logger.Error(ctx, "Failed to describe users [%s]: %+v", appVendorUserId, err)
	} else {
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.PassVendorNotifyTitle.GetDefaultMessage(appVendor.CompanyName),
			Content:     constants.PassVendorNotifyContent.GetDefaultMessage(users[0].GetUsername().GetValue(), appVendor.CompanyName),
			Owner:       sender.UserId,
			ContentType: constants.NfContentTypeVerify,
			Addresses:   []string{users[0].GetEmail().GetValue()},
		})
	}

	// send notifications
	nfclient.SendEmailNotification(ctx, emailNotifications)

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
	appVendorUserID, err := RejectVendorVerifyInfo(ctx, appVendorUserId, rejectMsg, approver)
	if err != nil {
		logger.Error(ctx, "Failed to reject vendorVerifyInfo [%s], %+v", appVendorUserID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	var emailNotifications []*models.EmailNotification
	users, err := accountclient.GetUsers(ctx, []string{appVendorUserId})
	if err != nil || len(users) != 1 {
		logger.Error(ctx, "Failed to describe users [%s]: %+v", appVendorUserId, err)
	} else {
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.RejectVendorNotifyTitle.GetDefaultMessage(appVendor.CompanyName),
			Content:     constants.RejectVendorNotifyContent.GetDefaultMessage(users[0].GetUsername().GetValue(), appVendor.CompanyName),
			Owner:       sender.UserId,
			ContentType: constants.NfContentTypeVerify,
			Addresses:   []string{users[0].GetEmail().GetValue()},
		})
	}
	// send notifications
	nfclient.SendEmailNotification(ctx, emailNotifications)

	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(appVendorUserID),
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
		logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], %+v", appVendorUserID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, appVendorUserID)
	}

	var vendor models.VendorVerifyInfo //need use a appvendor object to call function
	vendorPb := vendor.ParseVendor2Pb(ctx, appVendor)
	res := &pb.GetVendorVerifyInfoResponse{
		VendorVerifyInfo: vendorPb,
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
	isUrlFmt, err := VerifyUrl(ctx, url)

	if !isUrlFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.CompanyWebsite.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalUrlFormat, url)
	}

	email := req.AuthorizerEmail.GetValue()
	isEmailFmt, err := VerifyEmailFmt(ctx, email)

	if !isEmailFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerEmail.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalEmailFormat, email)
	}

	phone := req.AuthorizerPhone.GetValue()
	isPhoneFmt, err := VerifyPhoneFmt(ctx, phone)
	if !isPhoneFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerPhone.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalPhoneFormat, phone)
	}

	isBankAccountNumberFmt, err := VerifyBankAccountNumberFmt(ctx, req.BankAccountNumber.GetValue())
	if !isBankAccountNumberFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.BankAccountNumber.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorIllegalBankAccountNumberFormat, req.BankAccountNumber.GetValue())
	}
	return nil
}
