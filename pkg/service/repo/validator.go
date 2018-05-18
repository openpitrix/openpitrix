package repo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"path"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type QingStorCredential struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
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

	u, err := neturl.ParseRequestURI(url)
	if err != nil {
		return newErrorWithCode(ErrNotUrl, fmt.Errorf("url parse failed, %s", err))
	}

	switch repoType {
	case "s3":
		if u.Scheme != "s3" {
			return newErrorWithCode(ErrSchemeNotS3, fmt.Errorf("scheme is not s3"))
		}

		//fmt.Printf("%#v\n", u)
		m := compRegEx.FindStringSubmatch(u.Host + u.Path)
		//fmt.Printf("%#v\n", m)
		logger.Debug("Regexp result: %+v", m)

		if len(m) != 0 && len(m) == 4 {
			zone := m[1]
			host := m[2]
			bucket := m[3]

			var qc QingStorCredential
			err = jsonutil.Decode([]byte(credential), &qc)
			if err != nil {
				return newErrorWithCode(ErrCredentialNotJson, fmt.Errorf("json decode failed on credential, %+v", err))
			}

			if qc.AccessKeyId == "" {
				return newErrorWithCode(ErrNoAccessKeyId, fmt.Errorf("access_key_id not exist in credential"))
			}

			if qc.SecretAccessKey == "" {
				return newErrorWithCode(ErrNoSecretAccessKey, fmt.Errorf("secret_access_key not exist in credential"))
			}

			err = ValidateS3("https", host, qc.AccessKeyId, qc.SecretAccessKey, bucket, zone)
			if err != nil {
				return newErrorWithCode(ErrS3AccessDeny, fmt.Errorf("validate s3 failed, %+v", err))
			}
		} else {
			return newErrorWithCode(ErrUrlFormat, fmt.Errorf("url format error"))
		}
	case "http":
		if u.Scheme != "http" {
			return newErrorWithCode(ErrSchemeNotHttp, fmt.Errorf("scheme is not http"))
		}
		err := ValidateHTTP(u)
		if err != nil {
			return err
		}
	case "https":
		if u.Scheme != "https" {
			return newErrorWithCode(ErrSchemeNotHttps, fmt.Errorf("scheme is not https"))
		}
		err := ValidateHTTP(u)
		if err != nil {
			return err
		}
	default:
		return newErrorWithCode(ErrType, fmt.Errorf("type must be one of [s3, http, https]"))
	}

	return nil
}

func ValidateHTTP(u *neturl.URL) error {
	u.Path = path.Join(u.Path, "index.yaml")

	resp, err := http.Get(u.String())
	if err != nil {
		return newErrorWithCode(ErrHttpAccessDeny, fmt.Errorf("validate http failed, %+v", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newErrorWithCode(ErrHttpAccessDeny, fmt.Errorf("validate http failed, %+v", err))
	}

	var vals map[string]interface{}
	err = yaml.Unmarshal(body, &vals)
	if err != nil {
		return newErrorWithCode(ErrNotRepoUrl, fmt.Errorf("validate http failed, %+v", err))
	}
	return nil
}

func ValidateS3(scheme, host, accessKeyId, secretAccessKey, bucket, zone string) error {
	creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
	config := &aws.Config{
		Region:      aws.String(zone),
		Endpoint:    aws.String(fmt.Sprintf("%s://s3.%s.%s/%s/", scheme, zone, host, bucket)),
		Credentials: creds,
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	_, err = svc.ListBuckets(nil)
	if err != nil {
		return err
	}

	//fmt.Println(resp)
	return nil
}
