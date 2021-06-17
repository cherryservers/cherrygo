package cherrygo

import (
	"fmt"
	"strings"
)

const baseStoragePath = "/v1/projects/%s/storages"

type GetStorage interface {
	List(storageID string, projectID string) (Storage, *Response, error)
	Create(request *CreateStorage) (Storage, *Response, error)
	Delete(request *DeleteStorage) (Storage, *Response, error)
	Attach(request *AttachTo) (*Response, error)
	Detach(*DetachFrom) (*Response, error)
}

type Storage struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Href          string `json:"href"`
	Size          int    `json:"size"`
	AllowEditSize bool   `json:"allow_edit_size"`
	Unit          string `json:"unit"`
	Description   string `json:"description"`
	VlanID        string `json:"vlan_id"`
	VlanIP        string `json:"vlan_ip"`
	Initiator     string `json:"initiator"`
	DiscoveryIP   string `json:"discovery_ip"`
}

type StorageClient struct {
	client *Client
}

type CreateStorage struct {
	ProjectID   string `json:"project_id"`
	PlanID      string `json:"plan_id"`
	Description string `json:"description"`
	Size        string `json:"size"`
	Region      string `json:"region"`
}

type DeleteStorage struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
}

type AttachTo struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
	AttachTo  string `json:"attach_to"`
}

type DetachFrom struct {
	ProjectID string `json:"project_id"`
	StorageID string `json:"storage_id"`
}

func (s *StorageClient) List(storageID, projectID string) (Storage, *Response, error) {
	serverPath := strings.Join([]string{fmt.Sprintf(baseStoragePath, projectID), storageID}, "/")

	var trans Storage

	resp, err := s.client.MakeRequest("GET", serverPath, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *StorageClient) Create(request *CreateStorage) (Storage, *Response, error) {
	var trans Storage

	serverPath := fmt.Sprintf(baseStoragePath, request.ProjectID)

	resp, err := s.client.MakeRequest("POST", serverPath, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return trans, resp, err

}

func (s *StorageClient) Delete(request *DeleteStorage) (Storage, *Response, error) {
	var trans Storage

	serverPath := strings.Join([]string{fmt.Sprintf(baseStoragePath, request.ProjectID), request.StorageID}, "/")

	resp, err := s.client.MakeRequest("DELETE", serverPath, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return trans, resp, err
}

func (s *StorageClient) Attach(request *AttachTo) (*Response, error) {
	panic("not implement yet")
}
func (s *StorageClient) Detach(*DetachFrom) (*Response, error) {
	panic("not implement yet")
}
