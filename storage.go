package cherrygo

import (
	"fmt"
	"strings"
)

const baseStoragePath = "/v1/projects/%s/storages"

type GetStorage interface {
	List(projectID string, storageID string, opts *GetOptions) (BlockStorage, *Response, error)
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
	ProjectID   string `json:"project_id"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Region      string `json:"region"`
}

type DeleteStorage struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
}

type AttachTo struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
	AttachTo  int    `json:"attach_to"`
}

type AttachedTo struct {
	Href string `json:"href"`
}

type DetachFrom struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
}

func (s *StorageClient) List(projectID string, storageID string, opts *GetOptions) (BlockStorage, *Response, error) {
	path := strings.Join([]string{fmt.Sprintf(baseStoragePath, projectID), storageID}, "/")
	pathQuery := opts.WithQuery(path)

	var trans BlockStorage

	resp, err := s.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *StorageClient) Create(request *CreateStorage) (BlockStorage, *Response, error) {
	var trans BlockStorage

	serverPath := fmt.Sprintf(baseStoragePath, request.ProjectID)

	resp, err := s.client.MakeRequest("POST", serverPath, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return trans, resp, err

}

func (s *StorageClient) Delete(request *DeleteStorage) (*Response, error) {

	path := strings.Join([]string{fmt.Sprintf(baseStoragePath, request.ProjectID), request.StorageID}, "/")

	resp, err := s.client.MakeRequest("DELETE", path, request, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return resp, err
}

func (s *StorageClient) Attach(request *AttachTo) (BlockStorage, *Response, error) {
	var trans BlockStorage

	path := strings.Join([]string{fmt.Sprintf(baseStoragePath, request.ProjectID), request.StorageID, "attachments"}, "/")

	resp, err := s.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return trans, resp, err
}

func (s *StorageClient) Detach(request *DetachFrom) (*Response, error) {
	path := strings.Join([]string{fmt.Sprintf(baseStoragePath, request.ProjectID), request.StorageID, "attachments"}, "/")

	resp, err := s.client.MakeRequest("DELETE", path, nil, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return resp, err
}
