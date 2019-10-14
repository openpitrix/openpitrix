// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package internals3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	opconfig "openpitrix.io/openpitrix/pkg/config"
)

var aConf = opconfig.GetConf().Attachment

var creds = credentials.NewStaticCredentials(
	aConf.AccessKey,
	aConf.SecretKey,
	"")

var config = &aws.Config{
	Region:           aws.String("us-east-1"),
	Endpoint:         aws.String(aConf.Endpoint),
	DisableSSL:       aws.Bool(true),
	S3ForcePathStyle: aws.Bool(true),
	Credentials:      creds,
}

var Bucket = aws.String(aConf.BucketName)

var S3 *s3.S3

func init() {
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}
	S3 = s3.New(sess)
}
