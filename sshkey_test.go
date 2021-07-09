package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSSHKey_List(t *testing.T) {
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

		fmt.Fprint(writer, response)
	})

	sshKey, _, err := client.SSHKey.List("1", nil)
	if err != nil {
		t.Errorf("SSHKey.List returned %+v", err)
	}

	if !reflect.DeepEqual(sshKey, expected) {
		t.Errorf("SSHKey.List returned %+v, expected %+v", sshKey, expected)
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

		fmt.Fprint(writer, string(jsonBytes))
	})

	sshCreate := CreateSSHKey{
		Label: "test",
		Key:   "ssh-rsa AAAAB3NzaC1yc",
	}

	_, _, err := client.SSHKey.Create(&sshCreate)

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

		fmt.Fprint(writer)
	})

	sshDelete := DeleteSSHKey{
		ID: "1",
	}

	_, _, err := client.SSHKey.Delete(&sshDelete)

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

		fmt.Fprint(writer, string(jsonBytes))
	})

	sshUpdate := UpdateSSHKey{
		Label: "updated label",
		Key:   "ssh-rsa AAAAB3NzaC1ycupdated",
	}

	_, _, err := client.SSHKey.Update("1", &sshUpdate)

	if err != nil {
		t.Errorf("SSHKey.Update returned %+v", err)
	}
}
