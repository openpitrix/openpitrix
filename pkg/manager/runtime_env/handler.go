// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"context"

	"github.com/gocraft/dbr"
	"github.com/golang/protobuf/ptypes/wrappers"
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
	err := validateCreateRuntimeEnvRequest(req)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntimeEnv: %+v", err)
	}

	// create runtime env
	runtimeEnvId, err := p.createRuntimeEnv(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetRuntimeEnvUrl().GetValue(),
		s.UserId)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	// create labels
	err = p.createRuntimeEnvLabels(runtimeEnvId, req.Labels.GetValue())
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	// get response
	pbRuntimeEnv, err := p.getRuntimeEnvPbWithLabel(runtimeEnvId)
	if err != nil {
		logger.Errorf("CreateRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}
	res := &pb.CreateRuntimeEnvResponse{
		RuntimeEnv: pbRuntimeEnv,
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
		runtimeEnvIds, err = p.getRuntimeEnvIdsBySelectorString(req.Selector.GetValue())
		if err != nil {
			logger.Errorf("DescribeRuntimeEnvs: %+v", err)
			return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs: %+v", err)
		}
	}
	// get runtime envs
	pbRuntimeEnvs, count, err := p.getRuntimeEnvPbsWithoutLabelByReqAndId(req, runtimeEnvIds)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvs: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs: %+v", err)
	}
	// get runtime envs label
	pbRuntimeEnvs, err = p.getRuntimeEnvPbsLabel(pbRuntimeEnvs)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvs: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs %+v", err)
	}
	// get runtime env credential id
	pbRuntimeEnvs, err = p.getRuntimeEnvPbsRuntimeCredentialId(pbRuntimeEnvs)
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
	deleted, err := p.checkRuntimeEnvDeleted(runtimeEnvId)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("ModifyRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"ModifyRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
	}
	// update runtime env
	err = p.updateRuntimeEnv(req)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
	}

	// update runtime env label
	if req.Labels != nil {
		err := p.updateRuntimeEnvLabels(runtimeEnvId, req.Labels.GetValue())
		if err != nil {
			logger.Errorf("ModifyRuntimeEnv: %+v", err)
			return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", err)
		}
	}

	// get response
	pbRuntimeEnv, err := p.getRuntimeEnvPbWithLabel(runtimeEnvId)
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
	attached, err := p.checkRuntimeEnvAttached(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", err)
	}
	if attached {
		logger.Errorf("DeleteRuntimeEnv: runtime_env_credential has been attached [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnv: runtime_env_credential has been attached [%+v]", runtimeEnvId)
	}
	deleted, err := p.checkRuntimeEnvDeleted(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("DeleteRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
	}
	// deleted runtime env
	err = p.deleteRuntimeEnv(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", err)
	}

	// get runtime env
	pbRuntimeEnv, err := p.getRuntimeEnvPbWithLabel(runtimeEnvId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", err)
	}
	res := &pb.DeleteRuntimeEnvResponse{
		RuntimeEnv: pbRuntimeEnv,
	}
	return res, nil
}

func (p *Server) CreateRuntimeEnvCredential(ctx context.Context, req *pb.CreateRuntimeEnvCredentialRequset) (*pb.CreateRuntimeEnvCredentialResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeEnvCredential(req)
	if err != nil {
		logger.Errorf("CreateRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntimeEnvCredential: %+v", err)
	}
	// create runtime env credential
	runtimeEnvCredentialId, err := p.createRuntimeEnvCredential(
		req.Name.GetValue(),
		req.Description.GetValue(),
		s.UserId,
		req.Content,
	)
	if err != nil {
		logger.Errorf("CreateRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnvCredential: %+v", err)
	}
	// get runtime env credential
	pbRunTimeEnvCredential, err := p.getRuntimeEnvCredentialPbById(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("CreateRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnvCredential: %+v", err)
	}
	res := &pb.CreateRuntimeEnvCredentialResponse{
		RuntimeEnvCredential: pbRunTimeEnvCredential,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeEnvCredentials(ctx context.Context, req *pb.DescribeRuntimeEnvCredentialsRequset) (*pb.DescribeRuntimeEnvCredentialsResponse, error) {
	err := validateDescribeRuntimeEnvCredential(req)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvCredentials: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimeEnvCredentials: %+v", err)
	}
	pbRuntimeEnvCredential, count, err := p.getRuntimeEnvCredentialPbsByReq(req)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvCredentials: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvCredentials: %+v", err)
	}

	pbRuntimeEnvCredential, err = p.getRuntimeEnvCredentialPbsRuntimeEnvIds(pbRuntimeEnvCredential)
	if err != nil {
		logger.Errorf("DescribeRuntimeEnvCredentials: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvCredentials: %+v", err)
	}

	res := &pb.DescribeRuntimeEnvCredentialsResponse{
		RuntimeEnvCredentialSet: pbRuntimeEnvCredential,
		TotalCount:              count,
	}
	return res, nil
}

