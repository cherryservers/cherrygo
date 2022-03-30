package cherrygo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	mux       *http.ServeMux
	client    *Client
	server    *httptest.Server
	teamID    int
	projectID int
)

var authToken = "myToken"

func setup() {
	os.Setenv("CHERRY_AUTH_TOKEN", authToken)

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient()
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	teamID = 123
	projectID = 321
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestNewClient(t *testing.T) {
	setup()
	defer teardown()

	if client.BaseURL == nil || client.BaseURL.String() != server.URL {
		t.Errorf("NewClient BaseURL = %v, expected %v", client.BaseURL, server.URL)
	}

	if client.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", client.UserAgent, userAgent)
	}
}

func TestNewClientWithAuthVar(t *testing.T) {
	c, _ := NewClient(WithAuthToken(authToken))

	if c.AuthToken != authToken {
		t.Errorf("NewClient AuthToken = %v, expected %v", client.AuthToken, authToken)
	}
}

func TestErrorResponse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)

		fmt.Fprint(writer, `{
			"code": 400,
			"message": "Bad Request"
		}`)
	})

	_, err := client.MakeRequest(http.MethodGet, "/", nil, nil)

	expectedErr := "Error response from API: Bad Request (error code: 400)"
	if err.Error() != expectedErr {
		t.Fatalf("NewClient() expected error: %v, got %v", expectedErr, err)
	}
}

func TestCustomUserAgent(t *testing.T) {
	os.Setenv("CHERRY_AUTH_TOKEN", "token")

	ua := "testing/1.0"
	c, err := NewClient(WithUserAgent(ua))

	if err != nil {
		t.Fatalf("NewClient() unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := c.UserAgent; got != expected {
		t.Errorf("NewClient() UserAgent = %s; expected %s", got, expected)
	}
}
