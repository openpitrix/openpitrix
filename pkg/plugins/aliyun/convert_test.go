// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package aliyun

import "testing"

func TestConvertToInstanceType(t *testing.T) {
	type args struct {
		cpu    int
		memory int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: struct {
				cpu    int
				memory int
			}{cpu: 1, memory: 1 * G},
			want: "ecs.t5-lc1m1.small",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToInstanceType(tt.args.cpu, tt.args.memory)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToInstanceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToInstanceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToVolumeType(t *testing.T) {
	type args struct {
		volumeClass int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: struct{ volumeClass int }{volumeClass: 1},
			want: "cloud",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToVolumeType(tt.args.volumeClass)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToVolumeType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertToVolumeType() = %v, want %v", got, tt.want)
			}
		})
	}
}