func (p *Server) ModifyRuntimeEnvCredential(ctx context.Context, req *pb.ModifyRuntimeEnvCredentialRequest) (*pb.ModifyRuntimeEnvCredentialResponse, error) {
	// validate req
	err := validateModifyRuntimeEnvCredential(req)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ModifyRuntimeEnvCredential: %+v", err)
	}

	// check runtime env credential status
	runtimeEnvCredentialId := req.GetRuntimeEnvCredentialId().GetValue()
	deleted, err := p.checkRuntimeEnvCredentialDeleted(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnvCredential: %+v", err)
	}
	if deleted {
		logger.Errorf("ModifyRuntimeEnvCredential: runtime_env_credential has been deleted [%+v]", runtimeEnvCredentialId)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnvCredential: "+
			"runtime_env_credential has been deleted [%+v]", runtimeEnvCredentialId)
	}
	// update runtime env credential
	err = p.updateRuntimeEnvCredential(req)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnvCredential: %+v", err)
	}

	// get response
	pbRunTimeEnvCredential, err := p.getRuntimeEnvCredentialPbById(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("ModifyRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnvCredential: %+v", err)
	}
	res := &pb.ModifyRuntimeEnvCredentialResponse{
		RuntimeEnvCredential: pbRunTimeEnvCredential,
	}
	return res, nil
}

func (p *Server) DeleteRuntimeEnvCredential(ctx context.Context, req *pb.DeleteRuntimeEnvCredentialRequset) (*pb.DeleteRuntimeEnvCredentialResponse, error) {
	// validate req
	err := validateDeleteRuntimeEnvCredential(req)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DeleteRuntimeEnvCredential: %+v", err)
	}

	// check runtime env credential status
	runtimeEnvCredentialId := req.GetRuntimeEnvCredentialId().GetValue()
	attached, err := p.checkRuntimeEnvCredentialAttached(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnvCredential: %+v", err)
	}
	if attached {
		logger.Errorf("DeleteRuntimeEnvCredential: runtime_env_credential has been attached [%+v]", runtimeEnvCredentialId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnvCredential: runtime_env_credential has been attached [%+v]", runtimeEnvCredentialId)
	}
	deleted, err := p.checkRuntimeEnvCredentialDeleted(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnvCredential: %+v", err)
	}
	if deleted {
		logger.Errorf("DeleteRuntimeEnvCredential: runtime_env_credential has been deleted [%+v]", runtimeEnvCredentialId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntimeEnvCredential: runtime_env_credential has been deleted [%+v]", runtimeEnvCredentialId)
	}

	// update runtime env credential status
	err = p.deleteRuntimeEnvCredential(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: [%+v]", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnvCredential: [%+v]", err)
	}

	// get response
	pbRunTimeEnvCredential, err := p.getRuntimeEnvCredentialPbById(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("DeleteRuntimeEnvCredential: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnvCredential: %+v", err)
	}
	res := &pb.DeleteRuntimeEnvCredentialResponse{
		RuntimeEnvCredential: pbRunTimeEnvCredential,
	}
	return res, nil
}

