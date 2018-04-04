// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix drone service
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"openpitrix.io/libconfd"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/drone"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone"
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
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				info, err := client.GetInfo(context.Background(), &pbdrone.Empty{})
				if err != nil {
					logger.Fatal(err)
				}

				fmt.Println(jsonEncode(info))
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list enabled template resource",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				reply, err := client.GetTemplateFiles(context.Background(), &pbdrone.GetTemplateFilesRequest{
					Regexp: c.Args().First(),
				})
				if err != nil {
					logger.Fatal(err)
				}

				for _, file := range reply.Files {
					fmt.Println(file)
				}
				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from backend by keys",
			ArgsUsage: "key",

			Action: func(c *cli.Context) {
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				reply, err := client.GetValues(context.Background(), &pbdrone.GetValuesRequest{
					Keys: c.Args(),
				})
				if err != nil {
					logger.Fatal(err)
				}

				var maxLen = 1
				for _, k := range c.Args() {
					if len(k) > maxLen {
						maxLen = len(k)
					}
				}

				for _, k := range c.Args() {
					fmt.Printf("%-*s => %s\n", maxLen, k, reply.Values[k])
				}

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
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				cfg := libconfd.MustLoadConfig(c.String("config"))
				bcfg := libconfd.MustLoadBackendConfig(c.String("backend-config"))

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				_, err = client.StartConfd(context.Background(), &pbdrone.StartConfdRequest{
					ConfdConfig:        drone.To_pbdrone_ConfdConfig(cfg),
					ConfdBackendConfig: drone.To_pbdrone_ConfdBackendConfig(bcfg),
				})
				if err != nil {
					logger.Fatal(err)
				}

				fmt.Println("Done")
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Action: func(c *cli.Context) {
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				_, err = client.StopConfd(context.Background(), &pbdrone.Empty{})
				if err != nil {
					logger.Fatal(err)
				}

				fmt.Println("Done")
			},
		},
		{
			Name:  "confd-status",
			Usage: "get confd status",
			Action: func(c *cli.Context) {
				host := c.GlobalString("host")
				port := c.GlobalInt("port")

				client, conn, err := drone.DialDroneService(context.Background(), host, port)
				if err != nil {
					logger.Fatal(err)
				}
				defer conn.Close()

				reply, err := client.GetConfdStatus(context.Background(), &pbdrone.Empty{})
				if err != nil {
					logger.Fatal(err)
				}
				fmt.Println(jsonEncode(reply))
			},
		},

		{
			Name:  "serve",
			Usage: "run as drone service",
			Action: func(c *cli.Context) {
				id := drone.MakeDroneId(c.GlobalString("id-suffix"))
				port := c.GlobalInt("port")

				drone.Serve(
					drone.NewDefaultOptions(),
					drone.WithDrondId(id),
					drone.WithListenPort(port),
				)
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
