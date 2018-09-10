// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

type packageCmd struct {
	path        string
	version     string
	destination string

	out io.Writer
}

func newPackageCmd(out io.Writer) *cobra.Command {
	pkg := &packageCmd{out: out}

	cmd := &cobra.Command{
		Use:   "package [flags] [VERSION_PATH] [...]",
		Short: "package an app directory into an app archive",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("need at least one argument, the path to the app")
			}
			for i := 0; i < len(args); i++ {
				pkg.path = args[i]
				if err := pkg.run(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&pkg.version, "version", "", "set the version on the app to this semver version")
	f.StringVarP(&pkg.destination, "destination", "d", ".", "location to write the app.")

	return cmd
}

func (p *packageCmd) run() error {
	path, err := filepath.Abs(p.path)
	if err != nil {
		return err
	}

	ch, err := devkit.LoadDir(path)
	if err != nil {
		return err
	}

	// If version is set, modify the version.
	if len(p.version) != 0 {
		if err := setVersion(ch, p.version); err != nil {
			return err
		}
	}

	if filepath.Base(path) != ch.Metadata.Name {
		return fmt.Errorf("directory name [%s] and package.json name [%s] must match", filepath.Base(path), ch.Metadata.Name)
	}

	var dest string
	if p.destination == "." {
		// Save to the current working directory.
		dest, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		// Otherwise save to set destination
		dest = p.destination
	}

	name, err := devkit.Save(ch, dest)
	if err == nil {
		fmt.Fprintf(p.out, "Successfully packaged app and saved it to: %s\n", name)
	} else {
		return fmt.Errorf("failed to save: %s", err)
	}

	return err
}

func setVersion(version *opapp.OpApp, ver string) error {
	// Verify that version is a SemVer, and error out if it is not.
	if _, err := semver.NewVersion(ver); err != nil {
		return err
	}

	// Set the version field on the app.
	version.Metadata.Version = ver
	return nil
}
