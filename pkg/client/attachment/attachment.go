// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package attachmentclient

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewAttachmentManagerClient() (pb.AttachmentManagerClient, error) {
	conn, err := manager.NewClient(constants.AttachmentManagerHost, constants.AttachmentManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewAttachmentManagerClient(conn), err
}

func NewAttachmentServiceClient() (pb.AttachmentServiceClient, error) {
	conn, err := manager.NewClient(constants.AttachmentManagerHost, constants.AttachmentManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewAttachmentServiceClient(conn), err
}
