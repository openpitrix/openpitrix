// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gziputil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type ArchiveFiles map[string][]byte

func LoadArchive(in io.Reader, include ...string) (ArchiveFiles, error) {
	afs := make(ArchiveFiles)
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return afs, err
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
			return afs, err
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

		if len(include) > 0 && !stringutil.StringIn(n, include) {
			continue
		}

		if _, err := io.Copy(b, tr); err != nil {
			return afs, err
		}

		afs[n] = b.Bytes()
		b.Reset()
	}

	return afs, nil
}
