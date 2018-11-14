// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes/any"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

func Load(name string) (*opapp.OpApp, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		if validApp, err := IsAppDir(name); !validApp {
			return nil, err
		}
		return LoadDir(name)
	}
	return LoadFile(name)
}

// LoadFile loads from an archive file.
func LoadFile(name string) (*opapp.OpApp, error) {
	if fi, err := os.Stat(name); err != nil {
		return nil, err
	} else if fi.IsDir() {
		return nil, fmt.Errorf("cannot load a directory")
	}

	raw, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer raw.Close()

	return LoadArchive(raw)
}

// LoadArchive loads from a reader containing a compressed tar archive.
func LoadArchive(in io.Reader) (*opapp.OpApp, error) {
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return &opapp.OpApp{}, err
	}
	defer unzipped.Close()

	var files []opapp.BufferedFile
	tr := tar.NewReader(unzipped)
	for {
		b := bytes.NewBuffer(nil)
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &opapp.OpApp{}, err
		}

		if hd.FileInfo().IsDir() {
			// Use this instead of hd.Typeflag because we don't have to do any
			// inference chasing.
			continue
		}

		// Archive could contain \ if generated on Windows
		delimiter := "/"
		if strings.ContainsRune(hd.Name, '\\') {
			delimiter = "\\"
		}

		parts := strings.Split(hd.Name, delimiter)
		n := strings.Join(parts[1:], delimiter)

		// Normalize the path to the / delimiter
		n = strings.Replace(n, delimiter, "/", -1)

		if parts[0] == PackageJson {
			return nil, fmt.Errorf("[%s] not in base directory", PackageJson)
		}

		if _, err := io.Copy(b, tr); err != nil {
			return &opapp.OpApp{}, err
		}

		files = append(files, opapp.BufferedFile{Name: n, Data: b.Bytes()})
		b.Reset()
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files in app archive")
	}

	return LoadFiles(files)
}

func LoadDir(dir string) (*opapp.OpApp, error) {
	topdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// Just used for errors.
	c := &opapp.OpApp{}

	var files []opapp.BufferedFile
	topdir += string(filepath.Separator)

	err = filepath.Walk(topdir, func(name string, fi os.FileInfo, err error) error {
		n := strings.TrimPrefix(name, topdir)

		// Normalize to / since it will also work on Windows
		n = filepath.ToSlash(n)

		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(name)
		if err != nil {
			return fmt.Errorf("failed to read [%s]: %+v", n, err)
		}

		files = append(files, opapp.BufferedFile{Name: n, Data: data})
		return nil
	})
	if err != nil {
		return c, err
	}

	return LoadFiles(files)
}

// LoadFiles loads from in-memory files.
func LoadFiles(files []opapp.BufferedFile) (*opapp.OpApp, error) {
	c := &opapp.OpApp{}

	for _, f := range files {
		if f.Name == PackageJson {
			m, err := opapp.DecodePackageJson(f.Data)
			if err != nil {
				return c, err
			}
			c.Metadata = m
		} else if f.Name == ClusterJsonTmpl {
			c.ClusterConfTemplate = &opapp.ClusterConfTemplate{Raw: string(f.Data)}
		} else if f.Name == ConfigJson {
			m, err := opapp.DecodeConfigJson(f.Data)
			if err != nil {
				return c, err
			}
			c.ConfigTemplate = m
		} else {
			c.Files = append(c.Files, &any.Any{TypeUrl: f.Name, Value: f.Data})
		}
	}

	if c.Metadata == nil {
		return c, fmt.Errorf("missing file [%s]", PackageJson)
	}
	if c.ClusterConfTemplate == nil {
		return c, fmt.Errorf("missing file [%s]", ClusterJsonTmpl)
	}
	if c.ConfigTemplate == nil {
		return c, fmt.Errorf("missing file [%s]", ConfigJson)
	}
	if c.Metadata.Name == "" {
		return c, fmt.Errorf("failed to load [%s]: name must not be empty", PackageJson)
	}
	// Validate default config
	config := c.ConfigTemplate.GetDefaultConfig()
	err := opapp.ValidateClusterConfTmpl(c.ClusterConfTemplate, config)
	return c, err
}
