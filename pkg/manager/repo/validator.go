package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"openpitrix.io/openpitrix/pkg/logger"
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
		return errors.New("Visibility must be one of [public, private]")
	}

	return nil
}

func validate(repoType, url, credential, visibility string) error {
	err := validateVisibility(visibility)
	if err != nil {
		return err
	}

	u, err := neturl.ParseRequestURI(url)
	if err != nil {
		return fmt.Errorf("Url parse failed, %s", err)
	}

	switch repoType {
	case "s3":
		//fmt.Printf("%#v\n", u)
		m := compRegEx.FindStringSubmatch(u.Host + u.Path)
		//fmt.Printf("%#v\n", m)
		logger.Debugf("Regexp result: %+v", m)

		if len(m) != 0 && len(m) == 4 {
			zone := m[1]
			host := m[2]
			bucket := m[3]

			var qc QingStorCredential
			err = json.Unmarshal([]byte(credential), &qc)
			if err != nil {
				return fmt.Errorf("Json decode failed on credential, %+v", err)
			}

			if qc.AccessKeyId == "" || qc.SecretAccessKey == "" {
				return fmt.Errorf("access_key_id or secret_access_key not exist in credential")
			}

			err = ValidateS3(host, qc.AccessKeyId, qc.SecretAccessKey, bucket, zone)
			if err != nil {
				return fmt.Errorf("Validate qingstor failed, %+v", err)
			}
		} else {
			return errors.New("Url is not a bucket url of qingstor")
		}
	case "http":
		if u.Scheme != "http" {
			return errors.New("Scheme is not http")
		}
		err := ValidateHTTP(url)
		if err != nil {
			return fmt.Errorf("Validate http failed, %+v", err)
		}
	case "https":
		if u.Scheme != "https" {
			return errors.New("Scheme is not https")
		}
		err := ValidateHTTP(url)
		if err != nil {
			return fmt.Errorf("Validate https failed, %+v", err)
		}
	default:
		return fmt.Errorf("Type must be one of [s3, http, https]")
	}

	return nil
}

func ValidateHTTP(url string) error {
	_, err := http.Get(url)
	if err != nil {
		return err
	}
	return nil
}

func ValidateS3(host, access_key_id, secret_access_key, bucket, zone string) error {
	creds := credentials.NewStaticCredentials(access_key_id, secret_access_key, "")
	config := &aws.Config{
		Region:      aws.String(zone),
		Endpoint:    aws.String(fmt.Sprintf("http://s3.%s.%s/%s/", zone, host, bucket)),
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
