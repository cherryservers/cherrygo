// Package client provides an [*net/http.Client] wrapper that automatically
// retries failed requests when it is safe to do so.
package client

import (
	"io"
	"net/http"
	"time"
)

const (
	defaultMaxRetries                   = 5
	defaultExponentialBackoffBase       = 1 * time.Second
	defaultExponentialBackoffCap        = 30 * time.Second
	defaultExponentialBackoffMultiplier = 2
)

type requestDoer interface {
	// Do performs an HTTP request.
	// May not strictly adhere to RoundTripper/Client semantics.
	Do(*http.Request) (*http.Response, error)
}

// Client is a wrapper for http.Client that can retry
// on transient errors, with respect to request idempotency.
// Safe for concurrent use, just like [net/http.Client].
type Client struct {
	maxRetries int
	backoff    BackoffFunc
	debugDst   io.Writer

	requestDoer
	rootClient *http.Client
}

// Option is a client configuration option.
type Option func(*Client)

// WithMaxRetries sets a custom amount of maximum retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		c.maxRetries = n
	}
}

// WithBackoff sets a custom backoff generation function.
func WithBackoff(b BackoffFunc) Option {
	return func(c *Client) {
		c.backoff = b
	}
}

// WithHTTPClient sets a custom base HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cc *Client) {
		cc.rootClient = c
	}
}

// WithDebug enables debug mode, which dumps logs to w.
// No logs will be dumped if w is nil.
func WithDebug(w io.Writer) Option {
	return func(c *Client) {
		c.debugDst = w
	}
}

// New creates a new client.
func New(opts ...Option) *Client {
	client := Client{
		maxRetries: defaultMaxRetries,
		backoff: RateLimitedExponentialBackoff(
			ExponentialBackoffConfig{
				Base:       defaultExponentialBackoffBase,
				Cap:        defaultExponentialBackoffCap,
				Multiplier: defaultExponentialBackoffMultiplier,
			},
		),
		rootClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&client)
	}

	var c requestDoer = client.rootClient
	if client.debugDst != nil {
		c = &debugger{
			wrapped: client.rootClient,
			dst:     client.debugDst,
		}
	}

	client.requestDoer = &retrier{
		wrapped:    c,
		maxRetries: client.maxRetries,
		backoff:    client.backoff,
	}

	return &client
}
