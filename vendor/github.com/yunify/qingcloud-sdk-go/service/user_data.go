// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package service

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

var _ fmt.State
var _ time.Time

type UserDataService struct {
	Config     *config.Config
	Properties *UserDataServiceProperties
}

type UserDataServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) UserData(zone string) (*UserDataService, error) {
	properties := &UserDataServiceProperties{
		Zone: &zone,
	}

	return &UserDataService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/userdata/upload_userdata_attachment.html
func (s *UserDataService) UploadUserDataAttachment(i *UploadUserDataAttachmentInput) (*UploadUserDataAttachmentOutput, error) {
	if i == nil {
		i = &UploadUserDataAttachmentInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UploadUserDataAttachment",
		RequestMethod: "POST",
	}

	x := &UploadUserDataAttachmentOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type UploadUserDataAttachmentInput struct {
	AttachmentContent *string `json:"attachment_content" name:"attachment_content" location:"params"` // Required
	AttachmentName    *string `json:"attachment_name" name:"attachment_name" location:"params"`
}

func (v *UploadUserDataAttachmentInput) Validate() error {

	if v.AttachmentContent == nil {
		return errors.ParameterRequiredError{
			ParameterName: "AttachmentContent",
			ParentName:    "UploadUserDataAttachmentInput",
		}
	}

	return nil
}

type UploadUserDataAttachmentOutput struct {
	Message      *string `json:"message" name:"message"`
	Action       *string `json:"action" name:"action" location:"elements"`
	AttachmentID *string `json:"attachment_id" name:"attachment_id" location:"elements"`
	RetCode      *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
