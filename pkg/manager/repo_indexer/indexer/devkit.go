// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package indexer

import "openpitrix.io/openpitrix/pkg/pb"

type devkitIndexer struct {
	indexer
}

func NewDevkitIndexer(repo *pb.Repo) *devkitIndexer {
	return &devkitIndexer{
		indexer: indexer{repo: repo},
	}
}

func (i *devkitIndexer) IndexRepo() error {
	return nil
}
