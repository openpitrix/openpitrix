// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3/concurrency"
)

type Mutex struct {
	*concurrency.Mutex
}

func (etcd *Etcd) NewMutex(key string) (*Mutex, error) {
	session, err := concurrency.NewSession(etcd.Client)
	if err != nil {
		return nil, err
	}
	return &Mutex{concurrency.NewMutex(session, key)}, nil
}

// Lock locks the mutex with a cancelable context. If the context is canceled
// while trying to acquire the lock, the mutex tries to clean its stale lock entry.
func (m *Mutex) Lock(ctx context.Context) error {
	return m.Mutex.Lock(ctx)
}

func (m *Mutex) Unlock(ctx context.Context) error {
	return m.Mutex.Unlock(ctx)
}
