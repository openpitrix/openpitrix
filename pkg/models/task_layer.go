// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type TaskLayer struct {
	Tasks []*Task
	Child *TaskLayer
}

// WalkFunc is a callback type for use with TaskLayer.WalkTree
type WalkFunc func(parent *TaskLayer, current *TaskLayer)

func (t *TaskLayer) WalkTree(cb WalkFunc) {
	walkTaskLayerTree(nil, t, cb)
}

func (t *TaskLayer) IsLeaf() bool {
	if t.Child == nil {
		return true
	} else {
		return false
	}
}

func (t *TaskLayer) Leaf() *TaskLayer {
	current := t
	for {
		if current.IsLeaf() {
			return current
		} else {
			current = current.Child
		}
	}
}

func walkTaskLayerTree(parent *TaskLayer, current *TaskLayer, cb WalkFunc) {
	cb(parent, current)
	if current == nil || current.Child == nil {
		return
	} else {
		walkTaskLayerTree(current, current.Child, cb)
	}
}

func (t *TaskLayer) Append(target *TaskLayer) *TaskLayer {
	current := t.Leaf()
	if target == nil {
		return current
	} else {
		current.Child = target
		return current.Leaf()
	}
}
