// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix drone service
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
	app.Name = "openpitrix-drone"
	app.Usage = "openpitrix-drone provides drone service."
	app.Version = "0.0.0"

	app.UsageText = `openpitrix-drone [global options] command [options] [args...]

EXAMPLE:
   openpitrix-drone info
   openpitrix-drone list
   openpitrix-drone getv key
   openpitrix-drone confd-start
   openpitrix-drone serve
   openpitrix-drone tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "id-suffix",
			Value:  "", // drone@ip/suffix
			Usage:  "drone id suffix",
			EnvVar: "OPENPITRIX_DRONE_ID_SUFFIX",
		},
		cli.StringFlag{
			Name:   "dbpath",
			Value:  "drone.db",
			Usage:  "drone database path",
			EnvVar: "OPENPITRIX_DRONE_DBPATH",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			Usage:  "drone host ip",
			EnvVar: "OPENPITRIX_DRONE_HOST",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  constants.DroneServicePort,
			Usage:  "drone ip port",
			EnvVar: "OPENPITRIX_DRONE_PORT",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "drone log level (debug/info/warning/error/fatal/panic)",
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
			Usage: "show drone service info",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list enabled template resource",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				fmt.Println("TODO")
				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from backend by keys",
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
					EnvVar: "OPENPITRIX_CONFD_CONFIG_FILE",
				},
				cli.StringFlag{
					Name:   "backend-config",
					Value:  "confd-backend.toml",
					Usage:  "confd backend config file",
					EnvVar: "OPENPITRIX_CONFD_BACKEND_CONFIG_FILE",
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
			Usage: "run as drone service",
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
openpitrix-drone info
openpitrix-drone list

openpitrix-drone getv /
openpitrix-drone getv /key
openpitrix-drone getv / /key

openpitrix-drone confd-start
openpitrix-drone confd-stop
openpitrix-drone confd-status

openpitrix-drone serve

GOOS=windows openpitrix-drone list
LIBCONFD_GOOS=windows openpitrix-drone list
`
