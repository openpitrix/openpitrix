package repo

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/reporeader"
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

func validate(ctx context.Context, repoType, url, credential string, providers []string) error {
	var errCode uint32
	reader, err := reporeader.New(ctx, repoType, url, credential)
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
