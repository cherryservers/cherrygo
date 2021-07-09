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

	expected := []Projects{
		{
			ID:   321,
			Name: "My Project",
			Href: "/projects/321",
		},
		{
			ID:   322,
			Name: "My New Project",
			Href: "/projects/322",
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
