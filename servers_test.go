package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestServer_Get(t *testing.T) {
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
			RegionIso2: "LT",
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

		fmt.Fprint(writer, response)
	})

	server, _, err := client.Servers.Get(383531, nil)
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

		fmt.Fprint(writer, response)
	})

	power, _, err := client.Servers.PowerState(383531)
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
			RegionIso2: "LT",
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

		fmt.Fprint(writer, response)
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

	server, _, err := client.Servers.Create(&serverCreate)

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

		fmt.Fprint(writer)
	})

	_, _, err := client.Servers.Delete(383531)

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

		fmt.Fprint(writer, string(jsonBytes))
	})

	_, _, err := client.Servers.PowerOn(383531)

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

		fmt.Fprint(writer, string(jsonBytes))
	})

	_, _, err := client.Servers.PowerOff(383531)

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

		fmt.Fprint(writer, string(jsonBytes))
	})

	_, _, err := client.Servers.Reboot(383531)

	if err != nil {
		t.Errorf("Server.Reboot returned %+v", err)
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

		fmt.Fprint(writer, string(jsonBytes))
	})

	if _, _, err := client.Servers.ResetBMCPassword(383531); err != nil {
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

		fmt.Fprint(writer, string(jsonBytes))
	})

	tags := map[string]string{"env": "dev"}
	serverUpdate := UpdateServer{
		Tags:     &tags,
		Bgp:      false,
		Name:     "prod server",
		Hostname: "cherry.prod",
	}

	server, _, err := client.Servers.Update(383531, &serverUpdate)

	if err != nil {
		t.Errorf("Server.Update returned %+v", err)
	}

	if !reflect.DeepEqual(server, response) {
		t.Errorf("Server.List returned %+v, expected %+v", server, response)
	}
}
