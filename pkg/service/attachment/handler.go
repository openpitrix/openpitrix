// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachment

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

func (p *Server) CreateAttachment(ctx context.Context, req *pb.CreateAttachmentRequest) (*pb.CreateAttachmentResponse, error) {
	if len(req.AttachmentContent) == 0 {
		return nil, fmt.Errorf("content is empty")
	}
	var attachment = models.NewAttachment(req.AttachmentType.String())
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttachment).
		Record(attachment).
		Exec()
	if err != nil {
		return nil, err
	}

	err = putAttachmentFiles(ctx, attachment, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAttachmentResponse{
		AttachmentId: attachment.AttachmentId,
	}, nil
}

func (p *Server) AppendAttachment(ctx context.Context, req *pb.AppendAttachmentRequest) (*pb.AppendAttachmentResponse, error) {
	attachment, err := getAttachment(ctx, req.AttachmentId)
	if err != nil {
		return nil, err
	}

	err = putAttachmentFiles(ctx, attachment, req)
	if err != nil {
		return nil, err
	}
	return &pb.AppendAttachmentResponse{
		AttachmentId:   attachment.AttachmentId,
		AttachmentType: attachment.GetAttachmentType(),
	}, nil
}

func (p *Server) ReplaceAttachment(ctx context.Context, req *pb.ReplaceAttachmentRequest) (*pb.ReplaceAttachmentResponse, error) {
	attachment, err := getAttachment(ctx, req.AttachmentId)
	if err != nil {
		return nil, err
	}
	err = deleteAttachmentFiles(ctx, attachment)
	if err != nil {
		return nil, err
	}
	err = putAttachmentFiles(ctx, attachment, req)
	if err != nil {
		return nil, err
	}
	return &pb.ReplaceAttachmentResponse{
		AttachmentId:   attachment.AttachmentId,
		AttachmentType: attachment.GetAttachmentType(),
	}, nil
}

func (p *Server) GetAttachments(ctx context.Context, req *pb.GetAttachmentsRequest) (*pb.GetAttachmentsResponse, error) {
	attachments, err := getAttachments(ctx, req.AttachmentId)
	if err != nil {
		return nil, err
	}

	var res = &pb.GetAttachmentsResponse{
		Attachments: make(map[string]*pb.Attachment),
	}

	pbAtts, err := getAttachmentFiles(ctx, attachments, req)
	if err != nil {
		return nil, err
	}
	for _, pbAtt := range pbAtts {
		res.Attachments[pbAtt.AttachmentId] = pbAtt
	}

	return res, nil
}

func (p *Server) DeleteAttachments(ctx context.Context, req *pb.DeleteAttachmentsRequest) (*pb.DeleteAttachmentsResponse, error) {
	as, err := getAttachments(ctx, req.AttachmentId)
	if err != nil {
		return nil, err
	}
	if len(req.Filename) == 0 {
		err = removeAttachments(ctx, req.AttachmentId)
		if err != nil {
			return nil, err
		}
	}
	for _, a := range as {
		err = deleteAttachmentFiles(ctx, a, req.Filename...)
		if err != nil {
			return nil, err
		}
	}
	return &pb.DeleteAttachmentsResponse{
		AttachmentId: req.AttachmentId,
		Filename:     req.Filename,
	}, err
}
