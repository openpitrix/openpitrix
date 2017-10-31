// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// go test -enable-db-test=true

// start mysql in docker for test

package apps_test

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/config-v2"
	"openpitrix.io/openpitrix/pkg/handlers/apps"
)

var (
	tEnableDbTest = flag.Bool("enable-db-test", false, "enable database test")
	tConfigFile   = flag.String("config-file", "~/.openpitrix/config.yaml", "set config file path")
	tConfig       *config.Config
)

func TestMain(m *testing.M) {
	flag.Parse()
	tConfig, _ = config.Load(*tConfigFile)

	if *tEnableDbTest {
		// TODO(chai): start mysql

		// start service
		go apps.ListenAndServeAppsServer(tConfig)
		time.Sleep(time.Second * 3)
	}

	rv := m.Run()
	os.Exit(rv)
}

func tAssert(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Fatal("Assert failed")
	}
}

func tAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}
func TestAppDatabase(t *testing.T) {
	if !*tEnableDbTest {
		t.Skip()
	}
	tAssertf(t, tConfig != nil, "tConfig is nil")

	db, err := apps.OpenAppDatabase(tConfig)
	tAssertf(t, err == nil, "err: %v", err)
	db.TruncateTables()

	client := apps.NewAppClient(tConfig.Host, tConfig.Port)

	// apps is empty
	{
		items, err := client.GetApps(nil, 0, 0)
		tAssertf(t, err == nil, "err: %v", err)
		tAssert(t, len(items) == 0)
	}

	// create 1 app
	{
		item0 := apps.AppsItem{
			Id:          "app-12345678",
			Name:        "name-1",
			Description: "desc-1",
			Url:         "app1-url",
			Created:     time.Unix(time.Now().Unix(), 0),
		}

		err := client.CreateApp(nil, item0)
		tAssertf(t, err == nil, "err: %v", err)

		items, err := client.GetApps(nil, 0, 0)
		tAssertf(t, err == nil, "err: %v", err)
		tAssert(t, len(items) == 1)

		tAssert(t, items[0].Id == item0.Id)
		tAssert(t, items[0].Name == item0.Name)
		tAssert(t, items[0].Description == item0.Description)

		//tAssertf(t, reflect.DeepEqual(items[0], item0), "%v!= %v", items[0], item0)
	}

	// remove app
	{
		err := client.DeleteApp(nil, "app-12345678")
		tAssertf(t, err == nil, "err: %v", err)

		// now apps is empty
		items, err := client.GetApps(nil, 0, 0)
		tAssertf(t, err == nil, "err: %v", err)
		tAssert(t, len(items) == 0)
	}

	// remove a non exist app
	{
		err := client.DeleteApp(nil, "app-12345678")
		tAssert(t, err != nil) // must return a error
	}

	_ = client
}
