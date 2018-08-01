// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/config"
)

func getValidateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "validate_global_config [FILE]",
		Short: "Validate global config",
		Long:  "Validate global config from stdin or [FILE]",
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte
			var err error
			if len(args) == 0 {
				data, err = ioutil.ReadAll(os.Stdin)
			} else {
				data, err = ioutil.ReadFile(args[0])
			}
			if err != nil {
				return err
			}
			_, err = config.ParseGlobalConfig(data)
			return err
		},
	}

	//f := cmd.Flags()

	return cmd
}
