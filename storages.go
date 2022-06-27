package cherrygo

import (
	"fmt"
)

// StoragesService is an interface for interfacing with the Storages endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Storage
type StoragesService interface {
	List(projectID int, opts *GetOptions) ([]BlockStorage, *Response, error)
	Get(projectID int, storageID int, opts *GetOptions) (BlockStorage, *Response, error)
	Create(request *CreateStorage) (BlockStorage, *Response, error)
	Delete(request *DeleteStorage) (*Response, error)
	Attach(request *AttachTo) (BlockStorage, *Response, error)
	Detach(*DetachFrom) (*Response, error)
}

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
}

type StorageClient struct {
	client *Client
}

type CreateStorage struct {
	ProjectID   int    `json:"project_id"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Region      string `json:"region"`
}

type DeleteStorage struct {
	ProjectID int `json:"project_id"`
	StorageID int `json:"storage_id"`
}

type AttachTo struct {
	ProjectID int `json:"project_id"`
	StorageID int `json:"storage_id"`
	AttachTo  int `json:"attach_to"`
}

type AttachedTo struct {
	Href string `json:"href"`
}

type DetachFrom struct {
	ProjectID int `json:"project_id"`
	StorageID int `json:"storage_id"`
}

type StoragesClient struct {
	client *Client
}

func (c *StoragesClient) List(projectID int, opts *GetOptions) ([]BlockStorage, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/storages", baseProjectPath, projectID))

	var trans []BlockStorage
	resp, err := c.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *StoragesClient) Get(projectID int, storageID int, opts *GetOptions) (BlockStorage, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/storages/%d", baseProjectPath, projectID, storageID))

	var trans BlockStorage

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *StoragesClient) Create(request *CreateStorage) (BlockStorage, *Response, error) {
	var trans BlockStorage

	path := fmt.Sprintf("%s/%d/storages", baseProjectPath, request.ProjectID)

	resp, err := s.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err

}

func (s *StoragesClient) Delete(request *DeleteStorage) (*Response, error) {
	path := fmt.Sprintf("%s/%d/storages/%d", baseProjectPath, request.ProjectID, request.StorageID)

	resp, err := s.client.MakeRequest("DELETE", path, request, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return resp, err
}

func (s *StoragesClient) Attach(request *AttachTo) (BlockStorage, *Response, error) {
	var trans BlockStorage

	path := fmt.Sprintf("%s/%d/storages/%d/attachments", baseProjectPath, request.ProjectID, request.StorageID)

	resp, err := s.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *StoragesClient) Detach(request *DetachFrom) (*Response, error) {
	path := fmt.Sprintf("%s/%d/storages/%d/attachments", baseProjectPath, request.ProjectID, request.StorageID)

	resp, err := s.client.MakeRequest("DELETE", path, nil, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return resp, err
}
