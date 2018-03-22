package repo

import (
	"strings"
	"testing"
)

func TestValidate1(t *testing.T) {
	repoType := "s3"
	url := "https://s3.pek3a.qingstor.com/op-repo"
	credential := `{"access_key_id": "wiandianiaeudsadf8a33uffhufhud", "secret_access_key": "nduaufbuabfuebaufbaufaueuu"}`
	visibility := "public"

	err := validate(repoType, url, credential, visibility)

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
	url := "http://www.qingcloud.com"
	credential := ``
	visibility := "public"

	err := validate(repoType, url, credential, visibility)

	if err != nil {
		t.Error(err)
	}
}

func TestValidate3(t *testing.T) {
	repoType := "https"
	url := "https://www.qingcloud.com"
	credential := ``
	visibility := "public"

	err := validate(repoType, url, credential, visibility)

	if err != nil {
		t.Error(err)
	}
}
