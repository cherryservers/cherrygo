package cherrygo

import (
	"fmt"
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

	apiBody, err := os.ReadFile(filepath.Join(".", "testdata", "loadbalancer", "get.json"))
	require.NoError(t, err)

	setupGetWithOptsHandler(t, apiBody, "GET /v1/load-balancers/931318")

	got, _, err := testClient.LoadBalancers.Get(t.Context(), 931318, &GetOptions{Limit: 1})
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

func TestLoadBalancer_List(t *testing.T) {
	setup()
	defer teardown()

	apiBody, err := os.ReadFile(filepath.Join(".", "testdata", "loadbalancer", "list.json"))
	require.NoError(t, err)

	setupGetWithOptsHandler(t, apiBody, "GET /v1/projects/1/load-balancers")

	got, _, err := testClient.LoadBalancers.List(t.Context(), 1, &GetOptions{Limit: 1})
	require.NoError(t, err)

	assert.Equal(t, 931318, got[0].ID)
}

func TestLoadBalancer_Create(t *testing.T) {
	cases := []struct {
		name     string
		req      CreateLoadBalancer
		wantBody string
	}{
		{
			name: "only plan set",
			req: CreateLoadBalancer{
				Plan: "test-plan",
			},
			wantBody: "{\"slug\":\"test-plan\"}\n",
		},
		{
			name: "all fields set",
			req: CreateLoadBalancer{
				Name:   "test-name",
				Plan:   "test-plan",
				Region: "test-region",
				Cycle:  "test-cycle",
			},
			wantBody: "{\"name\":\"test-name\",\"slug\":\"test-plan\",\"region\":\"test-region\",\"cycle\":\"test-cycle\"}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()

			setupHandler(
				t,
				tc.wantBody,
				`{"id":1}`,
				http.StatusCreated,
				"POST /v1/projects/1/load-balancers",
			)

			got, _, err := testClient.LoadBalancers.Create(t.Context(), 1, tc.req)
			require.NoError(t, err)

			assert.Equal(t, 1, got.ID)
		})
	}
}

func TestLoadBalancer_Delete(t *testing.T) {
	setup()
	defer teardown()

	setupHandler(t, "", "", http.StatusNoContent, "DELETE /v1/load-balancers/1")

	_, err := testClient.LoadBalancers.Delete(t.Context(), 1)
	require.NoError(t, err)
}

func TestLoadBalancer_ListPlans(t *testing.T) {
	setup()
	defer teardown()

	apiBody, err := os.ReadFile(filepath.Join(".", "testdata", "loadbalancer", "plans.json"))
	require.NoError(t, err)

	setupGetWithOptsHandler(t, apiBody, "GET /v1/teams/1/load-balancer-plans")

	got, _, err := testClient.LoadBalancers.ListPlans(t.Context(), 1, &GetOptions{Limit: 1})
	require.NoError(t, err)

	wantAttribute := LoadBalancerPlanAttribute{
		Name:  "Number of vCPUs",
		Value: "1",
	}

	assert.Equal(t, "load_balancer_1", got[0].Slug)
	assert.Equal(t, float32(12.1), got[0].Pricing[0].Price)
	assert.Equal(t, "LT-Siauliai", got[0].Regions[0].Slug)
	assert.Equal(t, wantAttribute, got[0].Attributes[0])
}

func TestLoadBalancer_Update(t *testing.T) {
	cases := []struct {
		name     string
		req      UpdateLoadBalancer
		wantBody string
	}{
		{
			name:     "all fields omitted",
			wantBody: "{}\n",
		},
		{
			name: "all fields present",
			req: UpdateLoadBalancer{
				Name:                 "n",
				Plan:                 "p",
				StickyCookie:         "c",
				StickyEnabled:        new(bool),
				HTTPSRedirectEnabled: new(bool),
				ProxyEnabled:         new(bool),
				HealthCheckEnabled:   new(bool),
				HealthCheckPath:      "/",
				HealthCheckInterval:  1,
				HealthyThreshold:     1,
				UnhealthyThreshold:   1,
			},
			wantBody: "{\"name\":\"n\",\"slug\":\"p\"," +
				"\"sticky_cookie\":\"c\",\"sticky_enabled\":false," +
				"\"https_redirect_enabled\":false,\"proxy_enabled\":false," +
				"\"health_check_enabled\":false,\"health_check_path\":\"/\"," +
				"\"health_check_interval\":1,\"healthy_threshold\":1," +
				"\"unhealthy_threshold\":1}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()
			setupHandler(
				t,
				tc.wantBody,
				`{"id":1}`,
				http.StatusOK,
				"PUT /v1/load-balancers/1",
			)

			got, _, err := testClient.LoadBalancers.Update(t.Context(), 1, tc.req)
			require.NoError(t, err)

			assert.Equal(t, 1, got.ID)
		})
	}
}

