// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/config/test_config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

var p = &Server{&pi.Pi{}}
var tc = test_config.NewDbTestConfig("runtime")

func init() {
	p.Db = tc.GetDatabaseConn()
}

func TestServer_insertRuntimeEnvLabels_byCount(t *testing.T) {
	tc.CheckDbUnitTest(t)
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
		t.Fatalf("error runtime_env label count, should be 0")
	}
	err = p.insertRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId, nil)
	if err != nil {
		t.Fatal(err)
	}

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("runtime_env label count shold be 0")
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

func TestServer_getRuntimeEnvLabelsByEnvId(t *testing.T) {

	runtimeEnvLabels, err := p.getRuntimeEnvLabelsByEnvId()
	if err != nil {
		t.Fatal(err)
	}
	if len(runtimeEnvLabels) != 0 {
		t.Fatal("runtime_env label count shold be 0")
	}

	testRuntimeEnv1 := models.NewRuntimeEnv("test1", "test1", "http://openpitrix.io", "system")
	testRuntimeEnv2 := models.NewRuntimeEnv("test2", "test2", "http://openpitrix.io", "system")
	_, err = p.Db.InsertInto(models.RuntimeEnvTableName).
		Columns(models.RuntimeEnvColumns...).
		Record(testRuntimeEnv1).
		Record(testRuntimeEnv2).
		Exec()
	if err != nil {
		t.Fatal(err)
	}

	runtimeEnvTestMap1 := map[string]string{
		"zone": "pek3a",
		"team": "app",
		"hh":   "hh",
	}

	err = p.insertRuntimeEnvLabels(testRuntimeEnv1.RuntimeEnvId, runtimeEnvTestMap1)
	if err != nil {
		t.Fatal(err)
	}

	err = p.insertRuntimeEnvLabels(testRuntimeEnv2.RuntimeEnvId, runtimeEnvTestMap1)
	if err != nil {
		t.Fatal(err)
	}

	runtimeEnvLabels, err = p.getRuntimeEnvLabelsByEnvId(testRuntimeEnv1.RuntimeEnvId)
	if err != nil {
		t.Fatal(err)
	}
	if len(runtimeEnvLabels) != 3 {
		t.Fatal("runtime_env label count shold be 3")
	}
	for _, runtimeEnvLabel := range runtimeEnvLabels {
		if runtimeEnvLabel.RuntimeEnvId != testRuntimeEnv1.RuntimeEnvId {
			t.Fatalf("labels' runtime env id should be %+v", testRuntimeEnv1.RuntimeEnvId)
		}
		if _, ok := runtimeEnvTestMap1[runtimeEnvLabel.LabelKey]; !ok {
			t.Fatalf("faild to find label [%+v] in [%+v]", runtimeEnvLabels, testRuntimeEnv1)
		}
		if runtimeEnvTestMap1[runtimeEnvLabel.LabelKey] != runtimeEnvLabel.LabelValue {
			t.Fatalf("label [%+v] error,value should be [%+v]", runtimeEnvLabel, runtimeEnvTestMap1[runtimeEnvLabel.LabelKey])
		}
	}

	runtimeEnvLabels, err = p.getRuntimeEnvLabelsByEnvId(testRuntimeEnv1.RuntimeEnvId, testRuntimeEnv2.RuntimeEnvId)
	if err != nil {
		t.Fatal(runtimeEnvLabels)
	}
	if len(runtimeEnvLabels) != 6 {
		t.Fatal("runtime_env label count shold be 6")
	}
	for _, runtimeEnvLabel := range runtimeEnvLabels {
		if runtimeEnvLabel.RuntimeEnvId != testRuntimeEnv1.RuntimeEnvId &&
			runtimeEnvLabel.RuntimeEnvId != testRuntimeEnv2.RuntimeEnvId {
			t.Fatalf("labels' runtime env id should be %+v or %+v", testRuntimeEnv1.RuntimeEnvId, testRuntimeEnv2.RuntimeEnvId)
		}
		if _, ok := runtimeEnvTestMap1[runtimeEnvLabel.LabelKey]; !ok {
			t.Fatalf("faild to find label [%+v] in [%+v]", runtimeEnvLabels, testRuntimeEnv1)
		}
		if runtimeEnvTestMap1[runtimeEnvLabel.LabelKey] != runtimeEnvLabel.LabelValue {
			t.Fatalf("label [%+v] error,value should be [%+v]", runtimeEnvLabel, runtimeEnvTestMap1[runtimeEnvLabel.LabelKey])
		}
	}

	_, err = p.Db.DeleteFrom(models.RuntimeEnvLabelTableName).
		Where(db.Or(
			db.Eq(RuntimeEnvIdColumn, testRuntimeEnv1.RuntimeEnvId),
			db.Eq(RuntimeEnvIdColumn, testRuntimeEnv2.RuntimeEnvId))).Exec()
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Db.DeleteFrom(models.RuntimeEnvTableName).
		Where(db.Or(
			db.Eq(RuntimeEnvIdColumn, testRuntimeEnv1.RuntimeEnvId),
			db.Eq(RuntimeEnvIdColumn, testRuntimeEnv2.RuntimeEnvId))).Exec()
	if err != nil {
		t.Fatal(err)
	}
}

func TestServer_deleteRuntimeEnvLabels_byCount(t *testing.T) {
	testRuntimeEnv := models.NewRuntimeEnv("test1", "test1", "http://openpitrix.io", "system")
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
		t.Fatalf("error runtime_env label count, should be 0")
	}

	err = p.deleteRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId, nil)
	if err != nil {
		t.Fatal(err)
	}

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("error runtime_env label count, should be 0")
	}

	runtimeEnvTestMap := map[string]string{
		"zone": "pek3a",
		"team": "app",
		"hh":   "hh",
		"test": "test",
	}
	err = p.insertRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId, runtimeEnvTestMap)
	if err != nil {
		t.Fatal(err)
	}

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("error runtime_env label count, should be 4")
	}

	err = p.deleteRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"zone": "pek3a",
		})
	if err != nil {
		t.Fatal(err)
	}

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

	err = p.deleteRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"team": "appp",
		})
	if err != nil {
		t.Fatal(err)
	}

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

	err = p.deleteRuntimeEnvLabels(testRuntimeEnv.RuntimeEnvId,
		map[string]string{
			"team": "app",
			"hh":   "hh",
		})
	if err != nil {
		t.Fatal(err)
	}

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("error runtime_env label count, should be 1")
	}

	_, err = p.Db.DeleteFrom(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Exec()
	if err != nil {
		t.Fatal(err)
	}

	count, err = p.Db.Select("*").
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("error runtime_env label count, should be 0")
	}

	_, err = p.Db.DeleteFrom(models.RuntimeEnvTableName).
		Where(db.Eq(RuntimeEnvIdColumn, testRuntimeEnv.RuntimeEnvId)).
		Exec()
	if err != nil {
		t.Fatal(err)
	}

}
