// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package internals3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var creds = credentials.NewStaticCredentials(
	"openpitrixminioaccesskey",
	"openpitrixminiosecretkey",
	"")

var config = &aws.Config{
	Region:           aws.String("us-east-1"),
	Endpoint:         aws.String("http://openpitrix-minio:9000"),
	DisableSSL:       aws.Bool(true),
	S3ForcePathStyle: aws.Bool(true),
	Credentials:      creds,
}

var Bucket = aws.String("openpitrix-attachment")

var S3 *s3.S3

func init() {
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}
	S3 = s3.New(sess)
}
