// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package labelutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			args: "runtime=qingcloud&zone=pk3a&env=test",
		},
		{
			args: "runtime=kubernetes&env=dev",
		},
		{
			args: "runtime=kubernetes&team=openpitrix",
		},
		{
			args:    "runtime=&team=openpitrix",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args, got.ToString())
			}
		})
	}
}
