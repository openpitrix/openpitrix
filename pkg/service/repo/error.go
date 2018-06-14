package repo

const (
	ErrNotExpect = 901

	ErrVisibility        = 101
	ErrNotUrl            = 102
	ErrCredentialNotJson = 103
	ErrNoAccessKeyId     = 104
	ErrNoSecretAccessKey = 105
	ErrS3AccessDeny      = 106
	ErrUrlFormat         = 107
	ErrSchemeNotHttp     = 108
	ErrHttpAccessDeny    = 109
	ErrSchemeNotHttps    = 110
	ErrType              = 111
	ErrProviders         = 112
	ErrNotRepoUrl        = 113
	ErrSchemeNotS3       = 114
	ErrBadIndexYaml      = 115
)

type ErrorWithCode struct {
	code uint32
	msg  string
}

func newErrorWithCode(code uint32, err error) error {
	return &ErrorWithCode{
		code: code,
		msg:  err.Error(),
	}
}

func (e *ErrorWithCode) Error() string {
	return e.msg
}

func (e *ErrorWithCode) Code() uint32 {
	return e.code
}
