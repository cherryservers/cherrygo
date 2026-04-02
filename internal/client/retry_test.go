package client_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cherryservers/cherrygo/v3/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSpies(t *testing.T) (spyReadCloser, spyRoundTripper) {
	t.Helper()

	reqSpy := spyReadCloser{
		real: io.NopCloser(strings.NewReader("test")),
	}

	respSpy := spyRoundTripper{
		real: http.DefaultTransport,
	}

	return reqSpy, respSpy
}

type retryTestCase struct {
	title          string
	maxRetries     int
	method         string
	status         int
	fn             http.HandlerFunc
	reqCtx         context.Context
	backoffFn      client.BackoffFunc
	wantCalls      int
	wantErr        bool
	wantRespDrains int
	wantBody       string
}

func testRetry(t *testing.T, tc retryTestCase) {
	t.Run(tc.title, func(t *testing.T) {
		reqSpy, respSpy := setupSpies(t)

		httpClient := http.Client{
			Transport: &respSpy,
		}

		client := client.New(
			client.WithBackoff(tc.backoffFn),
			client.WithHTTPClient(&httpClient),
			client.WithMaxRetries(tc.maxRetries),
		)
		calls := 0

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tc.fn(w, r)
			calls++
		}))
		defer ts.Close()

		req, err := http.NewRequestWithContext(tc.reqCtx, tc.method, ts.URL, &reqSpy)
		require.NoError(t, err)

		resp, err := client.Do(req)

		assert.Equal(t, tc.wantRespDrains, respSpy.rc.closes, "Unexpected number of response drains.")
		assert.Equal(t, tc.wantRespDrains, respSpy.rc.reads, "Unexpected number of response drains.")
		assert.Equal(t, 1, reqSpy.closes, "The original request body should be closed exactly once.")
		assert.Equal(t, tc.wantCalls, calls, "Incorrect number HTTP server requests.")

		if tc.wantErr {
			assert.Error(t, err)
			assert.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)

			body, readErr := io.ReadAll(resp.Body)
			closeErr := resp.Body.Close()

			assert.NoError(t, readErr)
			assert.NoError(t, closeErr)
			assert.Equal(t, tc.wantBody, string(body))
		}
	})
}

func TestRetry(t *testing.T) {
	retryable := map[string][]int{
		"GET":    {408, 429, 502, 503, 504},
		"DELETE": {408, 429, 502, 503, 504},
		"PUT":    {408, 429, 502, 503, 504},
		"POST":   {429, 503},
		"PATCH":  {429, 503},
	}

	for method, statuses := range retryable {
		for _, status := range statuses {
			calls := 0
			testRetry(t, retryTestCase{
				title:      fmt.Sprintf("succeed on retry: method %q, status %d", method, status),
				maxRetries: 5,
				method:     method,
				status:     status,
				fn: func(w http.ResponseWriter, _ *http.Request) {
					if calls == 0 {
						w.WriteHeader(status)
						_, _ = fmt.Fprint(w, "error")
					} else {
						_, _ = fmt.Fprint(w, "test")
					}
					calls++
				},
				reqCtx:         t.Context(),
				backoffFn:      testBackoff(),
				wantCalls:      2,
				wantErr:        false,
				wantRespDrains: 1,
				wantBody:       "test",
			})
			testRetry(t, retryTestCase{
				title:      fmt.Sprintf("exceed max retries: method %q, status %d", method, status),
				maxRetries: 5,
				method:     method,
				status:     status,
				fn: func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(status)
					_, _ = fmt.Fprint(w, "error")
				},
				reqCtx:         t.Context(),
				backoffFn:      testBackoff(),
				wantCalls:      6,
				wantErr:        true,
				wantRespDrains: 6,
				wantBody:       "error",
			})
		}
	}

	unRetryable := map[string][]int{
		"POST":  {408, 502, 504},
		"PATCH": {408, 502, 504},
	}

	for method, statuses := range unRetryable {
		for _, status := range statuses {
			calls := 0
			testRetry(t, retryTestCase{
				title:      fmt.Sprintf("don't retry: method %q, status %d", method, status),
				maxRetries: 5,
				method:     method,
				status:     status,
				reqCtx:     t.Context(),
				backoffFn:  testBackoff(),
				fn: func(w http.ResponseWriter, _ *http.Request) {
					if calls == 0 {
						w.WriteHeader(status)
						_, _ = fmt.Fprint(w, "error")
					} else {
						_, _ = fmt.Fprint(w, "test")
					}
					calls++
				},
				wantCalls:      1,
				wantErr:        false,
				wantRespDrains: 0,
				wantBody:       "error",
			})
		}
	}
}

type errorSpyRoundTripper struct {
	err   error
	calls int
}

func (e *errorSpyRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	e.calls++
	return nil, e.err
}

type fakeNetError struct {
	transient bool
}

func (e fakeNetError) Error() string {
	return "test net error"
}

func (e fakeNetError) Temporary() bool {
	return e.transient
}

func (e fakeNetError) Timeout() bool {
	return e.transient
}

var _ net.Error = (*fakeNetError)(nil)

func TestRetryNetworkError(t *testing.T) {
	cases := []struct {
		transient bool
		wantCalls int
		title     string
	}{
		{
			title:     "don't retry non-transient net error",
			transient: false,
			wantCalls: 1,
		},
		{
			title:     "retry transient net error",
			transient: true,
			wantCalls: 6,
		},
	}

	for _, td := range cases {
		t.Run(td.title, func(t *testing.T) {
			reqSpy := spyReadCloser{
				real: io.NopCloser(strings.NewReader("test")),
			}

			fakeErr := fakeNetError{
				transient: td.transient,
			}

			rtSpy := errorSpyRoundTripper{
				err: fakeErr,
			}

			httpClient := http.Client{
				Transport: &rtSpy,
			}

			client := client.New(
				client.WithBackoff(testBackoff()),
				client.WithHTTPClient(&httpClient),
				client.WithMaxRetries(5),
			)

			req, err := http.NewRequest("GET", "fake-url", &reqSpy)
			require.NoError(t, err)

			resp, err := client.Do(req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, 1, reqSpy.closes, "The original request body should be closed exactly once.")
			assert.Equal(t, td.wantCalls, rtSpy.calls)
		})
	}
}

func TestRetryContextCancellation(t *testing.T) {
	calls := 0
	reqCtx, cancel := context.WithCancel(t.Context())

	var ctxCancelBackoff client.BackoffFunc = func(attempts int, _ *http.Response) time.Duration {
		if attempts < 2 {
			return time.Nanosecond
		}
		cancel()
		return time.Second
	}

	testRetry(t, retryTestCase{
		title:      "stop on request context cancellation",
		maxRetries: 5,
		method:     "GET",
		status:     429,
		fn: func(w http.ResponseWriter, _ *http.Request) {
			if calls < 3 {
				w.WriteHeader(429)
				_, _ = fmt.Fprint(w, "error")
			} else {
				_, _ = fmt.Fprint(w, "test")
			}
			calls++
		},
		reqCtx:         reqCtx,
		backoffFn:      ctxCancelBackoff,
		wantCalls:      3,
		wantErr:        true,
		wantRespDrains: 3,
		wantBody:       "",
	})
}
