// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

type createCmd struct {
	name    string
	out     io.Writer
	starter string
}

func newCreateCmd(out io.Writer) *cobra.Command {
	cc := &createCmd{out: out}

	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "create a new app with the given name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("the name of the new app is required")
			}
			cc.name = args[0]
			return cc.run()
		},
	}

	return cmd
}

func (c *createCmd) run() error {
	fmt.Fprintf(c.out, "Creating [%s]\n", c.name)

	appName := filepath.Base(c.name)
	cfile := &opapp.Metadata{
		Name:        appName,
		Description: "An OpenPitrix app",
		Version:     "0.1.0",
		AppVersion:  "1.0",
		ApiVersion:  devkit.ApiVersionV1,
	}

	_, err := devkit.Create(cfile, filepath.Dir(c.name))
	return err
}
