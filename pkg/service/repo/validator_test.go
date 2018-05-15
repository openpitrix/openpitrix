// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"strings"
	"testing"
)

func TestValidate1(t *testing.T) {
	repoType := "s3"
	url := "s3://s3.pek3a.qingstor.com/op-repo"
	credential := `{"access_key_id": "wiandianiaeudsadf8a33uffhufhud", "secret_access_key": "nduaufbuabfuebaufbaufaueuu"}`
	visibility := "public"
	providers := []string{"qingcloud"}

	err := validate(repoType, url, credential, visibility, providers)

	if err == nil {
		t.Errorf("expect error, because access_key_id and secret_access_key is wrong")
	}

	ok := strings.Contains(err.Error(), "InvalidAccessKeyId")
	if !ok {
		t.Error(err)
	}
}

func TestValidate2(t *testing.T) {
	repoType := "http"
	url := "https://kubernetes-charts.storage.googleapis.com"
	credential := ``
	visibility := "public"
	providers := []string{"qingcloud"}

	err := validate(repoType, url, credential, visibility, providers)

	if err == nil {
		t.Errorf("expect error, because type is not matched")
	}
}

func TestValidate3(t *testing.T) {
	repoType := "https"
	url := "https://kubernetes-charts.storage.googleapis.com"
	credential := ``
	visibility := "public"
	providers := []string{"qingcloud"}

	err := validate(repoType, url, credential, visibility, providers)

	if err != nil {
		t.Error(err)
	}
}
