package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestProjects_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []Project{
		{
			ID:   321,
			Name: "My Project",
			Href: "/projects/321",
			Bgp: ProjectBGP{
				Enabled:  true,
				LocalASN: 123,
			},
		},
		{
			ID:   322,
			Name: "My New Project",
			Href: "/projects/322",
			Bgp: ProjectBGP{
				Enabled:  false,
				LocalASN: 0,
			},
		},
	}

	mux.HandleFunc("/v1/teams/"+strconv.Itoa(teamID)+"/projects", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		fmt.Fprint(writer, response)
	})

	projects, _, err := client.Projects.List(teamID, nil)

	if err != nil {
		t.Errorf("Projects.List returned %+v", err)
	}

	if !reflect.DeepEqual(projects, expected) {
		t.Errorf("Projects.List returned %+v, expected %+v", projects, expected)
	}
}

func TestProject_Get(t *testing.T) {
	setup()
	defer teardown()

	expected := Project{
		ID:   projectID,
		Name: "My Project",
		Href: "/projects/321",
		Bgp: ProjectBGP{
			Enabled:  true,
			LocalASN: 123,
		},
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID), func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `{
			"id": 321,
			"name": "My Project",
			"href": "/projects/321",
			"bgp": {
				"enabled": true,
				"local_asn": 123
			}
		}`)
	})

	project, _, err := client.Projects.Get(projectID, nil)
	if err != nil {
		t.Errorf("Project.List returned %+v", err)
	}

	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Project.List returned %+v, expected %+v", project, expected)
	}
}

func TestProject_Create(t *testing.T) {
	setup()
	defer teardown()

	expected := Project{
		ID:   322,
		Name: "My Custom Project",
		Href: "/projects/322",
		Bgp: ProjectBGP{
			Enabled:  true,
			LocalASN: 123,
		},
	}

	requestBody := map[string]interface{}{
		"name": "My Custom Project",
	}

	mux.HandleFunc("/v1/teams/"+strconv.Itoa(teamID)+"/projects", func(writer http.ResponseWriter, request *http.Request) {
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

	projectCreate := CreateProject{
		Name: "My Custom Project",
	}

	project, _, err := client.Projects.Create(teamID, &projectCreate)
	if err != nil {
		t.Errorf("Project.Create returned %+v", err)
	}

	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Project.Create returned %+v, expected %+v", project, expected)
	}
}

func TestProject_Update(t *testing.T) {
	setup()
	defer teardown()

	requestBody := map[string]interface{}{
		"name": "My Updated Project",
		"bgp":  true,
	}

	mux.HandleFunc("/v1/projects/"+strconv.Itoa(projectID), func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)

		var v map[string]interface{}
		err := json.NewDecoder(request.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, requestBody) {
			t.Errorf("Request body\n sent %#v\n expected %#v", v, requestBody)
		}

		fmt.Fprint(writer, `{"id": 321}`)
	})

	name := "My Updated Project"
	bgp := true

	projectUpdate := UpdateProject{
		Name: &name,
		Bgp:  &bgp,
	}

	_, _, err := client.Projects.Update(projectID, &projectUpdate)
	if err != nil {
		t.Errorf("Project.Update returned %+v", err)
	}
}

func TestProject_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/projects/321", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)

		writer.WriteHeader(http.StatusNoContent)

		fmt.Fprint(writer)
	})

	_, err := client.Projects.Delete(projectID)

	if err != nil {
		t.Errorf("Project.Delete returned %+v", err)
	}
}

func TestProjectSSHKeys_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []SSHKey{
		{
			ID:          1,
			Label:       "test",
			Key:         "ssh-rsa AAAAB3NzaC1yc",
			Fingerprint: "fb:f0:21:33:e9:26:y3:2e:2e:b4:5c:8a:a6:26:64:ae",
			Updated:     "2021-04-20 16:40:54",
			Created:     "2021-04-20 13:40:43",
			Href:        "/ssh-keys/1",
		},
	}

	mux.HandleFunc("/v1/projects/123/ssh-keys", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)

		jsonBytes, _ := json.Marshal(expected)
		response := string(jsonBytes)

		fmt.Fprint(writer, response)
	})

	sshKeys, _, err := client.Projects.ListSSHKeys(123, nil)
	if err != nil {
		t.Errorf("Projects.ListSSHKeys returned %+v", err)
	}

	if !reflect.DeepEqual(sshKeys, expected) {
		t.Errorf("Projects.ListSSHKeys returned %+v, expected %+v", sshKeys, expected)
	}
}
