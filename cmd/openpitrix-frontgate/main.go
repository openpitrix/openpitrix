// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix frontgate service
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/service/metadata/frontgate"
)

func main() {
	app := cli.NewApp()
	app.Name = "openpitrix-frontgate"
	app.Usage = "openpitrix-frontgate provides frontgate service."
	app.Version = "0.0.0"

	app.UsageText = `openpitrix-frontgate [global options] command [options] [args...]

EXAMPLE:
   openpitrix-frontgate info
   openpitrix-frontgate list
   openpitrix-frontgate getv key
   openpitrix-frontgate setv key [value | @file]
   openpitrix-frontgate confd-start
   openpitrix-frontgate serve
   openpitrix-frontgate tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "id-suffix",
			Value:  "", // frontgate@ip/suffix
			Usage:  "frontgate id suffix",
			EnvVar: "OPENPITRIX_FRONTGATE_ID_SUFFIX",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  constants.FrontgateServicePort,
			Usage:  "frontgate host port",
			EnvVar: "OPENPITRIX_FRONTGATE_HOST",
		},
		cli.StringFlag{
			Name:   "pilot-host",
			Value:  constants.PilotManagerHost,
			Usage:  "pilot host ip",
			EnvVar: "OPENPITRIX_PILOT_HOST",
		},
		cli.IntFlag{
			Name:   "pilot-port",
			Value:  constants.PilotManagerPort,
			Usage:  "pilot host port",
			EnvVar: "OPENPITRIX_PILOT_HOST",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "frontgate log level (debug/info/warning/error/fatal/panic)",
			EnvVar: "OPENPITRIX_LOG_LEVEL",
		},
	}

	app.Before = func(c *cli.Context) error {
		flag.Parse()
		logger.SetLevelByString(c.GlobalString("log-level"))
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "info",
			Usage: "show frontgate service info",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list drone nodes",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from etcd by keys",
			ArgsUsage: "key",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:      "setv",
			Usage:     "set value to etcd",
			ArgsUsage: "key [value | @file]",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},

		{
			Name:  "confd-start",
			Usage: "start confd service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config",
					Value:  "confd.toml",
					Usage:  "confd config file",
					EnvVar: "OPENPITRIX_CONFD_CONFILE_FILE",
				},
				cli.StringFlag{
					Name:   "backend-config",
					Value:  "confd-backend.toml",
					Usage:  "confd backend config file",
					EnvVar: "OPENPITRIX_CONFD_BACKEND_CONFILE_FILE",
				},
			},

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Action: func(c *cli.Context) {
				fmt.Println("TODO")
			},
		},
		{
			Name:  "confd-status",
			Usage: "get confd status",
			Action: func(c *cli.Context) {
				fmt.Println("TODO")
			},
		},

		{
			Name:  "serve",
			Usage: "run as frontgate service",
			Action: func(c *cli.Context) {
				frontgate.Serve(nil)
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
	}

	app.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Fprintf(ctx.App.Writer, "not found '%v'!\n", command)
	}

	app.Run(os.Args)
}

func jsonEncode(m interface{}) string {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return ""
	}
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}

const tourTopic = `
openpitrix-frontgate info
openpitrix-frontgate list

openpitrix-frontgate getv /
openpitrix-frontgate getv /key
openpitrix-frontgate getv / /key

openpitrix-frontgate confd-start
openpitrix-frontgate confd-stop
openpitrix-frontgate confd-status

openpitrix-frontgate serve

GOOS=windows openpitrix-frontgate list
LIBCONFD_GOOS=windows openpitrix-frontgate list
`
