// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"openpitrix.io/openpitrix/pkg/util/stringutil"
	"openpitrix.io/openpitrix/test"
)

type Flag struct {
	*flag.FlagSet
}

type Cmd interface {
	GetActionName() string
	ParseFlag(f Flag)
	Run(out Out) error
}

var clientConfig = &test.ClientConfig{}

func init() {
	clientConfig = test.GetClientConfig()
	clientConfig.Debug = false
}

func newRootCmd(c string, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          c,
		Short:        "OpenPitrix cli tool",
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	cmd.AddCommand(getCobraCmds(AllCmd)...)
	flags.Parse(args)
	return cmd
}

func getCobraCmds(cmds []Cmd) (cobraCmds []*cobra.Command) {
	for _, cmd := range cmds {
		action := cmd.GetActionName()
		underscoreAction := stringutil.CamelCaseToUnderscore(action)
		run := cmd.Run
		c := &cobra.Command{
			Use:   fmt.Sprintf("%s [flags]", underscoreAction),
			Short: strings.Replace(underscoreAction, "_", " ", 0),
			RunE: func(c *cobra.Command, args []string) error {
				return run(Out{
					action: action,
					out:    c.OutOrStdout(),
				})
			},
		}
		f := c.Flags()
		cmd.ParseFlag(Flag{f})

		cobraCmds = append(cobraCmds, c)
	}
	return
}

func main() {
	cmd := newRootCmd(os.Args[0], os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
