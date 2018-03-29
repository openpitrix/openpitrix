// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"time"

	"github.com/Masterminds/semver"
)

// Versions is a list of versioned app references.
// Implements a sorter on App.
type Versions []*Version

// Len returns the length.
func (c Versions) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c Versions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c Versions) Less(a, b int) bool {
	// Failed parse pushes to the back.
	i, err := semver.NewVersion(c[a].Version)
	if err != nil {
		return true
	}
	j, err := semver.NewVersion(c[b].Version)
	if err != nil {
		return false
	}
	return i.LessThan(j)
}

// Version represents a app entry in the IndexFile
type Version struct {
	*Metadata
	URLs    []string  `json:"urls"`
	Created time.Time `json:"created,omitempty"`
	Removed bool      `json:"removed,omitempty"`
	Digest  string    `json:"digest,omitempty"`
}
