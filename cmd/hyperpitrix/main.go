// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix all in one
package main

import (
	"openpitrix.io/openpitrix/pkg/apigateway"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/service/app"
	"openpitrix.io/openpitrix/pkg/service/attachment"
	"openpitrix.io/openpitrix/pkg/service/category"
	"openpitrix.io/openpitrix/pkg/service/cluster"
	"openpitrix.io/openpitrix/pkg/service/isv"
	"openpitrix.io/openpitrix/pkg/service/job"
	"openpitrix.io/openpitrix/pkg/service/repo"
	"openpitrix.io/openpitrix/pkg/service/repo_indexer"
	"openpitrix.io/openpitrix/pkg/service/runtime"
	"openpitrix.io/openpitrix/pkg/service/runtime_provider"
	"openpitrix.io/openpitrix/pkg/service/task"
)

func getConf(database string) *config.Config {
	cfg := config.GetConf()
	cfg.Mysql.Database = database
	return cfg
}

func main() {
	go category.Serve(getConf("app"))
	go cluster.Serve(getConf("cluster"))
	go isv.Serve(getConf("isv"))
	go job.Serve(getConf("job"))
	go repo_indexer.Serve(getConf("repo"))
	go repo.Serve(getConf("repo"))
	go runtime.Serve(getConf("runtime"))
	go task.Serve(getConf("task"))
	go app.Serve(getConf("app"))
	go attachment.Serve(getConf("attachment"))
	go runtime_provider.Serve(getConf(""))

	apigateway.Serve(getConf(""))
}
