package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseRegionPath = "/v1/regions"

// RegionsService is an interface for interfacing with the the Images endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Regions
type RegionsService interface {
	List(ctx context.Context, opts *GetOptions) ([]Region, *Response, error)
	Get(ctx context.Context, region string, opts *GetOptions) (Region, *Response, error)
}

// Region fields
type Region struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Slug       string    `json:"slug,omitempty"`
	RegionIso2 string    `json:"region_iso_2,omitempty"`
	BGP        RegionBGP `json:"bgp,omitempty"`
	Location   string    `json:"location,omitempty"`
	Href       string    `json:"href,omitempty"`
}

type RegionsClient struct {
	client *Client
}

func (i *RegionsClient) List(ctx context.Context, opts *GetOptions) ([]Region, *Response, error) {
	path := opts.WithQuery(baseRegionPath)
	var trans []Region

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}

func (i *RegionsClient) Get(ctx context.Context, region string, opts *GetOptions) (Region, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s", baseRegionPath, region))
	var trans Region

	req, err := i.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Region{}, nil, err
	}

	resp, err := i.client.Do(req, &trans)
	return trans, resp, err
}
