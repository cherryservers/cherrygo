package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const teamsPath = "/v1/teams"

// TeamsService is an interface for interfacing with the Teams endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Teams
type TeamsService interface {
	List(ctx context.Context, opts *GetOptions) ([]Team, *Response, error)
	Get(ctx context.Context, teamID int, opts *GetOptions) (Team, *Response, error)
	Create(ctx context.Context, request *CreateTeam) (Team, *Response, error)
	Update(ctx context.Context, teamID int, request *UpdateTeam) (Team, *Response, error)
	Delete(ctx context.Context, teamID int) (*Response, error)
}

type Team struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Credit      Credit       `json:"credit,omitempty"`
	Billing     Billing      `json:"billing,omitempty"`
	Projects    []Project    `json:"projects,omitempty"`
	Memberships []Membership `json:"memberships,omitempty"`
	Href        string       `json:"href,omitempty"`
}

type Credit struct {
	Account   CreditDetails `json:"account,omitempty"`
	Promo     CreditDetails `json:"promo,omitempty"`
	Resources Resources     `json:"resources,omitempty"`
}

type CreditDetails struct {
	Remaining float32 `json:"remaining,omitempty"`
	Usage     float32 `json:"usage,omitempty"`
	Currency  string  `json:"currency,omitempty"`
}

type Resources struct {
	Pricing   Pricing       `json:"pricing,omitempty"`
	Remaining RemainingTime `json:"remaining,omitempty"`
}

type RemainingTime struct {
	Time int    `json:"time,omitempty"`
	Unit string `json:"unit,omitempty"`
}

type Pricing struct {
	Price     float32 `json:"price,omitempty"`
	UnitPrice float32 `json:"unit_price"`
	Taxed     bool    `json:"taxed,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	Unit      string  `json:"unit,omitempty"`
}

type Billing struct {
	Type        string `json:"type,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	CompanyCode string `json:"company_code,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Address1    string `json:"address_1,omitempty"`
	Address2    string `json:"address_2,omitempty"`
	CountryIso2 string `json:"country_iso_2,omitempty"`
	City        string `json:"city,omitempty"`
	Vat         Vat    `json:"vat,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

type Vat struct {
	Amount int    `json:"amount"`
	Number string `json:"number,omitempty"`
	Valid  bool   `json:"valid"`
}

type TeamsClient struct {
	client *Client
}

type CreateTeam struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type UpdateTeam struct {
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

// List func lists teams
func (c *TeamsClient) List(ctx context.Context, opts *GetOptions) ([]Team, *Response, error) {
	var trans []Team
	pathQuery := opts.WithQuery(teamsPath)

	req, err := c.client.NewRequest(ctx, http.MethodGet, pathQuery, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(req, &trans)
	return trans, resp, err
}

func (c *TeamsClient) Get(ctx context.Context, teamID int, opts *GetOptions) (Team, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", teamsPath, teamID))
	var trans Team

	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Team{}, nil, err
	}

	resp, err := c.client.Do(req, &trans)
	return trans, resp, err
}

func (c *TeamsClient) Create(ctx context.Context, request *CreateTeam) (Team, *Response, error) {
	path := teamsPath
	var trans Team

	req, err := c.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Team{}, nil, err
	}

	resp, err := c.client.Do(req, &trans)
	return trans, resp, err
}

func (c *TeamsClient) Update(ctx context.Context, teamID int, request *UpdateTeam) (Team, *Response, error) {
	path := fmt.Sprintf("%s/%d", teamsPath, teamID)
	var trans Team

	req, err := c.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return Team{}, nil, err
	}

	resp, err := c.client.Do(req, &trans)
	return trans, resp, err
}

func (c *TeamsClient) Delete(ctx context.Context, teamID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", teamsPath, teamID)

	req, err := c.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req, nil)
	return resp, err
}
