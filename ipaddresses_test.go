package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIpAddresses_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []IPAddress{
		{
			ID:            "e3f75899-1db3-b794-137f-78c5ee9096af",
			Address:       "5.199.171.0",
			AddressFamily: 4,
			Cidr:          "5.199.171.0/32",
			Type:          "floating-ip",
			RoutedTo: RoutedTo{
				ID:            "3ee8e5ce-4208-f437-7055-347e9e4e124e",
				Address:       "188.214.132.158",
				AddressFamily: 4,
				Cidr:          "188.214.132.128/25",
				Gateway:       "188.214.132.129",
				Type:          "primary-ip",
				Region: Region{
					ID:         1,
					Name:       "EU-Nord-1",
					RegionIso2: "LT",
					Href:       "/regions/1",
				},
			},
			PtrRecord: "ptr-a",
			ARecord:   "a-a",
			Href:      "/ips/e3f75899-1db3-b794-137f-78c5ee9096af",
		},
		{
			ID:            "e84d6ae8-573c-ecf9-a01d-afc57f95e910",
			Address:       "5.199.171.1",
			AddressFamily: 4,
			Cidr:          "5.199.171.1/32",
			Type:          "subnet",
			AssignedTo: AssignedTo{
				ID:       383531,
				Name:     "E5-1620v4",
				Href:     "/servers/383531",
				Hostname: "server-hostname",
				Image:    "Ubuntu 18.04 64bit",
				State:    "active",
			},
			PtrRecord: "ptr-b",
			ARecord:   "a-b",
			Href:      "/ips/e84d6ae8-573c-ecf9-a01d-afc57f95e910",
		},
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/ips", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	ips, _, err := testClient.IPAddresses.List(t.Context(), projectID, nil)
	if err != nil {
		t.Errorf("IPAddresses.List returned %+v", err)
	}

	if !reflect.DeepEqual(ips, expected) {
		t.Errorf("IPAddresses.List returned %+v, expected %+v", ips, expected)
	}
}

func TestIpAddress_Get(t *testing.T) {
	setup()
	defer teardown()

	ipUID := "af28c90d-e4a3-4999-a81a-7ae0d9e8eb88"
	expected := IPAddress{
		ID:            "e3f75899-1db3-b794-137f-78c5ee9096af",
		Address:       "5.199.171.0",
		AddressFamily: 4,
		Cidr:          "5.199.171.0/32",
		Type:          "floating-ip",
		RoutedTo: RoutedTo{
			ID:            "3ee8e5ce-4208-f437-7055-347e9e4e124e",
			Address:       "188.214.132.158",
			AddressFamily: 4,
			Cidr:          "188.214.132.128/25",
			Gateway:       "188.214.132.129",
			Type:          "primary-ip",
			Region: Region{
				ID:         1,
				Name:       "EU-Nord-1",
				RegionIso2: "LT",
				Location:   "Lithuania, Vilnius",
				Href:       "/regions/1",
			},
		},
		AssignedTo: AssignedTo{
			ID:       383531,
			Name:     "E5-1620v4",
			Href:     "/servers/383531",
			Hostname: "server-hostname",
			Image:    "Ubuntu 18.04 64bit",
			State:    "active",
		},
		PtrRecord: "ptr-r",
		ARecord:   "a-r",
		Href:      "/ips/e3f75899-1db3-b794-137f-78c5ee9096af",
	}

	mux.HandleFunc("/v1/ips/"+ipUID, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `{
			"id":"e3f75899-1db3-b794-137f-78c5ee9096af",
			"address":"5.199.171.0",
			"address_family":4,
			"cidr":"5.199.171.0/32",
			"type":"floating-ip",
			"routed_to":{
			   "id":"3ee8e5ce-4208-f437-7055-347e9e4e124e",
			   "address":"188.214.132.158",
			   "address_family":4,
			   "cidr":"188.214.132.128/25",
			   "gateway":"188.214.132.129",
			   "type":"primary-ip",
			   "region":{
				  "id":1,
				  "name":"EU-Nord-1",
				  "region_iso_2":"LT",
				  "href":"/regions/1",
				  "location":"Lithuania, Vilnius"
			   },
			   "ptr_record":"primry-ip-ptr",
			   "a_record":"primry-ip-a-record",
			   "ddos_scrubbing":false
			},
			"assigned_to":{
			   "id":383531,
			   "name":"E5-1620v4",
			   "href":"/servers/383531",
			   "hostname":"server-hostname",
			   "image":"Ubuntu 18.04 64bit",
			   "spot_instance":false,
			   "state":"active",
			   "created_at":"2021-06-14T06:31:09+00:00"
			},
			"ptr_record":"ptr-r",
			"a_record":"a-r",
			"ddos_scrubbing":false,
			"href":"/ips/e3f75899-1db3-b794-137f-78c5ee9096af"
		 }`)

		require.NoError(t, err)
	})

	ip, _, err := testClient.IPAddresses.Get(t.Context(), ipUID, nil)
	if err != nil {
		t.Errorf("IPAddress.List returned %+v", err)
	}

	if !reflect.DeepEqual(ip, expected) {
		t.Errorf("IPAddress.List returned %+v, expected %+v", ip, expected)
	}
}

