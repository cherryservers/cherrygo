package cherrygo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStorageClient_List(t *testing.T) {
	cherryServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(rw, `{
  "id": 0,
  "name": "string",
  "href": "string",
  "size": 0,
  "allow_edit_size": true,
  "unit": "string",
  "description": "string",
  "vlan_id": "string",
  "vlan_ip": "string",
  "initiator": "string",
  "discovery_ip": "string"
}`)
	}))
	defer cherryServer.Close()
	client, err := NewClientBase(cherryServer.Client(), "test")
	if err != nil {
		t.Error(err)
	}
	u, _ := url.Parse(cherryServer.URL)
	client.BaseURL = u
	storageClient := StorageClient{client}
	ebs, _, err := storageClient.List("10", "20")
	if err != nil {
		t.Error(err)
	}
	if ebs.ID != 0 {
		t.Errorf("Expected 0 got %d", ebs.ID)
	}
}
