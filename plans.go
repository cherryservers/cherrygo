package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const (
	teamPlanPath = "/v1/teams"
	basePlanPath = "/v1/plans"
)

// PlansService is an interface for interfacing with the Plan endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Plans
type PlansService interface {
	List(ctx context.Context, teamID int, opts *GetOptions) ([]Plan, *Response, error)
	GetBySlug(ctx context.Context, slug string, opts *GetOptions) (Plan, *Response, error)
	GetByID(ctx context.Context, id int, opts *GetOptions) (Plan, *Response, error)
}

// Plan data.
type Plan struct {
	ID               int                `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Slug             string             `json:"slug,omitempty"`
	Custom           bool               `json:"custom,omitempty"`
	Type             string             `json:"type,omitempty"`
	Specs            Specs              `json:"specs,omitempty"`
	Pricing          []Pricing          `json:"pricing,omitempty"`
	AvailableRegions []AvailableRegions `json:"available_regions,omitempty"`
	Category         string             `json:"category"`
	Softwares        []SoftwareImage    `json:"softwares"`
}

// SoftwareImage data.
type SoftwareImage struct {
	Image SoftwareImageSpecs `json:"image"`
}

// SoftwareImageSpecs data.
type SoftwareImageSpecs struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Specs specifies fields for specs
type Specs struct {
	Cpus      Cpus      `json:"cpus,omitempty"`
	Memory    Memory    `json:"memory,omitempty"`
	Storage   []Storage `json:"storage,omitempty"`
	Raid      Raid      `json:"raid,omitempty"`
	Nics      Nics      `json:"nics,omitempty"`
	Bandwidth Bandwidth `json:"bandwidth,omitempty"`
}

// Cpus fields
type Cpus struct {
	Count     int     `json:"count,omitempty"`
	Name      string  `json:"name,omitempty"`
	Cores     int     `json:"cores,omitempty"`
	Frequency float32 `json:"frequency,omitempty"`
	Unit      string  `json:"unit,omitempty"`
}

// Memory fields
type Memory struct {
	Count int    `json:"count,omitempty"`
	Total int    `json:"total,omitempty"`
	Unit  string `json:"unit,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Storage fields
type Storage struct {
	Count int     `json:"count,omitempty"`
	Name  string  `json:"name,omitempty"`
	Size  float32 `json:"size,omitempty"`
	Unit  string  `json:"unit,omitempty"`
}

// Raid fields
type Raid struct {
	Name string `json:"name,omitempty"`
}

// Nics fields
type Nics struct {
	Name string `json:"name,omitempty"`
}

// Bandwidth fields
type Bandwidth struct {
	Name string `json:"name,omitempty"`
}

// AvailableRegions data.
type AvailableRegions struct {
	*Region
	StockQty int `json:"stock_qty,omitempty"`
	SpotQty  int `json:"spot_qty,omitempty"`
}

// PlansClient makes plan related API requests.
type PlansClient struct {
	client *Client
}

// List func lists plans
func (p *PlansClient) List(ctx context.Context, teamID int, opts *GetOptions) ([]Plan, *Response, error) {
	basePath := basePlanPath
	if teamID != 0 {
		basePath = fmt.Sprintf("%s/%d/plans", teamPlanPath, teamID)
	}

	path := opts.WithQuery(basePath)
	var trans []Plan

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

func (p *PlansClient) get(ctx context.Context, path string) (Plan, *Response, error) {
	var trans Plan

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Plan{}, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

// GetByID retrieves server plan by ID.
func (p *PlansClient) GetByID(ctx context.Context, id int, opts *GetOptions) (Plan, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", basePlanPath, id))

	return p.get(ctx, path)
}

// GetBySlug retrieves server plan by slug.
func (p *PlansClient) GetBySlug(ctx context.Context, slug string, opts *GetOptions) (Plan, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s", basePlanPath, slug))

	return p.get(ctx, path)
}
