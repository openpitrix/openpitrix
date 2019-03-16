// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package retryutil

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"
)

func Retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	ctx := context.Background()
	return RetryWithContext(ctx, attempts, sleep, callback)
}

func RetryWithContext(ctx context.Context, attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= attempts {
			break
		}

		if sleep > 0 {
			time.Sleep(sleep)
		}

		logger.Warn(ctx, "Will retry %d because of error: %+v", i, err)
	}
	return fmt.Errorf("failed after %d attempts, error: %+v", attempts, err)
}
