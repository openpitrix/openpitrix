// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/provenance"
	"k8s.io/helm/pkg/repo"

	"openpitrix.io/openpitrix/pkg/repoiface/wrapper"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/util/stringutil"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/pb"
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

func (r *Reader) GetPackagesName(ctx context.Context) ([]string, error) {
	content, err := r.ReadFile(ctx, IndexYaml)
	if err != nil {
		return nil, err
	}
	var packagesName []string
	if stringutil.StringIn(constants.ProviderKubernetes, r.repo.GetProviders()) {
		var indexFile = new(repo.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		for pkgName := range indexFile.Entries {
			packagesName = append(packagesName, pkgName)
		}
	} else {
		var indexFile = new(opapp.IndexFile)
		err = yamlutil.Decode(content, indexFile)
		if err != nil {
			return nil, errors.Wrap(err, "decode yaml failed")
		}
		for pkgName := range indexFile.Entries {
			packagesName = append(packagesName, pkgName)
		}
	}
	return packagesName, nil
}

func (r *Reader) AddPackage(ctx context.Context, pkg []byte) error {
	exists, err := r.CheckFile(ctx, IndexYaml)
	if err != nil {
		return err
	}
	var content []byte
	if exists {
		content, err = r.ReadFile(ctx, IndexYaml)
		if err != nil {
			return err
		}
	}

	hash, err := provenance.Digest(bytes.NewReader(pkg))
	if err != nil {
		return err
	}
	if stringutil.StringIn(constants.ProviderKubernetes, r.repo.GetProviders()) {
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
