package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const (
	baseServerPath = "/v1/servers"
	endServersPath = "servers"
)

// ServersService is an interface for interfacing with the Server endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Servers
type ServersService interface {
	List(ctx context.Context, projectID int, opts *GetOptions) ([]Server, *Response, error)
	Get(ctx context.Context, serverID int, opts *GetOptions) (Server, *Response, error)
	PowerOff(ctx context.Context, serverID int) (Server, *Response, error)
	PowerOn(ctx context.Context, serverID int) (Server, *Response, error)
	Create(ctx context.Context, request *CreateServer) (Server, *Response, error)
	Delete(ctx context.Context, serverID int) (*Response, error)
	PowerState(ctx context.Context, serverID int) (PowerState, *Response, error)
	Reboot(ctx context.Context, serverID int) (Server, *Response, error)
	EnterRescueMode(ctx context.Context, serverID int, fields *RescueServerFields) (Server, *Response, error)
	ExitRescueMode(ctx context.Context, serverID int) (Server, *Response, error)
	Update(ctx context.Context, serverID int, request *UpdateServer) (Server, *Response, error)
	Reinstall(ctx context.Context, serverID int, fields *ReinstallServerFields) (Server, *Response, error)
	ListSSHKeys(ctx context.Context, serverID int, opts *GetOptions) ([]SSHKey, *Response, error)
	ResetBMCPassword(ctx context.Context, serverID int) (Server, *Response, error)
	ListCycles(ctx context.Context, opts *GetOptions) ([]ServerCycle, *Response, error)
	Upgrade(ctx context.Context, serverID int, plan string) (Server, *Response, error)
}

// Server response object
type Server struct {
	ID               int               `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Href             string            `json:"href,omitempty"`
	BMC              BMC               `json:"bmc,omitempty"`
	Hostname         string            `json:"hostname,omitempty"`
	Username         string            `json:"username,omitempty"`
	Password         string            `json:"password"`
	Image            string            `json:"image,omitempty"`
	DeployedImage    DeployedImage     `json:"deployed_image,omitempty"`
	SpotInstance     bool              `json:"spot_instance"`
	BGP              ServerBGP         `json:"bgp,omitempty"`
	Project          Project           `json:"project,omitempty"`
	Region           Region            `json:"region,omitempty"`
	State            string            `json:"state,omitempty"`
	Status           string            `json:"status,omitempty"`
	Plan             Plan              `json:"plan,omitempty"`
	AvailableRegions AvailableRegions  `json:"availableregions,omitempty"`
	Pricing          Pricing           `json:"pricing,omitempty"`
	IPAddresses      []IPAddress       `json:"ip_addresses,omitempty"`
	SSHKeys          []SSHKey          `json:"ssh_keys,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	Storage          BlockStorage      `json:"storage,omitempty"`
	Backup           BackupStorage     `json:"backup_storage,omitempty"`
	Created          string            `json:"created_at,omitempty"`
	TerminationDate  string            `json:"termination_date,omitempty"`
}

// BMC data.
type BMC struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// DeployedImage data.
type DeployedImage struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type reinstallRequest struct {
	ServerAction
	*ReinstallServerFields
}

// ReinstallServerFields holds the fields for a server reinstall request.
type ReinstallServerFields struct {
	Image           string   `json:"image"`
	Hostname        string   `json:"hostname,omitempty"`
	Password        string   `json:"password"`
	IPXE            string   `json:"ipxe,omitempty"`
	SSHKeys         []string `json:"ssh_keys,omitempty"`
	UserData        string   `json:"user_data,omitempty"`
	OSPartitionSize int      `json:"os_partition_size,omitempty"`
}

type rescueServer struct {
	ServerAction
	*RescueServerFields
}

// RescueServerFields holds the fields for a server rescue request.
type RescueServerFields struct {
	Password string `json:"password"`
}

// ServerAction fields for performed action on server
type ServerAction struct {
	Type string `json:"type"`
}

// PowerState fields
type PowerState struct {
	Power string `json:"power"`
}

// UpgradeServer action request body.
type UpgradeServer struct {
	ServerAction
	Plan string `json:"plan"`
}

