// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"context"

	"openpitrix.io/openpitrix/pkg/logger"
)

type callback func() error

func (etcd *Etcd) Dlock(ctx context.Context, key string, cb callback) error {
	logger.Debugf("Create dlock with key [%s]", key)
	mutex := etcd.NewMutex(key)
	err := mutex.Lock(ctx)
	if err != nil {
		logger.Fatalf("Dlock lock error: %+v", err)
		return err
	}
	defer mutex.Unlock(ctx)
	err = cb()
	return err
}
