// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import recipe "github.com/coreos/etcd/contrib/recipes"

type Queue struct {
	*recipe.Queue
}

func (etcd *Etcd) NewQueue(topic string) *Queue {
	return &Queue{recipe.NewQueue(etcd.Client, topic)}
}

func (q *Queue) Enqueue(val string) error {
	return q.Queue.Enqueue(val)
}

// Dequeue returns Enqueue()'d elements in FIFO order. If the
// queue is empty, Dequeue blocks until elements are available.
func (q *Queue) Dequeue() (string, error) {
	return q.Queue.Dequeue()
}
