package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseBackupPath = "/v1/backup-storages"

// BackupsService is an interface for interfacing with the the Backup Storage endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Backup-Storage
type BackupsService interface {
	ListPlans(ctx context.Context, opts *GetOptions) ([]BackupStoragePlan, *Response, error)
	ListBackups(ctx context.Context, projectID int, opts *GetOptions) ([]BackupStorage, *Response, error)
	Get(ctx context.Context, backupID int, opts *GetOptions) (BackupStorage, *Response, error)
	Create(ctx context.Context, request *CreateBackup) (BackupStorage, *Response, error)
	Update(ctx context.Context, request *UpdateBackupStorage) (BackupStorage, *Response, error)
	UpdateBackupMethod(ctx context.Context, request *UpdateBackupMethod) ([]BackupMethod, *Response, error)
	Delete(ctx context.Context, backupID int) (*Response, error)
}

// BackupsClient makes backup storage related API requests.
type BackupsClient struct {
	client *Client
}

// BackupStoragePlan data.
type BackupStoragePlan struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Slug          string    `json:"slug,omitempty"`
	SizeGigabytes int       `json:"size_gigabytes,omitempty"`
	Pricing       []Pricing `json:"pricing,omitempty"`
	Regions       []Region  `json:"regions,omitempty"`
	Href          string    `json:"href,omitempty"`
}

// BackupStorage data.
type BackupStorage struct {
	ID                   int            `json:"id,omitempty"`
	Status               string         `json:"status,omitempty"`
	State                string         `json:"state,omitempty"`
	PrivateIP            string         `json:"private_ip,omitempty"`
	PublicIP             string         `json:"public_ip,omitempty"`
	SizeGigabytes        int            `json:"size_gigabytes,omitempty"`
	UsedGigabytes        int            `json:"used_gigabytes,omitempty"`
	AttachedTo           AttachedTo     `json:"attached_to,omitempty"`
	Methods              []BackupMethod `json:"methods,omitempty"`
	AvailableIPAddresses []IPAddress    `json:"available_addresses,omitempty"`
	Rules                []Rule         `json:"rules,omitempty"`
	Plan                 Plan           `json:"plan,omitempty"`
	Pricing              Pricing        `json:"pricing,omitempty"`
	Region               Region         `json:"region,omitempty"`
	Href                 string         `json:"href,omitempty"`
}

// BackupMethod is a backup storage access method.
type BackupMethod struct {
	Name       string   `json:"name,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Port       int      `json:"port,omitempty"`
	Host       string   `json:"host,omitempty"`
	SSHKey     string   `json:"ssh_key,omitempty"`
	WhiteList  []string `json:"whitelist,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
	Processing bool     `json:"processing,omitempty"`
}

// Rule is the backup storage access method rule for an IP address.
type Rule struct {
	IPAddress      IPAddress      `json:"ip,omitempty"`
	EnabledMethods EnabledMethods `json:"methods,omitempty"`
}

// EnabledMethods for backup storage access.
type EnabledMethods struct {
	BORG bool `json:"borg,omitempty"`
	FTP  bool `json:"ftp,omitempty"`
	NFS  bool `json:"nfs,omitempty"`
	SMB  bool `json:"smb,omitempty"`
}

// CreateBackup is the body for backup storage creation request.
type CreateBackup struct {
	ServerID       int    `json:"server_id,omitempty"`
	BackupPlanSlug string `json:"slug"`
	RegionSlug     string `json:"region"`
	SSHKey         string `json:"ssh_key,omitempty"`
}

// UpdateBackupStorage is the body for a backup storage update request.
type UpdateBackupStorage struct {
	BackupStorageID int    `json:"id"`
	BackupPlanSlug  string `json:"slug,omitempty"`
	Password        string `json:"password,omitempty"`
	SSHKey          string `json:"ssh_key,omitempty"`
}

// UpdateBackupMethod is the body for a backup storage access method update request.
type UpdateBackupMethod struct {
	BackupStorageID  int      `json:"id"`
	BackupMethodName string   `json:"name"`
	Enabled          bool     `json:"enabled"`
	Whitelist        []string `json:"whitelist"`
}

// ListPlans lists backups storage plans.
func (s *BackupsClient) ListPlans(ctx context.Context, opts *GetOptions) ([]BackupStoragePlan, *Response, error) {
	var trans []BackupStoragePlan

	path := opts.WithQuery("/v1/backup-storage-plans")
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// ListBackups lists backup storage instances.
func (s *BackupsClient) ListBackups(ctx context.Context, projectID int, opts *GetOptions) ([]BackupStorage, *Response, error) {
	var trans []BackupStorage

	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/backup-storages", projectID))
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Get backup storage instance.
func (s *BackupsClient) Get(ctx context.Context, backupID int, opts *GetOptions) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseBackupPath, backupID))
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return BackupStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Create backup storage instance.
func (s *BackupsClient) Create(ctx context.Context, request *CreateBackup) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := fmt.Sprintf("/v1/servers/%d/backup-storages", request.ServerID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return BackupStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Update backup storage instance.
func (s *BackupsClient) Update(ctx context.Context, request *UpdateBackupStorage) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := fmt.Sprintf("%s/%d", baseBackupPath, request.BackupStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return BackupStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// UpdateBackupMethod updates backup storage instance access methods.
func (s *BackupsClient) UpdateBackupMethod(ctx context.Context, request *UpdateBackupMethod) ([]BackupMethod, *Response, error) {
	var trans []BackupMethod

	path := fmt.Sprintf("%s/%d/methods/%s", baseBackupPath, request.BackupStorageID, request.BackupMethodName)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, request)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Delete backup storage instance.
func (s *BackupsClient) Delete(ctx context.Context, backupID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseBackupPath, backupID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}
