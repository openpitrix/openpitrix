// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	completionLong = `
Output shell completion code for the specified shell (bash or zsh).
The shell code must be evaluated to provide interactive
completion of opctl commands.  This can be done by sourcing it from
the ~/.bash_profile or ~/.bashrc

Note for zsh users: zsh completions are only supported in versions of zsh >= 5.2`

	completionExample = `
# Installing bash completion on macOS using homebrew
## If running Bash 3.2 included with macOS
    brew install bash-completion
## or, if running Bash 4.1+
    brew install bash-completion@2
# Load the opctl completion code for bash into the current shell
    source <(opctl completion bash)
# Load the opctl completion code for zsh into the current shell
    source <(opctl completion zsh)
`
)

func getCompletionCmd() *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion code for the specified shell (bash or zsh)",

		Long:    completionLong,
		Example: completionExample,

		ValidArgs:             []string{"bash", "zsh"},
		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("shell not specified")
			}
			if len(args) > 1 {
				return fmt.Errorf("too many arguments, expected only the shell type")
			}
			shell := args[0]
			switch shell {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell type %q", shell)
			}
			return nil
		},
	}
	return completionCmd
}