func (p *Server) AttachCredentialToRuntimeEnv(ctx context.Context, req *pb.AttachCredentialToRuntimeEnvRequset) (*pb.AttachCredentialToRuntimeEnvResponse, error) {
	// validate req
	err := validateAttachCredentialToRuntimeEnv(req)
	if err != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "AttachCredentialToRuntimeEnv	: %+v", err)
	}

	// check runtime env status
	runtimeEnvId := req.RuntimeEnvId.GetValue()
	deleted, err := p.checkRuntimeEnvDeleted(runtimeEnvId)
	if err != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "AttachCredentialToRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("AttachCredentialToRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
		return nil, status.Errorf(codes.Internal,
			"AttachCredentialToRuntimeEnv: runtime_env has been deleted [%+v]", runtimeEnvId)
	}

	// check runtime env credential status
	runtimeEnvCredentialId := req.RuntimeEnvCredentialId.GetValue()
	deleted, err = p.checkRuntimeEnvCredentialDeleted(runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"AttachCredentialToRuntimeEnv: %+v", err)
	}
	if deleted {
		logger.Errorf("AttachCredentialToRuntimeEnv: runtime_env_credential has been deleted [%+v]", runtimeEnvCredentialId)
		return nil, status.Errorf(codes.Internal,
			"AttachCredentialToRuntimeEnv: runtime_env_credential has been deleted [%+v]",
			runtimeEnvCredentialId)
	}

	// check runtime env attached
	runtimeEnvAttachedCredential, err := p.getAttachedCredentialByEnvId(runtimeEnvId)
	if err != nil && err != dbr.ErrNotFound {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "AttachCredentialToRuntimeEnv: %+v", err)
	}
	if runtimeEnvAttachedCredential != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: runtime_env has been attached to %s", runtimeEnvAttachedCredential.RuntimeEnvCredentialId)
		return nil, status.Errorf(codes.Internal,
			"AttachCredentialToRuntimeEnv: runtime_env has been attached to %s",
			runtimeEnvAttachedCredential.RuntimeEnvCredentialId)
	}

	err = p.attachCredentialToRuntimeEnv(runtimeEnvId, runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "AttachCredentialToRuntimeEnv: %+v", err)
	}

	envAttachedCredential, err := p.getAttachedCredentialByEnvId(runtimeEnvId)
	if err != nil {
		logger.Errorf("AttachCredentialToRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal, "AttachCredentialToRuntimeEnv: failed to get runtimeEnvAttachedCredential [%+v]", err)
	}
	res := &pb.AttachCredentialToRuntimeEnvResponse{
		RuntimeEnvId:           &wrappers.StringValue{Value: envAttachedCredential.RuntimeEnvId},
		RuntimeEnvCredentialId: &wrappers.StringValue{Value: envAttachedCredential.RuntimeEnvCredentialId},
	}

	return res, nil
}

func (p *Server) DetachCredentialFromRuntimeEnv(ctx context.Context, req *pb.DetachCredentialFromRuntimeEnvRequset) (*pb.DetachCredentialFromRuntimeEnvResponse, error) {
	// validate req
	err := validateDetachCredentialFromRuntimeEnv(req)
	if err != nil {
		logger.Errorf("DetachCredentialFromRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DetachCredentialFromRuntimeEnv	: %+v", err)
	}

	// get runtimeEnvAttachedCredential to validate
	runtimeEnvId := req.RuntimeEnvId.GetValue()
	runtimeEnvAttachedCredential, err := p.getAttachedCredentialByEnvId(runtimeEnvId)
	if err != nil {
		logger.Errorf("DetachCredentialFromRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DetachCredentialFromRuntimeEnv: [%+v]", err)
	}
	runtimeEnvCredentialId := req.RuntimeEnvCredentialId.GetValue()
	if runtimeEnvAttachedCredential.RuntimeEnvCredentialId != runtimeEnvCredentialId {
		logger.Errorf("DetachCredentialFromRuntimeEnv: runtime_env_credential value not match")
		return nil, status.Errorf(codes.Internal,
			"DetachCredentialFromRuntimeEnv: runtime_env_credential value not match")
	}
	// detach credential from runtime_env
	err = p.detachCredentialFromRuntimeEnv(runtimeEnvId, runtimeEnvCredentialId)
	if err != nil {
		logger.Errorf("DetachCredentialFromRuntimeEnv: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DetachCredentialFromRuntimeEnv: %+v", err)
	}

	res := &pb.DetachCredentialFromRuntimeEnvResponse{
		RuntimeEnvId:           &wrappers.StringValue{Value: runtimeEnvId},
		RuntimeEnvCredentialId: &wrappers.StringValue{Value: runtimeEnvCredentialId},
	}

	return res, nil
}

