// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Indexer struct {
	*pi.Pi
	queue *etcd.Queue
}

func NewIndexer(pi *pi.Pi) *Indexer {
	return &Indexer{Pi: pi, queue: pi.Etcd.NewQueue("repo-indexer")}
}

func (i *Indexer) IndexRepo(repoId string) {

}
