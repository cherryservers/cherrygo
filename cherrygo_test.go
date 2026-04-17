package cherrygo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mux        *http.ServeMux
	testClient *Client
	server     *httptest.Server
	teamID     int
	projectID  int
)

var authToken = "myToken"

func setup() {
	os.Setenv("CHERRY_AUTH_TOKEN", authToken)

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	testClient, _ = NewClient()
	url, _ := url.Parse(server.URL)
	testClient.BaseURL = url
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

	if testClient.BaseURL == nil || testClient.BaseURL.String() != server.URL {
		t.Errorf("NewClient BaseURL = %v, expected %v", testClient.BaseURL, server.URL)
	}

	if testClient.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", testClient.UserAgent, userAgent)
	}
}

func TestNewClientWithAuthVar(t *testing.T) {
	c, _ := NewClient(WithAuthToken(authToken))

	if c.AuthToken != authToken {
		t.Errorf("NewClient AuthToken = %v, expected %v", testClient.AuthToken, authToken)
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

	req, err := testClient.NewRequest(t.Context(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	_, err = testClient.Do(req, nil)

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

func TestDebug(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprintf(w, `{"id": 1}`)
		require.NoError(t, err)
	})

	buf := &bytes.Buffer{}

	c, err := NewClient(WithAuthToken("HIDDEN"), WithDebug(buf), WithURL(server.URL))
	require.NoError(t, err)

	req, err := c.NewRequest(t.Context(), http.MethodGet, "/test", nil)
	require.NoError(t, err)

	_, err = c.Do(req, &struct {
		ID int `json:"id"`
	}{})
	require.NoError(t, err)

	got, err := io.ReadAll(buf)
	require.NoError(t, err)

	assert.Contains(t, string(got), "REQUEST")
	assert.Contains(t, string(got), "RESPONSE")
	assert.NotContains(t, string(got), "HIDDEN")
}
