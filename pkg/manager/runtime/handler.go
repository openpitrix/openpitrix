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

func (p *Server) CreateRuntimeEnv(ctx context.Context, req *pb.CreateRuntimeEnvRequest) (*pb.CreateRuntimeEnvResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(req)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntimeEnv: %+v", err)
	}

	// create runtime env
	runtimeId, err := p.createRuntime(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetRuntimeEnvUrl().GetValue(),
		s.UserId)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	// create labels
	err = p.createRuntimeLabels(runtimeId, req.Labels.GetValue())
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	// get response
	pbRuntime, err := p.getRuntimePbWithLabel(runtimeId)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}
	res := &pb.CreateRuntimeEnvResponse{
		RuntimeEnv: pbRuntime,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeEnvs(ctx context.Context, req *pb.DescribeRuntimeEnvsRequest) (*pb.DescribeRuntimeEnvsResponse, error) {
	// validate req
	err := validateDescribeRuntimeEnvRequest(req)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvs: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimeEnvs: %+v", err)
	}
	var runtimeEnvIds []string
	// get runtime env ids by selector
	if req.Selector != nil {
		runtimeEnvIds, err = p.getRuntimeIdsBySelectorString(req.Selector.GetValue())
		if err != nil {
			logger.Errorf("DescribeRuntimeEnvs: %+v", err)
			return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs: %+v", err)
		}
	}
	// get runtime envs
	pbRuntimeEnvs, count, err := p.getRuntimePbsWithoutLabelByReqAndId(req, runtimeEnvIds)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvs: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs: %+v", err)
	}
	// get runtime envs label
	pbRuntimeEnvs, err = p.getRuntimePbsLabel(pbRuntimeEnvs)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvs: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs %+v", err)
	}
	res := &pb.DescribeRuntimeEnvsResponse{
		RuntimeEnvSet: pbRuntimeEnvs,
		TotalCount:    count,
	}
	return res, nil
}

func (p *Server) ModifyRuntimeEnv(ctx context.Context, req *pb.ModifyRuntimeEnvRequest) (*pb.ModifyRuntimeEnvResponse, error) {
	// validate req
	err := validateModifyRuntimeEnvRequest(req)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ModifyRuntimeEnv: %+v", err)
	}
	// check runtime env can be modified
	runtimeEnvId := req.GetRuntimeEnvId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeEnvId)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("ModifyRuntimeEnv: runtime has been deleted [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"ModifyRuntimeEnv: runtime has been deleted [%+v]", runtimeEnvId)
	}
	// update runtime env
	err = p.updateRuntimeEnv(req)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
	}

	// update runtime env label
	if req.Labels != nil {
		err := p.updateRuntimeLabels(runtimeEnvId, req.Labels.GetValue())
		if err != nil {
			logger.Errorf("ModifyRuntimeEnv: %+v", err)
			return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
		}
	}

	// get response
	pbRuntimeEnv, err := p.getRuntimePbWithLabel(runtimeEnvId)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
	}
	res := &pb.ModifyRuntimeEnvResponse{
		RuntimeEnv: pbRuntimeEnv,
	}

	return res, nil
}

func (p *Server) DeleteRuntimeEnv(ctx context.Context, req *pb.DeleteRuntimeEnvRequest) (*pb.DeleteRuntimeEnvResponse, error) {
	// validate req
	err := validateDeleteRuntimeEnvRequest(req)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DeleteRuntimeEnvCredential: %+v", err)
	}

	// check runtime env can be deleted
	runtimeEnvId := req.GetRuntimeEnvId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("DeleteRuntimeEnv: runtime has been deleted [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnv: runtime has been deleted [%+v]", runtimeEnvId)
	}
	// deleted runtime env
	err = p.deleteRuntime(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", err)
	}

	// get runtime env
	pbRuntimeEnv, err := p.getRuntimePbWithLabel(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", err)
	}
	res := &pb.DeleteRuntimeEnvResponse{
		RuntimeEnv: pbRuntimeEnv,
	}
	return res, nil
}

