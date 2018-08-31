// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Masterminds/semver"
	"k8s.io/helm/pkg/urlutil"

	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

const APIVersionV1 = "v1"

type IndexFile struct {
	APIVersion string                `json:"api_version"`
	Generated  time.Time             `json:"generated"`
	Entries    map[string]OpVersions `json:"entries"`
	//PublicKeys []string            `json:"publicKeys,omitempty"`
}

// NewIndexFile initializes an index.
func NewIndexFile() *IndexFile {
	return &IndexFile{
		APIVersion: APIVersionV1,
		Generated:  time.Now(),
		Entries:    map[string]OpVersions{},
	}
}

// loadIndex loads an index file and does minimal validity checking.
//
// This will fail if API OpVersion is not set (ErrNoAPIVersion) or if the unmarshal fails.
func loadIndex(data []byte) (*IndexFile, error) {
	i := &IndexFile{}
	if err := yamlutil.Decode(data, i); err != nil {
		return i, err
	}
	i.SortEntries()
	return i, nil
}

// LoadIndexFile takes a file at the given path and returns an IndexFile object
func LoadIndexFile(path string) (*IndexFile, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadIndex(b)
}

// Add adds a file to the index
// This can leave the index in an unsorted state
func (i IndexFile) Add(md *Metadata, filename, baseURL, digest string) {
	u := filename
	if baseURL != "" {
		var err error
		_, file := filepath.Split(filename)
		u, err = urlutil.URLJoin(baseURL, file)
		if err != nil {
			u = filepath.Join(baseURL, file)
		}
	}
	cr := &OpVersion{
		URLs:     []string{u},
		Metadata: md,
		Digest:   digest,
		Created:  time.Now(),
	}
	if ee, ok := i.Entries[md.Name]; !ok {
		i.Entries[md.Name] = OpVersions{cr}
	} else {
		i.Entries[md.Name] = append(ee, cr)
	}
}

// Has returns true if the index has an entry for a app with the given name and exact version.
func (i IndexFile) Has(name, version string) bool {
	_, err := i.Get(name, version)
	return err == nil
}

// SortEntries sorts the entries by version in descending order.
//
// In canonical form, the individual version records should be sorted so that
// the most recent release for every version is in the 0th slot in the
// Entries.AppVersions array. That way, tooling can predict the newest
// version without needing to parse SemVers.
func (i IndexFile) SortEntries() {
	for _, versions := range i.Entries {
		sort.Sort(sort.Reverse(versions))
	}
}

// Get returns the AppVersion for the given name.
//
// If version is empty, this will return the app with the highest version.
func (i IndexFile) Get(name, version string) (*OpVersion, error) {
	vs, ok := i.Entries[name]
	if !ok {
		return nil, fmt.Errorf("no app name found")
	}
	if len(vs) == 0 {
		return nil, fmt.Errorf("no app version found")
	}

	var constraint *semver.Constraints
	if len(version) == 0 {
		constraint, _ = semver.NewConstraint("*")
	} else {
		var err error
		constraint, err = semver.NewConstraint(version)
		if err != nil {
			return nil, err
		}
	}

	for _, ver := range vs {
		test, err := semver.NewVersion(ver.Version)
		if err != nil {
			continue
		}

		if constraint.Check(test) {
			return ver, nil
		}
	}
	return nil, fmt.Errorf("no app version found for [%s-%s]", name, version)
}

// WriteFile writes an index file to the given destination path.
//
// The mode on the file is set to 'mode'.
func (i IndexFile) WriteFile(dest string, mode os.FileMode) error {
	b, err := yamlutil.Encode(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, b, mode)
}

// Merge merges the given index file into this index.
//
// This merges by name and version.
//
// If one of the entries in the given index does _not_ already exist, it is added.
// In all other cases, the existing record is preserved.
//
// This can leave the index in an unsorted state
func (i *IndexFile) Merge(f *IndexFile) {
	for _, cvs := range f.Entries {
		for _, cv := range cvs {
			if !i.Has(cv.Name, cv.Version) {
				e := i.Entries[cv.Name]
				i.Entries[cv.Name] = append(e, cv)
			}
		}
	}
}
