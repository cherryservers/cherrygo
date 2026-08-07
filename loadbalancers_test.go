package cherrygo

import (
	"io"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertContainsFunc[S ~[]E, E any](t *testing.T, s S, f func(E) bool, msg string) {
	assert.True(t, slices.ContainsFunc(s, f), msg)
}

func timeMustParse(t string) time.Time {
	tt, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err.Error())
	}
	return tt
}

func TestLoadBalancer_Get(t *testing.T) {
	setup()
	defer teardown()

	apiBody, err := os.Open(filepath.Join(".", "testdata", "loadbalancer.json"))
	require.NoError(t, err)
	defer func() {
		closeErr := apiBody.Close()
		t.Log(closeErr.Error())
	}()

	mux.HandleFunc("GET /v1/load-balancers/931318", func(w http.ResponseWriter, _ *http.Request) {
		_, handleErr := io.Copy(w, apiBody)
		require.NoError(t, handleErr)
	})

	got, _, err := testClient.LoadBalancers.Get(t.Context(), 931318, nil)
	require.NoError(t, err)

	wantAttribute := LoadBalancerPlanAttribute{
		Name:  "Number of vCPUs",
		Value: "1",
	}

	wantRule := LoadBalancerRule{
		ID:                  "9d4007cc-8109-11f1-a8ee-00163e7dabb3",
		SourcePort:          8765,
		DestinationPort:     6443,
		SourceProtocol:      "tcp",
		DestinationProtocol: "tcp",
		StickyEnabled:       false,
		Balance:             "roundrobin",
		Status:              "Active",
	}

	wantBackend := LoadBalancerBackend{
		Rule:         wantRule,
		IP:           netip.MustParseAddr("10.194.211.44"),
		Port:         6443,
		Status:       "Down",
		StatusDetail: "test",
		Weight:       99,
		Fails:        1,
		Downs:        1,
		Downtime:     62915,
	}

	wantCert := LoadBalancerCertificate{
		ID:      "9d557ed5-8048-11f1-a8ee-00163e7dabb3",
		CN:      "system:apiserver",
		Starts:  timeMustParse("2026-07-10T11:45:00+03:00"),
		Expires: timeMustParse("2027-07-10T11:45:00+03:00"),
	}

	wantHealthCheck := LoadBalancerHealthCheck{
		Enabled:            true,
		Path:               "/",
		Interval:           3000,
		HealthyThreshold:   3,
		UnhealthyThreshold: 3,
	}

	assert.Equal(t, 931318, got.ID)
	assert.Equal(t, "test", got.Name)
	assert.Equal(t, "deployed", got.Status)
	assert.Equal(t, "load_balancer_1", got.Plan.Slug)
	assert.Equal(t, float32(0.0166), got.Plan.Pricing[0].Price)
	assertContainsFunc(t, got.Plan.Regions, func(r Region) bool {
		return r.Slug == "LT-Siauliai"
	}, "plan regions don't contain \"LT-Siauliai\"")
	assert.Contains(t, got.Plan.Attributes, wantAttribute)
	assert.Equal(t, 927312, got.Servers[0].ID)
	assert.Contains(t, got.Rules, wantRule)
	assertContainsFunc(t, got.IPAdresses, func(ip IPAddress) bool {
		return ip.ID == "6a8132f6-215c-d83f-8259-6ec3fdc6185a"
	}, "doesn't have ip with id 6a8132f6-215c-d83f-8259-6ec3fdc6185a")
	assert.Equal(t, "LT-Siauliai", got.Region.Slug)
	assert.Contains(t, got.Backends, wantBackend)
	assert.Equal(t, "test", got.StickyCookie)
	assert.True(t, got.StickyEnabled)
	assert.True(t, got.HTTPSRedirectEnabled)
	assert.True(t, got.ProxyEnabled)
	assert.Contains(t, got.Certificates, wantCert)
	assert.Equal(t, float32(1.11), got.Pricing.Price)
	assert.Equal(t, wantHealthCheck, got.HealthCheck)
}
