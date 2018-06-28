// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix pilot service
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/pilot"
	"openpitrix.io/openpitrix/pkg/service/metadata/pilot/pilotutil"
	"openpitrix.io/openpitrix/pkg/util/pathutil"
	"openpitrix.io/openpitrix/pkg/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "pilot"
	app.Usage = "pilot provides pilot service."
	app.Version = version.GetVersionString()

	app.UsageText = `pilot [global options] command [options] [args...]

EXAMPLE:
   pilot gen-config
   pilot info
   pilot list
   pilot ping
   pilot exec
   pilot getv key
   pilot confd-info
   pilot confd-start
   pilot serve
   pilot send-task
   pilot tour`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "pilot-config.json",
			Usage:  "pilot config file",
			EnvVar: "OPENPITRIX_PILOT_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "gen-config",
			Usage: "gen default config",

			Action: func(c *cli.Context) {
				fmt.Println(pilot.NewDefaultConfigString())
				return
			},
		},

		{
			Name:  "info",
			Usage: "show pilot service info",

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				info, err := client.GetPilotConfig(context.Background(), &pbtypes.Empty{})
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
			Usage:     "list frontgate nodes",
			ArgsUsage: "[regexp]",

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				list, err := client.GetFrontgateList(context.Background(), &pbtypes.Empty{})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				re := c.Args().First()
				for _, v := range list.GetIdList() {
					if re == "" {
						fmt.Println(v)
						continue
					}
					matched, err := regexp.MatchString(re, v)
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					if matched {
						fmt.Println(v)
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
					Value: "pilot",
					Usage: "set endpoint type (pilot/frontgate/drone)",
				},

				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
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
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				switch s := c.String("endpoint-type"); s {
				case "pilot":
					_, err = client.PingPilot(context.Background(), &pbtypes.Empty{})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				case "frontgate":
					_, err = client.PingFrontgate(context.Background(), &pbtypes.FrontgateId{
						Id: c.String("frontgate-id"),
					})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				case "drone":
					_, err = client.PingDrone(
						context.Background(),
						&pbtypes.DroneEndpoint{
							FrontgateId: c.String("frontgate-id"),
							DroneIp:     c.String("drone-host"),
							DronePort:   int32(c.Int("drone-port")),
						},
					)
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
			Name:  "exec",
			Usage: "exec command on pilot/frontgate/drone service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "endpoint-type",
					Value: "frontgate",
					Usage: "set endpoint type (frontgate/drone)",
				},

				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
				},
				cli.StringFlag{
					Name:  "frontgate-node-id",
					Value: "frontgate-node-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "localhost",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
				cli.IntFlag{
					Name:  "timeout",
					Value: 3,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				switch s := c.String("endpoint-type"); s {
				case "frontgate":
					_, err = client.RunCommandOnFrontgateNode(context.Background(), &pbtypes.RunCommandOnFrontgateRequest{
						Endpoint: &pbtypes.FrontgateEndpoint{
							FrontgateId:     c.String("frontgate-id"),
							FrontgateNodeId: c.String("frontgate-node-id"),
						},
						Command:        strings.Join(c.Args(), " "),
						TimeoutSeconds: int32(c.Int("timeout")),
					})
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					fmt.Println("OK")
					return

				case "drone":
					_, err = client.RunCommandOnDrone(
						context.Background(),
						&pbtypes.RunCommandOnDroneRequest{
							Endpoint: &pbtypes.DroneEndpoint{
								FrontgateId: c.String("frontgate-id"),
								DroneIp:     c.String("drone-host"),
								DronePort:   int32(c.Int("drone-port")),
							},
							Command:        strings.Join(c.Args(), " "),
							TimeoutSeconds: int32(c.Int("timeout")),
						},
					)
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
			Name:  "confd-status",
			Usage: "get confd service status",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.IsConfdRunning(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				if reply.GetValue() {
					fmt.Printf("confd on frontgate(%s)/drone(%s:%d) is running\n",
						c.String("frontgate-id"), c.String("drone-host"), c.Int("drone-port"),
					)
				} else {
					fmt.Printf("confd on frontgate(%s)/drone(%s:%d) not running\n",
						c.String("frontgate-id"), c.String("drone-host"), c.Int("drone-port"),
					)
				}

				return
			},
		},

		{
			Name:  "confd-info",
			Usage: "get confd service config",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				reply, err := client.GetConfdConfig(context.Background(), &pbtypes.ConfdEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				fmt.Println(JSONString(reply))
				return
			},
		},

		{
			Name:  "confd-start",
			Usage: "start confd service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StartConfd(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println("Done")
				return
			},
		},
		{
			Name:  "confd-stop",
			Usage: "stop confd service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "frontgate-id",
					Value: "frontgate-001",
				},
				cli.StringFlag{
					Name:  "drone-host",
					Value: "",
				},
				cli.IntFlag{
					Name:  "drone-port",
					Value: constants.DroneServicePort,
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.StopConfd(context.Background(), &pbtypes.DroneEndpoint{
					FrontgateId: c.String("frontgate-id"),
					DroneIp:     c.String("drone-host"),
					DronePort:   int32(c.Int("drone-port")),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println("Done")
				return
			},
		},

		{
			Name:  "get-cmd-status",
			Usage: "get cmd status",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "task-id",
					Value: "",
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				status, err := client.GetSubtaskStatus(context.Background(), &pbtypes.SubTaskId{
					TaskId: c.String("task-id"),
				})
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println(status.Status)
				return
			},
		},
		{
			Name:  "send-task",
			Usage: "send task to pilot service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "task-file",
					Value: "task.json",
				},
			},

			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				// load task json
				task := func() *pbtypes.SubTaskMessage {
					data, err := ioutil.ReadFile(c.String("task-file"))
					if err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					p := new(pbtypes.SubTaskMessage)
					if err := json.Unmarshal(data, p); err != nil {
						logger.Critical("%+v", err)
						os.Exit(1)
					}
					return p
				}()

				client, conn, err := pilotutil.DialPilotService(
					context.Background(), cfg.Host, int(cfg.ListenPort),
				)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}
				defer conn.Close()

				_, err = client.HandleSubtask(context.Background(), task)
				if err != nil {
					logger.Critical("%+v", err)
					os.Exit(1)
				}

				fmt.Println("Done")
				return
			},
		},

		{
			Name:  "serve",
			Usage: "run as pilot service",
			Action: func(c *cli.Context) {
				cfgpath := pathutil.MakeAbsPath(c.GlobalString("config"))
				cfg := pilotutil.MustLoadPilotConfig(cfgpath)

				pilot.Serve(cfg)
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

func Atoi(s string, defaultValue int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultValue
}

const tourTopic = `
pilot gen-config

pilot info
pilot list

pilot getv /
pilot getv /key
pilot getv / /key

pilot confd-start
pilot confd-stop
pilot confd-status

pilot serve

GOOS=windows pilot list
LIBCONFD_GOOS=windows pilot list
`
