package client

import (
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// BackoffFunc generates backoff delays.
type BackoffFunc func(attempts int, resp *http.Response) time.Duration

// ExponentialBackoffConfig is the configuration for exponential backoff.
type ExponentialBackoffConfig struct {
	Base       time.Duration
	Cap        time.Duration
	Multiplier float64
}

// RateLimitedExponentialBackoff returns an backoff function
// that prioritizes `Retry-After` headers, defaulting to exponential backoff
// if they're not found. Adds partial jitter.
func RateLimitedExponentialBackoff(cfg ExponentialBackoffConfig) BackoffFunc {
	return func(attempts int, req *http.Response) time.Duration {
		backoff, ok := parseRetryAfter(req)
		if ok {
			return backoff
		}

		backoffSeconds := cfg.Base.Seconds() * math.Pow(cfg.Multiplier, float64(attempts))
		backoff = time.Second * time.Duration(backoffSeconds)
		backoff = jitter(backoff)

		return min(backoff, cfg.Cap)
	}
}

func parseRetryAfter(req *http.Response) (time.Duration, bool) {
	if req == nil {
		return 0, false
	}

	// Retry-After can be an integer with seconds or an HTTP date.
	d := req.Header.Get("Retry-After")

	seconds, err := strconv.Atoi(d)
	if err == nil {
		return time.Duration(seconds) * time.Second, true
	}

	date, err := http.ParseTime(d)
	if err == nil {
		return time.Until(date), true
	}

	return 0, false
}

func jitter(base time.Duration) time.Duration {
	return time.Duration(rand.Int64N(int64(base)/2)) + base/2
}
