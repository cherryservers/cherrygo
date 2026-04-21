package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseStoragePath = "/v1/storages"

// StoragesService is an interface for interfacing with the Storages endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Storage
type StoragesService interface {
	List(ctx context.Context, projectID int, opts *GetOptions) ([]BlockStorage, *Response, error)
	Get(ctx context.Context, storageID int, opts *GetOptions) (BlockStorage, *Response, error)
	Create(ctx context.Context, request *CreateStorage) (BlockStorage, *Response, error)
	Delete(ctx context.Context, storageID int) (*Response, error)
	Attach(ctx context.Context, request *AttachTo) (BlockStorage, *Response, error)
	Detach(ctx context.Context, storageID int) (*Response, error)
	Update(ctx context.Context, request *UpdateStorage) (BlockStorage, *Response, error)
}

// BlockStorage data.
type BlockStorage struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Href          string     `json:"href"`
	Size          int        `json:"size"`
	AllowEditSize bool       `json:"allow_edit_size"`
	Unit          string     `json:"unit"`
	Description   string     `json:"description,omitempty"`
	AttachedTo    AttachedTo `json:"attached_to,omitempty"`
	VlanID        string     `json:"vlan_id"`
	VlanIP        string     `json:"vlan_ip"`
	Initiator     string     `json:"initiator"`
	DiscoveryIP   string     `json:"discovery_ip"`
	Region        Region     `json:"region"`
}

// CreateStorage is the storage creation request body.
type CreateStorage struct {
	ProjectID   int    `json:"project_id"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Region      string `json:"region"`
}

// AttachTo is the storage attachment request body data.
type AttachTo struct {
	StorageID int `json:"storage_id"`
	AttachTo  int `json:"attach_to"`
}

// AttachedTo is the data of the instance the storage is attached to.
type AttachedTo struct {
	ID       int    `json:"id"`
	Hostname string `json:"hostname,omitempty"`
	Href     string `json:"href"`
}

// UpdateStorage is the request body for updating storage instances.
type UpdateStorage struct {
	StorageID   int    `json:"storage_id"`
	Size        int    `json:"size"`
	Description string `json:"description,omitempty"`
}

// StoragesClient makes storage related API requests.
type StoragesClient struct {
	client *Client
}

// List all project storages.
func (s *StoragesClient) List(ctx context.Context, projectID int, opts *GetOptions) ([]BlockStorage, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/storages", baseProjectPath, projectID))
	var trans []BlockStorage

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Get storage instance.
func (s *StoragesClient) Get(ctx context.Context, storageID int, opts *GetOptions) (BlockStorage, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseStoragePath, storageID))
	var trans BlockStorage

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return BlockStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Create storage instance.
func (s *StoragesClient) Create(ctx context.Context, request *CreateStorage) (BlockStorage, *Response, error) {
	var trans BlockStorage
	path := fmt.Sprintf("%s/%d/storages", baseProjectPath, request.ProjectID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return BlockStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Delete storage.
func (s *StoragesClient) Delete(ctx context.Context, storageID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseStoragePath, storageID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}

// Attach storage to server.
func (s *StoragesClient) Attach(ctx context.Context, request *AttachTo) (BlockStorage, *Response, error) {
	var trans BlockStorage
	path := fmt.Sprintf("%s/%d/attachments", baseStoragePath, request.StorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return BlockStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Detach storage from server.
func (s *StoragesClient) Detach(ctx context.Context, storageID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d/attachments", baseStoragePath, storageID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	return resp, err
}

// Update storage.
func (s *StoragesClient) Update(ctx context.Context, request *UpdateStorage) (BlockStorage, *Response, error) {
	var trans BlockStorage
	path := fmt.Sprintf("%s/%d", baseStoragePath, request.StorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return BlockStorage{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}
