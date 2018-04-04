// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package libconfd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
)

type Application struct {
	cfg    *Config
	client BackendClient
}

func NewApplication(cfg *Config, client BackendClient) *Application {
	return &Application{
		cfg:    cfg.Clone(),
		client: client,
	}
}

func (p *Application) List(re string) {
	_, paths, err := ListTemplateResource(p.cfg.ConfDir)
	if err != nil {
		logger.Fatal(err)
	}
	for _, s := range paths {
		basename := filepath.Base(s)
		if re == "" {
			fmt.Println(basename)
			continue
		}
		matched, err := regexp.MatchString(re, basename)
		if err != nil {
			logger.Fatal(err)
		}
		if matched {
			fmt.Println(basename)
		}
	}
}

func (p *Application) Info(names ...string) {
	if len(names) == 0 {
		_, paths, err := ListTemplateResource(p.cfg.ConfDir)
		if err != nil {
			logger.Fatal(err)
		}
		names = paths
	}
	for _, name := range names {
		if !strings.HasSuffix(name, ".toml") {
			name += ".toml"
		}
		tc, err := LoadTemplateResourceFile(p.cfg.ConfDir, name)
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(tc.TomlString())
	}
}

func (p *Application) Make(names ...string) {
	if len(names) == 0 {
		_, paths, err := ListTemplateResource(p.cfg.ConfDir)
		if err != nil {
			logger.Fatal(err)
		}
		names = paths
	}
	for _, name := range names {
		if !strings.HasSuffix(name, ".toml") {
			name += ".toml"
		}

		fmt.Print(filepath.Base(name), " ")

		tc, err := LoadTemplateResourceFile(p.cfg.ConfDir, name)
		if err != nil {
			logger.Fatal(err)
		}

		tcp := NewTemplateResourceProcessor(name, p.cfg, p.client, tc)

		cfg := p.cfg.Clone()
		cfg.Noop = true

		err = tcp.Process(&Call{Config: cfg, Client: p.client})
		if err != nil {
			logger.Fatal(err)
		}

		fmt.Println("done")
	}
}

func (p *Application) GetValues(keys ...string) {
	m, err := p.client.GetValues(keys)
	if err != nil {
		logger.Fatal(err)
	}

	var maxLen = 1
	for i := range keys {
		if len(keys[i]) > maxLen {
			maxLen = len(keys[i])
		}
	}

	for _, k := range keys {
		fmt.Printf("%-*s => %s\n", maxLen, k, m[k])
	}
}

func (p *Application) Run(opts ...Options) {
	service := NewProcessor()

	go func() {
		defer service.Close()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c

		fmt.Println("quit")
	}()

	service.Run(p.cfg, p.client, opts...)
}
