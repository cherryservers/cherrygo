package cherrygo

import (
	"fmt"
)

type GetStorages interface {
	List(projectID string, opts *GetOptions) ([]BlockStorage, *Response, error)
}

type StoragesClient struct {
	client *Client
}

func (c *StoragesClient) List(projectID string, opts *GetOptions) ([]BlockStorage, *Response, error) {
	path := fmt.Sprintf(baseStoragePath, projectID)
	pathQuery := opts.WithQuery(path)

	var trans []BlockStorage
	resp, err := c.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
