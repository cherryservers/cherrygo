package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
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

		fmt.Fprint(writer, response)
	})

	ips, _, err := client.IPAddresses.List(projectID, nil)

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
		fmt.Fprint(writer, `{
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
	})

	ip, _, err := client.IPAddresses.Get(ipUID, nil)
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
		fmt.Fprint(writer, response)
	})

	ipCreate := CreateIPAddress{
		Region:    "EU-Nord-1",
		PtrRecord: "ptr",
		ARecord:   "a",
		RoutedTo:  "3ee8e5ce-4208-f437-7055-347e9e4e124e",
		Tags:      &tags,
	}

	ipAddress, _, err := client.IPAddresses.Create(projectID, &ipCreate)
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

	ipId := "e3f75899-1db3-b794-137f-78c5ee9096af"
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

	mux.HandleFunc("/v1/ips/"+ipId, func(writer http.ResponseWriter, request *http.Request) {
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
		fmt.Fprint(writer, response)
	})

	ipUpdate := UpdateIPAddress{
		PtrRecord: "ptr-new",
		ARecord:   "a-new",
		Tags:      &tags,
	}

	ipAddress, _, err := client.IPAddresses.Update(ipId, &ipUpdate)
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

	ipId := "e3f75899-1db3-b794-137f-78c5ee9096af"

	mux.HandleFunc("/v1/ips/"+ipId, func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)

		writer.WriteHeader(http.StatusNoContent)

		fmt.Fprint(writer)
	})

	_, err := client.IPAddresses.Remove(ipId)

	if err != nil {
		t.Errorf("IPAddress.Remove returned %+v", err)
	}
}
