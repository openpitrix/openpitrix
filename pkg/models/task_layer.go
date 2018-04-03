// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type TaskLayer struct {
	Tasks []*Task
	Child *TaskLayer
}

// WalkFunc is a callback type for use with TaskLayer.WalkTree
type WalkFunc func(parent *TaskLayer, current *TaskLayer) error

func (t *TaskLayer) WalkTree(cb WalkFunc) error {
	return walkTaskLayerTree(nil, t, cb)
}

func (t *TaskLayer) Leaf() *TaskLayer {
	current := t
	for {
		if current.Child == nil {
			return current
		} else {
			current = current.Child
		}
	}
}

func walkTaskLayerTree(parent *TaskLayer, current *TaskLayer, cb WalkFunc) error {
	err := cb(parent, current)
	if err != nil {
		return err
	}

	if current.Child == nil {
		return nil
	} else {
		err = walkTaskLayerTree(current, current.Child, cb)
		return err
	}
}
