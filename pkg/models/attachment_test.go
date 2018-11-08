// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"testing"
	"time"
)

func Test_getAttachmentObjectPrefix(t *testing.T) {
	type args struct {
		attachmentId string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"att-xxyyzzqqwweerrttyyaassddff-iooihogwe"},
			want: "xx/yy/zz/qq/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAttachmentObjectPrefix(tt.args.attachmentId); got != tt.want {
				t.Errorf("getAttachmentObjectPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachment_RemoveObjectName(t *testing.T) {
	type fields struct {
		AttachmentId   string
		AttachmentType string
		CreateTime     time.Time
	}
	type args struct {
		filenameWithPrefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "",
			args:   args{"xx/yy/zz/qq/file/name.jpg"},
			want:   "file/name.jpg",
			fields: fields{AttachmentId: "att-xxyyzzqqwweerrttyyaassddff-iooihogwe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attachment{
				AttachmentId: tt.fields.AttachmentId,
				CreateTime:   tt.fields.CreateTime,
			}
			if got := a.RemoveObjectName(tt.args.filenameWithPrefix); got != tt.want {
				t.Errorf("Attachment.RemoveObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}