func (p *Server) getRuntimeEnvPbWithLabel(runtimeEnvId string) (*pb.RuntimeEnv, error) {
	runtimeEnv, err := p.getRuntimeEnv(runtimeEnvId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnv := models.RuntimeEnvToPb(runtimeEnv)
	runtimeEnvLabels, err := p.getRuntimeEnvLabelsByEnvId(runtimeEnvId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnv.Labels = models.RuntimeEnvLabelsToPbs(runtimeEnvLabels)

	return pbRuntimeEnv, nil
}

func (p *Server) createRuntimeEnv(name, description, url, userId string) (runtimeEnvId string, err error) {
	newRuntimeEnv := models.NewRuntimeEnv(name, description, url, userId)
	err = p.insertRuntimeEnv(*newRuntimeEnv)
	if err != nil {
		return "", err
	}
	return newRuntimeEnv.RuntimeEnvId, err
}

func (p *Server) createRuntimeEnvLabels(runtimeEnvId, labelString string) error {
	labelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	err = p.insertRuntimeEnvLabels(runtimeEnvId, labelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) getRuntimeEnvCredentialPbById(runtimeEnvCredentialId string) (*pb.RuntimeEnvCredential, error) {
	runtimeEnvCredential, err := p.getRuntimeEnvCredential(runtimeEnvCredentialId)
	if err != nil {
		return nil, err
	}
	pbRuntimeEnvCredential := models.RuntimeEnvCredentialToPb(runtimeEnvCredential)
	return pbRuntimeEnvCredential, nil
}

func (p *Server) getRuntimeEnvPbsLabel(pbRuntimeEnvs []*pb.RuntimeEnv) ([]*pb.RuntimeEnv, error) {
	var runtimeEnvIds []string
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		runtimeEnvIds = append(runtimeEnvIds, pbRuntimeEnv.RuntimeEnvId.GetValue())
	}
	runtimeEnvLabels, err := p.getRuntimeEnvLabelsByEnvId(runtimeEnvIds...)
	if err != nil {
		return nil, err
	}
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		for _, runtimeEnvLabel := range runtimeEnvLabels {
			if pbRuntimeEnv.RuntimeEnvId.GetValue() == runtimeEnvLabel.RuntimeEnvId {
				pbRuntimeEnv.Labels = append(pbRuntimeEnv.Labels, models.RuntimeEnvLabelToPb(runtimeEnvLabel))
			}
		}
	}
	return pbRuntimeEnvs, nil
}

func (p *Server) getRuntimeEnvPbsRuntimeCredentialId(pbRuntimeEnvs []*pb.RuntimeEnv) ([]*pb.RuntimeEnv, error) {
	var runtimeEnvIds []string
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		runtimeEnvIds = append(runtimeEnvIds, pbRuntimeEnv.RuntimeEnvId.GetValue())
	}
	runtimeEnvAttachedCredentials, err := p.getAttachedCredentialsByEnvIds(runtimeEnvIds)
	if err != nil {
		return nil, err
	}
	for _, pbRuntimeEnv := range pbRuntimeEnvs {
		for _, runtimeEnvAttachedCredential := range runtimeEnvAttachedCredentials {
			if pbRuntimeEnv.RuntimeEnvId.GetValue() == runtimeEnvAttachedCredential.RuntimeEnvId {
				pbRuntimeEnv.RuntimeEnvCredentialId = utils.ToProtoString(runtimeEnvAttachedCredential.RuntimeEnvCredentialId)
				break
			}
		}
	}
	return pbRuntimeEnvs, nil
}

func (p *Server) getRuntimeEnvCredentialPbsRuntimeEnvIds(pbrRuntimeEnvCredentials []*pb.RuntimeEnvCredential) ([]*pb.RuntimeEnvCredential, error) {
	var runtimeEnvCredentialIds []string
	for _, pbrRuntimeEnvCredential := range pbrRuntimeEnvCredentials {
		runtimeEnvCredentialIds = append(runtimeEnvCredentialIds, pbrRuntimeEnvCredential.RuntimeEnvCredentialId.GetValue())
	}
	runtimeEnvAttachedCredentials, err := p.getAttachedCredentialsByCredentialIds(runtimeEnvCredentialIds)
	if err != nil {
		return nil, err
	}
	for _, pbrRuntimeEnvCredential := range pbrRuntimeEnvCredentials {
		for _, runtimeEnvAttachedCredential := range runtimeEnvAttachedCredentials {
			if pbrRuntimeEnvCredential.RuntimeEnvCredentialId.GetValue() == runtimeEnvAttachedCredential.RuntimeEnvCredentialId {
				pbrRuntimeEnvCredential.RuntimeEnvId = append(pbrRuntimeEnvCredential.RuntimeEnvId, runtimeEnvAttachedCredential.RuntimeEnvId)
			}
		}
	}
	return pbrRuntimeEnvCredentials, nil
}

