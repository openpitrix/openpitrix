// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix drone service
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone/droneutil"
	"openpitrix.io/openpitrix/pkg/util/pathutil"
	"openpitrix.io/openpitrix/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone"
	app.Usage = "drone provides drone service."
	app.Version = version.GetVersionString()

	app.UsageText = `drone [global options] command [options] [args...]

EXAMPLE:
   drone gen-config
   drone info
   drone list
   drone ping
   drone getv key
   drone confd-start
   drone serve
   drone tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "drone-config.json",
			Usage:  "drone config file",
			EnvVar: "OPENPITRIX_DRONE_CONFIG",
		},
		cli.StringFlag{
			Name:   "config-confd",
			Value:  "confd-config.json",
			Usage:  "drone confd file (ignored if missing)",
			EnvVar: "OPENPITRIX_CONFD_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "debug",
			Usage:  "debug app",
			Hidden: true,

			Action: func(c *cli.Context) {
				return
			},
		},

		{
			Name:  "gen-config",
			Usage: "gen default config",

			Action: func(c *cli.Context) {
				fmt.Println(drone.NewDefaultConfigString())
				return
			},
		},

		{
			Name:  "info",
			Usage: "show drone service info",

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				info, err := client.GetDroneConfig(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println(JSONString(info))
				return
			},
		},
		{
			Name:      "list",
			Usage:     "list enabled template resource",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.GetTemplateFiles(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				re := c.Args().First()
				for _, file := range reply.GetValueList() {
					if re == "" {
						fmt.Println(file)
						continue
					}
					matched, err := regexp.MatchString(re, file)
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					if matched {
						fmt.Println(file)
					}
				}
				return
			},
		},
		{
			Name:  "ping",
			Usage: "ping pilot/frontgate/drone service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "endpoint-type",
					Value: "drone",
					Usage: "set endpoint type (pilot/frontgate/drone)",
				},

				cli.StringFlag{
					Name:  "frontgate-host",
					Value: "localhost",
				},
				cli.IntFlag{
					Name:  "frontgate-port",
					Value: constants.FrontgateServicePort,
				},

				cli.StringFlag{
					Name:  "drone-host",
					Value: "localhost",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				switch s := c.String("endpoint-type"); s {
				case "pilot":
					_, err = client.PingPilot(context.Background(), &pbtypes.FrontgateEndpoint{
						NodeIp:   c.String("frontgate-host"),
						NodePort: int32(c.Int("frontgate-port")),
					})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}

					fmt.Println("OK")
					return

				case "frontgate":
					_, err = client.PingFrontgate(context.Background(), &pbtypes.FrontgateEndpoint{
						NodeIp:   c.String("frontgate-host"),
						NodePort: int32(c.Int("frontgate-port")),
					})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}

					fmt.Println("OK")
					return

				case "drone":
					_, err = client.PingDrone(context.Background(), &pbtypes.Empty{})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}

					fmt.Println("OK")
					return

				default:
					logger.Critical("unknown endpoint type: %s\n", s)
					os.Exit(1)
				}

				return
			},
		},

		{
			Name:      "getv",
			Usage:     "get values from backend by keys",
			ArgsUsage: "key",

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.GetValues(context.Background(), &pbtypes.StringList{
					ValueList: c.Args(),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				var maxLen = 1
				for _, k := range c.Args() {
					if len(k) > maxLen {
						maxLen = len(k)
					}
				}

				for _, k := range c.Args() {
					fmt.Printf("%-*s => %s\n", maxLen, k, reply.GetValueMap()[k])
				}

				return
			},
		},

		{
			Name:  "confd-status",
			Usage: "get confd service status",
			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.IsConfdRunning(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				if reply.GetValue() {
					fmt.Println("confd is running")
				} else {
					fmt.Println("confd not running")
				}
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
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StartConfd(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println("Done")
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := droneutil.MustLoadDroneConfig(cfgpath)

				client, conn, err := droneutil.DialDroneService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StopConfd(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println("Done")
			},
		},

		{
			Name:  "serve",
			Usage: "run as drone service",
			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfgConfdPath := pathutil.MakeAbsPath(c.GlobalString("config-confd"))

				cfg := droneutil.MustLoadDroneConfig(cfgpath)
				cfgManager := drone.NewConfigManager(cfgpath, cfg)

				confdServer := drone.NewConfdServer(cfgConfdPath)
				if cfgConfd, _ := droneutil.LoadConfdConfig(cfgConfdPath); cfgConfd != nil {
					cfg := cfgManager.Get()
					fnHookKeyAdjuster := func(absKey string) (realKey string) {
						if absKey == "/"+cfg.ConfdSelfHost || strings.HasPrefix(absKey, "/"+cfg.ConfdSelfHost+"/") {
							return absKey
						}

						if absKey == "/self" {
							return "/" + cfg.ConfdSelfHost
						}
						if strings.HasPrefix(absKey, "/self/") {
							return "/" + cfg.ConfdSelfHost + absKey[len("/self/")-1:]
						}
						return absKey
					}

					confdServer.SetConfig(cfgConfd, fnHookKeyAdjuster)
				}

				drone.Serve(cfgManager, confdServer)
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

func JSONString(m interface{}) string {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return ""
	}
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}

const tourTopic = `
drone gen-config

drone info
drone list
drone ping

drone getv /
drone getv /key
drone getv / /key

drone confd-start
drone confd-stop

drone serve

GOOS=windows drone list
LIBCONFD_GOOS=windows drone list
`
