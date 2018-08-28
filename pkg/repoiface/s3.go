// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type S3Interface struct {
	url             *neturl.URL
	accessKeyId     string
	secretAccessKey string
	bucket          string
	config          *aws.Config
}

type S3Credential struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

var (
	compRegEx = regexp.MustCompile(`^s3\.(?P<zone>.+)\.(?P<host>.+\..+)/(?P<bucket>.+)/?$`)
)

func NewS3Interface(ctx context.Context, u *neturl.URL, credential string) (*S3Interface, error) {
	var s3Credential S3Credential
	err := jsonutil.Decode([]byte(credential), &s3Credential)
	if err != nil {
		return nil, ErrDecodeJsonFailed
	}

	if s3Credential.AccessKeyId == "" {
		return nil, ErrEmptyAccessKeyId
	}

	if s3Credential.SecretAccessKey == "" {
		return nil, ErrEmptySecretAccessKey
	}

	accessKeyId := s3Credential.AccessKeyId
	secretAccessKey := s3Credential.SecretAccessKey

	m := compRegEx.FindStringSubmatch(u.Host + u.Path)
	if len(m) != 0 && len(m) == 4 {
		zone := m[1]
		host := m[2]
		bucket := m[3]
		creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
		config := &aws.Config{
			Region:      aws.String(zone),
			Endpoint:    aws.String(fmt.Sprintf("%s://s3.%s.%s", "https", zone, host)),
			Credentials: creds,
		}
		return &S3Interface{
			url:             u,
			accessKeyId:     accessKeyId,
			secretAccessKey: secretAccessKey,
			config:          config,
			bucket:          bucket,
		}, nil
	} else {
		logger.Error(ctx, "Repo url [%s] regex test failed", u.String())
		return nil, ErrParseUrlFailed
	}
}

func (i *S3Interface) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	sess, err := session.NewSession(i.config)
	if err != nil {
		logger.Error(ctx, "Connect to s3 failed: %+v", err)
		return nil, ErrGetIndexYamlFailed
	}

	svc := s3.New(sess)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		logger.Error(ctx, "Failed to get s3 repo [%+v] index.yaml, error: %+v", i, err)
		return nil, ErrGetIndexYamlFailed
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, ErrGetIndexYamlFailed
	}
	return body, nil
}

func (i *S3Interface) WriteFile(ctx context.Context, filename string, data []byte) error {
	return nil
}
