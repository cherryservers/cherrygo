package cherrygo

import (
	"fmt"
)

type GetStorages interface {
	List(projectID string) ([]BlockStorage, *Response, error)
}

type StoragesClient struct {
	client *Client
}

func (c *StoragesClient) List(projectID string) ([]BlockStorage, *Response, error) {
	storagePath := fmt.Sprintf(baseStoragePath, projectID)

	var trans []BlockStorage
	resp, err := c.client.MakeRequest("GET", storagePath, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
