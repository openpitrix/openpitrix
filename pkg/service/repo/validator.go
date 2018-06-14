package repo

import (
	"fmt"
	"regexp"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/reporeader"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type QingStorCredential struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type IndexYaml struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Entries    map[string]interface{} `yaml:"entries"`
	Generated  string                 `yaml:"generated"`
}

var (
	compRegEx = regexp.MustCompile(`^s3\.(?P<zone>.+)\.(?P<host>.+\..+)/(?P<bucket>.+)/?$`)
)

func validateVisibility(visibility string) error {
	switch visibility {
	case "public":
	case "private":
	default:
		return fmt.Errorf("visibility must be one of [public, private]")
	}

	return nil
}

func validateProviders(providers []string) error {
	if len(providers) == 0 {
		return fmt.Errorf("providers must be provided")
	}

	for _, provider := range providers {
		if !stringutil.StringIn(provider, []string{constants.ProviderKubernetes, constants.ProviderQingCloud}) {
			return fmt.Errorf("provider must be in range [kubernetes, qingcloud]")
		}
	}

	return nil
}

func validate(repoType, url, credential, visibility string, providers []string) error {
	err := validateVisibility(visibility)
	if err != nil {
		return newErrorWithCode(ErrVisibility, err)
	}

	err = validateProviders(providers)
	if err != nil {
		return newErrorWithCode(ErrProviders, err)
	}

	var errCode uint32
	reader, err := reporeader.New(repoType, url, credential)
	if err != nil {
		switch err {
		case reporeader.ErrParseUrlFailed:
			errCode = ErrUrlFormat
		case reporeader.ErrDecodeJsonFailed:
			errCode = ErrCredentialNotJson
		case reporeader.ErrEmptyAccessKeyId:
			errCode = ErrNoAccessKeyId
		case reporeader.ErrEmptySecretAccessKey:
			errCode = ErrNoSecretAccessKey
		case reporeader.ErrInvalidType:
			errCode = ErrType
		case reporeader.ErrSchemeNotMatched:
			switch repoType {
			case constants.TypeHttp:
				errCode = ErrSchemeNotHttp
			case constants.TypeHttps:
				errCode = ErrSchemeNotHttps
			case constants.TypeS3:
				errCode = ErrSchemeNotS3
			}
		}
		return newErrorWithCode(errCode, err)
	}

	body, err := reader.GetIndexYaml()
	if err != nil {
		switch err {
		case reporeader.ErrGetIndexYamlFailed:
			switch repoType {
			case constants.TypeHttp, constants.TypeHttps:
				errCode = ErrHttpAccessDeny
			case constants.TypeS3:
				errCode = ErrS3AccessDeny
			}

		}
		return newErrorWithCode(errCode, err)
	}

	var y IndexYaml
	err = yamlutil.Decode(body, &y)
	if err != nil {
		errCode = ErrBadIndexYaml
		return newErrorWithCode(errCode, err)
	}

	return nil
}
