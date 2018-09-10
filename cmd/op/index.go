// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

type repoIndexCmd struct {
	dir   string
	url   string
	out   io.Writer
	merge string
}

func newIndexCmd(out io.Writer) *cobra.Command {
	index := &repoIndexCmd{out: out}

	cmd := &cobra.Command{
		Use:   "index [flags] [DIR]",
		Short: "generate an index file given a directory containing packaged app",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgsLength(len(args), "path to a directory"); err != nil {
				return err
			}

			index.dir = args[0]

			return index.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&index.url, "url", "", "url of app repo")
	f.StringVar(&index.merge, "merge", "", "merge the generated index into the given index")

	return cmd
}

func (i *repoIndexCmd) run() error {
	path, err := filepath.Abs(i.dir)
	if err != nil {
		return err
	}

	return index(path, i.url, i.merge)
}

func index(dir, url, mergeTo string) error {
	out := filepath.Join(dir, "index.yaml")

	i, err := devkit.IndexDirectory(dir, url)
	if err != nil {
		return err
	}
	if mergeTo != "" {
		// if index.yaml is missing then create an empty one to merge into
		var i2 *opapp.IndexFile
		if _, err := os.Stat(mergeTo); os.IsNotExist(err) {
			i2 = opapp.NewIndexFile()
			i2.WriteFile(mergeTo, 0755)
		} else {
			i2, err = opapp.LoadIndexFile(mergeTo)
			if err != nil {
				return fmt.Errorf("merge failed: %s", err)
			}
		}
		i.Merge(i2)
	}
	i.SortEntries()
	return i.WriteFile(out, 0755)
}
