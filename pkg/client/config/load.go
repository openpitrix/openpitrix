// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	flag "github.com/spf13/pflag"
)

func AddFlag(f *flag.FlagSet, s *string) {
	f.StringVarP(s, "config", "f", defaultConfigFile(),
		"specify config file of your credentials")
}

func defaultConfigFile() string {
	const f = "config.json"
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "openpitrix", f)
	}
	return filepath.Join(guessUnixHomeDir(), ".openpitrix", f)
}

func guessUnixHomeDir() string {
	// Prefer $HOME over user.Current due to glibc bug: golang.org/issue/13470
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	// Else, fall back to user.Current:
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return ""
}
