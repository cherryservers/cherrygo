package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestStorage_List(t *testing.T) {
	setup()
	defer teardown()

	expected := BlockStorage{
		ID:            123,
		Name:          "name",
		Href:          "/storages/123",
		Size:          256,
		AllowEditSize: true,
		Unit:          "GB",
		Description:   "string",
		AttachedTo: AttachedTo{
			Href: "/servers/1",
		},
		VlanID:      "1",
		VlanIP:      "1.1.1.1",
		Initiator:   "com.cherryservers:initiator",
		DiscoveryIP: "1.1.1.1",
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `{
			"id": 123,
			"name": "name",
			"href": "/storages/123",
			"size": 256,
			"allow_edit_size": true,
			"unit": "GB",
			"description": "string",
			"attached_to": {
				"href": "/servers/1"
			},
			"vlan_id": "1",
			"vlan_ip": "1.1.1.1",
			"initiator": "com.cherryservers:initiator",
			"discovery_ip": "1.1.1.1"
		}`)
	})

	storage, _, err := client.Storage.List(strconv.Itoa(projectID), "123")
	if err != nil {
		t.Errorf("Storage.List returned %+v", err)
	}

	if !reflect.DeepEqual(storage, expected) {
		t.Errorf("Storage.List returned %+v, expected %+v", storage, expected)
	}
}

func TestStorage_Create(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"project_id":  "321",
		"description": "desc",
		"size":        521.00,
		"region":      "EU-Nord-1",
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/storages", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		fmt.Fprint(writer, `{"id": 123}`)
	})

	createStorage := CreateStorage{
		ProjectID:   "321",
		Description: "desc",
		Size:        521,
		Region:      "EU-Nord-1",
	}

	_, _, err := client.Storage.Create(&createStorage)
	if err != nil {
		t.Errorf("Storage.List returned %+v", err)
	}
}

func TestStorage_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer)
	})

	deleteStorage := DeleteStorage{
		ProjectID: "321",
		StorageID: "123",
	}

	_, err := client.Storage.Delete(&deleteStorage)
	if err != nil {
		t.Errorf("Storage.Delete returned %+v", err)
	}
}

func TestStorage_Attach(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"project_id": "321",
		"storage_id": "123",
		"attach_to":  float64(1234),
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/storages/123/attachments", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		fmt.Fprint(writer, `{"id": 123}`)
	})

	attachStorage := AttachTo{
		ProjectID: "321",
		StorageID: "123",
		AttachTo:  1234,
	}

	_, _, err := client.Storage.Attach(&attachStorage)
	if err != nil {
		t.Errorf("Storage.Attach returned %+v", err)
	}
}

func TestStorage_Detach(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID)+"/storages/123/attachments", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer)
	})

	detachStorage := DetachFrom{
		ProjectID: "321",
		StorageID: "123",
	}

	_, err := client.Storage.Detach(&detachStorage)
	if err != nil {
		t.Errorf("Storage.Detach returned %+v", err)
	}
}
