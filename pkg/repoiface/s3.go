// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"path"
	"strings"

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
	prefix          string
	config          *aws.Config
}

type S3Credential struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

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

	// e.g. s3://s3.us-east-2.amazonaws.com/my-openpitrix
	creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
	var region, endpoint, bucket, prefix string
	bucket, prefix = getBucketPrefix(u.Path)

	if strings.HasPrefix(u.Host, "s3.") {
		region = strings.Split(u.Host, ".")[1]
		endpoint = fmt.Sprintf("https://%s", u.Host)
	} else {
		// If using alternative s3 endpoint (e.g. Minio) default region to us-east-1
		region = "us-east-1"
		endpoint = fmt.Sprintf("http://%s", u.Host)
	}

	config := &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		DisableSSL:       aws.Bool(strings.HasPrefix(endpoint, "http://")),
		S3ForcePathStyle: aws.Bool(endpoint != ""),
		Credentials:      creds,
	}
	return &S3Interface{
		url:             u,
		accessKeyId:     accessKeyId,
		secretAccessKey: secretAccessKey,
		config:          config,
		bucket:          bucket,
		prefix:          prefix,
	}, nil
}

func (i *S3Interface) getService(ctx context.Context) (*s3.S3, error) {
	sess, err := session.NewSession(i.config)
	if err != nil {
		logger.Error(ctx, "Connect to s3 [%s] failed: %+v", i.url, err)
		return nil, err
	}
	svc := s3.New(sess)
	return svc, err
}

func (i *S3Interface) CheckFile(ctx context.Context, filename string) (bool, error) {
	svc, err := i.getService(ctx)
	if err != nil {
		return false, err
	}

	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(path.Join(i.prefix, GetFileName(filename))),
	})
	if err != nil {
		logger.Error(ctx, "Failed to read file [%s] from s3 [%s], error: %+v", filename, i.url, err)
		return false, nil
	}

	return true, nil
}

func (i *S3Interface) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	svc, err := i.getService(ctx)
	if err != nil {
		return nil, err
	}

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(path.Join(i.prefix, GetFileName(filename))),
	})
	if err != nil {
		logger.Error(ctx, "Failed to read file [%s] from s3 [%s], error: %+v", filename, i.url, err)
		return nil, err
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (i *S3Interface) DeleteFile(ctx context.Context, filename string) error {
	svc, err := i.getService(ctx)
	if err != nil {
		return err
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(path.Join(i.prefix, GetFileName(filename))),
	})
	if err != nil {
		logger.Error(ctx, "Failed to delete file [%s] from s3 [%s], error: %+v", filename, i.url, err)
		return err
	}

	return nil
}

func (i *S3Interface) WriteFile(ctx context.Context, filename string, data []byte) error {
	svc, err := i.getService(ctx)
	if err != nil {
		return err
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(i.bucket),
		Key:    aws.String(path.Join(i.prefix, GetFileName(filename))),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		logger.Error(ctx, "Failed to write file [%s] to s3 [%s], error: %+v", filename, i.url, err)
		return err
	}

	return nil
}

func (i *S3Interface) CheckRead(ctx context.Context) error {
	svc, err := i.getService(ctx)
	if err != nil {
		return err
	}
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(i.bucket),
	})
	if err != nil {
		logger.Error(ctx, "Failed to get bucket info from s3 [%s], error: %+v", i.url, err)
	}
	return err
}

func (i *S3Interface) CheckWrite(ctx context.Context) error {
	return i.WriteFile(ctx, ".openpitrix.test", []byte{})
}
