package cherrygo

import (
	"fmt"
	"strings"
)

const baseServersPath = "/v1/projects"
const endServersPath = "servers"

// GetServers interface metodas isgauti team'sus
type GetServers interface {
	List(projectID string, opts *GetOptions) ([]Server, *Response, error)
}

// IPAddresses fields
type IPAddresses struct {
	ID            string            `json:"id,omitempty"`
	Address       string            `json:"address,omitempty"`
	AddressFamily int               `json:"address_family,omitempty"`
	Cidr          string            `json:"cidr,omitempty"`
	Gateway       string            `json:"gateway,omitempty"`
	Type          string            `json:"type,omitempty"`
	Region        Region            `json:"region,omitempty"`
	RoutedTo      RoutedTo          `json:"routed_to,omitempty"`
	AssignedTo    AssignedTo        `json:"assigned_to,omitempty"`
	TargetedTo    AssignedTo        `json:"targeted_to,omitempty"`
	Project       Project           `json:"project,omitempty"`
	PtrRecord     string            `json:"ptr_record,omitempty"`
	ARecord       string            `json:"a_record,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Href          string            `json:"href,omitempty"`
}

// RoutedTo fields
type RoutedTo struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address,omitempty"`
	AddressFamily int    `json:"address_family,omitempty"`
	Cidr          string `json:"cidr,omitempty"`
	Gateway       string `json:"gateway,omitempty"`
	Type          string `json:"type,omitempty"`
	Region        Region `json:"region,omitempty"`
}

// AssignedTo fields
type AssignedTo struct {
	ID       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Href     string  `json:"href,omitempty"`
	Hostname string  `json:"hostname,omitempty"`
	Image    string  `json:"image,omitempty"`
	Region   Region  `json:"region,omitempty"`
	State    string  `json:"state,omitempty"`
	Pricing  Pricing `json:"pricing,omitempty"`
}

// ServersClient paveldi client
type ServersClient struct {
	client *Client
}

// List func lists teams
func (s *ServersClient) List(projectID string, opts *GetOptions) ([]Server, *Response, error) {
	//root := new(teamRoot)

	//serversIDString := strconv.Itoa(projectID)

	path := strings.Join([]string{baseServersPath, projectID, endServersPath}, "/")
	pathQuery := opts.WithQuery(path)

	var trans []Server
	//resp := t.client.Bumba()
	//log.Println("\nFROM LIST1: ", root.Teams)
	resp, err := s.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
