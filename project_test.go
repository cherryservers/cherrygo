package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestProject_List(t *testing.T) {
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

	project, _, err := client.Project.List(strconv.Itoa(projectID))
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

	project, _, err := client.Project.Create(teamID, &projectCreate)
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

	projectUpdate := UpdateProject{
		Name: "My Updated Project",
		Bgp:  true,
	}

	_, _, err := client.Project.Update(strconv.Itoa(projectID), &projectUpdate)
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

	projectDelete := DeleteProject{
		ID: strconv.Itoa(projectID),
	}

	_, _, err := client.Project.Delete(strconv.Itoa(projectID), &projectDelete)

	if err != nil {
		t.Errorf("Project.Delete returned %+v", err)
	}
}
