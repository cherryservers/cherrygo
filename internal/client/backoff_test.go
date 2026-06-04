package client_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/cherryservers/cherrygo/v3/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff(t *testing.T) {
	fn := client.RateLimitedExponentialBackoff(client.ExponentialBackoffConfig{
		Base:       2 * time.Second,
		Cap:        60 * time.Second,
		Multiplier: 2,
	})

	wantDelaysMin := []time.Duration{1, 2, 4, 8, 16, 32, 60, 60}
	wantDelaysMax := []time.Duration{2, 4, 8, 16, 32, 60, 60, 60}

	gotDelays := make([]time.Duration, 0, len(wantDelaysMax))
	gotDelaysAgain := make([]time.Duration, 0, len(wantDelaysMax))

	for i := range wantDelaysMax {
		wantDelaysMax[i] = wantDelaysMax[i] * time.Second
		wantDelaysMin[i] = wantDelaysMin[i] * time.Second

		gotDelays = append(gotDelays, fn(i, nil))
		gotDelaysAgain = append(gotDelaysAgain, fn(i, nil))
	}

	for i := range gotDelays {
		assert.LessOrEqual(t, wantDelaysMin[i], gotDelays[i], "Backoff delay too small.")
		assert.GreaterOrEqual(t, wantDelaysMax[i], gotDelays[i], "Backoff delay too big.")

		assert.LessOrEqual(t, wantDelaysMin[i], gotDelaysAgain[i], "Backoff delay too small.")
		assert.GreaterOrEqual(t, wantDelaysMax[i], gotDelaysAgain[i], "Backoff delay too big.")
	}
	assert.NotEqual(t, gotDelays, gotDelaysAgain, "Duplicate delays generated.")
}

func TestBackoffRetryAfter(t *testing.T) {
	now := time.Now()

	cases := []struct {
		retryAfter   string
		wantDelayMin time.Duration
		wantDelayMax time.Duration
	}{
		{
			retryAfter:   "1",
			wantDelayMin: 1 * time.Second,
			wantDelayMax: 1 * time.Second,
		},
		{
			retryAfter:   now.Add(time.Second * 30).UTC().Format(http.TimeFormat),
			wantDelayMin: 25 * time.Second,
			wantDelayMax: 30 * time.Second,
		},
		{
			retryAfter:   "",
			wantDelayMin: 1 * time.Second,
			wantDelayMax: 2 * time.Second,
		},
	}

	fn := client.RateLimitedExponentialBackoff(client.ExponentialBackoffConfig{
		Base:       2 * time.Second,
		Cap:        60 * time.Second,
		Multiplier: 2,
	})

	for _, td := range cases {
		t.Run(td.retryAfter, func(t *testing.T) {
			head := make(http.Header)
			head.Add("Retry-After", td.retryAfter)
			resp := &http.Response{
				Header: head,
			}
			gotDelay := fn(0, resp)

			assert.LessOrEqual(t, td.wantDelayMin, gotDelay)
			assert.GreaterOrEqual(t, td.wantDelayMax, gotDelay)
		})
	}
}
