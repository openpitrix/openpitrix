// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/libconfd"
	_ "openpitrix.io/libconfd/backends/etcdv3"
)

func main() {
	app := cli.NewApp()
	app.Name = "miniconfd"
	app.Usage = "miniconfd is simple confd, only support toml/etcd backend."
	app.Version = "0.1.0"

	app.UsageText = `miniconfd [global options] command [options] [args...]

EXAMPLE:
   miniconfd list
   miniconfd info
   miniconfd make target
   miniconfd getv key
   miniconfd tour

   miniconfd run`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "confd.toml",
			Usage:  "miniconfd config file",
			EnvVar: "MINICONFD_CONFILE_FILE",
		},
		cli.StringFlag{
			Name:   "backend-config",
			Value:  "confd-backend.toml",
			Usage:  "miniconfd backend config file",
			EnvVar: "MINICONFD_BACKEND_CONFILE_FILE",
		},
	}

	app.Before = func(context *cli.Context) error {
		flag.Parse()
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "list",
			Usage:     "list enabled template resource",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				cfg := libconfd.MustLoadConfig(c.GlobalString("config"))

				backendConfig := libconfd.MustLoadBackendConfig(c.GlobalString("backend-config"))
				backendClient := libconfd.MustNewBackendClient(backendConfig)

				libconfd.NewApplication(cfg, backendClient).List(c.Args().First())
				return
			},
		},
		{
			Name:      "info",
			Usage:     "show template resource info",
			ArgsUsage: "[name...]",

			Action: func(c *cli.Context) {
				cfg := libconfd.MustLoadConfig(c.GlobalString("config"))

				backendConfig := libconfd.MustLoadBackendConfig(c.GlobalString("backend-config"))
				backendClient := libconfd.MustNewBackendClient(backendConfig)

				libconfd.NewApplication(cfg, backendClient).Info(c.Args()...)
				return
			},
		},

		{
			Name:      "make",
			Usage:     "make template target, not run any command",
			ArgsUsage: "[target...]",

			Action: func(c *cli.Context) {
				cfg := libconfd.MustLoadConfig(c.GlobalString("config"))

				backendConfig := libconfd.MustLoadBackendConfig(c.GlobalString("backend-config"))
				backendClient := libconfd.MustNewBackendClient(backendConfig)

				libconfd.NewApplication(cfg, backendClient).Make(c.Args()...)
				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from backend by keys",
			ArgsUsage: "key",

			Action: func(c *cli.Context) {
				cfg := libconfd.MustLoadConfig(c.GlobalString("config"))

				backendConfig := libconfd.MustLoadBackendConfig(c.GlobalString("backend-config"))
				backendClient := libconfd.MustNewBackendClient(backendConfig)

				libconfd.NewApplication(cfg, backendClient).GetValues(c.Args()...)
				return
			},
		},

		{
			Name:  "tour",
			Usage: "show more examples",
			Action: func(c *cli.Context) {
				fmt.Println(tourTopic)
			},
		},

		{
			Name:  "run",
			Usage: "run confd service",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "once",
					Usage: "run with onetime flag",
				},
				cli.BoolFlag{
					Name:  "noop",
					Usage: "run with noop flag",
				},
				cli.BoolFlag{
					Name:  "watch",
					Usage: "run with watch mode",
				},
			},

			Action: func(c *cli.Context) {
				cfg := libconfd.MustLoadConfig(c.GlobalString("config"))

				backendConfig := libconfd.MustLoadBackendConfig(c.GlobalString("backend-config"))
				backendClient := libconfd.MustNewBackendClient(backendConfig)

				libconfd.NewApplication(cfg, backendClient).Run(
					func(cfg *libconfd.Config) {
						cfg.Onetime = c.Bool("once")
					},
					func(cfg *libconfd.Config) {
						cfg.Noop = c.Bool("noop")
					},
					func(cfg *libconfd.Config) {
						cfg.Watch = c.Bool("watch")
					},
				)
				return
			},
		},
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

const tourTopic = `
miniconfd list
miniconfd info simple

miniconfd make simple
miniconfd make simple.windows

miniconfd getv /
miniconfd getv /key
miniconfd getv / /key

miniconfd run
miniconfd run -once
miniconfd run -noop
miniconfd run -once -noop

GOOS=windows miniconfd list
LIBCONFD_GOOS=windows miniconfd list
`
