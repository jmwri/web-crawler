package internal

import "time"

// BackoffFunc returns the duration to sleep before making an attempt
type BackoffFunc func(attempt int) time.Duration

// SimpleBackoff adds a small increment for each attempt
func SimpleBackoff(attempt int) time.Duration {
	increment := time.Millisecond * 500
	return time.Duration(attempt-1) * increment
}
