// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

)

//server struct
var ss struct{
	server 		*Server
	ctx 		context.Context
	cancel 		context.CancelFunc
}

var ENVs = map[string]map[string]string{
	"etcd": {
		"OPENPITRIX_ETCD_ENDPOINTS": "127.0.0.1:12379",
	},
	"mysql": {
		"OPENPITRIX_MYSQL_HOST":     "127.0.0.1",
		"OPENPITRIX_MYSQL_PORT":     "13306",
		"OPENPITRIX_MYSQL_USER":     "root",
		"OPENPITRIX_MYSQL_PASSWORD": "password",
		"OPENPITRIX_MYSQL_DATABASE": "mbing",
	},
}

func TestMain(m *testing.M){

	//check if enable testing
	if !*tTestingEnvEnabled {
		fmt.Println("testing env disabled")
		os.Exit(1)
	}

	//setup envs
	for _, typeEnvs := range ENVs{
		for k, v := range typeEnvs{
			os.Setenv(k, v)
		}
	}

	//setup server
	InitGlobelSetting()
	ss.server, _ = NewServer()
	ss.ctx, ss.cancel = context.WithTimeout(context.Background(), time.Second)

	//run testing
	flag.Parse()
	exitCode := m.Run()

	//stop testing and exit
	ss.cancel()
	os.Exit(exitCode)
}

func TestNewServer(t *testing.T) {
	succInfo := "TestNewServer Passed, server: %v, context: %v, cancleFunc: %v"
	t.Logf(succInfo, ss.server, ss.ctx, ss.cancel)
}

