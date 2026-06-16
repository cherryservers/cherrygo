package cherrygo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"
	"unicode"

	"github.com/cherryservers/cherrygo/v3/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []Server{{
		ID:   383531,
		Name: "E5-1620v4",
		Href: "/servers/383531",
		BMC: BMC{
			User:     "kuser",
			Password: "d564!h#4s8",
		},
		Hostname: "server-hostname",
		Username: "user",
		Password: "hjas345dgf54",
		Image:    "Ubuntu 18.04 64bit",
		Region: Region{
			ID:         1,
			Name:       "EU-Nord-1",
			RegionISO2: "LT",
			Href:       "/regions/1",
		},
		BGP: ServerBGP{
			Enabled: true,
		},
		State: "active",
		Tags:  map[string]string{"env": "dev"},
	}}

	mux.HandleFunc("GET /v1/projects/123/servers", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	servers, _, err := testClient.Servers.List(t.Context(), 123, nil)
	require.NoError(t, err)

	assert.Equal(t, expected, servers)
}

func TestServer_Get(t *testing.T) {
	setup()
	defer teardown()

	expected := Server{
		ID:   383531,
		Name: "E5-1620v4",
		Href: "/servers/383531",
		BMC: BMC{
			User:     "kuser",
			Password: "d564!h#4s8",
		},
		Hostname: "server-hostname",
		Username: "user",
		Password: "hjas345dgf54",
		Image:    "Ubuntu 18.04 64bit",
		Region: Region{
			ID:         1,
			Name:       "EU-Nord-1",
			RegionISO2: "LT",
			Href:       "/regions/1",
		},
		BGP: ServerBGP{
			Enabled: true,
		},
		State: "active",
		Tags:  map[string]string{"env": "dev"},
	}

	mux.HandleFunc("/v1/servers/383531", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	server, _, err := testClient.Servers.Get(t.Context(), 383531, nil)
	if err != nil {
		t.Errorf("Servers.Get returned %+v", err)
	}

	if !reflect.DeepEqual(server, expected) {
		t.Errorf("Servers.Get returned %+v, expected %+v", server, expected)
	}
}

func TestServer_PowerState(t *testing.T) {
	setup()
	defer teardown()

	expected := PowerState{
		Power: "on",
	}

	mux.HandleFunc("/v1/servers/383531", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	power, _, err := testClient.Servers.PowerState(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.PowerState returned %+v", err)
	}

	if !reflect.DeepEqual(power, expected) {
		t.Errorf("Server.PowerState returned %+v, expected %+v", power, expected)
	}
}

func TestServer_Create(t *testing.T) {
	setup()
	defer teardown()

	expected := Server{
		ID:       383531,
		Name:     "E5-1620v4",
		Href:     "/servers/383531",
		Hostname: "server-hostname",
		Image:    "Ubuntu 18.04 64bit",
		Region: Region{
			ID:         1,
			Name:       "EU-Nord-1",
			RegionISO2: "LT",
			Href:       "/regions/1",
		},
		State: "active",
		Tags:  map[string]string{"env": "dev"},
	}

	requestBody := map[string]interface{}{
		"plan":         "e5_1620v4",
		"project_id":   float64(projectID),
		"hostname":     "server-hostname",
		"image":        "ubuntu_22_04",
		"region":       "eu_nord_1",
		"ssh_keys":     []interface{}{"1", "2", "3"},
		"ip_addresses": []interface{}{"e3f75899-1db3-b794-137f-78c5ee9096af"},
		"user_data":    "dXNlcl9kYXRh",
		"tags":         map[string]interface{}{"env": "dev"},
		"spot_market":  false,
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/servers", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err = fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	tags := map[string]string{"env": "dev"}
	serverCreate := CreateServer{
		Plan:        "e5_1620v4",
		ProjectID:   projectID,
		Hostname:    "server-hostname",
		Image:       "ubuntu_22_04",
		Region:      "eu_nord_1",
		SSHKeys:     []string{"1", "2", "3"},
		IPAddresses: []string{"e3f75899-1db3-b794-137f-78c5ee9096af"},
		UserData:    "dXNlcl9kYXRh",
		Tags:        &tags,
	}

	server, _, err := testClient.Servers.Create(t.Context(), &serverCreate)
	if err != nil {
		t.Errorf("Server.Create returned %+v", err)
	}

	if !reflect.DeepEqual(server, expected) {
		t.Errorf("Server.Create returned %+v, expected %+v", server, expected)
	}
}

func TestServer_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/383531", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)

		writer.WriteHeader(http.StatusNoContent)

		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.Servers.Delete(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.Delete returned %+v", err)
	}
}

func TestServer_PowerOn(t *testing.T) {
	setup()
	defer teardown()

	expected := map[string]interface{}{
		"type": "power_on",
	}

	response := Server{
		ID: 383531,
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, expected)
		}

		jsonBytes, _ := json.Marshal(response)

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.PowerOn(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.PowerOn returned %+v", err)
	}
}

func TestServer_PowerOff(t *testing.T) {
	setup()
	defer teardown()

	expected := map[string]interface{}{
		"type": "power_off",
	}

	response := Server{
		ID: 383531,
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, expected)
		}

		jsonBytes, _ := json.Marshal(response)

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.PowerOff(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.PowerOff returned %+v", err)
	}
}

