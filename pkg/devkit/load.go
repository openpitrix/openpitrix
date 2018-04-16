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

	"openpitrix.io/openpitrix/pkg/devkit/app"
)

func Load(name string) (*app.App, error) {
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
func LoadFile(name string) (*app.App, error) {
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
func LoadArchive(in io.Reader) (*app.App, error) {
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return &app.App{}, err
	}
	defer unzipped.Close()

	var files []app.BufferedFile
	tr := tar.NewReader(unzipped)
	for {
		b := bytes.NewBuffer(nil)
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &app.App{}, err
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
			return nil, fmt.Errorf("%s not in base directory", PackageJson)
		}

		if _, err := io.Copy(b, tr); err != nil {
			return &app.App{}, err
		}

		files = append(files, app.BufferedFile{Name: n, Data: b.Bytes()})
		b.Reset()
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files in app archive")
	}

	return LoadFiles(files)
}

func LoadDir(dir string) (*app.App, error) {
	topdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// Just used for errors.
	c := &app.App{}

	var files []app.BufferedFile
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
			return fmt.Errorf("error reading [%s]: %s", n, err)
		}

		files = append(files, app.BufferedFile{Name: n, Data: data})
		return nil
	})
	if err != nil {
		return c, err
	}

	return LoadFiles(files)
}

// LoadFiles loads from in-memory files.
func LoadFiles(files []app.BufferedFile) (*app.App, error) {
	c := &app.App{}

	for _, f := range files {
		if f.Name == PackageJson {
			m, err := app.UnmarshalMetadata(f.Data)
			if err != nil {
				return c, err
			}
			c.Metadata = m
		} else if f.Name == ClusterJsonTmpl {
			c.ClusterConfTemplate = &app.ClusterConfTemplate{Raw: string(f.Data)}
		} else if f.Name == ConfigJson {
			m, err := app.UnmarshalConfigTemplate(f.Data)
			if err != nil {
				return c, err
			}
			c.ConfigTemplate = m
		} else {
			c.Files = append(c.Files, f)
		}
	}

	if c.Metadata == nil {
		return c, fmt.Errorf("version metadata (package) missing")
	}
	if c.Metadata.Name == "" {
		return c, fmt.Errorf("invalid version (package): name must not be empty")
	}
	// Validate default config
	config := c.ConfigTemplate.GetDefaultConfig()
	err := app.ValidateClusterConfTmpl(c.ClusterConfTemplate, &config)
	return c, err
}
