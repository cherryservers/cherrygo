package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHKey_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []SSHKey{{
		ID:          1,
		Label:       "test",
		Key:         "ssh-rsa AAAAB3NzaC1yc",
		Fingerprint: "fb:f0:21:33:e9:26:y3:2e:2e:b4:5c:8a:a6:26:64:ae",
		Updated:     "2021-04-20 16:40:54",
		Created:     "2021-04-20 13:40:43",
		Href:        "/ssh-keys/1",
	}}

	mux.HandleFunc("GET /v1/ssh-keys", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	sshKey, _, err := testClient.SSHKeys.List(t.Context(), nil)
	require.NoError(t, err)

	assert.Equal(t, expected, sshKey)
}

func TestSSHKey_Get(t *testing.T) {
	setup()
	defer teardown()

	expected := SSHKey{
		ID:          1,
		Label:       "test",
		Key:         "ssh-rsa AAAAB3NzaC1yc",
		Fingerprint: "fb:f0:21:33:e9:26:y3:2e:2e:b4:5c:8a:a6:26:64:ae",
		Updated:     "2021-04-20 16:40:54",
		Created:     "2021-04-20 13:40:43",
		Href:        "/ssh-keys/1",
	}

	mux.HandleFunc("/v1/ssh-keys/1", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		_, err := fmt.Fprint(writer, response)
		require.NoError(t, err)
	})

	sshKey, _, err := testClient.SSHKeys.Get(t.Context(), 1, nil)
	if err != nil {
		t.Errorf("SSHKey.Get returned %+v", err)
	}

	if !reflect.DeepEqual(sshKey, expected) {
		t.Errorf("SSHKey.Get returned %+v, expected %+v", sshKey, expected)
	}
}

func TestSSHKey_Create(t *testing.T) {
	setup()
	defer teardown()

	expected := SSHKey{
		ID:          1,
		Label:       "test",
		Key:         "ssh-rsa AAAAB3NzaC1yc",
		Fingerprint: "fb:f0:21:33:e9:26:y3:2e:2e:b4:5c:8a:a6:26:64:ae",
		Updated:     "2021-04-20 16:40:54",
		Created:     "2021-04-20 13:40:43",
		Href:        "/ssh-keys/1",
	}

	requestBody := map[string]interface{}{
		"label": "test",
		"key":   "ssh-rsa AAAAB3NzaC1yc",
	}

	mux.HandleFunc("/v1/ssh-keys", func(writer http.ResponseWriter, request *http.Request) {
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

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	sshCreate := CreateSSHKey{
		Label: "test",
		Key:   "ssh-rsa AAAAB3NzaC1yc",
	}

	_, _, err := testClient.SSHKeys.Create(t.Context(), &sshCreate)
	if err != nil {
		t.Errorf("SSHKey.Create returned %+v", err)
	}
}

func TestSSHKey_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/ssh-keys/1", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)

		writer.WriteHeader(http.StatusNoContent)

		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, _, err := testClient.SSHKeys.Delete(t.Context(), 1)
	if err != nil {
		t.Errorf("SSHKey.Delete returned %+v", err)
	}
}

func TestSSHKey_Update(t *testing.T) {
	setup()
	defer teardown()

	expected := SSHKey{
		ID:    1,
		Label: "updated label",
	}

	requestBody := map[string]interface{}{
		"label": "updated label",
		"key":   "ssh-rsa AAAAB3NzaC1ycupdated",
	}

	mux.HandleFunc("/v1/ssh-keys/1", func(writer http.ResponseWriter, request *http.Request) {
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

		_, err = fmt.Fprint(writer, string(jsonBytes))
		require.NoError(t, err)
	})

	label := "updated label"
	key := "ssh-rsa AAAAB3NzaC1ycupdated"
	sshUpdate := UpdateSSHKey{
		Label: &label,
		Key:   &key,
	}

	_, _, err := testClient.SSHKeys.Update(t.Context(), 1, &sshUpdate)
	if err != nil {
		t.Errorf("SSHKey.Update returned %+v", err)
	}
}
