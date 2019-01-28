// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type EmailNotification struct {
	Title       string
	Content     string
	Owner       string
	ContentType string
	Addresses   []string
}