func TestIpAddress_Create(t *testing.T) {
	setup()
	defer teardown()

	tags := map[string]string{"env": "dev"}
	expected := IPAddress{
		ID:            "e3f75899-1db3-b794-137f-78c5ee9096af",
		Address:       "5.199.171.0",
		AddressFamily: 4,
		Cidr:          "5.199.171.0/32",
		Type:          "floating-ip",
		Region: Region{
			ID:         1,
			Name:       "EU-Nord-1",
			RegionIso2: "LT",
			Href:       "/regions/1",
		},
		RoutedTo: RoutedTo{
			ID:            "3ee8e5ce-4208-f437-7055-347e9e4e124es",
			Address:       "188.214.132.158",
			AddressFamily: 4,
			Cidr:          "188.214.132.128/25",
			Gateway:       "188.214.132.129",
			Type:          "primary-ip",
			Region: Region{
				ID:         1,
				Name:       "EU-Nord-1",
				RegionIso2: "LT",
				Href:       "/regions/1",
			},
		},
		PtrRecord: "ptr-r",
		ARecord:   "a-r",
		Tags:      &tags,
		Href:      "/ips/e3f75899-1db3-b794-137f-78c5ee9096af",
	}

	requestBody := map[string]interface{}{
		"region":     "EU-Nord-1",
		"ptr_record": "ptr",
		"a_record":   "a",
		"routed_to":  "3ee8e5ce-4208-f437-7055-347e9e4e124e",
		"tags":       map[string]interface{}{"env": "dev"},
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/ips", func(writer http.ResponseWriter, request *http.Request) {
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

	ipCreate := CreateIPAddress{
		Region:    "EU-Nord-1",
		PtrRecord: "ptr",
		ARecord:   "a",
		RoutedTo:  "3ee8e5ce-4208-f437-7055-347e9e4e124e",
		Tags:      &tags,
	}

	ipAddress, _, err := testClient.IPAddresses.Create(t.Context(), projectID, &ipCreate)
	if err != nil {
		t.Errorf("IPAddress.Create returned %+v", err)
	}

	if !reflect.DeepEqual(ipAddress, expected) {
		t.Errorf("IPAddress.Create returned %+v, expected %+v", ipAddress, expected)
	}
}

func TestIpAddress_Update(t *testing.T) {
	setup()
	defer teardown()

	ipID := "e3f75899-1db3-b794-137f-78c5ee9096af"
	tags := map[string]string{"env": "dev"}
	expected := IPAddress{
		ID:        "e3f75899-1db3-b794-137f-78c5ee9096af",
		PtrRecord: "ptr-new",
		ARecord:   "a-new",
		Tags:      &tags,
	}

	requestBody := map[string]interface{}{
		"ptr_record": "ptr-new",
		"a_record":   "a-new",
		"tags":       map[string]interface{}{"env": "dev"},
	}

	mux.HandleFunc("/v1/ips/"+ipID, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)

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

	ipUpdate := UpdateIPAddress{
		PtrRecord: "ptr-new",
		ARecord:   "a-new",
		Tags:      &tags,
	}

	ipAddress, _, err := testClient.IPAddresses.Update(t.Context(), ipID, &ipUpdate)
	if err != nil {
		t.Errorf("IPAddress.Update returned %+v", err)
	}

	if !reflect.DeepEqual(ipAddress, expected) {
		t.Errorf("IPAddress.Update returned %+v, expected %+v", ipAddress, expected)
	}
}

func TestIpAddress_Delete(t *testing.T) {
	setup()
	defer teardown()

	ipID := "e3f75899-1db3-b794-137f-78c5ee9096af"

	mux.HandleFunc("/v1/ips/"+ipID, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)

		writer.WriteHeader(http.StatusNoContent)

		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.IPAddresses.Remove(t.Context(), ipID)
	if err != nil {
		t.Errorf("IPAddress.Remove returned %+v", err)
	}
}

func TestIPAddress_Assign(t *testing.T) {
	setup()
	defer teardown()

	assignRequest := AssignIPAddress{
		ServerID: 123,
	}

	mux.HandleFunc("PUT /v1/ips/abc123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		v := new(AssignIPAddress)
		err := json.NewDecoder(r.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, assignRequest, *v)

		_, err = fmt.Fprint(w, `{"id": "abc123", "address": "127.0.0.1"}`)
		require.NoError(t, err)
	})

	ip, _, err := testClient.IPAddresses.Assign(t.Context(), "abc123", &assignRequest)
	require.NoError(t, err)

	assert.Equal(t, "abc123", ip.ID)
	assert.Equal(t, "127.0.0.1", ip.Address)
}

func TestIPAddress_Unassign(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /v1/ips/abc123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		v := new(UpdateIPAddress)
		err := json.NewDecoder(r.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, "0", v.TargetedTo)

		_, err = fmt.Fprint(w, `{"id": "abc123", "address": "127.0.0.1"}`)
		require.NoError(t, err)
	})

	_, err := testClient.IPAddresses.Unassign(t.Context(), "abc123")
	require.NoError(t, err)
}
