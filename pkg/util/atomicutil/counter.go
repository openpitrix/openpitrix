// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package atomicutil

import "sync/atomic"

type Counter int32

func (c *Counter) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(c), n)
}
func (c *Counter) Get() int32 {
	return atomic.LoadInt32((*int32)(c))
}