// CreateServer fields for ordering new server
type CreateServer struct {
	ProjectID       int                `json:"project_id"`
	Plan            string             `json:"plan"`
	Hostname        string             `json:"hostname,omitempty"`
	Image           string             `json:"image,omitempty"`
	Region          string             `json:"region"`
	SSHKeys         []string           `json:"ssh_keys,omitempty"`
	IPAddresses     []string           `json:"ip_addresses,omitempty"`
	UserData        string             `json:"user_data,omitempty"`
	Tags            *map[string]string `json:"tags,omitempty"`
	SpotInstance    bool               `json:"spot_market"`
	OSPartitionSize int                `json:"os_partition_size,omitempty"`
	IPXE            string             `json:"ipxe,omitempty"`
	StorageID       int                `json:"storage_id,omitempty"`
	Cycle           string             `json:"cycle,omitempty"`
	DiscountCode    string             `json:"discount,omitempty"`
}

// UpdateServer fields for updating a server with specified tags
type UpdateServer struct {
	Name     string             `json:"name,omitempty"`
	Hostname string             `json:"hostname,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
	BGP      *bool              `json:"bgp,omitempty"`
}

// ServerCycle data.
type ServerCycle struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ServersClient makes server related API requests.
type ServersClient struct {
	client *Client
}

// List func lists teams
func (s *ServersClient) List(ctx context.Context, projectID int, opts *GetOptions) ([]Server, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/servers", projectID))
	var trans []Server

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Get server.
func (s *ServersClient) Get(ctx context.Context, serverID int, opts *GetOptions) (Server, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseServerPath, serverID))
	var trans Server

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

func (s *ServersClient) action(ctx context.Context, serverID int, serverAction ServerAction) (Server, *Response, error) {
	var trans Server
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, serverAction)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// PowerOff function turns server off
func (s *ServersClient) PowerOff(ctx context.Context, serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "power_off",
	}

	return s.action(ctx, serverID, action)
}

// PowerOn function turns server on
func (s *ServersClient) PowerOn(ctx context.Context, serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "power_on",
	}

	return s.action(ctx, serverID, action)
}

// Reboot function restarts desired server
func (s *ServersClient) Reboot(ctx context.Context, serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "reboot",
	}

	return s.action(ctx, serverID, action)
}

// EnterRescueMode on server.
func (s *ServersClient) EnterRescueMode(ctx context.Context, serverID int, fields *RescueServerFields) (Server, *Response, error) {
	var trans Server
	request := &rescueServer{ServerAction{Type: "enter-rescue-mode"}, fields}
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// ExitRescueMode on server.
func (s *ServersClient) ExitRescueMode(ctx context.Context, serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "exit-rescue-mode",
	}

	return s.action(ctx, serverID, action)
}

// ResetBMCPassword for bare metal server.
func (s *ServersClient) ResetBMCPassword(ctx context.Context, serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "reset-bmc-password",
	}

	return s.action(ctx, serverID, action)
}

// Reinstall server OS.
func (s *ServersClient) Reinstall(ctx context.Context, serverID int, fields *ReinstallServerFields) (Server, *Response, error) {
	var trans Server
	request := &reinstallRequest{ServerAction{Type: "reinstall"}, fields}
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Upgrade virtual server plan.
func (s *ServersClient) Upgrade(ctx context.Context, serverID int, plan string) (Server, *Response, error) {
	var trans Server
	request := &UpgradeServer{
		ServerAction: ServerAction{Type: "upgrade"},
		Plan:         plan,
	}
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// PowerState retrieves server power state.
func (s *ServersClient) PowerState(ctx context.Context, serverID int) (PowerState, *Response, error) {
	path := fmt.Sprintf("%s/%d?fields=power", baseServerPath, serverID)
	var trans PowerState

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return PowerState{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Create server.
func (s *ServersClient) Create(ctx context.Context, request *CreateServer) (Server, *Response, error) {
	var trans Server
	path := fmt.Sprintf("/v1/projects/%d/servers", request.ProjectID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Update server.
func (s *ServersClient) Update(ctx context.Context, serverID int, request *UpdateServer) (Server, *Response, error) {
	var trans Server
	path := fmt.Sprintf("%s/%d", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Delete server.
func (s *ServersClient) Delete(ctx context.Context, serverID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}

// ListSSHKeys list SSH keys assigned to the server.
func (s *ServersClient) ListSSHKeys(ctx context.Context, serverID int, opts *GetOptions) ([]SSHKey, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/ssh-keys", baseServerPath, serverID))
	var trans []SSHKey

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// ListCycles lists available billing cycles.
func (s *ServersClient) ListCycles(ctx context.Context, opts *GetOptions) ([]ServerCycle, *Response, error) {
	path := opts.WithQuery("cycles")
	var trans []ServerCycle

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}
