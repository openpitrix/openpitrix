// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package archiveutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path"
	"strings"
)

type Files map[string][]byte

func Load(in io.Reader) (Files, error) {
	var files = make(Files)
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return files, err
	}
	defer unzipped.Close()

	tr := tar.NewReader(unzipped)
	for {
		b := bytes.NewBuffer(nil)
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return files, err
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

		if _, err := io.Copy(b, tr); err != nil {
			return files, err
		}

		files[n] = b.Bytes()
		b.Reset()
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files in app archive")
	}
	return files, err
}

func Save(files map[string][]byte, prefix ...string) ([]byte, error) {
	f := new(bytes.Buffer)
	zipper := gzip.NewWriter(f)

	// Wrap in tar writer
	twriter := tar.NewWriter(zipper)
	p := path.Join(prefix...)
	for name, content := range files {
		if err := writeToTar(twriter, path.Join(p, name), content); err != nil {
			return nil, err
		}
	}

	twriter.Close()
	zipper.Close()

	return f.Bytes(), nil
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
