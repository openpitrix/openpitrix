// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"testing"
)

type Docker struct {
	Name       string
	Image      string
	Port       int
	Volume     map[string]string
	WorkDir    string
	Entrypoint string
	t          *testing.T
}

func NewDocker(t *testing.T, name, image string) *Docker {
	return &Docker{
		Name:   name,
		Image:  image,
		t:      t,
		Volume: make(map[string]string),
	}
}

func (d *Docker) Setup() string {
	d.t.Log(d.Teardown())

	s := fmt.Sprintf("docker run -i -d --name='%s' --entrypoint='%s'", d.Name, d.Entrypoint)
	for src, dst := range d.Volume {
		// " -v src:dst"
		s = fmt.Sprintf("%s -v %s:%s", s, src, dst)
	}
	if len(d.WorkDir) > 0 {
		s = fmt.Sprintf("%s -w %s", s, d.WorkDir)
	}
	if d.Port > 0 {
		s = fmt.Sprintf("%s -p %d:%d", s, d.Port, d.Port)
	}
	s = fmt.Sprintf("%s %s sh", s, d.Image)
	return ExecCmd(d.t, s)
}

func (d *Docker) Exec(cmd string) string {
	return ExecCmd(d.t, fmt.Sprintf("docker exec -i %s %s", d.Name, cmd))
}

func (d *Docker) ExecD(cmd string) string {
	return ExecCmd(d.t, fmt.Sprintf("docker exec -i -d %s %s", d.Name, cmd))
}

func (d *Docker) Teardown() string {
	return ExecCmd(d.t, fmt.Sprintf("docker rm -f %s || true", d.Name))
}
