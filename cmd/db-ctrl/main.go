// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/db/ctrl"
	"openpitrix.io/openpitrix/pkg/logger"
)

func main() {
	flag := config.GetFlagSet()
	var (
		schemaPath string
		isCleanup  bool
	)
	flag.StringVar(&schemaPath, "schema-path", "", "specify schema path")
	flag.StringVar(&schemaPath, "s", "", "specify schema path")
	flag.BoolVar(&isCleanup, "cleanup", false, "specify to cleanup database")
	flag.BoolVar(&isCleanup, "C", false, "specify to cleanup database")
	config.ParseFlag()
	if len(schemaPath) == 0 {
		currentFile, err := os.Executable()
		if err != nil {
			panic(err)
		}
		currentFilePath := filepath.Dir(currentFile)
		schemaPath = filepath.Join(currentFilePath, "schema")
		log.Printf("Unspecified schema path, using default [%s]", schemaPath)
	} else {
		log.Printf("Specified schema path [%s]", schemaPath)
	}
	fileInfo, err := os.Stat(schemaPath)
	if os.IsNotExist(err) {
		log.Fatalf("The schema path [%s] is not exist", schemaPath)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("The schema path [%s] is not dir", schemaPath)
	}
	cfg := config.LoadConf()
	logger.Disable()
	if isCleanup {
		ctrl.Cleanup(cfg)
	} else {
		ctrl.Start(cfg, schemaPath)
	}
}
