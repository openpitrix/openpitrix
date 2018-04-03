// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) CreateRuntime(ctx context.Context, req *pb.CreateRuntimeRequest) (*pb.CreateRuntimeResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(req)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntime: %+v", err)
	}

	// create runtime env
	runtimeId, err := p.createRuntime(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetRuntimeUrl().GetValue(),
		s.UserId)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// create labels
	err = p.createRuntimeLabels(runtimeId, req.Labels.GetValue())
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// get response
	pbRuntime, err := p.getRuntimePbWithLabel(runtimeId)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}
	res := &pb.CreateRuntimeResponse{
		Runtime: pbRuntime,
	}
	return res, nil
}

func (p *Server) DescribeRuntimes(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimesResponse, error) {
	// validate req
	err := validateDescribeRuntimeRequest(req)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimes: %+v", err)
	}
	var runtimeIds []string
	// get runtime env ids by selector
	if req.Selector != nil {
		runtimeIds, err = p.getRuntimeIdsBySelectorString(req.Selector.GetValue())
		if err != nil {
			logger.Errorf("DescribeRuntimes: %+v", err)
			return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
		}
	}
	// get runtime envs
	pbRuntimes, count, err := p.getRuntimePbsWithoutLabelByReqAndId(req, runtimeIds)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
	}
	// get runtime envs label
	pbRuntimes, err = p.getRuntimePbsLabel(pbRuntimes)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes %+v", err)
	}
	res := &pb.DescribeRuntimesResponse{
		RuntimeSet: pbRuntimes,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) ModifyRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) (*pb.ModifyRuntimeResponse, error) {
	// validate req
	err := validateModifyRuntimeRequest(req)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ModifyRuntime: %+v", err)
	}
	// check runtime env can be modified
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}
	if deleted {
		logger.Errorf("ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// update runtime env
	err = p.updateRuntime(req)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}

	// update runtime env label
	if req.Labels != nil {
		err := p.updateRuntimeLabels(runtimeId, req.Labels.GetValue())
		if err != nil {
			logger.Errorf("ModifyRuntime: %+v", err)
			return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
		}
	}

	// get response
	pbRuntime, err := p.getRuntimePbWithLabel(runtimeId)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}
	res := &pb.ModifyRuntimeResponse{
		Runtime: pbRuntime,
	}

	return res, nil
}

func (p *Server) DeleteRuntime(ctx context.Context, req *pb.DeleteRuntimeRequest) (*pb.DeleteRuntimeResponse, error) {
	// validate req
	err := validateDeleteRuntimeRequest(req)
	if err != nil {
		logger.Errorf("DeleteRuntimeCredential: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DeleteRuntimeCredential: %+v", err)
	}

	// check runtime env can be deleted
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: %+v", err)
	}
	if deleted {
		logger.Errorf("DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// deleted runtime env
	err = p.deleteRuntime(runtimeId)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}

	// get runtime env
	pbRuntime, err := p.getRuntimePbWithLabel(runtimeId)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}
	res := &pb.DeleteRuntimeResponse{
		Runtime: pbRuntime,
	}
	return res, nil
}

func (p *Server) CreateRuntimeCredential(context.Context, *pb.CreateRuntimeCredentialRequset) (*pb.CreateRuntimeCredentialResponse, error) {
	return nil, nil
}
func (p *Server) DescribeRuntimeCredentials(context.Context, *pb.DescribeRuntimeCredentialsRequset) (*pb.DescribeRuntimeCredentialsResponse, error) {
	return nil, nil
}
func (p *Server) ModifyRuntimeCredential(context.Context, *pb.ModifyRuntimeCredentialRequest) (*pb.ModifyRuntimeCredentialResponse, error) {
	return nil, nil
}
func (p *Server) DeleteRuntimeCredential(context.Context, *pb.DeleteRuntimeCredentialRequset) (*pb.DeleteRuntimeCredentialResponse, error) {
	return nil, nil
}

