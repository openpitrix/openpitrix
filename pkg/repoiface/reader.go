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
	"openpitrix.io/openpitrix/pkg/logger"
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
		p, err := chartutil.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			logger.Error(ctx, "Failed to load package, error: %+v", err)
			return nil, err
		}

		pkgVersion := wrapper.HelmVersionWrapper{ChartVersion: &repo.ChartVersion{Metadata: p.Metadata}}
		return pkgVersion, nil
	} else {
		p, err := devkit.LoadArchive(bytes.NewReader(pkg))
		if err != nil {
			logger.Error(ctx, "Failed to load package, error: %+v", err)
			return nil, err
		}

		pkgVersion := wrapper.OpVersionWrapper{OpVersion: &opapp.OpVersion{Metadata: p.Metadata}}
		return pkgVersion, nil
	}
}

func (r *Reader) GetAppVersions(ctx context.Context) (AppVersions, error) {
	content, err := r.getIndexYaml(ctx)
	var appVersions = make(AppVersions)
	if r.isK8s() {
		var indexFile = new(repo.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		for appName, versions := range indexFile.Entries {
			var vs = make(Versions)
			for _, version := range versions {
				vs[version.Version] = version.AppVersion
			}
			appVersions[appName] = vs
		}
	} else {
		var indexFile = new(opapp.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		for appName, versions := range indexFile.Entries {
			var vs = make(Versions)
			for _, version := range versions {
				vs[version.Version] = version.AppVersion
			}
			appVersions[appName] = vs
		}
	}
	return appVersions, nil
}

func (r *Reader) AddPackage(ctx context.Context, pkg []byte) error {
	content, err := r.getIndexYaml(ctx)
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
	}
	if err != nil {
		return err
	}
	return r.WriteFile(ctx, IndexYaml, content)
}

func (r *Reader) DeletePackage(ctx context.Context, appName, version string) error {
	content, err := r.getIndexYaml(ctx)
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
