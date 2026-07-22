// Package backoff provides backoff calculation functions.
package backoff

import (
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// Func generates backoff delays.
type Func func(attempts int, resp *http.Response) time.Duration

// ExponentialBackoffConfig is the configuration for exponential backoff.
type ExponentialBackoffConfig struct {
	Base       time.Duration
	Cap        time.Duration
	Multiplier float64
}

// RateLimitedExponentialBackoff returns an backoff function
// that prioritizes `Retry-After` headers, defaulting to exponential backoff
// if they're not found. Adds partial jitter.
func RateLimitedExponentialBackoff(cfg ExponentialBackoffConfig) Func {
	exp := ExponentialBackoff(cfg)
	return func(attempts int, resp *http.Response) time.Duration {
		backoff, ok := parseRetryAfter(resp)
		if ok {
			return backoff
		}

		return exp(attempts, resp)
	}
}

// ExponentialBackoff returns an exponential backoff with partial jitter.
func ExponentialBackoff(cfg ExponentialBackoffConfig) Func {
	return func(attempts int, _ *http.Response) time.Duration {
		backoffSeconds := cfg.Base.Seconds() * math.Pow(cfg.Multiplier, float64(attempts))

		// Guard against overflow on high iterations.
		if math.IsNaN(backoffSeconds) ||
			math.IsInf(backoffSeconds, 0) ||
			backoffSeconds <= 0 ||
			backoffSeconds >= cfg.Cap.Seconds() {
			return cfg.Cap
		}
		backoff := time.Duration(backoffSeconds * float64(time.Second))
		return jitter(backoff)
	}
}

func parseRetryAfter(resp *http.Response) (time.Duration, bool) {
	if resp == nil {
		return 0, false
	}

	// Retry-After can be an integer with seconds or an HTTP date.
	d := resp.Header.Get("Retry-After")

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
