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

func getAttachment(ctx context.Context, attachmentId string) (*models.Attachment, error) {
	var a models.Attachment
	_, err := pi.Global().DB(ctx).
		Select(models.AttachmentColumns...).
		From(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentId)).
		Load(&a)
	return &a, err
}

func removeAttachments(ctx context.Context, attachmentIds []string) error {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableAttachment).
		Where(db.Eq(constants.ColumnAttachmentId, attachmentIds)).
		Exec()
	return err
}

func listAttachmentFilenames(ctx context.Context, attachment *models.Attachment) ([]string, error) {
	// with prefix
	var filenames []string
	output, err := internals3.S3.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: internals3.Bucket,
		Prefix: aws.String(attachment.GetObjectPrefix()),
	})
	if err != nil {
		return nil, err
	}
	for _, o := range output.Contents {
		if o.Key != nil {
			filenames = append(filenames, attachment.RemoveObjectName(*o.Key))
		}
	}
	return filenames, nil
}

func deleteAttachmentFiles(ctx context.Context, attachment *models.Attachment, filename ...string) error {
	// filenames with prefix
	var filenames []string
	var err error
	// prepare object keys with prefix
	if len(filename) == 0 {
		filenames, err = listAttachmentFilenames(ctx, attachment)
		if err != nil {
			return err
		}
	} else {
		// with prefix
		for _, f := range filename {
			filenames = append(filenames, f)
		}
	}

	// delete object
	for _, filename := range filenames {
		_, err := internals3.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: internals3.Bucket,
			Key:    aws.String(attachment.GetObjectName(filename)),
		})
		if err != nil {
			if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
				continue
			}
			return err
		}
	}
	return nil
}

type contents interface {
	GetAttachmentContent() map[string][]byte
}

func putAttachmentFiles(ctx context.Context, attachment *models.Attachment, contents contents) error {
	for filename, content := range contents.GetAttachmentContent() {
		_, err := internals3.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: internals3.Bucket,
			Key:    aws.String(attachment.GetObjectName(filename)),
			Body:   bytes.NewReader(content),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type getAttachmentReq interface {
	GetFilename() []string
	GetIgnoreContent() bool
}

func getFile(ctx context.Context, attachment *models.Attachment, filename string) (*s3.GetObjectOutput, error) {
	return internals3.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: internals3.Bucket,
		Key:    aws.String(attachment.GetObjectName(filename)),
	})
}

func getAttachmentFiles(ctx context.Context, attachments []*models.Attachment, req getAttachmentReq) ([]*pb.Attachment, error) {
	var err error
	var pbAtts []*pb.Attachment
	for _, a := range attachments {
		var attachmentContent = make(map[string][]byte)
		// filenames with prefix
		var filenames []string

		// prepare object keys with prefix
		if len(req.GetFilename()) == 0 {
			filenames, err = listAttachmentFilenames(ctx, a)
			if err != nil {
				return nil, err
			}
		} else {
			for _, filename := range req.GetFilename() {
				filenames = append(filenames, filename)
			}
		}

		// get object content
		for _, filename := range filenames {
			var content []byte
			if req.GetIgnoreContent() {
				attachmentContent[filename] = content
				continue
			}
			output, err := getFile(ctx, a, filename)
			if err != nil {
				if e, ok := err.(awserr.Error); ok && e.Code() == s3.ErrCodeNoSuchKey {
					continue
				}
				return nil, err
			}
			content, err = ioutil.ReadAll(output.Body)
			if err != nil {
				return nil, err
			}
			attachmentContent[filename] = content
		}

		var pbAttachment = models.AttachmentToPb(a)
		pbAttachment.AttachmentContent = attachmentContent
		pbAtts = append(pbAtts, pbAttachment)
	}
	return pbAtts, nil
}