func TestServer_Reboot(t *testing.T) {
	setup()
	defer teardown()

	expected := map[string]interface{}{
		"type": "reboot",
	}

	response := Server{
		ID: 383531,
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, expected)
		}

		jsonBytes, _ := json.Marshal(response)

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.Reboot(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.Reboot returned %+v", err)
	}
}

func TestServer_EnterRescueMode(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"type":     "enter-rescue-mode",
		"password": "abcdef",
	}

	expected := Server{
		ID:     383531,
		Status: "rescue mode",
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err = fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.EnterRescueMode(t.Context(), 383531, &RescueServerFields{Password: "abcdef"})
	if err != nil {
		t.Errorf("Server.EnterRescueMode returned %+v", err)
	}
}

func TestServer_ExitRescueMode(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"type": "exit-rescue-mode",
	}

	expected := Server{
		ID:     383531,
		Status: "deployed",
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err = fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.ExitRescueMode(t.Context(), 383531)
	if err != nil {
		t.Errorf("Server.ExitRescueMode returned %+v", err)
	}
}

func TestServersClient_ResetBMCPassword(t *testing.T) {
	setup()
	defer teardown()

	expected := map[string]interface{}{
		"type": "reset-bmc-password",
	}

	response := Server{
		ID: 383531,
	}

	mux.HandleFunc("/v1/servers/383531/actions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		if err := json.NewDecoder(request.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, expected)
		}

		jsonBytes, _ := json.Marshal(response)

		_, err := fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	if _, _, err := testClient.Servers.ResetBMCPassword(t.Context(), 383531); err != nil {
		t.Errorf("Servers.ResetBMCPassword returned %+v", err)
	}
}

func TestServer_Update(t *testing.T) {
	setup()
	defer teardown()

	response := Server{
		ID:       383531,
		Name:     "prod server",
		Hostname: "cherry.prod",
		BGP: ServerBGP{
			Enabled: false,
		},
		Tags: map[string]string{"env": "dev"},
	}

	mux.HandleFunc("/v1/servers/383531", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		expected := map[string]interface{}{
			"name":     "prod server",
			"hostname": "cherry.prod",
			"tags":     map[string]interface{}{"env": "dev"},
			"bgp":      false,
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, expected)
		}

		jsonBytes, _ := json.Marshal(response)

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	tags := map[string]string{"env": "dev"}
	serverUpdate := UpdateServer{
		Tags:     &tags,
		BGP:      new(bool),
		Name:     "prod server",
		Hostname: "cherry.prod",
	}

	server, _, err := testClient.Servers.Update(t.Context(), 383531, &serverUpdate)
	if err != nil {
		t.Errorf("Server.Update returned %+v", err)
	}

	if !reflect.DeepEqual(server, response) {
		t.Errorf("Server.List returned %+v, expected %+v", server, response)
	}
}