func TestLoadBalancer_Reset(t *testing.T) {
	setup()
	defer teardown()

	setupHandler(
		t,
		"{\"type\":\"reset\"}\n",
		`{"id":1}`,
		http.StatusOK,
		"POST /v1/load-balancers/1/actions",
	)

	got, _, err := testClient.LoadBalancers.Reset(t.Context(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, got.ID)
}

func TestLoadBalancer_GetRule(t *testing.T) {
	setup()
	defer teardown()

	apiResp, err := os.ReadFile(filepath.Join(".", "testdata", "loadbalancer", "rule.json"))
	require.NoError(t, err)

	setupGetWithOptsHandler(t, apiResp, "GET /v1/load-balancers/1/rules/9d4007cc-8109-11f1-a8ee-00163e7dabb3")

	got, _, err := testClient.LoadBalancers.GetRule(
		t.Context(),
		1,
		"9d4007cc-8109-11f1-a8ee-00163e7dabb3",
		&GetOptions{Limit: 1})
	require.NoError(t, err)

	assert.Equal(t, "9d4007cc-8109-11f1-a8ee-00163e7dabb3", got.ID)
	assert.Equal(t, 8765, got.SourcePort)
	assert.Equal(t, 6443, got.DestinationPort)
	assert.Equal(t, "tcp", got.SourceProtocol)
	assert.Equal(t, "tcp", got.DestinationProtocol)
	assert.False(t, got.StickyEnabled)
	assert.Equal(t, "roundrobin", got.Balance)
	assert.Equal(t, "Active", got.Status)
}

func TestLoadBalancer_ListRules(t *testing.T) {
	setup()
	defer teardown()

	setupGetWithOptsHandler(t, []byte(`[{"id":"a"}]`), "GET /v1/load-balancers/1/rules")

	got, _, err := testClient.LoadBalancers.ListRules(t.Context(), 1, &GetOptions{Limit: 1})
	require.NoError(t, err)

	assert.Equal(t, "a", got[0].ID)
}

func TestLoadBalancer_CreateRule(t *testing.T) {
	cases := []struct {
		name        string
		req         CreateLoadBalancerRule
		wantReqBody string
	}{
		{
			name:        "no params",
			wantReqBody: "{\"source_protocol\":\"\",\"source_port\":0,\"destination_protocol\":\"\",\"destination_port\":0}\n",
		},
		{
			name: "all params",
			req: CreateLoadBalancerRule{
				SourceProtocol:      "test-proto",
				SourcePort:          1,
				DestinationProtocol: "test-proto",
				DestinationPort:     1,
				CertificateID:       "1",
				Balance:             "test-balance",
			},
			wantReqBody: "{\"source_protocol\":\"test-proto\",\"source_port\":1,\"destination_protocol\":\"test-proto\",\"destination_port\":1,\"certificate_id\":\"1\",\"balance\":\"test-balance\"}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()

			setupHandler(
				t,
				tc.wantReqBody,
				`[{"id":"a"}]`,
				http.StatusOK,
				"POST /v1/load-balancers/1/rules",
			)

			got, _, err := testClient.LoadBalancers.CreateRule(t.Context(), 1, tc.req)
			require.NoError(t, err)

			assert.Equal(t, "a", got[0].ID)
		})
	}
}

func TestLoadBalancer_DeleteRule(t *testing.T) {
	setup()
	defer teardown()

	setupHandler(t, "", "", 204, "DELETE /v1/load-balancers/1/rules/a")

	_, err := testClient.LoadBalancers.DeleteRule(t.Context(), 1, "a")
	require.NoError(t, err)
}

