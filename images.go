package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseImagePath = "/v1/plans"

// ImagesService is an interface for interfacing with the the Images endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Images
type ImagesService interface {
	List(ctx context.Context, plan string, opts *GetOptions) ([]Image, *Response, error)
}

type Image struct {
	ID      int       `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Pricing []Pricing `json:"pricing,omitempty"`
}

type ImagesClient struct {
	client *Client
}

// List func lists images
func (i *ImagesClient) List(ctx context.Context, plan string, opts *GetOptions) ([]Image, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s/images", baseImagePath, plan))
	var trans []Image

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
