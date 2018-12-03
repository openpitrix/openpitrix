// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachment

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

var (
	_ pb.AttachmentManagerServer = &Server{}
	_ pb.AttachmentServiceServer = &Server{}
)

func (p *Server) CreateAttachment(ctx context.Context, req *pb.CreateAttachmentRequest) (*pb.CreateAttachmentResponse, error) {
	if len(req.AttachmentContent) == 0 {
		return nil, fmt.Errorf("content is empty")
	}
	var attachment = models.NewAttachment()
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
		AttachmentId: attachment.AttachmentId,
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
		AttachmentId: attachment.AttachmentId,
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

func (p *Server) GetAttachment(ctx context.Context, req *pb.GetAttachmentRequest) (*pb.GetAttachmentResponse, error) {
	var content = &pb.GetAttachmentResponse{}
	if len(req.AttachmentId) == 0 || len(req.Filename) == 0 {
		return content, nil
	}
	attachment, err := getAttachment(ctx, req.AttachmentId)
	if err != nil {
		return content, err
	}
	output, err := getFile(ctx, attachment, req.Filename)
	if err != nil {
		if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
			return content, nil
		}
		return content, err
	}
	file, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	content.Content = file
	if output.ETag != nil {
		content.Etag = *output.ETag
	}
	return content, nil
}
