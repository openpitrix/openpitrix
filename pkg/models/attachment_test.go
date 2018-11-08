// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "testing"

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
