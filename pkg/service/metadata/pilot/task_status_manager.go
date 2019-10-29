// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"sync"

	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type TaskStatusManager struct {
	db map[string]pbtypes.SubTaskStatus
	sync.Mutex
}

func NewTaskStatusManager() *TaskStatusManager {
	return &TaskStatusManager{
		db: make(map[string]pbtypes.SubTaskStatus),
	}
}

func (p *TaskStatusManager) GetStatus(id string) (v pbtypes.SubTaskStatus, ok bool) {
	p.Lock()
	defer p.Unlock()

	v, ok = p.db[id]
	return
}

func (p *TaskStatusManager) PutStatus(v pbtypes.SubTaskStatus) {
	p.Lock()
	defer p.Unlock()

	p.db[v.TaskId] = v
}
