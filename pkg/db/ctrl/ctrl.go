// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ctrl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
)

const CreateDb = "CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;"
const DropDb = "DROP DATABASE %s;"
const CreateDbCtrlTable = "CREATE TABLE %s (id int PRIMARY KEY, version int);"
const DbCtrlTableName = "db_ctrl"

type DbCtrl struct {
	Id      int   `json:"id"`
	Version int64 `json:"version"`
}

func (d *DbCtrl) getVersionStr() string {
	return fmt.Sprintf("%.6d", d.Version)
}

func getDbCtrlVersion(session *db.Database) (dbCtrl *DbCtrl, err error) {
	err = session.Select("*").From(DbCtrlTableName).Where("id = 1").LoadOne(&dbCtrl)
	if err != nil {
		log.Printf("Failed to query from [db_ctrl]: %+v", err)
	}
	return
}

func updateDbCtrlVersion(session *db.Database, version int64) (err error) {
	_, err = session.Update(DbCtrlTableName).Where("id = 1").Set("version", version).Exec()
	return
}

func createDbCtrlRecord(session *db.Database) {
	_, err := session.Exec(fmt.Sprintf(CreateDbCtrlTable, DbCtrlTableName))
	if err != nil {
		log.Fatalf("Failed to create table [db_ctrl]: %+v", err)
	}
	_, err = session.InsertInto(DbCtrlTableName).Columns("id", "version").Values(1, 0).Exec()
	if err != nil {
		log.Fatalf("Failed to insert record into [db_ctrl]: %+v", err)
	}
}

func execFile(session *db.Database, path string) {
	log.Printf("Executing file [%s]", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	querys := strings.Split(string(b), ";")
	for _, query := range querys {
		if len(strings.TrimSpace(query)) == 0 {
			continue
		}
		_, err := session.Exec(query)
		if err != nil {
			log.Printf("Execute file [%s] error: %+v", path, err)
		}
	}
}

func Init(m config.MysqlConfig, schemaPath string) {
	session := openDatabase(m)
	defer session.Close()
	dbCtrl, err := getDbCtrlVersion(session)
	if err == nil {
		log.Printf("Got current db ctrl version [%s], skip init tables", dbCtrl.getVersionStr())
		return
	}
	initFilePath := filepath.Join(schemaPath, "init.sql")
	execFile(session, initFilePath)
	createDbCtrlRecord(session)
}

func Upgrade(m config.MysqlConfig, schemaPath string) {
	session := openDatabase(m)
	defer session.Close()
	upgradeDirPath := filepath.Join(schemaPath, "upgrade")
	files, err := ioutil.ReadDir(upgradeDirPath)
	if err != nil {
		panic(err)
	}
	dbCtrl, err := getDbCtrlVersion(session)
	if err == nil {
		log.Printf("Got current db ctrl version [%s], skip some upgrade tables", dbCtrl.getVersionStr())
	}
	for _, file := range files {
		fileVersion := strings.Split(file.Name(), "_")[0]
		version, err := strconv.ParseInt(fileVersion, 10, 64)
		if err != nil {
			panic(err)
		}
		if version <= dbCtrl.Version {
			continue
		}
		path := filepath.Join(upgradeDirPath, file.Name())
		execFile(session, path)
		err = updateDbCtrlVersion(session, version)
		if err != nil {
			log.Fatalf("Failed to update db ctrl version to [%s]", fileVersion)
		} else {
			log.Printf("Update db ctrl version to [%s]", fileVersion)
		}
	}
}

func openTempDatabase(m config.MysqlConfig) *dbr.Session {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", m.User, m.Password, m.Host, m.Port)
	tempDb, err := dbr.Open("mysql", dsn, nil)
	if err != nil {
		panic(err)
	}
	return tempDb.NewSession(nil)
}

func createDatabase(m config.MysqlConfig) {
	session := openTempDatabase(m)
	_, err := session.Exec(fmt.Sprintf(CreateDb, m.Database))
	if err != nil {
		panic(err)
	}
}

func dropDatabase(m config.MysqlConfig) {
	session := openTempDatabase(m)
	_, err := session.Exec(fmt.Sprintf(DropDb, m.Database))
	if err != nil {
		panic(err)
	}
	log.Printf("Drop database [%s] success", m.Database)
}

func openDatabase(m config.MysqlConfig) *db.Database {
	database, err := db.OpenDatabase(m)
	if err != nil {
		panic(err)
	}
	return database
}

func Cleanup(cfg *config.Config) {
	m := cfg.Mysql
	dropDatabase(m)
}

func Start(cfg *config.Config, schemaPath string) {
	m := cfg.Mysql
	createDatabase(m)
	Init(m, schemaPath)
	Upgrade(m, schemaPath)
	log.Print("Done, start http server")
	HttpServe(m, schemaPath)
}

func HttpServe(m config.MysqlConfig, schemaPath string) {
	session := openDatabase(m)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dbCtrlVersion, err := getDbCtrlVersion(session)
		if err != nil {
			fmt.Fprintf(w, "{\"err\": \"%+v\"}", err)
			return
		}
		b, err := json.Marshal(dbCtrlVersion)
		if err != nil {
			fmt.Fprintf(w, "{\"err\": \"%+v\"}", err)
			return
		}
		s := string(b)
		log.Printf("Current db ctrl version [%s]", s)
		fmt.Fprintf(w, s)
	})
	http.HandleFunc("/upgrade", func(w http.ResponseWriter, r *http.Request) {
		Upgrade(m, schemaPath)
		w.Write([]byte("{\"msg\":\"done\"}"))
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", constants.DbCtrlPort), nil))
}
