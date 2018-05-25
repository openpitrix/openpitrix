// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
)

var (
	flagConfigFile  = flag.String("config", "config.json", "set config file")
	flagCheckConfig = flag.Bool("check-config", false, "check config file")
	flagCheckServer = flag.Bool("check-server", false, "check json response")
)

type Config struct {
	ListenPort int    `json:"listen_port"`
	Key        string `json:"key"`

	Msg0 string `json:"msg0"`
	Msg1 string `json:"msg1"`
	Msg2 string `json:"msg2"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	if *flagCheckConfig {
		checkConfigFile(*flagConfigFile)
		return
	}

	cfg := MustLoadConfig(*flagConfigFile)
	if *flagCheckServer {
		checkJsonResponse(cfg)
		return
	}

	addr := fmt.Sprintf(":%d", cfg.ListenPort)
	fmt.Printf("Please visit http://localhost:%d/\n", cfg.ListenPort)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s\n", JsonEncode(cfg))
		log.Printf("config: %s\n", JsonEncode(cfg))
	})
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func MustLoadConfig(path string) *Config {
	var p Config
	if err := JsonLoad(path, &p); err != nil {
		log.Fatalf("MustLoadConfig: JsonLoad: %v\n", err)
	}
	if err := IsValidConfig(&p); err != nil {
		log.Fatalf("MustLoadConfig: invalid config: %v\n", err)
	}
	return &p
}

func IsValidConfig(cfg *Config) error {
	if cfg.ListenPort <= 0 {
		return fmt.Errorf("cfg: invalid port: %d", cfg.ListenPort)
	}
	return nil
}

func JsonLoad(filename string, m interface{}) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		return
	}
	return
}

func JsonEncode(m interface{}) []byte {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return nil
	}
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return data
}

func checkConfigFile(path string) {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return
	}
	var p Config
	if err := JsonLoad(path, &p); err != nil {
		log.Fatalf("checkConfigFile: JsonLoad: %v\n", err)
	}
	if err := IsValidConfig(&p); err != nil {
		log.Fatalf("checkConfigFile: invalid config: %v\n", err)
	}

	fmt.Println("OK")
}

func checkJsonResponse(cfg *Config) {
	var got Config

	addr := fmt.Sprintf("http://%s:%d", getLocalIP(), cfg.ListenPort)
	if err := getJsonByURL(addr, &got); err != nil {
		log.Fatalf("checkJsonResponse: getJsonByURL: %v\n", err)
	}

	if !reflect.DeepEqual(cfg, &got) {
		log.Printf("checkJsonResponse: expect = %v\n", cfg)
		log.Printf("checkJsonResponse: got = %v\n", &got)
		log.Fatal()
	}

	fmt.Println("OK")
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func getJsonByURL(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
