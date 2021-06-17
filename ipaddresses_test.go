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

	expected := []IPAddresses{
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

	ips, _, err := client.IPAddresses.List(strconv.Itoa(projectID))

	if err != nil {
		t.Errorf("IPAddresses.List returned %+v", err)
	}

	if !reflect.DeepEqual(ips, expected) {
		t.Errorf("IPAddresses.List returned %+v, expected %+v", ips, expected)
	}
}
