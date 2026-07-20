package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []BlockStorage{{
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
		VLANID:      "1",
		VLANIP:      "1.1.1.1",
		Initiator:   "com.cherryservers:initiator",
		DiscoveryIP: "1.1.1.1",
	}}

	mux.HandleFunc("GET /v1/projects/123/storages", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `[{
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
		}]`)
		require.NoError(t, err)
	})

	storages, _, err := testClient.Storages.List(t.Context(), 123, nil)
	if err != nil {
		t.Errorf("Storages.List returned %+v", err)
	}

	if !reflect.DeepEqual(storages, expected) {
		t.Errorf("Storages.List returned %+v, expected %+v", storages, expected)
	}
}

func TestStorage_Get(t *testing.T) {
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
		VLANID:      "1",
		VLANIP:      "1.1.1.1",
		Initiator:   "com.cherryservers:initiator",
		DiscoveryIP: "1.1.1.1",
	}

	mux.HandleFunc("/v1/storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `{
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
		require.NoError(t, err)
	})

	storage, _, err := testClient.Storages.Get(t.Context(), 123, nil)
	if err != nil {
		t.Errorf("Storages.Get returned %+v", err)
	}

	if !reflect.DeepEqual(storage, expected) {
		t.Errorf("Storages.Get returned %+v, expected %+v", storage, expected)
	}
}

func TestStorage_Create(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
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

		_, err = fmt.Fprint(writer, `{"id": 123}`)
		require.NoError(t, err)
	})

	createStorage := CreateStorage{
		Description: "desc",
		Size:        521,
		Region:      "EU-Nord-1",
	}

	_, _, err := testClient.Storages.Create(t.Context(), 321, &createStorage)
	if err != nil {
		t.Errorf("Storages.Create returned %+v", err)
	}
}

func TestStorage_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.Storages.Delete(t.Context(), 123)
	if err != nil {
		t.Errorf("Storages.Delete returned %+v", err)
	}
}

func TestStorage_Attach(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"attach_to": float64(1234),
	}

	mux.HandleFunc("/v1/storages/123/attachments", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		_, err = fmt.Fprint(writer, `{"id": 123}`)
		require.NoError(t, err)
	})

	attachStorage := AttachTo{
		AttachTo: 1234,
	}

	_, _, err := testClient.Storages.Attach(t.Context(), 123, &attachStorage)
	if err != nil {
		t.Errorf("Storages.Attach returned %+v", err)
	}
}

func TestStorage_Detach(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/storages/123/attachments", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.Storages.Detach(t.Context(), 123)
	if err != nil {
		t.Errorf("Storages.Detach returned %+v", err)
	}
}

func TestStorage_Update(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"size":        float64(500),
		"description": "volume 1",
	}

	mux.HandleFunc("/v1/storages/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		_, err = fmt.Fprint(writer, `{"id": 123}`)
		require.NoError(t, err)
	})

	updateStorage := UpdateStorage{
		Size:        500,
		Description: "volume 1",
	}

	_, _, err := testClient.Storages.Update(t.Context(), 123, &updateStorage)
	if err != nil {
		t.Errorf("Storages.Update returned %+v", err)
	}
}