func (p *Server) CreateRuntimeEnvCredential(context.Context, *pb.CreateRuntimeEnvCredentialRequset) (*pb.CreateRuntimeEnvCredentialResponse, error) {
	return nil, nil
}
func (p *Server) DescribeRuntimeEnvCredentials(context.Context, *pb.DescribeRuntimeEnvCredentialsRequset) (*pb.DescribeRuntimeEnvCredentialsResponse, error) {
	return nil, nil
}
func (p *Server) ModifyRuntimeEnvCredential(context.Context, *pb.ModifyRuntimeEnvCredentialRequest) (*pb.ModifyRuntimeEnvCredentialResponse, error) {
	return nil, nil
}
func (p *Server) DeleteRuntimeEnvCredential(context.Context, *pb.DeleteRuntimeEnvCredentialRequset) (*pb.DeleteRuntimeEnvCredentialResponse, error) {
	return nil, nil
}
func (p *Server) AttachCredentialToRuntimeEnv(context.Context, *pb.AttachCredentialToRuntimeEnvRequset) (*pb.AttachCredentialToRuntimeEnvResponse, error) {
	return nil, nil
}
func (p *Server) DetachCredentialFromRuntimeEnv(context.Context, *pb.DetachCredentialFromRuntimeEnvRequset) (*pb.DetachCredentialFromRuntimeEnvResponse, error) {
	return nil, nil
}

func (p *Server) getRuntimePbWithLabel(runtimeId string) (*pb.RuntimeEnv, error) {
	runtimeEnv, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnv := models.RuntimeToPb(runtimeEnv)
	runtimeEnvLabels, err := p.getRuntimeLabelsById(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnv.Labels = models.RuntimeLabelsToPbs(runtimeEnvLabels)

	return pbRuntimeEnv, nil
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

func (p *Server) getRuntimeCredentialPbById(runtimeCredentialId string) (*pb.RuntimeEnvCredential, error) {
	runtimeEnvCredential, err := p.getRuntimeCredential(runtimeCredentialId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnvCredential := models.RuntimeCredentialToPb(runtimeEnvCredential)
	return pbRuntimeEnvCredential, nil
}

func (p *Server) getRuntimePbsLabel(pbRuntimeEnvs []*pb.RuntimeEnv) ([]*pb.RuntimeEnv, error) {
	var runtimeEnvIds []string
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		runtimeEnvIds = append(runtimeEnvIds, pbRuntimeEnv.RuntimeEnvId.GetValue())
	}
	runtimeEnvLabels, err := p.getRuntimeLabelsById(runtimeEnvIds...)
	if err != nil {
		return nil, err
	}
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		for _, runtimeEnvLabel := range runtimeEnvLabels {
			if pbRuntimeEnv.RuntimeEnvId.GetValue() == runtimeEnvLabel.RuntimeId {
				pbRuntimeEnv.Labels = append(pbRuntimeEnv.Labels, models.RuntimeLabelToPb(runtimeEnvLabel))
			}
		}
	}
	return pbRuntimeEnvs, nil
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
	req *pb.DescribeRuntimeEnvsRequest, runtimeIds []string) (
	runtimeEnvPbs []*pb.RuntimeEnv, count uint32, err error) {
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
	runtimeEnvPbs = models.RuntimeEnvToPbs(runtimes)
	return runtimeEnvPbs, count, nil
}

func (p *Server) getRuntimeCredentialPbsByReq(req *pb.DescribeRuntimeEnvCredentialsRequset) (
	runtimeEnvCredentialPbs []*pb.RuntimeEnvCredential, count uint32, err error) {
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	filterCondition := manager.BuildFilterConditions(req, models.RuntimeCredentialTableName)
	runtimeEnvCredentials, count, err := p.getRuntimeCredentialsByFilterCondition(filterCondition, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	runtimeEnvCredentialPbs = models.RuntimeEnvCredentialToPbs(runtimeEnvCredentials)
	return runtimeEnvCredentialPbs, count, nil
}

func (p *Server) updateRuntimeEnv(req *pb.ModifyRuntimeEnvRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := p.updateRuntimeByMap(req.RuntimeEnvId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeCredential(req *pb.ModifyRuntimeEnvCredentialRequest) error {
	attributes := manager.BuildUpdateAttributes(
		req, NameColumn, DescriptionColumn)
	attributes[RuntimeCredentialContentColumn] = models.RuntimeCredentialContentMapToString(req.Content)
	err := p.updateRuntimeCredentialByMap(req.RuntimeEnvCredentialId.GetValue(), attributes)
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
	runtimeEnvCredentialId string, err error) {

	newRunTimeEnvCredential := models.NewRuntimeCredential(name, description, userId, content)
	err = p.insertRuntimeCredential(*newRunTimeEnvCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeEnvCredential.RuntimeCredentialId, nil
}
