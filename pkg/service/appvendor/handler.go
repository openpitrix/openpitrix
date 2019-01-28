// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"math"

	accountclient "openpitrix.io/openpitrix/pkg/client/account"
	amclient "openpitrix.io/openpitrix/pkg/client/am"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	nfclient "openpitrix.io/openpitrix/pkg/client/notification"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
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
		adminUsers, err := amclient.GetRoleUsers(ctx, []string{constants.RoleGlobalAdmin})
		if err != nil {
			logger.Error(ctx, "Failed to describe role [%s] users: %+v", constants.RoleGlobalAdmin, err)
		} else {
			for _, adminUser := range adminUsers {
				emailNotifications = append(emailNotifications, &models.EmailNotification{
					Title:       constants.SubmitVendorNotifyAdminTitle.GetDefaultMessage(vendor.CompanyName),
					Content:     constants.SubmitVendorNotifyAdminContent.GetDefaultMessage(adminUser, vendor.CompanyName),
					Owner:       sender.UserId,
					ContentType: constants.NfContentTypeVerify,
					Addresses:   []string{adminUser.Email},
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

	_, err = PassVendorVerifyInfo(ctx, appVendorUserId)
	if err != nil {
		logger.Error(ctx, "Failed to pass vendorVerifyInfo [%s], %+v", appVendorUserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

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
	appClient, err := appclient.NewAppManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendorStatisticses []*models.VendorStatistics
	vendors, vendorCount, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	/*============================================================================================================*/
	//To get ClusterCountTotal
	var clusterCntAll4AllPages int32
	var clusterCntMonth4AllPages int32
	for _, vendor := range vendors {
		//step1:Get real appCnt for each vendor
		var vendorStatistics models.VendorStatistics
		appVendorUserId := vendor.UserId

		pbApps, appCnt, err := appClient.DescribeAppsWithAppVendorUserId(ctx, appVendorUserId, db.DefaultSelectLimit, 0)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}

		//step2:if the real appCnt is smaller than db.DefaultSelectLimit,there is only one page apps,and the rows of this one page is length of pbApps.
		//Just accumulate the clusterCnt4SingleApp for each app.
		if appCnt <= int32(db.DefaultSelectLimit) {
			for _, pbApp := range pbApps {
				_, clusterCntAll4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), false, db.DefaultSelectLimit, 0)
				_, clusterCntMonth4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), true, db.DefaultSelectLimit, 0)
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
				}
				clusterCntAll4AllPages = clusterCntAll4AllPages + clusterCntAll4SingleApp
				clusterCntMonth4AllPages = clusterCntMonth4AllPages + clusterCntMonth4SingleApp
			}

		} else {
			//step3:if the real appCnt is bigger than db.DefaultSelectLimit(200),there are more than 1 page Apps.
			//Should accumulate the clusterCnt4SingleApp for each apps ,then accumulate the number for each page.
			pages := int(math.Ceil(float64(appCnt / db.DefaultSelectLimit)))
			for i := 1; i <= pages; i++ {
				offset := db.DefaultSelectLimit * i
				pbApps, _, err := appClient.DescribeAppsWithAppVendorUserId(ctx, appVendorUserId, db.DefaultSelectLimit, uint32(offset))
				if err != nil {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
				}

				var clusterCntAll4OnePage int32
				var clusterCntMonth4OnePage int32
				for _, pbApp := range pbApps {
					_, clusterCntAll4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), false, db.DefaultSelectLimit, uint32(offset))
					_, clusterCntMonth4SingleApp, err := clusterClient.DescribeClustersWithAppId(ctx, pbApp.AppId.GetValue(), true, db.DefaultSelectLimit, uint32(offset))
					if err != nil {
						return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
					}
					clusterCntAll4OnePage = clusterCntAll4OnePage + clusterCntAll4SingleApp
					clusterCntMonth4OnePage = clusterCntMonth4OnePage + clusterCntMonth4SingleApp
				}
				clusterCntAll4AllPages = clusterCntAll4AllPages + clusterCntAll4OnePage
				clusterCntMonth4AllPages = clusterCntMonth4AllPages + clusterCntMonth4OnePage
			}

		}

		/*============================================================================================================*/
		vendorStatistics.UserId = vendor.UserId

		vendorStatistics.CompanyName = vendor.CompanyName
		vendorStatistics.ActiveAppCount = int32(appCnt)
		vendorStatistics.ClusterCountTotal = clusterCntAll4AllPages
		vendorStatistics.ClusterCountMonth = clusterCntMonth4AllPages

		vendorStatisticses = append(vendorStatisticses, &vendorStatistics)

	}

	var vendorStatistics models.VendorStatistics //need use a vendorStatistics object to call function
	vendorStatisticsPbSet := vendorStatistics.ParseVendorStatisticsSet2PbSet(ctx, vendorStatisticses)

	res := &pb.DescribeVendorStatisticsResponse{
		VendorVerifyStatisticsSet: vendorStatisticsPbSet,
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
