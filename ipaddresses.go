package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseIPPath = "/v1/ips"

// IPAddressesService is an interface for interfacing with the the Server endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Ip-Addresses
type IPAddressesService interface {
	List(ctx context.Context, projectID int, opts *GetOptions) ([]IPAddress, *Response, error)
	Get(ctx context.Context, ipID string, opts *GetOptions) (IPAddress, *Response, error)
	Create(ctx context.Context, projectID int, request *CreateIPAddress) (IPAddress, *Response, error)
	Remove(ctx context.Context, ipID string) (*Response, error)
	Update(ctx context.Context, ipID string, request *UpdateIPAddress) (IPAddress, *Response, error)
	Assign(ctx context.Context, ipID string, request *AssignIPAddress) (IPAddress, *Response, error)
	Unassign(ctx context.Context, ipID string) (*Response, error)
}

// IPAddress data.
type IPAddress struct {
	ID            string             `json:"id,omitempty"`
	Address       string             `json:"address,omitempty"`
	AddressFamily int                `json:"address_family,omitempty"`
	CIDR          string             `json:"cidr,omitempty"`
	Gateway       string             `json:"gateway,omitempty"`
	Type          string             `json:"type,omitempty"`
	Region        Region             `json:"region,omitempty"`
	RoutedTo      RoutedTo           `json:"routed_to,omitempty"`
	AssignedTo    AssignedTo         `json:"assigned_to,omitempty"`
	TargetedTo    AssignedTo         `json:"targeted_to,omitempty"`
	Project       Project            `json:"project,omitempty"`
	PTRRecord     string             `json:"ptr_record,omitempty"`
	ARecord       string             `json:"a_record,omitempty"`
	Tags          *map[string]string `json:"tags,omitempty"`
	DDoSScrubbing bool               `json:"ddos_scrubbing,omitempty"`
	Href          string             `json:"href,omitempty"`
}

// RoutedTo fields
type RoutedTo struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address,omitempty"`
	AddressFamily int    `json:"address_family,omitempty"`
	CIDR          string `json:"cidr,omitempty"`
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

// IPsClient makes IP address related API requests.
type IPsClient struct {
	client *Client
}

// CreateIPAddress fields for adding addition IP address
type CreateIPAddress struct {
	Region        string             `json:"region,omitempty"`
	PTRRecord     string             `json:"ptr_record,omitempty"`
	ARecord       string             `json:"a_record,omitempty"`
	RoutedTo      string             `json:"routed_to,omitempty"`
	AssignedTo    string             `json:"assigned_to,omitempty"`
	TargetedTo    string             `json:"targeted_to,omitempty"`
	Tags          *map[string]string `json:"tags,omitempty"`
	DDoSScrubbing bool               `json:"ddos_scrubbing,omitempty"`
}

// UpdateIPAddress fields for updating IP address
type UpdateIPAddress struct {
	PTRRecord  string             `json:"ptr_record,omitempty"`
	ARecord    string             `json:"a_record,omitempty"`
	RoutedTo   string             `json:"routed_to,omitempty"`
	AssignedTo string             `json:"assigned_to,omitempty"`
	TargetedTo string             `json:"targeted_to,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}

// AssignIPAddress is the IP address assignment request body.
// Subnet type IP addresses can be only assigned to a server.
// Floating IP address can be assigned directly to a server or routed to subnet type IP address.
type AssignIPAddress struct {
	ServerID int    `json:"targeted_to,omitempty"`
	RoutedTo string `json:"routed_to,omitempty"`
}

// List func lists ip addresses
func (i *IPsClient) List(ctx context.Context, projectID int, opts *GetOptions) ([]IPAddress, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/ips", baseProjectPath, projectID))
	var trans []IPAddress

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

// Get IP address.
func (i *IPsClient) Get(ctx context.Context, ipID string, opts *GetOptions) (IPAddress, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s", baseIPPath, ipID))
	var trans IPAddress

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return IPAddress{}, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

// Create function orders new floating IP address
func (i *IPsClient) Create(ctx context.Context, projectID int, request *CreateIPAddress) (IPAddress, *Response, error) {
	var trans IPAddress
	path := fmt.Sprintf("%s/%d/ips", baseProjectPath, projectID)

	req, err := i.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return IPAddress{}, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

// Update function updates existing IP address
func (i *IPsClient) Update(ctx context.Context, ipID string, request *UpdateIPAddress) (IPAddress, *Response, error) {
	var trans IPAddress
	path := fmt.Sprintf("%s/%s", baseIPPath, ipID)

	req, err := i.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return IPAddress{}, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

// Remove function removes existing project IP address
func (i *IPsClient) Remove(ctx context.Context, ipID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", baseIPPath, ipID)

	req, err := i.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req, nil)
	return resp, err
}

// Assign IP address.
func (i *IPsClient) Assign(ctx context.Context, ipID string, request *AssignIPAddress) (IPAddress, *Response, error) {
	var trans IPAddress
	path := fmt.Sprintf("%s/%s", baseIPPath, ipID)

	req, err := i.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return IPAddress{}, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

// Unassign IP address.
func (i *IPsClient) Unassign(ctx context.Context, ipID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", baseIPPath, ipID)
	request := UpdateIPAddress{
		TargetedTo: "0",
	}

	req, err := i.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req, nil)
	return resp, err
}
