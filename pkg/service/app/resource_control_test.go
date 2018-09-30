// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
)

func Test_getAppVersionStatus(t *testing.T) {
	type args struct {
		defaultStatus string
		currentStatus string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{constants.StatusActive, constants.StatusDraft},
			want: constants.StatusActive,
		},
		{
			args: args{constants.StatusDraft, constants.StatusDraft},
			want: constants.StatusDraft,
		},
		{
			args: args{constants.StatusDraft, constants.StatusActive},
			want: constants.StatusActive,
		},
		{
			args: args{constants.StatusDraft, constants.StatusDeleted},
			want: constants.StatusDraft,
		},
		{
			args: args{constants.StatusDraft, constants.StatusSuspended},
			want: constants.StatusSuspended,
		},
		{
			args: args{constants.StatusActive, constants.StatusActive},
			want: constants.StatusActive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAppVersionStatus(tt.args.defaultStatus, tt.args.currentStatus); got != tt.want {
				t.Errorf("getAppVersionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
