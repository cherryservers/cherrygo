package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"slices"
	"time"
)

const bodyReadLimit = 4096

type retrier struct {
	wrapped    requestDoer
	maxRetries int
	backoff    BackoffFunc
}

// Do executes requests, retrying unsuccessful ones, when it safe to do so.
//
// Requests are passed to the wrapped [requestDoer], which is expected to act like a
// [net/http.Client]. If the response status code indicates success or is unsafe to
// retry, returns that response with a nil error. If the request context
// expires or the retry attempt limit is reached, the response will be nil and
// an error will be returned.
func (r *retrier) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	var lastErr error

	// The original request is discarded, make sure its body is closed.
	defer func() {
		if req.Body != nil {
			_ = req.Body.Close()
		}
	}()

	for attempts := 0; attempts < r.maxRetries+1; attempts++ {
		clone, err := cloneRequest(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to clone request: %w", err)
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		resp, err := r.wrapped.Do(clone)
		if err != nil {
			if !isTransient(err) || !isIdempotent(clone.Method) {
				return nil, err
			}
			lastErr = err
		} else {
			if isSuccessful(resp.StatusCode) || !safeToRetry(resp.StatusCode, clone.Method) {
				return resp, nil
			}
			lastErr = fmt.Errorf("bad status: %q", resp.Status)

			// Try to drain and close the body, so the connection is freed.
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, bodyReadLimit))
			_ = resp.Body.Close()
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(r.backoff(attempts, resp)):
		}
	}

	return nil, fmt.Errorf("max retries %d exceeded, last attempt: %w", r.maxRetries, lastErr)
}

func cloneRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	c := req.Clone(ctx)
	if req.Body == nil {
		return c, nil
	}

	if req.GetBody == nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		_ = req.Body.Close()

		req.GetBody = func() (io.ReadCloser, error) {
			buf := bytes.NewBuffer(body)
			return io.NopCloser(buf), nil
		}

		req.Body, err = req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("failed to get request body: %w", err)
		}
	}

	var err error

	c.Body, err = req.GetBody()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func isTransient(err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				return true
			}
		}
	}

	return false
}

func isSuccessful(status int) bool {
	return status >= 200 && status < 300
}

func isIdempotent(method string) bool {
	return method != http.MethodConnect &&
		method != http.MethodPost &&
		method != http.MethodPatch
}

func safeToRetry(status int, method string) bool {
	idempotentRetryable := []int{
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
	nonIdempotentRetryable := []int{
		http.StatusServiceUnavailable,
		http.StatusTooManyRequests,
	}

	if isIdempotent(method) {
		return slices.Contains(idempotentRetryable, status)
	}
	return slices.Contains(nonIdempotentRetryable, status)
}
