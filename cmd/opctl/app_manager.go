// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"openpitrix.io/openpitrix/test"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/models"
)

type CreateAppCmd struct {
	*models.OpenpitrixCreateAppRequest
}

func NewCreateAppCmd() Cmd {
	return &CreateAppCmd{
		&models.OpenpitrixCreateAppRequest{},
	}
}

func (*CreateAppCmd) GetActionName() string {
	return "CreateApp"
}

func (c *CreateAppCmd) ParseFlag(f Flag) {
	f.StringVarP(&c.ChartName, "chart_name", "C", "", "chart_name")
	f.StringVarP(&c.Description, "description", "d", "", "description")
	f.StringVarP(&c.Home, "home", "H", "", "home")
	f.StringVarP(&c.Icon, "icon", "i", "", "icon")
	f.StringVarP(&c.Maintainers, "maintainers", "m", "", "maintainers")
	f.StringVarP(&c.Name, "name", "n", "", "name")
	f.StringVarP(&c.Owner, "owner", "o", "", "owner")
	f.StringVarP(&c.Readme, "readme", "R", "", "readme")
	f.StringVarP(&c.RepoID, "repo_id", "r", "", "repo_id")
	f.StringVarP(&c.Screenshots, "screenshots", "s", "", "screenshots")
	f.StringVarP(&c.Sources, "sources", "S", "", "sources")
}

func (c *CreateAppCmd) Run(out Out) error {
	params := app_manager.NewCreateAppParams()
	params.WithBody(c.OpenpitrixCreateAppRequest)

	out.WriteRequest(params)

	client := test.GetClient(clientConfig)
	res, err := client.AppManager.CreateApp(params)
	if err != nil {
		return err
	}

	out.WriteResponse(res.Payload)

	return nil
}

type DescribeAppCmd struct {
	*app_manager.DescribeAppsParams
}

func NewDescribeAppCmd() Cmd {
	return &DescribeAppCmd{
		app_manager.NewDescribeAppsParams(),
	}
}

func (*DescribeAppCmd) GetActionName() string {
	return "DescribeApp"
}

func (c *DescribeAppCmd) ParseFlag(f Flag) {
	c.SearchWord = new(string)
	c.Limit = new(int64)
	c.Offset = new(int64)
	f.StringSliceVarP(&c.AppID, "app_id", "a", []string{""}, "")
	f.StringSliceVarP(&c.ChartName, "chart_name", "c", []string{""}, "")
	f.StringSliceVarP(&c.Name, "name", "n", []string{""}, "")
	f.StringSliceVarP(&c.Owner, "owner", "o", []string{""}, "")
	f.StringSliceVarP(&c.Status, "status", "s", []string{""}, "")
	f.StringSliceVarP(&c.RepoID, "repo_id", "r", []string{""}, "")
	f.StringVarP(c.SearchWord, "search_word", "S", "", "")
	f.Int64VarP(c.Limit, "limit", "L", 20, "")
	f.Int64VarP(c.Offset, "offset", "O", 0, "")
}

func (c *DescribeAppCmd) Run(out Out) error {

	out.WriteRequest(c.DescribeAppsParams)

	client := test.GetClient(clientConfig)
	res, err := client.AppManager.DescribeApps(c.DescribeAppsParams)
	if err != nil {
		return err
	}

	out.WriteResponse(res.Payload)

	return nil
}
