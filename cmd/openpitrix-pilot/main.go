// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix pilot service
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
)

func main() {
	app := cli.NewApp()
	app.Name = "openpitrix-pilot"
	app.Usage = "openpitrix-pilot provides pilot service."
	app.Version = "0.1.0"

	app.UsageText = `openpitrix-pilot [global options] command [options] [args...]

EXAMPLE:
   openpitrix-pilot info
   openpitrix-pilot list
   openpitrix-pilot getv key
   openpitrix-pilot confd-start
   openpitrix-pilot serve
   openpitrix-pilot tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "id-suffix",
			Value:  "", // pilot@ip/suffix
			Usage:  "pilot id suffix",
			EnvVar: "OPENPITRIX_PILOT_ID_SUFFIX",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  constants.PilotManagerPort,
			Usage:  "pilot ip port",
			EnvVar: "OPENPITRIX_PILOT_PORT",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "pilot log level (debug/info/warning/error/fatal/panic)",
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
			Usage: "show pilot service info",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list frontgate nodes",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from frontgate etcd by keys",
			ArgsUsage: "key",

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
					EnvVar: "pilot_CONFD_CONFILE_FILE",
				},
				cli.StringFlag{
					Name:   "backend-config",
					Value:  "confd-backend.toml",
					Usage:  "confd backend config file",
					EnvVar: "pilot_CONFD_BACKEND_CONFILE_FILE",
				},
			},

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:  "confd-status",
			Usage: "get confd status",
			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},

		{
			Name:  "serve",
			Usage: "run as pilot service",
			Action: func(c *cli.Context) {
				fmt.Println("TODO")
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
openpitrix-pilot info
openpitrix-pilot list

openpitrix-pilot getv /
openpitrix-pilot getv /key
openpitrix-pilot getv / /key

openpitrix-pilot confd-start
openpitrix-pilot confd-stop
openpitrix-pilot confd-status

openpitrix-pilot serve

GOOS=windows openpitrix-pilot list
LIBCONFD_GOOS=windows openpitrix-pilot list
`
