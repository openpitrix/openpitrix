// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func RunInDocker() bool {
	if s, err := pkgGetThisCgroupDir("devices"); err == nil {
		return strings.HasPrefix(s, "/docker/")
	}
	if pkgFileExists("/.dockerenv") || pkgFileExists("/.dockerinit") {
		return true
	}
	return false
}

func UseDockerLinkedEnvironmentVariables() {
	// mysql database
	if v, ok := os.LookupEnv("OPENPITRIX_DB_" + "MYSQL_PORT_3306_TCP_PORT"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_DB_PORT", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_DB_ENV_" + "MYSQL_DATABASE"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_DB_DBNAME", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_DB_ENV_" + "MYSQL_ROOT_PASSWORD"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_DB_ROOTPASSWORD", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_DB_ENV_" + "MYSQL_USER"); ok && v != "" {
		// unused
	}
	if v, ok := os.LookupEnv("OPENPITRIX_DB_ENV_" + "MYSQL_PASSWORD"); ok && v != "" {
		// unused
	}

	// openpitrix service port
	if v, ok := os.LookupEnv("OPENPITRIX_APP_ENV_" + "OPENPITRIX_CONFIG_APP_PORT"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_APP_PORT", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_RUNTIME_ENV_" + "OPENPITRIX_CONFIG_RUNTIME_PORT"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_RUNTIME_PORT", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_CLUSTER_ENV_" + "OPENPITRIX_CONFIG_CLUSTER_PORT"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_CLUSTER_PORT", v)
	}
	if v, ok := os.LookupEnv("OPENPITRIX_REPO_ENV_" + "OPENPITRIX_CONFIG_REPO_PORT"); ok && v != "" {
		os.Setenv("OPENPITRIX_CONFIG_REPO_PORT", v)
	}
}

func pkgFileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

// https://github.com/docker/libcontainer/blob/master/cgroups/utils.go

// Returns the relative path to the cgroup docker is running in.
func pkgGetThisCgroupDir(subsystem string) (string, error) {
	f, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return pkgParseCgroupFile(subsystem, f)
}

func pkgParseCgroupFile(subsystem string, r io.Reader) (string, error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}

		text := s.Text()
		parts := strings.Split(text, ":")

		for _, subs := range strings.Split(parts[1], ",") {
			if subs == subsystem {
				return parts[2], nil
			}
		}
	}

	return "", errors.New("not found")
}
