package cherrygo

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTeam_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/teams/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer)
	})

	_, err := client.Teams.Delete(123)
	if err != nil {
		t.Errorf("Teams.Delete returned %+v", err)
	}
}
