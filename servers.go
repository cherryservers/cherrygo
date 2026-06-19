package cherrygo

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/netip"
	"time"
)

const (
	baseServerPath = "/v1/servers"
	endServersPath = "servers"
)

// ServerStatus is the pseudo-enum for server statuses.
//
// Used as a constraint for server status parameters.
// It is not used as a field type in response structs and
// makes no guarantees about enumerating all possible status values.
type ServerStatus int

const (
	// StatusDeployed status is used to indicate an active server
	// deployment. This is generally the status to watch for when
	// provisioning a new server, except when using custom installation
	// procedures, in which case, see `Allocated`.
	StatusDeployed ServerStatus = iota

	// StatusAllocated status is used to indicate an active server
	// deployment, when custom installation procedures are used, i.e.
	// iPXE, since these don't go through the full standard deployment
	// process.
	StatusAllocated
)

var serverStatusValues = map[ServerStatus]string{
	StatusDeployed:  "deployed",
	StatusAllocated: "allocated",
}

func (ss ServerStatus) String() string {
	if v, ok := serverStatusValues[ss]; ok {
		return v
	}

	return fmt.Sprintf("invalid ServerStatus(%d)", ss)
}

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
	AllowBMCAccess(ctx context.Context, serverID int, ip4 string) (Server, *Response, error)
	WaitForStatus(ctx context.Context, serverID int, status ServerStatus) (Server, *Response, error)
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

	// IP is the address at which the BMC can be reached.
	IP netip.Addr `json:"ip,omitzero"`

	// AllowedIP is the address that is whitelisted for BMC access.
	AllowedIP netip.Addr `json:"allowed_ip,omitzero"`

	Expires time.Time `json:"expires,omitzero"`
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

type allowBMCAccess struct {
	ServerAction
	AllowedIP string `json:"allowed_ip,omitempty"`
}

// CreateServer fields for ordering new server
type CreateServer struct {
	ProjectID int    `json:"project_id"`
	Plan      string `json:"plan"`

	// PrebuiltID allows selecting a pre-assembled plan variant.
	// Requires Plan to be set as well.
	PrebuiltID int `json:"prebuilt_id,omitempty"`

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

// AllowBMCAccess allows BMC/IPMI access from the specified IPv4 address for a limited duration.
// If ip4 is empty, no whitelist will be used, i.e. all addresses will be allowed.
func (s *ServersClient) AllowBMCAccess(ctx context.Context, serverID int, ip4 string) (Server, *Response, error) {
	var srv Server
	body := &allowBMCAccess{
		ServerAction: ServerAction{Type: "create-console-access"},
		AllowedIP:    ip4,
	}
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return Server{}, nil, err
	}

	resp, err := s.client.Do(req, &srv)
	return srv, resp, err
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

// WaitForStatus blocks until server reaches specified status.
func (s *ServersClient) WaitForStatus(ctx context.Context, serverID int, status ServerStatus) (Server, *Response, error) {
	if s.client.pollBackoff == nil {
		return Server{}, nil, errors.New("nil client pollBackoff function")
	}

	attempt := 0
	for {
		server, resp, err := s.Get(ctx, serverID, nil)
		if err != nil {
			return Server{}, nil, err
		}

		if server.Status == status.String() {
			return server, resp, nil
		}

		if resp == nil {
			return Server{}, nil, fmt.Errorf("nil response from GET server %d", serverID)
		}

		select {
		case <-time.After(s.client.pollBackoff(attempt, resp.Response)):
			attempt++
		case <-ctx.Done():
			return Server{}, nil, ctx.Err()
		}
	}
}

// GeneratePassword generates a password that matches Cherry Servers secure password
// criteria in a cryptographically secure way.
func GeneratePassword() (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		all       = lowercase + uppercase + digits
		length    = 20
	)
	password := make([]byte, length)

	var charset string
	for i := range length {
		switch i {
		case 0:
			// Ensure there is at least one lower-case alphabetical character.
			charset = lowercase
		case 1:
			// Ensure there is at least one upper-case alphabetical
			// character that is not first.
			charset = uppercase
		case 2:
			// Ensure there is at least one digit that is not last.
			charset = digits
		default:
			charset = all
		}
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[idx.Int64()]
	}
	return string(password), nil
}