func (p *Server) getRuntimeEnvIdsBySelectorString(selectorString string) ([]string, error) {
	selectorMap, err := SelectorStringToMap(selectorString)
	if err != nil {
		return nil, err
	}
	runtimeEnvIds, err := p.getRuntimeEnvIdsBySelectorMap(selectorMap)
	if err != nil {
		return nil, err
	}
	return runtimeEnvIds, nil
}

func (p *Server) getRuntimeEnvPbsWithoutLabelByReqAndId(
	req *pb.DescribeRuntimeEnvsRequest, runtimeEnvIds []string) (
	runtimeEnvPbs []*pb.RuntimeEnv, count uint32, err error) {
	// build filter condition
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	filterCondition := manager.BuildFilterConditions(req, models.RuntimeEnvTableName)
	if len(runtimeEnvIds) > 0 {
		if filterCondition == nil {
			filterCondition = db.Eq(RuntimeEnvIdColumn, runtimeEnvIds)
		} else {
			filterCondition = db.And(filterCondition, db.Eq(RuntimeEnvIdColumn, runtimeEnvIds))
		}
	} else if len(runtimeEnvIds) == 0 && req.Selector != nil {
		return nil, 0, nil
	}
	//get runtime envs
	runtimeEnvs, count, err := p.getRuntimeEnvsByFilterCondition(filterCondition, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	runtimeEnvPbs = models.RuntimeEnvToPbs(runtimeEnvs)
	return runtimeEnvPbs, count, nil
}

func (p *Server) getRuntimeEnvCredentialPbsByReq(req *pb.DescribeRuntimeEnvCredentialsRequset) (
	runtimeEnvCredentialPbs []*pb.RuntimeEnvCredential, count uint32, err error) {
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	filterCondition := manager.BuildFilterConditions(req, models.RuntimeEnvCredentialTableName)
	runtimeEnvCredentials, count, err := p.getRuntimeEnvCredentialsByFilterCondition(filterCondition, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	runtimeEnvCredentialPbs = models.RuntimeEnvCredentialToPbs(runtimeEnvCredentials)
	return runtimeEnvCredentialPbs, count, nil
}

func (p *Server) updateRuntimeEnv(req *pb.ModifyRuntimeEnvRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := p.updateRuntimeEnvByMap(req.RuntimeEnvId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeEnvCredential(req *pb.ModifyRuntimeEnvCredentialRequest) error {
	attributes := manager.BuildUpdateAttributes(
		req, NameColumn, DescriptionColumn)
	attributes[RuntimeEnvCredentialContentColumn] = models.RuntimeEnvCredentialContentMapToString(req.Content)
	err := p.updateRuntimeEnvCredentialByMap(req.RuntimeEnvCredentialId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeEnvLabels(runtimeEnvId string, labelString string) error {
	newLabelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	oldRuntimeEnvLabels, err := p.getRuntimeEnvLabelsByEnvId(runtimeEnvId)
	if err != nil {
		return err
	}
	oldLabelMap := LabelStructToMap(oldRuntimeEnvLabels)
	additionLabelMap, deletionLabelMap := LabelMapDiff(oldLabelMap, newLabelMap)
	err = p.deleteRuntimeEnvLabels(runtimeEnvId, deletionLabelMap)
	if err != nil {
		return err
	}
	err = p.insertRuntimeEnvLabels(runtimeEnvId, additionLabelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) createRuntimeEnvCredential(name, description, userId string, content map[string]string) (
	runtimeEnvCredentialId string, err error) {

	newRunTimeEnvCredential := models.NewRuntimeEnvCredential(name, description, userId, content)
	err = p.insertRuntimeEnvCredential(*newRunTimeEnvCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeEnvCredential.RuntimeEnvCredentialId, nil
}

func (p *Server) attachCredentialToRuntimeEnv(runtimeEnvId, runtimeEnvCredentialId string) error {
	newRuntimeEnvAttachedCredential := models.NewRuntimeEnvAttachedCredential(runtimeEnvId, runtimeEnvCredentialId)
	err := p.insertRuntimeEnvAttachedCredential(*newRuntimeEnvAttachedCredential)
	if err != nil {
		return err
	}
	return nil
}
