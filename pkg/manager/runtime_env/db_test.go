package runtime_env

import (
	"fmt"
	"testing"

	"github.com/koding/multiconfig"

	"openpitrix.io/openpitrix/pkg/config/test_config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

var p = &Server{&pi.Pi{}}

var dbSessionSuccess = false
var testConfig *test_config.OpTestConfig

func init() {
	testConfig = test_config.LoadConf()
	if testConfig.DbTest {
		db, err := test_config.OpenDatabase(testConfig.Db)
		if err != nil {
			logger.Fatalf("failed to open database %+v", testConfig.Db)
		}
		err = db.Ping()
		if err != nil {
			logger.Fatalf("failed to ping database %+v", testConfig.Db)
		}
		dbSessionSuccess = true
		p.Db = db
	}
}

func checkDbTest(t *testing.T) {
	if !testConfig.DbTest {
		fmt.Println("run db unit tests by set environment variables:")
		loader := &multiconfig.EnvironmentLoader{}
		loader.PrintEnvs(new(test_config.OpTestConfig))
		t.Skip()
	}
}

func TestInsertRuntimeEnvLabels_byCount(t *testing.T) {
	checkDbTest(t)
	if !dbSessionSuccess {
		t.Fatalf("failed to open database")
	}

	testRuntimeEnv := models.NewRuntimeEnv("test", "test", "http://openpitrix.io", "system")
	_, err := p.Db.InsertInto(models.RuntimeEnvTableName).
		Columns(models.RuntimeEnvColumns...).
		Record(testRuntimeEnv).
		Exec()
	if err != nil {
		t.Fatal(err)
	}
	count, err := p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("error runtime_env label count, should be %v", count)
	}

	p.insertRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"openpitrix": "test",
			"env":        "test"})

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("error runtime_env label count, should be 2")
	}

	p.insertRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"runtime": "qingcloud",
		})
	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()

	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("error runtime_env label count, should be 3")
	}

	p.insertRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"zone": "pek3a",
			"team": "app",
			"hh":   "hh",
		})
	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()

	if err != nil {
		t.Fatal(err)
	}
	if count != 6 {
		t.Fatalf("error runtime_env label count, should be 6")
	}

	p.Db.DeleteFrom(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId))
	if err != nil {
		t.Fatal(err)
	}
}
