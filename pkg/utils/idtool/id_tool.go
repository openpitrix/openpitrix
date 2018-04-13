// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package idtool

import (
	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
}

// format likes: B6BZVN3mOPvx
func GetUuid(prefix string) string {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + i
}
