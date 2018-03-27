package test_config

import (
	"openpitrix.io/openpitrix/pkg/config"
)

func NewTestDbConfig(database string) config.MysqlConfig {
	return config.MysqlConfig{
		Host:     "localhost",
		Port:     "13306",
		User:     "root",
		Password: "password",
		Database: database,
	}
}

func NewTestEtcdEndpoints() []string {
	return []string{"localhost:12379"}
}
