package utils

import (
	"fmt"
	"time"
)

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

// WaitForSpecific wait a function return true.
func WaitForSpecific(f func() bool, timeout time.Duration, waitInterval time.Duration) error {
	return WaitForSpecificOrError(func() (bool, error) {
		return f(), nil
	}, timeout, waitInterval)
}

// WaitFor wait a function return true.
func WaitFor(f func() bool) error {
	return WaitForSpecific(f, 180*time.Second, 3*time.Second)
}
