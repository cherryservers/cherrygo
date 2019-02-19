package cherrygo

import (
	"log"
	"strconv"
	"strings"
)

const baseIPSPath = "/v1/projects"
const endIPSPath = "ips"

// GetIPS interface metodas isgauti team'sus
type GetIPS interface {
	List(projectID int) ([]IPAddresses, *Response, error)
}

// IPSClient paveldi client
type IPSClient struct {
	client *Client
}

// List func lists teams
func (i *IPSClient) List(projectID int) ([]IPAddresses, *Response, error) {
	//root := new(teamRoot)

	projectIDString := strconv.Itoa(projectID)

	ipsPath := strings.Join([]string{baseIPSPath, projectIDString, endIPSPath}, "/")

	var trans []IPAddresses
	//resp := t.client.Bumba()
	//log.Println("\nFROM LIST1: ", root.Teams)
	resp, err := i.client.MakeRequest("GET", ipsPath, nil, &trans)
	if err != nil {
		log.Fatal(err)
	}

	return trans, resp, err
}
