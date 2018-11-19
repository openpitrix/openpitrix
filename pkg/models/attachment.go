// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAttachmentId() string {
	return fmt.Sprintf(
		"att-%s-%s",
		idutil.GetAttachmentPrefix(),
		idutil.GetUuid(""),
	)
}

// "att-xxyyzzqqwweerrttyyaassddff-iooihogwe" => "xx/yy/zz/qq/"
func getAttachmentObjectPrefix(attachmentId string) string {
	return fmt.Sprintf(
		"%s/%s/%s/%s/",
		attachmentId[4:6], attachmentId[6:8], attachmentId[8:10], attachmentId[10:12])
}

type Attachment struct {
	AttachmentId   string
	AttachmentType string
	CreateTime     time.Time
}

func (a Attachment) GetObjectName(filename string) string {
	return a.GetObjectPrefix() + filename
}

func (a Attachment) GetObjectPrefix() string {
	return getAttachmentObjectPrefix(a.AttachmentId)
}

func (a Attachment) GetAttachmentType() pb.AttachmentType {
	return pb.AttachmentType(pb.AttachmentType_value[a.AttachmentType])
}

var AttachmentColumns = db.GetColumnsFromStruct(&Attachment{})

func NewAttachment(t string) *Attachment {
	return &Attachment{
		AttachmentId:   NewAttachmentId(),
		AttachmentType: t,
		CreateTime:     time.Now(),
	}
}

func AttachmentToPb(attachment *Attachment) *pb.Attachment {
	pbAttachment := pb.Attachment{}
	pbAttachment.AttachmentId = attachment.AttachmentId
	pbAttachment.AttachmentType = attachment.GetAttachmentType()
	pbAttachment.CreateTime = pbutil.ToProtoTimestamp(attachment.CreateTime)
	return &pbAttachment
}

func AttachmentsToPbs(attachments []*Attachment) (pbAttachments []*pb.Attachment) {
	for _, attachment := range attachments {
		pbAttachments = append(pbAttachments, AttachmentToPb(attachment))
	}
	return
}
