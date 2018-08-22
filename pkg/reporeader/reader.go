// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package reporeader

import (
	"context"
	"fmt"
	neturl "net/url"
	"regexp"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Err error

var (
	ErrGetIndexYamlFailed   Err = fmt.Errorf("get index.yaml failed")
	ErrParseUrlFailed       Err = fmt.Errorf("parse url failed")
	ErrDecodeJsonFailed     Err = fmt.Errorf("decode json failed")
	ErrEmptyAccessKeyId     Err = fmt.Errorf("access key id is empty")
	ErrEmptySecretAccessKey Err = fmt.Errorf("secret access key is empty")
	ErrSchemeNotMatched     Err = fmt.Errorf("scheme not matched")
	ErrInvalidType          Err = fmt.Errorf("invalid repo type")
)

type Reader interface {
	GetIndexYaml() ([]byte, error)
}

type S3Credential struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

var (
	compRegEx = regexp.MustCompile(`^s3\.(?P<zone>.+)\.(?P<host>.+\..+)/(?P<bucket>.+)/?$`)
)

const IndexYaml = "index.yaml"

func New(ctx context.Context, repoType, url, credential string) (Reader, error) {
	u, err := neturl.ParseRequestURI(url)
	if err != nil {
		logger.Error(ctx, "Parse url [%s] failed, error: %+v", url, err)
		return nil, ErrParseUrlFailed
	}

	switch repoType {
	case constants.TypeS3:
		m := compRegEx.FindStringSubmatch(u.Host + u.Path)
		logger.Debug(ctx, "Repo url [%s] regexp result: %+v", url, m)

		if len(m) != 0 && len(m) == 4 {
			zone := m[1]
			host := m[2]
			bucket := m[3]

			var qc S3Credential
			err = jsonutil.Decode([]byte(credential), &qc)
			if err != nil {
				return nil, ErrDecodeJsonFailed
			}

			if qc.AccessKeyId == "" {
				return nil, ErrEmptyAccessKeyId
			}

			if qc.SecretAccessKey == "" {
				return nil, ErrEmptySecretAccessKey
			}

			return NewS3Reader(ctx, u, qc.AccessKeyId, qc.SecretAccessKey, zone, host, bucket), nil
		} else {
			logger.Error(ctx, "Repo url [%s] regex test failed", url)
			return nil, ErrParseUrlFailed
		}
	case constants.TypeHttp:
		if u.Scheme != constants.TypeHttp {
			return nil, ErrSchemeNotMatched
		}
		return NewHttpReader(ctx, u), nil
	case constants.TypeHttps:
		if u.Scheme != constants.TypeHttps {
			return nil, ErrSchemeNotMatched
		}
		return NewHttpReader(ctx, u), nil
	default:
		return nil, ErrInvalidType
	}

}
