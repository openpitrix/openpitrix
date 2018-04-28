// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Replace word.
//
// Example:
//	go run replace.go -dir=. -ext=.go -old=old -new=new
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	flagSourceDir = flag.String("dir", ".", "Set dir.")
	flagOldWord   = flag.String("old", "", "Set old word.")
	flagNewWord   = flag.String("new", "", "Set new word.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage of %s: -old=old -new=new\n",
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	total := 0
	filepath.Walk(*flagSourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("filepath.Walk: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if fixWords(path) {
			fmt.Printf("fix %s\n", path)
			total++
		}
		return nil
	})

	fmt.Printf("total %d\n", total)
}

func fixWords(path string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("ioutil.ReadFile(%s): %v", path, err)
	}
	if !bytes.Contains(data, []byte(*flagOldWord)) {
		return false
	}
	data = bytes.Replace(data, []byte(*flagOldWord), []byte(*flagNewWord), -1)
	if err = ioutil.WriteFile(path, data, 0666); err != nil {
		log.Fatalf("ioutil.WriteFile(%s): %v", path, err)
	}
	return true
}
