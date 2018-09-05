// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"time"

	"github.com/Masterminds/semver"
)

// OpVersions is a list of versioned app references.
// Implements a sorter on OpApp.
type OpVersions []*OpVersion

// Len returns the length.
func (c OpVersions) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c OpVersions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c OpVersions) Less(a, b int) bool {
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

// OpVersion represents a app entry in the IndexFile
type OpVersion struct {
	*Metadata
	URLs    []string  `json:"urls"`
	Created time.Time `json:"created,omitempty"`
	Removed bool      `json:"removed,omitempty"`
	Digest  string    `json:"digest,omitempty"`
}

func (h OpVersion) GetUrls() []string { return h.URLs }
