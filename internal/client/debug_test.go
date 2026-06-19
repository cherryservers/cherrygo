package client_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cherryservers/cherrygo/v3/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprintf(w, `{"id": 1}`)
		require.NoError(t, err)
	}))
	defer ts.Close()

	buf := &bytes.Buffer{}

	client := client.New(client.WithDebug(buf))
	req, err := http.NewRequest("GET", ts.URL, strings.NewReader("test"))
	require.NoError(t, err)

	req.Header.Add("Authorization", "SECRET")

	_, err = client.Do(req)
	require.NoError(t, err)

	got, err := io.ReadAll(buf)
	require.NoError(t, err)

	assert.Contains(t, string(got), "REQUEST")
	assert.Contains(t, string(got), "RESPONSE")
	assert.Contains(t, string(got), "***")
	assert.NotContains(t, string(got), "SECRET")
}
