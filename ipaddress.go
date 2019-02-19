package cherrygo

import (
	"log"
	"strconv"
	"strings"
)

// GetIP interface metodas isgauti team'sus
type GetIP interface {
	List(projectID int, ipID string) (IPAddresses, *Response, error)
	Create(projectID int, request *CreateIPAddress) (IPAddresses, *Response, error)
	Remove(projectID int, request *RemoveIPAddress) (IPAddresses, *Response, error)
	Update(projectID int, ipID string, request *UpdateIPAddress) (IPAddresses, *Response, error)
}

// IPClient paveldi client
type IPClient struct {
	client *Client
}

// CreateIPAddress fields for adding addition IP address
type CreateIPAddress struct {
	Type       string `json:"type", omitempty`
	Region     string `json:"region", omitempty`
	PtrRecord  string `json:"ptr_record", omitempty`
	ARecord    string `json:"a_record", omitempty`
	RoutedTo   string `json:"routed_to", omitempty`
	AssignedTo string `json:"assigned_to", omitempty`
}

// UpdateIPAddress fields for updating IP address
type UpdateIPAddress struct {
	PtrRecord  string `json:"ptr_record", omitempty`
	ARecord    string `json:"a_record", omitempty`
	RoutedTo   string `json:"routed_to", omitempty`
	AssignedTo string `json:"assigned_to", omitempty`
}

// RemoveIPAddress fields for removing IP address
type RemoveIPAddress struct {
	ID string `json:"id", omitempty`
}

// List func lists teams
func (i *IPClient) List(projectID int, ipID string) (IPAddresses, *Response, error) {
	//root := new(teamRoot)

	projectIDString := strconv.Itoa(projectID)

	ipsPath := strings.Join([]string{baseIPSPath, projectIDString, endIPSPath, ipID}, "/")

	var trans IPAddresses

	resp, err := i.client.MakeRequest("GET", ipsPath, nil, &trans)
	if err != nil {
		log.Fatal(err)
	}

	return trans, resp, err
}

// Create function orders new floating IP address
func (i *IPClient) Create(projectID int, request *CreateIPAddress) (IPAddresses, *Response, error) {

	var trans IPAddresses

	projectIDString := strconv.Itoa(projectID)

	ipAddressPath := strings.Join([]string{baseIPSPath, projectIDString, endIPSPath}, "/")

	resp, err := i.client.MakeRequest("POST", ipAddressPath, request, &trans)
	if err != nil {
		log.Fatal(err)
	}
	return trans, resp, err

}

// Update function updates existing floating IP address
func (i *IPClient) Update(projectID int, ipID string, request *UpdateIPAddress) (IPAddresses, *Response, error) {

	var trans IPAddresses

	projectIDString := strconv.Itoa(projectID)

	ipAddressPath := strings.Join([]string{baseIPSPath, projectIDString, endIPSPath, ipID}, "/")

	resp, err := i.client.MakeRequest("PUT", ipAddressPath, request, &trans)
	if err != nil {
		log.Fatal(err)
	}

	return trans, resp, err
}

// Remove function remove existing floating IP address
func (i *IPClient) Remove(projectID int, request *RemoveIPAddress) (IPAddresses, *Response, error) {

	var trans IPAddresses

	projectIDString := strconv.Itoa(projectID)

	ipAddressPath := strings.Join([]string{baseIPSPath, projectIDString, endIPSPath, request.ID}, "/")

	resp, err := i.client.MakeRequest("DELETE", ipAddressPath, request, &trans)
	if err != nil {
		log.Fatal(err)
	}
	return trans, resp, err
}
