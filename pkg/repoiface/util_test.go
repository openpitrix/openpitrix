// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"testing"
)

func Test_getBucketPrefix(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantBucket string
		wantPrefix string
	}{
		{
			"",
			args{"/bucket_name/prefix/long/name"},
			"bucket_name",
			"prefix/long/name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBucket, gotPrefix := getBucketPrefix(tt.args.path)
			if gotBucket != tt.wantBucket {
				t.Errorf("getBucketPrefix() gotBucket = %v, want %v", gotBucket, tt.wantBucket)
			}
			if gotPrefix != tt.wantPrefix {
				t.Errorf("getBucketPrefix() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
		})
	}
}