func TestServer_UpdateBodyFieldsOmitted(t *testing.T) {
	setup()
	defer teardown()

	bod := UpdateServer{}

	mux.HandleFunc("PUT /v1/servers/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		got, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, "{}", string(got))

		_, err = fmt.Fprint(w, `{"id": 123}`)
		require.NoError(t, err)
	})

	_, _, err := testClient.Servers.Update(t.Context(), 123, &bod)
	require.NoError(t, err)
}

func TestServer_Reinstall(t *testing.T) {
	setup()
	defer teardown()

	reinstallRequest := ReinstallServerFields{
		Image:           "test-img",
		Hostname:        "test-host",
		Password:        "test-pass",
		IPXE:            "test-ipxe",
		SSHKeys:         []string{"123"},
		UserData:        "test-user-data",
		OSPartitionSize: 1,
	}

	mux.HandleFunc("POST /v1/servers/123/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		v := new(ReinstallServerFields)
		err := json.NewDecoder(r.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, reinstallRequest, *v)

		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprint(w, `{"id": 123, "deployed_image": {"slug": "test-img"}}`)
		require.NoError(t, err)
	})

	server, _, err := testClient.Servers.Reinstall(t.Context(), 123, &reinstallRequest)
	require.NoError(t, err)

	assert.Equal(t, 123, server.ID)
	assert.Equal(t, "test-img", server.DeployedImage.Slug)
}

func TestServer_ListSSHKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /v1/servers/123/ssh-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, `[{"id": 123, "label": "test-key"}]`)
		require.NoError(t, err)
	})

	keys, _, err := testClient.Servers.ListSSHKeys(t.Context(), 123, nil)
	require.NoError(t, err)

	assert.Equal(t, 123, keys[0].ID)
	assert.Equal(t, "test-key", keys[0].Label)
}

func TestServer_ListCycles(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /cycles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, `[{"id": 123, "name": "test-name", "slug": "test-slug"}]`)
		require.NoError(t, err)
	})

	cycles, _, err := testClient.Servers.ListCycles(t.Context(), nil)
	require.NoError(t, err)

	assert.Equal(t, 123, cycles[0].ID)
	assert.Equal(t, "test-name", cycles[0].Name)
	assert.Equal(t, "test-slug", cycles[0].Slug)
}

func TestServer_Upgrade(t *testing.T) {
	setup()
	defer teardown()

	want := UpgradeServer{
		ServerAction: ServerAction{Type: "upgrade"},
		Plan:         "test-plan",
	}

	mux.HandleFunc("POST /v1/servers/123/actions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		w.WriteHeader(http.StatusCreated)

		v := new(UpgradeServer)
		err := json.NewDecoder(r.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, want, *v)

		_, err = fmt.Fprint(w, `{"id": 123, "plan":{"id": 123, "slug": "test-plan"}}`)
		require.NoError(t, err)
	})

	server, _, err := testClient.Servers.Upgrade(t.Context(), 123, "test-plan")
	require.NoError(t, err)

	assert.Equal(t, 123, server.ID)
	assert.Equal(t, want.Plan, server.Plan.Slug)
}

func TestServer_WaitForStatusSucceedsWhenServerDeployedBeforeContextExpires(t *testing.T) {
	mux := http.NewServeMux()
	apiServer := httptest.NewServer(mux)
	defer apiServer.Close()

	var pollF backoff.Func = func(_ int, _ *http.Response) time.Duration {
		return 0
	}

	client, err := NewClient(WithAPIKey(
		"fakeKey"),
		WithPollBackoff(pollF),
		WithURL(apiServer.URL),
	)
	require.NoError(t, err)

	pollCount := 0
	mux.HandleFunc("GET /v1/servers/123", func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		if pollCount > 2 {
			_, err = fmt.Fprint(w, `{"id": 123, "status": "deployed"}`)
		} else {
			_, err = fmt.Fprint(w, `{"id": 123, "status": "deploying"}`)
		}
		require.NoError(t, err)
	})

	srv, resp, err := client.Servers.WaitForStatus(t.Context(), 123, StatusDeployed)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "deployed", srv.Status)
	assert.Equal(t, 3, pollCount)
}

