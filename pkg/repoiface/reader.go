// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface/wrapper"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type Reader struct {
	RepoInterface
	repo *pb.Repo
}

const IndexYaml = "index.yaml"

func NewReader(ctx context.Context, r *pb.Repo) (*Reader, error) {
	iface, err := New(ctx, r.Type.GetValue(), r.Url.GetValue(), r.Credential.GetValue())
	reader := Reader{
		repo:          r,
		RepoInterface: iface,
	}
	return &reader, err
}

type AppVersions map[string]Versions

type Versions map[string]string

func (r *Reader) getIndexYaml(ctx context.Context) ([]byte, error) {
	var content []byte
	exists, err := r.CheckFile(ctx, IndexYaml)
	if err != nil {
		return content, err
	}
	if exists {
		return r.ReadFile(ctx, IndexYaml)
	}
	return content, nil
}

func (r *Reader) isK8s() bool {
	return stringutil.StringIn(constants.ProviderKubernetes, r.repo.GetProviders())
}

func (r *Reader) LoadPackage(ctx context.Context, pkg []byte) (wrapper.VersionInterface, error) {
	if r.isK8s() {
		return LoadPackage(ctx, Helm, pkg)
	} else {
		return LoadPackage(ctx, Vmbased, pkg)
	}
}

func (r *Reader) GetAppVersions(ctx context.Context) (AppVersions, error) {
	index, err := r.GetIndex(ctx)
	if err != nil {
		return nil, err
	}
	var appVersions = make(AppVersions)
	entries := index.GetEntries()
	for appName, vs := range entries {
		var versions = make(Versions)
		for _, v := range vs {
			versions[v.GetVersion()] = v.GetAppVersion()
		}
		appVersions[appName] = versions
	}
	return appVersions, nil
}

func (r *Reader) GetIndex(ctx context.Context) (wrapper.IndexInterface, error) {
	content, err := r.getIndexYaml(ctx)
	if err != nil {
		return nil, err
	}
	if r.isK8s() {
		var indexFile = new(repo.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		return &wrapper.HelmIndexWrapper{IndexFile: indexFile}, nil
	} else {
		var indexFile = new(opapp.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		return &wrapper.OpIndexWrapper{IndexFile: indexFile}, nil
	}
}

func (r *Reader) AddPackage(ctx context.Context, pkg []byte) error {
	content, err := r.getIndexYaml(ctx)
	if err != nil {
		return err
	}
	hash, err := provenance.Digest(bytes.NewReader(pkg))
	if err != nil {
		return err
	}
	if r.isK8s() {
		var indexFile = repo.NewIndexFile()
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return errors.Wrap(err, "decode yaml failed")
		}
		app, err := chartutil.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			return err
		}
		w := wrapper.HelmVersionWrapper{ChartVersion: &repo.ChartVersion{Metadata: app.Metadata}}
		indexFile.Add(app.Metadata, w.GetPackageName(), "", hash)
		content, err = yamlutil.Encode(indexFile)
		if err != nil {
			return err
		}
	} else {
		var indexFile = opapp.NewIndexFile()
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return errors.Wrap(err, "decode yaml failed")
		}
		app, err := devkit.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			return err
		}
		w := wrapper.OpVersionWrapper{OpVersion: &opapp.OpVersion{Metadata: app.Metadata}}
		indexFile.Add(app.Metadata, w.GetPackageName(), "", hash)
		content, err = yamlutil.Encode(indexFile)
		if err != nil {
			return err
		}
	}
	return r.WriteFile(ctx, IndexYaml, content)
}

func (r *Reader) DeletePackage(ctx context.Context, appName, version string) error {
	content, err := r.getIndexYaml(ctx)
	if err != nil {
		return err
	}
	pkgName := fmt.Sprintf("%s-%s.tgz", appName, version)

	if r.isK8s() {
		var indexFile = repo.NewIndexFile()
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return errors.Wrap(err, "decode yaml failed")
		}
		if versions, ok := indexFile.Entries[appName]; ok {
			var newVersions repo.ChartVersions
			for _, v := range versions {
				if v.Version != version {
					newVersions = append(newVersions, v)
				}
			}
			indexFile.Entries[appName] = newVersions
			if len(newVersions) == 0 {
				delete(indexFile.Entries, appName)
			}
		}
		content, err = yamlutil.Encode(indexFile)
	} else {
		var indexFile = opapp.NewIndexFile()
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return errors.Wrap(err, "decode yaml failed")
		}
		if versions, ok := indexFile.Entries[appName]; ok {
			var newVersions opapp.OpVersions
			for _, v := range versions {
				if v.Version != version {
					newVersions = append(newVersions, v)
				}
			}
			indexFile.Entries[appName] = newVersions
			if len(newVersions) == 0 {
				delete(indexFile.Entries, appName)
			}
		}
		content, err = yamlutil.Encode(indexFile)
	}
	if err != nil {
		return err
	}
	err = r.WriteFile(ctx, IndexYaml, content)
	if err != nil {
		return err
	}
	return r.DeleteFile(ctx, pkgName)
}
