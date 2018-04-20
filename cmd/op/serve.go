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
)

type serveCmd struct {
	out      io.Writer
	url      string
	address  string
	repoPath string
}

func newServeCmd(out io.Writer) *cobra.Command {
	srv := &serveCmd{out: out}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start a local http web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return srv.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&srv.repoPath, "repo-path", "", "local directory path from which to serve apps")
	f.StringVar(&srv.address, "address", devkit.DefaultServeAddr, "address to listen on")
	f.StringVar(&srv.url, "url", "", "external URL of app repository")

	return cmd
}

func (s *serveCmd) run() error {
	repoPath, err := filepath.Abs(s.repoPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return err
	}

	fmt.Fprintln(s.out, "Regenerating index. This may take a moment.")
	if len(s.url) > 0 {
		err = index(repoPath, s.url, "")
	} else {
		err = index(repoPath, "http://"+s.address, "")
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(s.out, "Now serving you on [http://%s/]\n", s.address)
	return devkit.StartLocalRepo(repoPath, s.address)
}
