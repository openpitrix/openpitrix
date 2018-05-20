// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package reporeader

import (
	"fmt"
	"io/ioutil"
	neturl "net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"openpitrix.io/openpitrix/pkg/logger"
)

type S3Reader struct {
	url             *neturl.URL
	accessKeyId     string
	secretAccessKey string
	bucket          string
	config          *aws.Config
}

func NewS3Reader(url *neturl.URL, accessKeyId, secretAccessKey, zone, host, bucket string) *S3Reader {
	creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
	config := &aws.Config{
		Region:      aws.String(zone),
		Endpoint:    aws.String(fmt.Sprintf("%s://s3.%s.%s/%s/", "http", zone, host, bucket)),
		Credentials: creds,
	}
	return &S3Reader{
		url:             url,
		accessKeyId:     accessKeyId,
		secretAccessKey: secretAccessKey,
		config:          config,
		bucket:          bucket,
	}
}

func (s *S3Reader) GetIndexYaml() ([]byte, error) {
	sess, err := session.NewSession(s.config)
	if err != nil {
		logger.Error("Connect to s3 failed: %+v", err)
		return nil, ErrGetIndexYamlFailed
	}

	svc := s3.New(sess)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(IndexYaml),
	})
	if err != nil {
		logger.Error("Failed to get s3 repo [%+v] index.yaml, error: %+v", s, err)
		return nil, ErrGetIndexYamlFailed
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, ErrGetIndexYamlFailed
	}
	return body, nil
}
