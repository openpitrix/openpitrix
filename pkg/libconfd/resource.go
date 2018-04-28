// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type _TemplateResourceConfig struct {
	TemplateResource TemplateResource `toml:"template"`
}

// TemplateResource is the representation of a parsed template resource.
type TemplateResource struct {
	Src           string      `toml:"src" json:"src"`
	Dest          string      `toml:"dest" json:"dest"`
	Prefix        string      `toml:"prefix" json:"prefix"`
	Keys          []string    `toml:"keys" json:"keys"`
	Mode          string      `toml:"mode" json:"mode"`
	Gid           int         `toml:"gid" json:"gid"`
	Uid           int         `toml:"uid" json:"uid"`
	CheckCmd      string      `toml:"check_cmd" json:"check_cmd"`
	ReloadCmd     string      `toml:"reload_cmd" json:"reload_cmd"`
	FileMode      os.FileMode `toml:"file_mode" json:"file_mode"`
	PGPPrivateKey []byte      `toml:"pgp_private_key" json:"pgp_private_key"`
}

var _LIBCONFD_GOOS = func() string {
	if s := os.Getenv("LIBCONFD_GOOS"); s != "" {
		return s
	}
	if s := os.Getenv("GOOS"); s != "" {
		return s
	}
	return runtime.GOOS
}()

func isTemplateResourceFileShouldBeBuilt(abspath string) bool {
	basename := filepath.Base(abspath)

	if strings.HasPrefix(basename, "_") {
		return false
	}
	if strings.HasPrefix(basename, ".") {
		return false
	}
	if !strings.HasSuffix(basename, ".toml") {
		return false
	}

	switch {
	case strings.HasSuffix(basename, ".darwin.toml"):
		if _LIBCONFD_GOOS != "darwin" {
			return false
		}
	case strings.HasSuffix(basename, ".linux.toml"):
		if _LIBCONFD_GOOS != "linux" {
			return false
		}
	case strings.HasSuffix(basename, ".windows.toml"):
		if _LIBCONFD_GOOS != "windows" {
			return false
		}
	}

	if data, err := ioutil.ReadFile(abspath); err == nil {
		if bytes.Contains(data, []byte("# +build ignore")) {
			return false
		}
		if bytes.Contains(data, []byte("# +build !"+_LIBCONFD_GOOS)) {
			return false
		}
		if bytes.Contains(data, []byte("# +build "+_LIBCONFD_GOOS)) {
			return true
		}
	}
	return true
}

func ListTemplateResource(confdir string) ([]*TemplateResource, []string, error) {
	if !dirExists(confdir) {
		return nil, nil, fmt.Errorf("confdir '%s' does not exist", confdir)
	}

	globpaths, err := filepath.Glob(filepath.Join(confdir, "*.toml"))
	if err != nil {
		return nil, nil, err
	}

	var paths []string
	for _, s := range globpaths {
		if isTemplateResourceFileShouldBeBuilt(s) {
			paths = append(paths, s)
		}
	}

	var lastError error
	var tcs = make([]*TemplateResource, len(paths))

	for i, s := range paths {
		tcs[i], err = LoadTemplateResourceFile(confdir, s)
		if err != nil {
			lastError = err
		}
	}
	if lastError != nil {
		return tcs, paths, lastError
	}

	return tcs, paths, nil
}

func LoadTemplateResourceFile(confdir, name string) (*TemplateResource, error) {
	if !filepath.IsAbs(name) {
		name = filepath.Join(confdir, "conf.d", name)
	}

	p := &_TemplateResourceConfig{
		TemplateResource: TemplateResource{
			Gid: -1,
			Uid: -1,
		},
	}

	_, err := toml.DecodeFile(name, p)
	if err != nil {
		return nil, err
	}

	return &p.TemplateResource, nil
}

func (p *TemplateResource) TomlString() string {
	q := _TemplateResourceConfig{
		TemplateResource: *p,
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(q); err != nil {
		GetLogger().Panic(err)
	}
	return buf.String()
}

func (p *TemplateResource) SaveFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(p.TomlString())
	if err != nil {
		return err
	}

	return nil
}

func (p *TemplateResource) getAbsKeys() []string {
	s := make([]string, len(p.Keys))
	for i, k := range p.Keys {
		s[i] = path.Join(p.Prefix, k)
	}
	return s
}
