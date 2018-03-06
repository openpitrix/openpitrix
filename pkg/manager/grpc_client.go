// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/logger"
)

func NewClient(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				logger.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				logger.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()
	return conn, err
}
