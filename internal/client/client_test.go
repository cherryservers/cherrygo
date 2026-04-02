package client_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cherryservers/cherrygo/v3/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBackoff() client.BackoffFunc {
	return func(_ int, _ *http.Response) time.Duration {
		return time.Nanosecond
	}
}

func TestClientOneShotSuccess(t *testing.T) {
	client := client.New()
	count := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "test")
		count++
	}))
	defer ts.Close()

	spy := spyReadCloser{
		real: io.NopCloser(strings.NewReader("test")),
	}

	req, err := http.NewRequest("GET", ts.URL, &spy)
	require.NoError(t, err)

	resp, doErr := client.Do(req)
	body, readErr := io.ReadAll(resp.Body)
	closeErr := resp.Body.Close()

	assert.NoError(t, doErr)
	assert.NoError(t, readErr)
	assert.NoError(t, closeErr)
	assert.Equal(t, "test\n", string(body))
	assert.Equal(t, 1, count)
	assert.Equal(t, 1, spy.closes)
}
