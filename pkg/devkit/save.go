// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

// Reference: vendor/k8s.io/helm/pkg/chartutil/save.go:97

func Save(c *opapp.OpApp, outDir string) (string, error) {
	// Create archive
	if fi, err := os.Stat(outDir); err != nil {
		return "", err
	} else if !fi.IsDir() {
		return "", fmt.Errorf("location [%s] is not a directory", outDir)
	}

	if c.Metadata == nil {
		return "", fmt.Errorf("no [%s] data", PackageJson)
	}

	cfile := c.Metadata
	if cfile.Name == "" {
		return "", fmt.Errorf("no package name specified (%s)", PackageJson)
	} else if cfile.Version == "" {
		return "", fmt.Errorf("no app version specified (%s)", PackageJson)
	}

	filename := cfile.GetPackageName()
	filename = filepath.Join(outDir, filename)
	if stat, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filename), 0755); !os.IsExist(err) {
			return "", err
		}
	} else if !stat.IsDir() {
		return "", fmt.Errorf("[%s] is not a directory", filepath.Dir(filename))
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	// Wrap in gzip writer
	zipper := gzip.NewWriter(f)

	// Wrap in tar writer
	twriter := tar.NewWriter(zipper)
	rollback := false
	defer func() {
		twriter.Close()
		zipper.Close()
		f.Close()
		if rollback {
			os.Remove(filename)
		}
	}()

	if err := writeTarContents(twriter, c, ""); err != nil {
		rollback = true
	}
	return filename, err
}

func writeTarContents(out *tar.Writer, c *opapp.OpApp, prefix string) error {
	base := filepath.Join(prefix, c.Metadata.Name)

	cdata, err := json.Marshal(c.Metadata)
	if err != nil {
		return err
	}
	if err := writeToTar(out, base+"/"+PackageJson, cdata); err != nil {
		return err
	}

	if c.ConfigTemplate != nil && len(c.ConfigTemplate.Raw) > 0 {
		if err := writeToTar(out, base+"/"+ConfigJson, []byte(c.ConfigTemplate.Raw)); err != nil {
			return err
		}
	}
	if c.ClusterConfTemplate != nil && len(c.ClusterConfTemplate.Raw) > 0 {
		if err := writeToTar(out, base+"/"+ClusterJsonTmpl, []byte(c.ClusterConfTemplate.Raw)); err != nil {
			return err
		}
	}

	// Save files
	for _, f := range c.Files {
		n := filepath.Join(base, f.TypeUrl)
		if err := writeToTar(out, n, f.Value); err != nil {
			return err
		}
	}

	return nil
}

// writeToTar writes a single file to a tar archive.
func writeToTar(out *tar.Writer, name string, body []byte) error {
	h := &tar.Header{
		Name: name,
		Mode: 0755,
		Size: int64(len(body)),
	}
	if err := out.WriteHeader(h); err != nil {
		return err
	}
	if _, err := out.Write(body); err != nil {
		return err
	}
	return nil
}

func savePackageJson(filename string, metadata *opapp.Metadata) error {
	out, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0644)
}
