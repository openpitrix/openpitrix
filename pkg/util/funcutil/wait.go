// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package funcutil

import (
	"fmt"
	"time"
)

// Reference: https://github.com/yunify/qingcloud-sdk-go/blob/f1702ea44452e2d6b6048608c070237cd45088d6/utils/wait.go
// TimeoutError An Error represents a timeout error.
type TimeoutError struct {
	timeout time.Duration
}

// Error message
func (e *TimeoutError) Error() string { return fmt.Sprintf("Wait timeout [%s] ", e.timeout) }

// Timeout duration
func (e *TimeoutError) Timeout() time.Duration { return e.timeout }

// NewTimeoutError create a new TimeoutError
func NewTimeoutError(timeout time.Duration) *TimeoutError {
	return &TimeoutError{timeout: timeout}
}

// WaitForSpecificOrError wait a function return true or error.
func WaitForSpecificOrError(f func() (bool, error), timeout time.Duration, waitInterval time.Duration) error {
	stop, err := f()
	if err != nil {
		return err
	}
	if stop {
		return nil
	}
	ticker := time.NewTicker(waitInterval)
	defer ticker.Stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			stop, err := f()
			if err != nil {
				return err
			}
			if stop {
				return nil
			}
		case <-timer.C:
			return NewTimeoutError(timeout)
		}
	}
}
