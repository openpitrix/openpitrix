// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	"testing"
)

var ctx = context.TODO()

func TestValidate1(t *testing.T) {
	repoType := "s3"
	url := "s3://s3.pek3a.qingstor.com/op-repo"
	credential := `{"access_key_id": "wiandianiaeudsadf8a33uffhufhud", "secret_access_key": "nduaufbuabfuebaufbaufaueuu"}`
	providers := []string{"qingcloud"}

	err := validate(ctx, repoType, url, credential, providers)

	if err == nil {
		t.Errorf("expect error, because access_key_id and secret_access_key is wrong")
	}
}

func TestValidate2(t *testing.T) {
	repoType := "http"
	url := "https://kubernetes-charts.storage.googleapis.com"
	credential := ``
	providers := []string{"qingcloud"}

	err := validate(ctx, repoType, url, credential, providers)

	if err == nil {
		t.Errorf("expect error, because type is not matched")
	}
}

func TestValidate3(t *testing.T) {
	repoType := "https"
	url := "https://kubernetes-charts.storage.googleapis.com"
	credential := ``
	providers := []string{"qingcloud"}

	err := validate(ctx, repoType, url, credential, providers)

	if err != nil {
		t.Error(err)
	}
}

func TestValidate4(t *testing.T) {
	repoType := "https"
	url := "https://kubernetes1-charts.storage.googleapis.com"
	credential := ``
	providers := []string{"qingcloud"}

	err := validate(ctx, repoType, url, credential, providers)

	if err == nil {
		t.Error("error expect, because this is a bad url")
	}
}

func TestValidate5(t *testing.T) {
	repoType := "https"
	url := "https://baidu.com"
	credential := ``
	providers := []string{"qingcloud"}

	err := validate(ctx, repoType, url, credential, providers)

	if err == nil {
		t.Error("error expect, because this is a bad url")
	}
}