func (p *Server) getRuntimePbWithLabel(runtimeId string) (*pb.Runtime, error) {
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime := models.RuntimeToPb(runtime)
	runtimeLabels, err := p.getRuntimeLabelsById(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime.Labels = models.RuntimeLabelsToPbs(runtimeLabels)

	return pbRuntime, nil
}

func (p *Server) createRuntime(name, description, url, userId string) (runtimeId string, err error) {
	newRuntime := models.NewRuntime(name, description, url, userId)
	err = p.insertRuntime(*newRuntime)
	if err != nil {
		return "", err
	}
	return newRuntime.RuntimeId, err
}

func (p *Server) createRuntimeLabels(runtimeId, labelString string) error {
	labelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(runtimeId, labelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) getRuntimeCredentialPbById(runtimeCredentialId string) (*pb.RuntimeCredential, error) {
	runtimeCredential, err := p.getRuntimeCredential(runtimeCredentialId)
	if err != nil {
		return nil, err
	}
	pbRuntimeCredential := models.RuntimeCredentialToPb(runtimeCredential)
	return pbRuntimeCredential, nil
}

func (p *Server) getRuntimePbsLabel(pbRuntimes []*pb.Runtime) ([]*pb.Runtime, error) {
	var runtimeIds []string
	for _, pbRuntime := range pbRuntimes {
		runtimeIds = append(runtimeIds, pbRuntime.RuntimeId.GetValue())
	}
	runtimeLabels, err := p.getRuntimeLabelsById(runtimeIds...)
	if err != nil {
		return nil, err
	}
	for _, pbRuntime := range pbRuntimes {
		for _, runtimeLabel := range runtimeLabels {
			if pbRuntime.RuntimeId.GetValue() == runtimeLabel.RuntimeId {
				pbRuntime.Labels = append(pbRuntime.Labels, models.RuntimeLabelToPb(runtimeLabel))
			}
		}
	}
	return pbRuntimes, nil
}

func (p *Server) getRuntimeIdsBySelectorString(selectorString string) ([]string, error) {
	selectorMap, err := SelectorStringToMap(selectorString)
	if err != nil {
		return nil, err
	}
	runtimeIds, err := p.getRuntimeIdsBySelectorMap(selectorMap)
	if err != nil {
		return nil, err
	}
	return runtimeIds, nil
}

func (p *Server) getRuntimePbsWithoutLabelByReqAndId(
	req *pb.DescribeRuntimesRequest, runtimeIds []string) (
	runtimePbs []*pb.Runtime, count uint32, err error) {
	// build filter condition
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	filterCondition := manager.BuildFilterConditions(req, models.RuntimeTableName)
	if len(runtimeIds) > 0 {
		if filterCondition == nil {
			filterCondition = db.Eq(RuntimeIdColumn, runtimeIds)
		} else {
			filterCondition = db.And(filterCondition, db.Eq(RuntimeIdColumn, runtimeIds))
		}
	} else if len(runtimeIds) == 0 && req.Selector != nil {
		return nil, 0, nil
	}
	// get runtime
	runtimes, count, err := p.getRuntimesByFilterCondition(filterCondition, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	runtimePbs = models.RuntimeToPbs(runtimes)
	return runtimePbs, count, nil
}

func (p *Server) getRuntimeCredentialPbsByReq(req *pb.DescribeRuntimeCredentialsRequset) (
	runtimeCredentialPbs []*pb.RuntimeCredential, count uint32, err error) {
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	filterCondition := manager.BuildFilterConditions(req, models.RuntimeCredentialTableName)
	runtimeCredentials, count, err := p.getRuntimeCredentialsByFilterCondition(filterCondition, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	runtimeCredentialPbs = models.RuntimeCredentialToPbs(runtimeCredentials)
	return runtimeCredentialPbs, count, nil
}

func (p *Server) updateRuntime(req *pb.ModifyRuntimeRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := p.updateRuntimeByMap(req.RuntimeId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeCredential(req *pb.ModifyRuntimeCredentialRequest) error {
	attributes := manager.BuildUpdateAttributes(
		req, NameColumn, DescriptionColumn)
	attributes[RuntimeCredentialContentColumn] = models.RuntimeCredentialContentMapToString(req.Content)
	err := p.updateRuntimeCredentialByMap(req.RuntimeCredentialId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeLabels(runtimeId string, labelString string) error {
	newLabelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	oldRuntimeLabels, err := p.getRuntimeLabelsById(runtimeId)
	if err != nil {
		return err
	}
	oldLabelMap := LabelStructToMap(oldRuntimeLabels)
	additionLabelMap, deletionLabelMap := LabelMapDiff(oldLabelMap, newLabelMap)
	err = p.deleteRuntimeLabels(runtimeId, deletionLabelMap)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(runtimeId, additionLabelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) createRuntimeCredential(name, description, userId string, content map[string]string) (
	runtimeCredentialId string, err error) {

	newRunTimeCredential := models.NewRuntimeCredential(name, description, userId, content)
	err = p.insertRuntimeCredential(*newRunTimeCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeCredential.RuntimeCredentialId, nil
}
