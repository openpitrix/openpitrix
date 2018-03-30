// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// readFileStat return a fileInfo describing the named file.
func readFileStat(name string) (fi fileInfo, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return
	}

	fi.Mode = stats.Mode()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return
	}

	fi.Md5 = fmt.Sprintf("%x", h.Sum(nil))
	return fi, nil
}
