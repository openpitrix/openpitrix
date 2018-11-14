// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachment

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"openpitrix.io/openpitrix/pkg/client/internals3"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

func getAttachments(ctx context.Context, attachmentIds []string) ([]*models.Attachment, error) {
	var as []*models.Attachment
	_, err := pi.Global().DB(ctx).
		Select(models.AttachmentColumns...).
		From(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentIds)).
		Load(&as)
	return as, err
}

func removeAttachments(ctx context.Context, attachmentIds []string) error {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentIds)).
		Exec()
	return err
}

func (p *Server) UploadAttachment(ctx context.Context, req *pb.UploadAttachmentRequest) (*pb.UploadAttachmentResponse, error) {
	var attachment *models.Attachment
	if len(req.AttachmentId) == 0 {
		attachment = models.NewAttachment()
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableAttachment).
			Record(attachment).
			Exec()
		if err != nil {
			return nil, err
		}
	} else {
		var a models.Attachment
		err := pi.Global().DB(ctx).
			Select(models.AttachmentColumns...).
			From(constants.TableAttachment).
			Where(db.Eq(constants.ColumnAttachmentId, req.AttachmentId)).
			LoadOne(&a)
		if err != nil {
			return nil, err
		}
		attachment = &a
	}
	for filename, content := range req.AttachmentContent {
		_, err := internals3.S3.PutObject(&s3.PutObjectInput{
			Bucket: internals3.Bucket,
			Key:    aws.String(attachment.GetObjectName(filename)),
			Body:   bytes.NewReader(content),
		})
		if err != nil {
			return nil, err
		}
	}
	return &pb.UploadAttachmentResponse{
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

	for _, a := range attachments {
		var attachmentContent = make(map[string][]byte)
		// filenames with prefix
		var filenames []string

		// prepare object keys with prefix
		if len(req.Filename) == 0 {
			output, err := internals3.S3.ListObjects(&s3.ListObjectsInput{
				Bucket: internals3.Bucket,
				Prefix: aws.String(a.GetObjectPrefix()),
			})
			if err != nil {
				return nil, err
			}
			for _, o := range output.Contents {
				if o.Key != nil {
					filenames = append(filenames, *o.Key)
				}
			}
		} else {
			for _, filename := range req.Filename {
				filenames = append(filenames, a.GetObjectName(filename))
			}
		}

		// get object content
		for _, filename := range filenames {
			output, err := internals3.S3.GetObject(&s3.GetObjectInput{
				Bucket: internals3.Bucket,
				Key:    aws.String(filename),
			})
			if err != nil {
				if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
					continue
				}
				return nil, err
			}
			content, err := ioutil.ReadAll(output.Body)
			if err != nil {
				return nil, err
			}
			attachmentContent[filename] = content
		}

		var pbAttachment = models.AttachmentToPb(a)
		pbAttachment.AttachmentContent = attachmentContent
		res.Attachments[a.AttachmentId] = pbAttachment
	}
	return res, nil
}

func (p *Server) DeleteAttachments(ctx context.Context, req *pb.DeleteAttachmentsRequest) (*pb.DeleteAttachmentsResponse, error) {
	as, err := getAttachments(ctx, req.AttachmentId)
	if err != nil {
		return nil, err
	}
	for _, a := range as {
		// filenames with prefix
		var filenames []string
		// prepare object keys with prefix
		if len(req.Filename) == 0 {
			output, err := internals3.S3.ListObjects(&s3.ListObjectsInput{
				Bucket: internals3.Bucket,
				Prefix: aws.String(a.GetObjectPrefix()),
			})
			if err != nil {
				return nil, err
			}
			for _, o := range output.Contents {
				if o.Key != nil {
					filenames = append(filenames, *o.Key)
				}
			}
		} else {
			for _, filename := range req.Filename {
				filenames = append(filenames, a.GetObjectName(filename))
			}
		}

		// delete object
		for _, filename := range filenames {
			_, err := internals3.S3.DeleteObject(&s3.DeleteObjectInput{
				Bucket: internals3.Bucket,
				Key:    aws.String(filename),
			})
			if err != nil {
				if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
					continue
				}
				return nil, err
			}
		}
	}
	if len(req.Filename) == 0 {
		err = removeAttachments(ctx, req.AttachmentId)
	}
	return &pb.DeleteAttachmentsResponse{
		AttachmentId: req.AttachmentId,
		Filename:     req.Filename,
	}, err
}
