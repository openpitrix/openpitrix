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

func TestURLJoin(t *testing.T) {
	type args struct {
		repoUrl  string
		fileName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"http://some.com", "right"}, "http://some.com/right"},
		{"", args{"http://some.com", "right/right"}, "http://some.com/right/right"},
		{"", args{"http://some.com", "http://right.com"}, "http://right.com"},
		{"", args{"http://some.com/right", "../left"}, "http://some.com/left"},
		{"", args{"http://some.com/subdir/", "index.yaml"}, "http://some.com/subdir/index.yaml"},
		{"", args{"http://some.com/subdir", "index.yaml"}, "http://some.com/subdir/index.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLJoin(tt.args.repoUrl, tt.args.fileName); got != tt.want {
				t.Errorf("URLJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
