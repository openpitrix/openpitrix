package test_config

import (
	"flag"
	"fmt"
	"os"

	"github.com/gocraft/dbr"
	"github.com/koding/multiconfig"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
)

var defaultEventReceiver = db.EventReceiver{}

func LoadConf() *OpTestConfig {
	flag.CommandLine.Parse(os.Args[1:])
	config := new(OpTestConfig)
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{DefaultTagName: "default"},
		&multiconfig.EnvironmentLoader{},
	)
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		logger.Panicf("Failed to load testConfig: %+v", err)
		panic(err)
	}
	logger.SetLevelByString(config.Log.Level)
	logger.Debugf("LoadConf: %+v", config)
	return config
}

type DbConfig struct {
	Type     string `default:"mysql"`
	Host     string `default:"127.0.0.1"`
	Port     string `default:"23306"`
	User     string `default:"root"`
	Password string `default:"password"`
	Database string `default:"openpitrix"`
}

type LogConfig struct {
	Level string `default:"debug"` // debug, info, warn, error, fatal
}

type OpTestConfig struct {
	Log    LogConfig
	Db     DbConfig
	DbTest bool
}

func (m *DbConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}

func OpenDatabase(cfg DbConfig) (*db.Database, error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	switch cfg.Type {
	case "mysql":
		conn, err := dbr.Open(cfg.Type, cfg.GetUrl()+"?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", &defaultEventReceiver)
		if err != nil {
			return nil, err
		}
		return &db.Database{Session: conn.NewSession(nil)}, nil
	}
	return nil, fmt.Errorf("unknown database type")
}
