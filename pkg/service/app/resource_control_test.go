// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
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

func Test_matchPackageFailedError(t *testing.T) {
	var err error
	var res *pb.ValidatePackageResponse

	err = fmt.Errorf("failed to decode package.json: invalid character")
	res = &pb.ValidatePackageResponse{}
	matchPackageFailedError(err, res)
	require.Equal(t, "decode failed, invalid character", res.ErrorDetails["package.json"])

	err = fmt.Errorf("error reading package.json: invalid character")
	res = &pb.ValidatePackageResponse{}
	matchPackageFailedError(err, res)
	require.Equal(t, "invalid character", res.ErrorDetails["package.json"])

	err = fmt.Errorf("missing file [package.json]")
	res = &pb.ValidatePackageResponse{}
	matchPackageFailedError(err, res)
	require.Equal(t, "not found", res.ErrorDetails["package.json"])
}
