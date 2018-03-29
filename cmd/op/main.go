// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newRootCmd(c string, args []string) *cobra.Command {
	// Thanks to Helm Project (https://helm.sh/).
	cmd := &cobra.Command{
		Use:          c,
		Short:        "The devkit for OpenPitrix.",
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()

	cmd.AddCommand(
		// app commands
		newCreateCmd(out),
		//newLintCmd(out),
		newPackageCmd(out),
		newIndexCmd(out),
		newServeCmd(out),
	)
	flags.Parse(args)
	return cmd
}

func main() {
	cmd := newRootCmd(os.Args[0], os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkArgsLength(argsReceived int, requiredArgs ...string) error {
	expectedNum := len(requiredArgs)
	if argsReceived != expectedNum {
		arg := "arguments"
		if expectedNum == 1 {
			arg = "argument"
		}
		return fmt.Errorf("this command needs %v %s: %s", expectedNum, arg, strings.Join(requiredArgs, ", "))
	}
	return nil
}