func TestLoadBalancer_ListServers(t *testing.T) {
	setup()
	defer teardown()

	setupGetWithOptsHandler(t, []byte(`[{"id":1}]`), "GET /v1/load-balancers/1/servers")

	got, _, err := testClient.LoadBalancers.ListServers(
		t.Context(),
		1,
		&GetOptions{Limit: 1},
	)

	require.NoError(t, err)
	assert.Equal(t, 1, got[0].ID)
}

func TestLoadBalancer_AddServer(t *testing.T) {
	cases := []struct {
		name        string
		req         AddLoadBalancerServer
		wantReqBody string
	}{
		{
			name:        "no params",
			wantReqBody: "{\"server_id\":0}\n",
		},
		{
			name: "all params",
			req: AddLoadBalancerServer{
				ServerID:     1,
				ServerWeight: 1,
			},
			wantReqBody: "{\"server_id\":1,\"server_weight\":1}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()

			setupHandler(
				t,
				tc.wantReqBody,
				`[{"id":1}]`,
				200,
				"POST /v1/load-balancers/1/servers",
			)

			got, _, err := testClient.LoadBalancers.AddServer(t.Context(), 1, tc.req)
			require.NoError(t, err)
			assert.Equal(t, 1, got[0].ID)
		})
	}
}

func TestLoadBalancer_DeleteServer(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("DELETE /v1/load-balancers/1/servers", func(w http.ResponseWriter, r *http.Request) {
		handleErr := r.ParseForm()
		require.NoError(t, handleErr)

		serverID := r.Form.Get("server_id")
		assert.Equal(t, "1", serverID)

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := testClient.LoadBalancers.DeleteServer(t.Context(), 1, 1)
	require.NoError(t, err)
}

func TestLoadBalancer_AddCertificate(t *testing.T) {
	cases := []struct {
		name        string
		req         AddLoadBalancerCertificate
		wantReqBody string
	}{
		{
			name:        "no params",
			wantReqBody: "{\"key\":\"\",\"certificate\":\"\"}\n",
		},
		{
			name: "all params",
			req: AddLoadBalancerCertificate{
				PrivateKey:  "a",
				Certificate: "a",
			},
			wantReqBody: "{\"key\":\"a\",\"certificate\":\"a\"}\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer teardown()

			setupHandler(
				t,
				tc.wantReqBody,
				`[{"id": "a"}]`,
				http.StatusOK,
				"POST /v1/load-balancers/1/certificates",
			)

			got, _, err := testClient.LoadBalancers.AddCertificate(
				t.Context(), 1, tc.req)
			require.NoError(t, err)

			assert.Equal(t, "a", got[0].ID)
		})
	}
}

func TestLoadBalancer_ListCertificates(t *testing.T) {
	setup()
	defer teardown()

	respBody, err := os.ReadFile(filepath.Join("testdata", "loadbalancer", "certificates.json"))
	require.NoError(t, err)

	setupGetWithOptsHandler(t, respBody, "GET /v1/load-balancers/1/certificates")

	got, _, err := testClient.LoadBalancers.ListCertificates(t.Context(), 1, &GetOptions{Limit: 1})
	require.NoError(t, err)

	assert.Equal(t, "86807cc4-9a2e-11f1-a8ee-00163e7dabb3", got[0].ID)
	assert.Equal(t, "kube-apiserver", got[0].CN)
	assert.Equal(t, "2026-07-10 11:45:00 +0300 EEST", got[0].Starts.String())
	assert.Equal(t, "2027-07-10 11:45:00 +0300 EEST", got[0].Expires.String())
}

func setupGetWithOptsHandler(t *testing.T, body []byte, handlePattern string) {
	mux.HandleFunc(handlePattern, func(w http.ResponseWriter, r *http.Request) {
		handleErr := r.ParseForm()
		require.NoError(t, handleErr)
		limit := r.Form.Get("limit")
		assert.Equal(t, "1", limit)

		_, handleErr = fmt.Fprint(w, string(body))
		require.NoError(t, handleErr)
	})
}

func setupHandler(t *testing.T, wantReqBody, respBody string, status int, handlePattern string) {
	mux.HandleFunc(handlePattern, func(w http.ResponseWriter, r *http.Request) {
		if wantReqBody != "" {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, wantReqBody, string(body))
		}
		w.WriteHeader(status)
		if respBody != "" {
			_, err := fmt.Fprint(w, respBody)
			require.NoError(t, err)
		}
	})
}