func TestServer_WaitForStatusFailsWhenContextCancelled(t *testing.T) {
	mux := http.NewServeMux()
	apiServer := httptest.NewServer(mux)
	defer apiServer.Close()

	ctx, cancel := context.WithCancel(t.Context())

	var pollF backoff.Func = func(_ int, _ *http.Response) time.Duration {
		if ctx.Err() == nil {
			return 0
		}
		return time.Second * 5
	}

	client, err := NewClient(WithAPIKey(
		"fakeKey"),
		WithPollBackoff(pollF),
		WithURL(apiServer.URL),
	)
	require.NoError(t, err)

	pollCount := 0
	mux.HandleFunc("GET /v1/servers/123", func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		if pollCount > 2 {
			_, err = fmt.Fprint(w, `{"id": 123, "status": "deployed"}`)
		} else if pollCount > 1 {
			cancel()
			_, err = fmt.Fprint(w, `{"id": 123, "status": "deploying"}`)
		} else {
			_, err = fmt.Fprint(w, `{"id": 123, "status": "deploying"}`)
		}
		require.NoError(t, err)
	})

	srv, resp, err := client.Servers.WaitForStatus(ctx, 123, StatusDeployed)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "", srv.Status)
	assert.Equal(t, 2, pollCount)
}

func TestGeneratePasswordGeneratesValidCherryServersPasswords(t *testing.T) {
	const minLen = 16
	var failures []string

	for range 1000 {
		pw, err := GeneratePassword()
		require.NoError(t, err)

		if (len([]rune(pw))) < minLen {
			failures = append(
				failures,
				fmt.Sprintf("password %q is too short, needs to be at least %d length", pw, minLen),
			)
		}

		var hasLower, hasNonFirstUpper, hasNonLastDigit bool

		for i, r := range pw {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				failures = append(failures, fmt.Sprintf("password %q has forbidden symbol %q", pw, r))
			}

			if unicode.IsLower(r) {
				hasLower = true
			}

			if i != 0 && unicode.IsUpper(r) {
				hasNonFirstUpper = true
			}

			if i != len([]rune(pw))-1 && unicode.IsDigit(r) {
				hasNonLastDigit = true
			}
		}

		if !hasLower {
			failures = append(failures, fmt.Sprintf("password %q doesn't have a lower case letter", pw))
		}
		if !hasNonFirstUpper {
			failures = append(
				failures,
				fmt.Sprintf("password %q doesn't have an upper case letter that is not the first symbol", pw))
		}
		if !hasNonLastDigit {
			failures = append(
				failures,
				fmt.Sprintf("password %q doesn't have a digit that is not the last symbol", pw))
		}
	}

	assert.Empty(t, failures)
}

func TestServer_AllowBMCAccess(t *testing.T) {
	setup()
	defer teardown()

	apiResponse, err := os.Open(filepath.Join("testdata", "allow_bmc.json"))
	require.NoError(t, err)
	defer func() {
		_ = apiResponse.Close()
	}()

	wantReqBody := allowBMCAccess{
		ServerAction: ServerAction{Type: "create-console-access"},
		AllowedIP:    "123.123.123.123",
	}

	mux.HandleFunc("POST /v1/servers/899712/actions", func(w http.ResponseWriter, r *http.Request) {
		gotReqBody := allowBMCAccess{}
		err := json.NewDecoder(r.Body).Decode(&gotReqBody)
		require.NoError(t, err)

		assert.Equal(t, wantReqBody, gotReqBody)

		w.WriteHeader(http.StatusCreated)
		_, _ = io.ReadAll(io.TeeReader(apiResponse, w))
		_, _ = apiResponse.Seek(0, io.SeekStart)
	})

	got, resp, err := testClient.Servers.AllowBMCAccess(t.Context(), 899712, "123.123.123.123")
	require.NoError(t, err)
	require.NotNil(t, resp)

	var want Server
	err = json.NewDecoder(apiResponse).Decode(&want)
	require.NoError(t, err)

	assert.Equal(t, want.BMC, got)
}
