// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	. "openpitrix.io/openpitrix/cmd/opctl/common"
)

type Gen struct {
	swagger spec.Swagger
}

func toString(a spec.StringOrArray) string {
	if len(a) == 0 {
		return ""
	}
	return a[0]
}

func getName(str string) string {
	return strings.Replace(str, ".", "_", -1)
}

func (g *Gen) GetCmdFromOperation(op *spec.Operation) Cmd {
	var c = Cmd{}
	c.Action = op.ID
	c.Description = op.Summary
	c.Service = op.Tags[0]
	c.Query = make(map[string]Param)
	c.Body = make(map[string]Param)
	if op.Security != nil && len(op.Security) == 0 {
		c.Insecurity = true
	}
	for _, p := range op.Parameters {
		if p.In == "query" {
			var t string
			if p.Type == "array" {
				t += "[]"
				if p.Items != nil {
					t += p.Items.Type
				}
			} else if p.Format == "int32" {
				t = "int32"
			} else if p.Format == "int64" {
				t = "int64"
			} else {
				t = p.Type
			}
			c.Query[getName(p.Name)] = Param{
				Shorthand: "",
				Help:      p.Description,
				Type:      t,
				Default:   p.Default,
			}
		} else {
			if p.Schema != nil {
				def := strings.Split(p.Schema.Ref.String(), "/")
				defName := def[len(def)-1]
				schema := g.swagger.Definitions[defName]

				for name, s := range schema.Properties {
					var t string
					if s.Type.Contains("array") {
						t += "[]"
						if s.Items != nil && s.Items.Schema != nil {
							t += toString(s.Items.Schema.Type)
						}
					} else if s.Format == "byte" {
						t += "byte"
					} else {
						t += toString(s.Type)
					}
					c.Body[getName(name)] = Param{
						Shorthand: "",
						Help:      s.Description,
						Type:      t,
						Default:   s.Default,
					}
				}
			}
		}
	}
	return c
}

func (g *Gen) Parse(content []byte) {
	err := g.swagger.UnmarshalJSON(content)
	if err != nil {
		Error(err, "unmarshal stdin")
	}
	var cmds Cmds
	for _, path := range g.swagger.Paths.Paths {
		var ops = []*spec.Operation{
			path.Get,
			path.Put,
			path.Post,
			path.Delete,
			path.Options,
			path.Head,
			path.Patch,
		}
		for _, op := range ops {
			if op != nil {
				cmd := g.GetCmdFromOperation(op)
				cmds = append(cmds, cmd)
			}
		}
	}
	sort.Sort(cmds)
	t, err := yaml.Marshal(cmds)
	if err != nil {
		Error(err, "marshal yaml")
	}
	os.Stdout.Write(t)
}

func main() {
	var filePath string
	flag.StringVarP(&filePath, "file", "f", "", "")
	flag.Parse()

	var content []byte
	var err error
	if filePath != "" {
		content, err = ioutil.ReadFile(filePath)
		if err != nil {
			Error(err, fmt.Sprintf("read file [%s]", filePath))
		}
	} else {
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			Error(err, "read stdin")
		}
	}

	g := Gen{
		swagger: spec.Swagger{},
	}
	g.Parse(content)
}
