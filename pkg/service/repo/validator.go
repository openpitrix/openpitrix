package repo

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/repoiface"
)

func validate(ctx context.Context, repoType, url, credential string) error {
	var errCode uint32
	r, err := repoiface.New(ctx, repoType, url, credential)
	if err != nil {
		switch err {
		case repoiface.ErrParseUrlFailed:
			errCode = ErrUrlFormat
		case repoiface.ErrDecodeJsonFailed:
			errCode = ErrCredentialNotJson
		case repoiface.ErrEmptyAccessKeyId:
			errCode = ErrNoAccessKeyId
		case repoiface.ErrEmptySecretAccessKey:
			errCode = ErrNoSecretAccessKey
		case repoiface.ErrInvalidType:
			errCode = ErrType
		case repoiface.ErrSchemeNotMatched:
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

	err = r.CheckRead(ctx)
	if err != nil {
		switch repoType {
		case constants.TypeHttp:
			errCode = ErrHttpAccessDeny
		case constants.TypeHttps:
			errCode = ErrHttpAccessDeny
		case constants.TypeS3:
			errCode = ErrS3AccessDeny
		}
		return newErrorWithCode(errCode, err)
	}

	return nil
}
