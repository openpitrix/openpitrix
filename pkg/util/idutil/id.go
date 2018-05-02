// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package idutil

import (
	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
}

func newId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// format likes: B6BZVN3mOPvx
func GetUuid(prefix string) string {
	id := newId()
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

const Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"

// format likes: 300m50zn91nwz5
func GetUuid36(prefix string) string {
	id := newId()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36
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
