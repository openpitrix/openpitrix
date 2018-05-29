// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// fileInfo describes a configuration file and is returned by readFileStat.
type fileInfo struct {
	Uid  uint32
	Gid  uint32
	Mode os.FileMode
	Md5  string
}

// ensureFileDir ensure file's dir is exist.
func ensureFileDir(file string) error {
	dir := filepath.Dir(file)
	if dir == "" {
		return nil
	}
	if dirExists(dir) {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return false
	}
	return true
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func fileNotExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}

func getFuncName(skips ...int) string {
	var skip = 1
	if len(skips) > 0 {
		skip = skips[0]
	}

	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	name := runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

func strInStrList(s string, ss []string) bool {
	for _, t := range ss {
		if s == t {
			return true
		}
	}
	return false
}
