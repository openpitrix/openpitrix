// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package labelutil

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		labelString string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			args: args{
				"runtime=qingcloud&zone=pk3a&env=test",
			},
			want: map[string][]string{
				"runtime": {"qingcloud"},
				"zone":    {"pk3a"},
				"env":     {"test"},
			},
		},
		{
			args: args{
				"runtime=kubernetes&env=dev",
			},
			want: map[string][]string{
				"runtime": {"kubernetes"},
				"env":     {"dev"},
			},
		},
		{
			args: args{
				"runtime=kubernetes&team=openpitrix",
			},
			want: map[string][]string{
				"runtime": {"kubernetes"},
				"team":    {"openpitrix"},
			},
		},
		{
			args: args{
				"runtime=&team=openpitrix",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.labelString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
